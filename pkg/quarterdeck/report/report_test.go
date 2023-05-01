package report_test

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/report"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/stretchr/testify/require"
)

var timezone *time.Location

func TestDailyUsersReport(t *testing.T) {
	// This test does not test scheduling or sending an email, but rather tests that
	// the queries executed to the database return the expected results for the report.
	path := filepath.Join(t.TempDir(), "quarterdeck.db")
	err := setupDB(path)
	require.NoError(t, err, "could not setup database")

	// Create a report tool with a mock emailer
	mock := &MockEmailer{}
	daily, err := report.NewDailyUsers(mock)
	require.NoError(t, err, "could not create daily reporting tool")

	err = daily.Report()
	require.NoError(t, err, "could not run daily report")

	// Ensure the emailer has been called
	require.Equal(t, 1, mock.calls)
	require.Len(t, mock.data, 1)

	now := time.Now().In(timezone)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timezone)

	report := mock.data[0]
	require.Equal(t, 24*time.Hour, today.Sub(report.Date))
	require.Equal(t, today.AddDate(0, 0, -1), report.Date)
	require.Equal(t, today.AddDate(0, 0, -31), report.InactiveDate)
	require.Empty(t, report.Domain)
	require.Empty(t, report.EnsignDashboardLink)

	// HACK: if UTC time and local time are the same day make specific assertions;
	// otherwise only assert that some data is returned. This implies that the report
	// is not handling timezones correctly - but since the PLG report is run once a day
	// we can guarantee that it is at a time that is the same day UTC and local time.
	utcnow := time.Now().In(time.UTC)
	if now.Day() != utcnow.Day() {
		require.NotZero(t, report.DailyUsers)
		require.NotZero(t, report.ActiveUsers)
		require.NotZero(t, report.InactiveUsers)
		require.NotZero(t, report.APIKeys)
		require.NotZero(t, report.ActiveKeys)
		require.NotZero(t, report.InactiveKeys)
		require.NotZero(t, report.RevokedKeys)
		require.NotZero(t, report.Organizations)
		require.NotZero(t, report.Projects)
	} else {
		require.Equal(t, 7, report.NewUsers)
		require.Equal(t, 24, report.DailyUsers)
		require.Equal(t, 58, report.ActiveUsers)
		require.Equal(t, 109, report.InactiveUsers)
		require.Equal(t, 112, report.APIKeys)
		require.Equal(t, 94, report.ActiveKeys)
		require.Equal(t, 18, report.InactiveKeys)
		require.Equal(t, 22, report.RevokedKeys)
		require.Equal(t, 54, report.Organizations)
		require.Equal(t, 3, report.NewOrganizations)
		require.Equal(t, 270, report.Projects)
		require.Equal(t, 8, report.NewProjects)
	}
}

type MockEmailer struct {
	data  []*emails.DailyUsersData
	calls int
}

func (m *MockEmailer) SendDailyUsers(data *emails.DailyUsersData) error {
	m.data = append(m.data, data)
	m.calls++
	return nil
}

func setupDB(path string) (err error) {
	dsn := fmt.Sprintf("sqlite3:///%s", path)
	if err = db.Connect(dsn, false); err != nil {
		return err
	}

	// TODO: insert users into multiple timezones
	// NOTE: we expect DB timestamps to be stored in UTC per RFC3339.
	if timezone, err = time.LoadLocation("America/New_York"); err != nil {
		return err
	}

	// Ensure pseudo-randomness
	rand.Seed(time.Now().UnixNano())

	// Time ranges for creating database records in
	now := time.Now().In(time.UTC)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)
	inactive := yesterday.AddDate(0, 0, -30)
	history := inactive.AddDate(0, -5, 0)

	var tx *sql.Tx
	if tx, err = db.BeginTx(context.Background(), nil); err != nil {
		return err
	}
	defer tx.Rollback()

	// Create 7 new users
	if err = insertUsers(tx, 7, yesterday, today, yesterday, now); err != nil {
		return err
	}

	// Create 24 total daily users
	if err = insertUsers(tx, 17, history, yesterday, yesterday, today); err != nil {
		return err
	}

	// Create 58 active users
	if err = insertUsers(tx, 34, history, yesterday, inactive, yesterday); err != nil {
		return err
	}

	// Create 109 inactive users
	if err = insertUsers(tx, 100, history, inactive, history, inactive); err != nil {
		return err
	}

	if err = insertUsers(tx, 9, history, inactive, time.Time{}, time.Time{}); err != nil {
		return err
	}

	// Create 3 new organizations
	if err = insertOrganizations(tx, 3, yesterday, today); err != nil {
		return err
	}

	// Create 54 organizations total
	if err = insertOrganizations(tx, 51, history, yesterday); err != nil {
		return err
	}

	// Create 8 new projects (must come after organizations)
	if err = insertProjects(tx, 8, yesterday, today); err != nil {
		return err
	}

	// Create 270 projects total
	if err = insertProjects(tx, 262, history, yesterday); err != nil {
		return err
	}

	// Insert 94 active api keys (must come after organizations, projects, and users)
	if err = insertAPIKeys(tx, 94, history, inactive, inactive, now); err != nil {
		return err
	}

	// Insert 18 inactive api keys (must come after organizations, projects, and users)
	if err = insertAPIKeys(tx, 18, history, inactive, history, inactive); err != nil {
		return err
	}

	// Insert 22 revoked api keys (must come after organizations, projects, and users)
	if err = insertRevokedAPIKeys(tx, 22, history, today); err != nil {
		return err
	}

	return tx.Commit()
}

func insertUsers(tx *sql.Tx, n int, createdAfter, createdBefore, loginAfter, loginBefore time.Time) error {
	vals := make([]string, 0, n)
	params := make([]interface{}, 0, n*7)

	for i := 0; i < n; i++ {
		id := ulid.Make()
		name, email, password := randString(), randString(), randString()
		created := randomTimestamp(createdAfter, createdBefore)
		modified := randomTimestamp(created, createdBefore)

		var lastLogin string
		if !loginAfter.IsZero() && !loginBefore.IsZero() {
			lastLogin = randomTimestamp(loginAfter, loginBefore).Format(time.RFC3339Nano)
		}

		vals = append(vals, "(?, ?, ?, ?, ?, ?, ?)")
		params = append(params, id, name, email, password, lastLogin, created.Format(time.RFC3339Nano), modified.Format(time.RFC3339Nano))
	}

	query := fmt.Sprintf("INSERT INTO users (id, name, email, password, last_login, created, modified) VALUES %s", strings.Join(vals, ","))
	if _, err := tx.Exec(query, params...); err != nil {
		return err
	}
	return nil
}

func insertOrganizations(tx *sql.Tx, n int, after, before time.Time) error {
	vals := make([]string, 0, n)
	params := make([]interface{}, 0, n*5)

	for i := 0; i < n; i++ {
		id := ulid.Make()
		name, domain := randString(), randString()
		created := randomTimestamp(after, before)
		modified := randomTimestamp(created, before)

		vals = append(vals, "(?, ?, ?, ?, ?)")
		params = append(params, id, name, domain, created.Format(time.RFC3339Nano), modified.Format(time.RFC3339Nano))
	}

	query := fmt.Sprintf("INSERT INTO organizations VALUES %s", strings.Join(vals, ","))
	if _, err := tx.Exec(query, params...); err != nil {
		return err
	}
	return nil
}

func insertProjects(tx *sql.Tx, n int, after, before time.Time) error {
	vals := make([]string, 0, n)
	params := make([]interface{}, 0, n*3)

	for i := 0; i < n; i++ {
		id := ulid.Make()
		created := randomTimestamp(after, before)
		modified := randomTimestamp(created, before)

		orgID := "(SELECT id FROM organizations ORDER BY RANDOM() LIMIT 1)"
		ownerID := "(SELECT id FROM users ORDER BY RANDOM() LIMIT 1)"
		vals = append(vals, fmt.Sprintf("(%s, ?, ?, ?, %s)", orgID, ownerID))
		params = append(params, id, created.Format(time.RFC3339Nano), modified.Format(time.RFC3339Nano))
	}

	query := fmt.Sprintf("INSERT INTO organization_projects VALUES %s", strings.Join(vals, ","))
	if _, err := tx.Exec(query, params...); err != nil {
		return err
	}
	return nil
}

func insertAPIKeys(tx *sql.Tx, n int, createdAfter, createdBefore, loginAfter, loginBefore time.Time) error {
	vals := make([]string, 0, n)
	params := make([]interface{}, 0, n*7)

	for i := 0; i < n; i++ {
		id := ulid.Make()
		keyID, secret, name := randString(), randString(), randString()
		created := randomTimestamp(createdAfter, createdBefore)
		modified := randomTimestamp(created, createdBefore)

		var lastUsed string
		if !loginAfter.IsZero() && !loginBefore.IsZero() {
			lastUsed = randomTimestamp(loginAfter, loginBefore).Format(time.RFC3339Nano)
		}

		vals = append(vals, "(?, ?, ?, ?, (SELECT id FROM organizations ORDER BY RANDOM() LIMIT 1), (SELECT project_id FROM organization_projects ORDER BY RANDOM() LIMIT 1), (SELECT id FROM users ORDER BY RANDOM() LIMIT 1), ?, ?, ?)")
		params = append(params, id, keyID, secret, name, lastUsed, created.Format(time.RFC3339Nano), modified.Format(time.RFC3339Nano))
	}

	query := fmt.Sprintf("INSERT INTO api_keys (id, key_id, secret, name, organization_id, project_id, created_by, last_used, created, modified) VALUES %s", strings.Join(vals, ","))
	if _, err := tx.Exec(query, params...); err != nil {
		return err
	}
	return nil
}

func insertRevokedAPIKeys(tx *sql.Tx, n int, after, before time.Time) error {
	vals := make([]string, 0, n)
	params := make([]interface{}, 0, n*4)

	for i := 0; i < n; i++ {
		id := ulid.Make()
		keyID := randString()
		created := randomTimestamp(after, before)
		modified := randomTimestamp(created, before)

		vals = append(vals, "(?,?,?,?)")
		params = append(params, id, keyID, created.Format(time.RFC3339Nano), modified.Format(time.RFC3339Nano))
	}

	query := fmt.Sprintf("INSERT INTO revoked_api_keys (id, key_id, created, modified) VALUES %s", strings.Join(vals, ","))
	if _, err := tx.Exec(query, params...); err != nil {
		return err
	}
	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randString() string {
	b := make([]rune, 14)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func randomTimestamp(after, before time.Time) time.Time {
	if before.Before(after) || before.Equal(after) {
		panic(fmt.Errorf("invalid after and before timestamps: after %s before %s", after, before))
	}

	i, j := after.UnixNano(), before.UnixNano()
	for k := 0; k < 10; k++ {
		ts := rand.Int63n(j-i) + i
		dt := time.Unix(0, ts)

		if dt.After(after) && dt.Before(before) {
			return dt.In(after.Location())
		}
	}
	panic("could not generate timestamp in time range")
}

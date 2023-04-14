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
	require.Equal(t, 0, report.NewUsers)
	require.Equal(t, 0, report.DailyUsers)
	require.Equal(t, 0, report.ActiveUsers)
	require.Equal(t, 0, report.InactiveUsers)
	require.Equal(t, 0, report.APIKeys)
	require.Equal(t, 0, report.ActiveKeys)
	require.Equal(t, 0, report.InactiveKeys)
	require.Equal(t, 0, report.RevokedKeys)
	require.Equal(t, 54, report.Organizations)
	require.Equal(t, 3, report.NewOrganizations)
	require.Equal(t, 0, report.Projects)
	require.Equal(t, 0, report.NewProjects)
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

	// Time ranges for creating database records in
	now := time.Now().In(timezone)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timezone)
	yesterday := today.AddDate(0, 0, -1)
	inactive := yesterday.AddDate(0, 0, -30)
	history := inactive.AddDate(0, -5, 0)

	var tx *sql.Tx
	if tx, err = db.BeginTx(context.Background(), nil); err != nil {
		return err
	}
	defer tx.Rollback()

	// Create 3 new organizations
	if err = insertOrganizations(tx, 3, yesterday, today); err != nil {
		return err
	}

	// Create 54 organizations total
	if err = insertOrganizations(tx, 51, history, yesterday); err != nil {
		return err
	}

	return tx.Commit()
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
		params = append(params, id, name, domain, created, modified)

		fmt.Println(created)
	}

	query := fmt.Sprintf("INSERT INTO organizations VALUES %s", strings.Join(vals, ","))
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
	i, j := after.UnixNano(), before.UnixNano()
	ts := rand.Int63n(j-i) + i
	return time.Unix(0, ts)
}

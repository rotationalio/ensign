package report

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog/log"
)

// The DailyUsers report is generated once per day at the specified schedule using a
// sleep-based cron mechanism. When the report is generated an email is sent to the
// Rotational admins with the report contents.
type DailyUsers struct {
	emailer      DailyUsersEmailer
	lastRun      time.Time
	nextRun      time.Time
	pollInterval time.Duration
	timezone     *time.Location
	done         chan struct{}
}

type DailyUsersEmailer interface {
	SendDailyUsers(*emails.DailyUsersData) error
}

func NewDailyUsers(emailer DailyUsersEmailer) (report *DailyUsers, err error) {
	report = &DailyUsers{
		emailer:      emailer,
		pollInterval: 15 * time.Minute,
		done:         make(chan struct{}),
	}

	// TODO: make timezone configurable
	if report.timezone, err = time.LoadLocation("America/New_York"); err != nil {
		return nil, fmt.Errorf("could not parse timezone: %w", err)
	}

	// TODO: make schedule configurable
	// Run every morning at 6 am in the specified time zone
	now := time.Now().In(report.timezone)
	report.nextRun = time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, report.timezone)
	if now.After(report.nextRun) {
		report.nextRun = report.nextRun.AddDate(0, 0, 1)
	}

	return report, nil
}

// Run the daily users report tool which checks if it needs to run every 15 minutes
// (e.g. the poll interval) and if the next run time is after the current time it runs
// the report, otherwise it continues to sleep.
func (r *DailyUsers) Run() error {
	ticker := time.NewTicker(r.pollInterval)
	log.Info().Msg("daily users report scheduler started")

	for {
		select {
		case <-r.done:
			return nil
		case ts := <-ticker.C:
			if err := r.Scheduler(ts); err != nil {
				// Since this is going to stop the reporting tool; report with fatal level.
				sentry.Fatal(nil).Err(err).Msg("daily report scheduler terminated")
				return err
			}
		}
	}
}

// Shutdown the report runner gracefully; this will ensure that a report that is
// currently being generated will finish before shutting down.
func (r *DailyUsers) Shutdown() error {
	r.done <- struct{}{}
	return nil
}

// Scheduler determines if the report needs to be run and if it does, it runs the report
// and schedules the next run (updating the lastRun timestamp). We expect the ts to be
// passed in from the ticker at run; but if a zero-valued timestamp is passed in, then
// the current timestamp is used. If this method returns an error it is assumed to be
// fatal (e.g. so bad it should shut down the reporting tool).
func (r *DailyUsers) Scheduler(ts time.Time) (err error) {
	log.Debug().Time("last_run", r.lastRun).Time("next_run", r.nextRun).Msg("daily users report scheduler")
	if ts.IsZero() {
		ts = time.Now().In(r.timezone)
	}

	// This shouldn't happen but better safe than sorry.
	if ts.Before(r.lastRun) || ts.Equal(r.lastRun) {
		return ErrBeforeLastRun
	}

	if ts.Before(r.nextRun) {
		// We haven't reached the time for the next run so skip (no error)
		return nil
	}

	// At this point it's either nextRun or after nextRun so it's time to run the report.
	r.lastRun = time.Now()
	err = r.Report()

	// Schedule the next run
	r.nextRun = r.nextRun.AddDate(0, 0, 1)
	log.Info().Time("last_run", r.lastRun).Time("next_run", r.nextRun).Msg("daily users report has been run")
	return err
}

// Runs the daily users report and emails the admins. If an error is returned from
// report it is expected to be fatal; all other errors should simply be logged.
func (r *DailyUsers) Report() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		// Being unable to connect to the database is a fatal error.
		return err
	}
	defer tx.Rollback()

	// Create the report data for sending the report to admins.
	now := time.Now().In(r.timezone)
	yesterday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, r.timezone).AddDate(0, 0, -1)

	report := &emails.DailyUsersData{
		Date:         yesterday,
		InactiveDate: yesterday.AddDate(0, 0, -30),
	}

	day := sql.Named("day", report.Date.Format("2006-01-02"))
	inactive := sql.Named("inactive", report.InactiveDate.Format("2006-01-02"))

	// New Users
	if err = tx.QueryRow("SELECT count(id) FROM users WHERE date(created) == date(:day)", day).Scan(&report.NewUsers); err != nil {
		return err
	}

	// Daily Users
	// NOTE: if they user is logged in today that doesn't necessarily mean they were
	// logged in yesterday, but the only way to track this is with Prometheus.
	if err = tx.QueryRow("SELECT count(id) FROM users WHERE date(last_login) >= date(:day)", day).Scan(&report.DailyUsers); err != nil {
		return err
	}

	// Active Users
	if err = tx.QueryRow("SELECT count(id) FROM users WHERE date(last_login) >= date(:inactive)", inactive).Scan(&report.ActiveUsers); err != nil {
		return err
	}

	// Inactive Users
	if err = tx.QueryRow("SELECT count(id) FROM users WHERE last_login == '' || date(last_login) < date(:inactive)", inactive).Scan(&report.InactiveUsers); err != nil {
		return err
	}

	// API Keys
	if err = tx.QueryRow("SELECT count(id) FROM api_keys").Scan(&report.APIKeys); err != nil {
		return err
	}

	// Revoked API Keys
	if err = tx.QueryRow("SELECT count(id) FROM revoked_api_keys").Scan(&report.RevokedKeys); err != nil {
		return err
	}

	// Active API Keys
	if err = tx.QueryRow("SELECT count(id) FROM api_keys WHERE date(last_used) >= date(:inactive)", inactive).Scan(&report.ActiveKeys); err != nil {
		return err
	}

	// Inactive API Keys
	if err = tx.QueryRow("SELECT count(id) FROM api_keys WHERE last_used == '' || date(last_used) < date(:inactive)", inactive).Scan(&report.InactiveKeys); err != nil {
		return err
	}

	// New Organizations
	if err = tx.QueryRow("SELECT count(id) FROM organizations WHERE date(created) == date(:day)", day).Scan(&report.NewOrganizations); err != nil {
		return err
	}

	// Organizations
	if err = tx.QueryRow("SELECT count(id) FROM organizations").Scan(&report.Organizations); err != nil {
		return err
	}

	// New Projects
	if err = tx.QueryRow("SELECT count(*) FROM organization_projects WHERE date(created) == date(:day)", day).Scan(&report.NewProjects); err != nil {
		return err
	}

	// Projects
	if err = tx.QueryRow("SELECT count(*) FROM organization_projects").Scan(&report.Projects); err != nil {
		return err
	}

	// Commit the transaction to conclude it (errors not fatal).
	if err = tx.Commit(); err != nil {
		sentry.Error(nil).Err(err).Msg("could not commit daily users report readonly tx")
	}

	// Email the report to the admins; if the email fails, log it (not fatal)
	if err = r.emailer.SendDailyUsers(report); err != nil {
		sentry.Error(nil).Err(err).Msg("could not send daily report to admins")
	}
	return nil
}

// Get the last time the daily report was run.
func (r *DailyUsers) LastRun() time.Time {
	return r.lastRun
}

// Get the timestamp that the next report should run (within 15 minutes).
func (r *DailyUsers) NextRun() time.Time {
	return r.nextRun
}

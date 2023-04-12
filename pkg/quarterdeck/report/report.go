package report

import (
	"fmt"
	"time"

	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog/log"
)

// The DailyUsers report is generated once per day at the specified schedule using a
// sleep-based cron mechanism. When the report is generated an email is sent to the
// Rotational admins with the report contents.
type DailyUsers struct {
	lastRun      time.Time
	nextRun      time.Time
	pollInterval time.Duration
	timezone     *time.Location
	done         chan struct{}
}

func NewDailyUsers() (report *DailyUsers, err error) {
	report = &DailyUsers{
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
	for {
		select {
		case <-r.done:
			return nil
		case ts := <-ticker.C:
			if err := r.Scheduler(ts); err != nil {
				// Since this is going to stop the reporting tool; report with fatal level.
				sentry.Fatal(nil).Err(err).Msg("daily report scheduler terminated")
				return nil
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
func (r *DailyUsers) Report() error {
	// TODO: query the database.
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

package quarterdeck_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rotationalio/ensign/pkg/quarterdeck"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSendDailyUsersEmail(t *testing.T) {
	// NOTE: in most cases this test will be skipped, but if you set the
	// $QUARTERDECK_TEST_SENDING_EMAILS environment variable then this test will attempt
	// to send the daily users email report. You must also have the $SENDGRID_API_KEY
	// and $QUARTERDECK_ADMIN_EMAIL environment variables set for this to work.

	// NOTE: if you place a .env file in this directory alongside the test file, it
	// will be read, making it simpler to run tests and set environment variables.
	godotenv.Load()

	if os.Getenv("QUARTERDECK_TEST_SENDING_EMAILS") == "" {
		t.Skip("skip quarterdeck send emails test")
	}

	// Create a test configuration for a quarterdeck server
	// HACK: this is a bit fragile but it sets up the sendgrid client.
	conf, err := config.Config{
		Maintenance:  false,
		BindAddr:     "127.0.0.1:0",
		Mode:         gin.TestMode,
		LogLevel:     logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:   false,
		AllowOrigins: []string{"http://localhost:3000"},
		EmailURL: config.URLConfig{
			Base:   "http://localhost:3000",
			Invite: "/invite",
			Verify: "/verify",
		},
		SendGrid: emails.Config{
			APIKey:     os.Getenv("SENDGRID_API_KEY"),
			FromEmail:  "enson@rotational.io",
			AdminEmail: os.Getenv("QUARTERDECK_ADMIN_EMAIL"),
			Testing:    false,
		},
		Database: config.DatabaseConfig{
			URL:      "sqlite3:///" + filepath.Join(t.TempDir(), "test.db"),
			ReadOnly: false,
		},
		Token: config.TokenConfig{
			Keys: map[string]string{
				"01GE6191AQTGMCJ9BN0QC3CCVG": "testdata/01GE6191AQTGMCJ9BN0QC3CCVG.pem",
				"01GE62EXXR0X0561XD53RDFBQJ": "testdata/01GE62EXXR0X0561XD53RDFBQJ.pem",
			},
			Audience:        "http://localhost:3000",
			Issuer:          "http://quarterdeck.test/",
			AccessDuration:  10 * time.Minute,
			RefreshDuration: 20 * time.Minute,
			RefreshOverlap:  -10 * time.Minute,
		},
		RateLimit: config.RateLimitConfig{
			PerSecond: 60.00,
			Burst:     120,
			TTL:       5 * time.Minute,
		},
	}.Mark()
	assert.NoError(t, err, "test configuration is invalid")

	srv, err := quarterdeck.New(conf)
	assert.NoError(t, err, "could not create the quarterdeck api server from the test configuration")

	err = srv.Setup()
	assert.NoError(t, err, "could not setup quarterdeck api server from test config")

	// Send the email
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	inactive := today.AddDate(0, 0, -30)

	data := &emails.DailyUsersData{
		Date:                today,
		InactiveDate:        inactive,
		Domain:              "ensign.localhost",
		EnsignDashboardLink: "http://grafana.localhost",
		NewUsers:            41,
		DailyUsers:          17,
		ActiveUsers:         1032,
		InactiveUsers:       21,
		APIKeys:             96,
		InactiveKeys:        32,
		ActiveKeys:          64,
		RevokedKeys:         102,
		Organizations:       109,
		NewOrganizations:    32,
		Projects:            3312,
		NewProjects:         41,
	}

	err = srv.SendDailyUsers(data)
	assert.NoError(t, err, "could not send email to users")
}

package emails_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/emails/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

// If the eyeball flag is set, then the tests will write MIME emails to the testdata directory.
var eyeball = flag.Bool("eyeball", false, "Generate MIME emails for eyeball testing")

// This suite mocks the SendGrid email client to verify that email metadata is
// populated corectly and emails can be marshaled into bytes for transmission.
func TestEmailSuite(t *testing.T) {
	suite.Run(t, &EmailTestSuite{})
}

type EmailTestSuite struct {
	suite.Suite
	conf emails.Config
}

func (suite *EmailTestSuite) SetupSuite() {
	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog MUST be false otherwise this will be overriden
	logger.Discard()

	suite.conf = emails.Config{
		Testing:    true,
		FromEmail:  "service@example.com",
		AdminEmail: "admin@example.com",
		Archive:    filepath.Join("fixtures", "emails"),
	}
}

func (suite *EmailTestSuite) BeforeTest(suiteName, testName string) {
	setupMIMEDir(suite.T())
}

func (suite *EmailTestSuite) AfterTest(suiteName, testName string) {
	mock.Reset()
}

func (suite *EmailTestSuite) TearDownSuite() {
	logger.ResetLogger()
}

// If eyeball testing is enabled, this removes and recreates the eyeball directory for
// this test.
func setupMIMEDir(t *testing.T) {
	if *eyeball {
		path := filepath.Join("testdata", fmt.Sprintf("eyeball%s", t.Name()))
		err := os.RemoveAll(path)
		require.NoError(t, err)
		err = os.MkdirAll(path, 0755)
		require.NoError(t, err)
	}
}

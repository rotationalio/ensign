package quarterdeck_test

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func (suite *quarterdeckTestSuite) TestSecurityTxt() {
	require := suite.Require()

	// Create a basic HTTP request rather than use the Quarterdeck client to ensure the
	// headers and data returned are the expected values.
	req, err := http.NewRequest(http.MethodGet, suite.srv.URL()+"/.well-known/security.txt", nil)
	require.NoError(err, "could not create basic http request")

	rep, err := http.DefaultClient.Do(req)
	require.NoError(err, "could not execute basic http request")
	defer rep.Body.Close()

	// The content-type must be text/plain for the security.text file
	require.Equal("text/plain; charset=utf-8", rep.Header.Get("Content-Type"))

	// Ensure the content returned is as expected
	expected, err := os.ReadFile("testdata/security.txt")
	require.NoError(err, "could not read testdata/security.txt fixture")

	actual, err := io.ReadAll(rep.Body)
	require.NoError(err, "could not read the body response from the server")
	require.Equal(expected, actual, "the security.txt file does not match the testdata fixture")

	// TODO: check that the GPG signature is valid
	// This can be done on the command line with: gpg --auto-key-retrieve --verify --output -
	// Not sure how to do this in golang tests, however.

	// Error if the security.txt has expired - if this test fails, regenerate the security.txt file!
	lines := strings.Split(string(actual), "\n")
	found := false
	for _, line := range lines {
		if strings.HasPrefix(line, "Expires:") {
			found = true
			parts := strings.Split(line, " ")
			require.Len(parts, 2, "could not split expires directive")

			expires, err := time.Parse(time.RFC3339, parts[1])
			require.NoError(err, "could not parse expires timestamp")

			require.True(time.Now().Before(expires), "the security.txt file has expired, regenerate it!")
		}
	}
	require.True(found, "could not find expires line")
}

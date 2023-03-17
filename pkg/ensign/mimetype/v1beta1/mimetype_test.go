package mimetype_test

import (
	"strings"
	"testing"

	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		s        string
		expected mimetype.MIME
	}{
		{"application/json", mimetype.MIME_APPLICATION_JSON},
		{"user/format-5", mimetype.MIME_USER_SPECIFIED5},
		{"APPLICATION/XML", mimetype.MIME_APPLICATION_XML},
		{"  TEXT/csv  ", mimetype.MIME_TEXT_CSV},
		{"application/json; charset=utf-8", mimetype.MIME_APPLICATION_JSON_UTF8},
		{"application/atom+xml; charset=latin1", mimetype.MIME_APPLICATION_ATOM},
		{"application/vnd.myapp.type+xml", mimetype.MIME_APPLICATION_XML},
		{"text/atom+xml", mimetype.MIME_APPLICATION_ATOM},
		{"application/csv", mimetype.MIME_TEXT_CSV},
		{"text/vnd.myapp.type+xml", mimetype.MIME_APPLICATION_XML},
	}

	for _, tc := range testCases {
		mime, err := mimetype.Parse(tc.s)
		require.NoError(t, err, "could not parse %q", tc.s)
		require.Equal(t, tc.expected, mime, "expected mimetype not returned")
		require.Equal(t, tc.expected, mimetype.MustParse(tc.s), "expected mimetype not returned")
	}

	// Test bad cases
	testCases = []struct {
		s        string
		expected mimetype.MIME
	}{
		{"text/svg+png", mimetype.MIME_TEXT_PLAIN},
		{"image/png", mimetype.MIME_APPLICATION_OCTET_STREAM},
		{"application/vnd.myapp.type", mimetype.MIME_APPLICATION_OCTET_STREAM},
	}

	for _, tc := range testCases {
		mime, err := mimetype.Parse(tc.s)
		require.Error(t, err, "expected unknown mimetype error")
		require.Equal(t, mimetype.MIME_UNSPECIFIED, mime, "unexpected mime returned")
		require.Equal(t, tc.expected, mimetype.MustParse(tc.s), "default mime type unexpected")
	}
}

func TestStrings(t *testing.T) {
	// Ensure that the mimetype strings are used
	for key, val := range mimetype.MIMEType_name {
		require.Equal(t, val, mimetype.MIME(key).MimeType())
	}
}

func TestCoverage(t *testing.T) {
	// Ensure that all of the protocol buffer mimetypes are defined in Go
	for key := range mimetype.MIME_name {
		require.Contains(t, mimetype.MIMEType_name, key, "the protobuf MIME_name contains a value that is not in MIMEType_name")
	}

	for key := range mimetype.MIMEType_name {
		require.Contains(t, mimetype.MIME_name, key, "the MIMEType_name contains a value that is not specified by the protocol buffers")
	}

	// Get sets of all the values in the value maps
	mimeVals := make(map[int32]struct{})
	mimetypeVals := make(map[int32]struct{})

	for _, val := range mimetype.MIME_value {
		require.Contains(t, mimetype.MIMEType_name, val, "the protobuf MIME_value contains a value that is not in MIMEType_name")
		mimeVals[val] = struct{}{}
	}

	// Ensure that all of the values are represented in the name maps
	for _, val := range mimetype.MIMEType_value {
		require.Contains(t, mimetype.MIME_name, val, "the MIMEType_value contains a value that is not in the protocol buffers")
		mimetypeVals[val] = struct{}{}
	}

	// Ensure that all the values are defined in Go
	for key, val := range mimetype.MIME_value {
		require.Contains(t, mimetypeVals, val, "the protobuf MIME_value contains %q that is not in MIMEType_values", key)
		mimeVals[val] = struct{}{}
	}

	for key, val := range mimetype.MIMEType_value {
		require.Contains(t, mimeVals, val, "the MIMEType_value contains %q that is not in the protocol buffers", key)
		mimetypeVals[val] = struct{}{}
	}

	// Ensure that the hand-coded values are complete (no need to check generated code)
	for key, val := range mimetype.MIMEType_name {
		require.Contains(t, mimetype.MIMEType_value, val, "missing %q from value map", val)
		require.Equal(t, mimetype.MIMEType_value[val], key, "string mismatch in MIMEType_value map")
	}

	for key, val := range mimetype.MIMEType_value {
		require.Contains(t, mimetype.MIMEType_name, val, "missing %q from name map", val)
		require.Equal(t, mimetype.MIMEType_name[val], key, "string mismatch in MIMEType_name map")
	}

	// One last sanity check that should be true based on the above tests.
	require.Equal(t, len(mimetype.MIME_name), len(mimetype.MIMEType_name), "the MIMETYPE_name data structure is not synchronized with the protocol buffers")
}

func TestPrefixes(t *testing.T) {
	prefixes := map[string]struct{}{
		"application": {},
		"text":        {},
		"user":        {},
	}

	// Ensure prefix spelling is correct
	for key := range mimetype.MIMEType_value {
		parts := strings.Split(key, "/")
		require.Len(t, parts, 2, "mimetype %q should only conain one / sep", key)
		require.Contains(t, prefixes, parts[0], "unknown prefix %q", parts[0])
	}

	for _, key := range mimetype.MIMEType_name {
		parts := strings.Split(key, "/")
		require.Len(t, parts, 2, "mimetype %q should only conain one / sep", key)
		require.Contains(t, prefixes, parts[0], "unknown prefix %q", parts[0])
	}
}

func TestRegexp(t *testing.T) {
	testCases := []string{
		"application/json",
		"application/ld+json",
		"application/json; charset=utf8",
		"application/protobuf; msg=trisa.v1beta1.SecureEnvelope counterparty=trisa.example.com",
		"application/atom+xml; charset=utf-8",
		"application/vnd.google.protobuf",
		"application/vnd.myapp.type+xml",
		"application/x-protobuf",
		"text/plain",
		"text/html; charset=utf8",
		"text/tsv+csv",
		"text/tsv+csv; charset=utf8",
		"user/format-0",
		"user/format-853",
	}

	for _, tc := range testCases {
		require.True(t, mimetype.MIMERegExp.MatchString(tc), "%q does not match mime regexp", tc)
		result := mimetype.Components(tc)
		for _, key := range []string{"prefix", "mime", "ext", "pairs"} {
			require.Contains(t, result, key, "missing required component")
		}
	}
}

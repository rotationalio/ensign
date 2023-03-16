package mimetype

import (
	"fmt"
	"regexp"
	"strings"
)

// Parse a string mimetype into a mimetype constant. If the given mimetype is unknown
// then an error is returned. Parse returns best effort mimetypes. For example if the
// mimetype is application/vnd.myapp.type+xml then application/xml is returned. This
// method is case and whitespace insensitive.
func Parse(s string) (MIME, error) {
	// Case and whitespace insensitivity
	s = strings.ToLower(strings.TrimSpace(s))

	// Check if the mimetype is in the map
	if mime, ok := MIMEType_value[s]; ok {
		return MIME(mime), nil
	}

	// try alternative formulations of the mimetype
	components := Components(s)

	// remove charset or any key/value pairs from mime type along with the extension
	if mime, ok := MIMEType_value[fmt.Sprintf("%s/%s", components["prefix"], components["mime"])]; ok {
		return MIME(mime), nil
	}

	// remove pairs but keep extension
	if mime, ok := MIMEType_value[fmt.Sprintf("%s/%s+%s", components["prefix"], components["mime"], components["ext"])]; ok {
		return MIME(mime), nil
	}

	// try just the extension
	if mime, ok := MIMEType_value[fmt.Sprintf("%s/%s", components["prefix"], components["ext"])]; ok {
		return MIME(mime), nil
	}

	// try the alternate prefix
	var prefix string
	switch components["prefix"] {
	case "application":
		prefix = "text"
	case "text":
		prefix = "application"
	}

	if prefix != "" {
		// try new prefix with extension
		if mime, ok := MIMEType_value[fmt.Sprintf("%s/%s+%s", prefix, components["mime"], components["ext"])]; ok {
			return MIME(mime), nil
		}

		// try alternative prefix without extension
		if mime, ok := MIMEType_value[fmt.Sprintf("%s/%s", prefix, components["mime"])]; ok {
			return MIME(mime), nil
		}

		// try alternative prefix using extension only
		if mime, ok := MIMEType_value[fmt.Sprintf("%s/%s", prefix, components["ext"])]; ok {
			return MIME(mime), nil
		}
	}

	return MIME_UNSPECIFIED, fmt.Errorf("unknown mimetype %q", s)
}

// MustParse does not return an error but returns a default mimetype; text/plain if the
// mimetype prefix is text otherwise application/octet-stream.
func MustParse(s string) MIME {
	if mime, err := Parse(s); err == nil {
		return mime
	}

	s = strings.ToLower(strings.TrimSpace(s))
	components := Components(s)

	switch components["prefix"] {
	case "text":
		return MIME_TEXT_PLAIN
	default:
		return MIME_APPLICATION_OCTET_STREAM
	}
}

// Enum value maps for mimetype strings
var (
	MIMEType_name = map[int32]string{
		0:    "application/octet-stream",
		1:    "text/plain",
		2:    "text/csv",
		3:    "text/html",
		4:    "text/calendar",
		50:   "application/json",
		51:   "application/ld+json",
		52:   "application/jsonlines",
		53:   "application/ubjson",
		54:   "application/bson",
		100:  "application/xml",
		101:  "application/atom+xml",
		252:  "application/msgpack",
		253:  "application/parquet",
		254:  "application/avro",
		255:  "application/protobuf",
		512:  "application/pdf",
		513:  "application/java-archive",
		514:  "application/python-pickle",
		1000: "user/format-0",
		1001: "user/format-1",
		1002: "user/format-2",
		1003: "user/format-3",
		1004: "user/format-4",
		1005: "user/format-5",
		1006: "user/format-6",
		1007: "user/format-7",
		1008: "user/format-8",
		1009: "user/format-9",
	}
	MIMEType_value = map[string]int32{
		"application/octet-stream":  0,
		"text/plain":                1,
		"text/csv":                  2,
		"text/html":                 3,
		"text/calendar":             4,
		"application/json":          50,
		"application/ld+json":       51,
		"application/jsonlines":     52,
		"application/ubjson":        53,
		"application/bson":          54,
		"application/xml":           100,
		"application/atom+xml":      101,
		"application/msgpack":       252,
		"application/parquet":       253,
		"application/avro":          254,
		"application/protobuf":      255,
		"application/pdf":           512,
		"application/java-archive":  513,
		"application/python-pickle": 514,
		"user/format-0":             1000,
		"user/format-1":             1001,
		"user/format-2":             1002,
		"user/format-3":             1003,
		"user/format-4":             1004,
		"user/format-5":             1005,
		"user/format-6":             1006,
		"user/format-7":             1007,
		"user/format-8":             1008,
		"user/format-9":             1009,
	}
	MIMERegExp = regexp.MustCompile(`^(?P<prefix>application|user|text)\/(?P<mime>[\w\-_\.]+)(\+(?P<ext>[\w-]+))?(?P<pairs>;(\s+([\w\.\-_]+=[\w\.\-_]+))+)?$`)
)

// Returns the MimeType name as defined by the IETF specification.
func (x MIME) MimeType() string {
	return MIMEType_name[int32(x.Number())]
}

// Components returns the mimetype parts in the following format: prefix/mime+ext; pairs
// Note that the pairs will include the leading ; and all spaces between key/values.
func Components(s string) map[string]string {
	result := make(map[string]string)
	match := MIMERegExp.FindStringSubmatch(s)
	if len(match) == 0 {
		return result
	}

	for i, name := range MIMERegExp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}

// MimeType constant aliases for more readable Go code.
const (
	Unspecified             = MIME_UNSPECIFIED
	Unknown                 = MIME_UNKNOWN
	TextPlain               = MIME_TEXT_PLAIN
	TextCSV                 = MIME_TEXT_CSV
	TextHTML                = MIME_TEXT_HTML
	TextCalendar            = MIME_TEXT_CALENDAR
	ApplicationOctetStream  = MIME_APPLICATION_OCTET_STREAM
	ApplicationJSON         = MIME_APPLICATION_JSON
	ApplicationJSONUTF8     = MIME_APPLICATION_JSON_UTF8
	ApplicationJSONLD       = MIME_APPLICATION_JSON_LD
	AppplicationJSONLines   = MIME_APPLICATION_JSON_LINES
	ApplicationUBJSON       = MIME_APPLICATION_UBJSON
	ApplicationBSON         = MIME_APPLICATION_BSON
	ApplicationXML          = MIME_APPLICATION_XML
	ApplicationAtom         = MIME_APPLICATION_ATOM
	ApplicationMsgPack      = MIME_APPLICATION_MSGPACK
	ApplicationParquet      = MIME_APPLICATION_PARQUET
	ApplicationAvro         = MIME_APPLICATION_AVRO
	ApplicationProtobuf     = MIME_APPLICATION_PROTOBUF
	ApplicationPDF          = MIME_APPLICATION_PDF
	ApplicationJavaArchive  = MIME_APPLICATION_JAVA_ARCHIVE
	ApplicationPythonPickle = MIME_APPLICATION_PYTHON_PICKLE
	UserSpecified0          = MIME_USER_SPECIFIED0
	UserSpecified1          = MIME_USER_SPECIFIED1
	UserSpecified2          = MIME_USER_SPECIFIED2
	UserSpecified3          = MIME_USER_SPECIFIED3
	UserSpecified4          = MIME_USER_SPECIFIED4
	UserSpecified5          = MIME_USER_SPECIFIED5
	UserSpecified6          = MIME_USER_SPECIFIED6
	UserSpecified7          = MIME_USER_SPECIFIED7
	UserSpecified8          = MIME_USER_SPECIFIED8
	UserSpecified9          = MIME_USER_SPECIFIED9
)

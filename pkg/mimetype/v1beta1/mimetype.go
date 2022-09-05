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

// String replaces the protocol buffer String() method to return a human-friendly
// representation of the mimetype rather than the enum variable name.
// NOTE: must delete the protocol buffer String() method everytime they're generated.
func (x MIME) String() string {
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

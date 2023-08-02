package tenant

import (
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	mimetype "github.com/rotationalio/go-ensign/mimetype/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Return fake events
func fixtureEvents() (events []*api.Event) {
	events = make([]*api.Event, 0)

	// Plaintext event
	events = append(events, &api.Event{
		Data:     []byte("hello world"),
		Metadata: map[string]string{},
		Mimetype: mimetype.TextPlain,
		Type: &api.Type{
			Name:         "Message",
			MajorVersion: 1,
		},
		Created: timestamppb.Now(),
	})

	// CSV event
	events = append(events, &api.Event{
		Data:     []byte("hello,world"),
		Metadata: map[string]string{},
		Mimetype: mimetype.TextCSV,
		Type: &api.Type{
			Name:         "Spreadsheet",
			MajorVersion: 1,
			MinorVersion: 1,
		},
		Created: timestamppb.Now(),
	})

	// HTML event
	events = append(events, &api.Event{
		Data:     []byte("<html><body><h1>Hello World</h1></body></html>"),
		Metadata: map[string]string{},
		Mimetype: mimetype.TextHTML,
		Type: &api.Type{
			Name:         "Webpage",
			MajorVersion: 1,
			PatchVersion: 1,
		},
		Created: timestamppb.Now(),
	})

	// JSON events
	events = append(events, &api.Event{
		Data:     []byte(`{"price": 334.11, "symbol": "MSFT", "timestamp": 1690899527135, "volume": 50}`),
		Metadata: map[string]string{},
		Mimetype: mimetype.ApplicationJSON,
		Type: &api.Type{
			Name:         "StockQuote",
			MinorVersion: 1,
		},
		Created: timestamppb.Now(),
	})
	events = append(events, &api.Event{
		Data:     []byte(`{"price": 320.23, "symbol": "APPL", "timestamp": 1690899527135, "volume": 25}`),
		Metadata: map[string]string{},
		Mimetype: mimetype.ApplicationJSON,
		Type: &api.Type{
			Name:         "StockQuote",
			MinorVersion: 1,
		},
		Created: timestamppb.Now(),
	})
	events = append(events, &api.Event{
		Data:     []byte(`{"price": 335.12, "symbol": "MSFT", "timestamp": 1690899527135, "volume": 40}`),
		Metadata: map[string]string{},
		Mimetype: mimetype.ApplicationJSON,
		Type: &api.Type{
			Name:         "StockQuote",
			MinorVersion: 1,
		},
		Created: timestamppb.Now(),
	})

	// XML event
	events = append(events, &api.Event{
		Data:     []byte(`<note><to>Arthur</to><from>Marvin</from><heading>Life</heading><body>Don't Panic!</body></note>`),
		Metadata: map[string]string{},
		Mimetype: mimetype.ApplicationXML,
		Type: &api.Type{
			Name:         "Note",
			PatchVersion: 1,
		},
		Created: timestamppb.Now(),
	})

	// msgpack events
	events = append(events, &api.Event{
		Data:     []byte("\x81\xa4name\xa3Bob\xa3age\x18\x1e"),
		Metadata: map[string]string{},
		Mimetype: mimetype.ApplicationMsgPack,
		Type: &api.Type{
			Name:         "Person",
			MajorVersion: 1,
			MinorVersion: 3,
		},
		Created: timestamppb.Now(),
	})
	events = append(events, &api.Event{
		Data:     []byte("\x81\xa4name\xa3Alice\xa3age\x18\x1e"),
		Metadata: map[string]string{},
		Mimetype: mimetype.ApplicationMsgPack,
		Type: &api.Type{
			Name:         "Person",
			MajorVersion: 1,
			MinorVersion: 3,
		},
		Created: timestamppb.Now(),
	})

	// protobuf events
	events = append(events, &api.Event{
		Data:     []byte("\x0a\x03Bob\x10\x1e"),
		Metadata: map[string]string{},
		Mimetype: mimetype.ApplicationProtobuf,
		Type: &api.Type{
			Name:         "Person",
			MajorVersion: 2,
			PatchVersion: 1,
		},
		Created: timestamppb.Now(),
	})
	events = append(events, &api.Event{
		Data:     []byte("\x0a\x05Alice\x10\x1e"),
		Metadata: map[string]string{},
		Mimetype: mimetype.ApplicationProtobuf,
		Type: &api.Type{
			Name:         "Person",
			MajorVersion: 2,
			PatchVersion: 1,
		},
	})

	return events
}

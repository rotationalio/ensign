package api_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/oklog/ulid/v2"
	. "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHash(t *testing.T) {
	// For the same event, but different policies, each hashing method should return
	// a different hash signature of the fixture.
	hashes := make(map[string]struct{})
	event := createRandomEvent(
		mimetype.ApplicationMsgPack,
		&Type{Name: "RandomData", MajorVersion: 1, MinorVersion: 2, PatchVersion: 3},
		map[string]string{
			"FirstKey":  "rand",
			"SecondKey": "rand",
			"ThirdKey":  "rand",
			"foo":       "bar",
			"color":     "blue",
		},
	)

	policies := []*Deduplication{
		{Strategy: Deduplication_STRICT},
		{Strategy: Deduplication_DATAGRAM},
		{Strategy: Deduplication_KEY_GROUPED, Keys: []string{"FirstKey", "SecondKey"}},
		{Strategy: Deduplication_UNIQUE_KEY, Keys: []string{"FirstKey", "SecondKey"}},
	}

	for _, policy := range policies {
		sig, err := event.Hash(policy)
		require.NoError(t, err, "could not create a hash of the event")
		hashes[base64.RawStdEncoding.EncodeToString(sig)] = struct{}{}
	}

	require.Len(t, hashes, len(policies), "expected a unique hash for each policy, one of the hashes is duplicated")
}

func TestHashStrict(t *testing.T) {
	fixtures, err := loadFixtures()
	require.NoError(t, err, "could not load event fixtures for strict hashing tests")

	testCases := []struct {
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			&EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			fixtures[0], []byte{0xb, 0xa6, 0x75, 0xb8, 0x91, 0xd7, 0xee, 0xc8, 0xa4, 0xe1, 0xf, 0x84, 0x7d, 0x8a, 0x14, 0x9b}, nil,
		},
		{
			fixtures[1], []byte{0xe6, 0xab, 0x9c, 0x80, 0xd, 0xec, 0xa8, 0x2a, 0x5d, 0xaf, 0x5a, 0xbd, 0x99, 0x16, 0xb3, 0x66}, nil,
		},
		{
			fixtures[2], []byte{0x32, 0x1e, 0xf3, 0xfb, 0x10, 0x56, 0xb0, 0x12, 0x2d, 0x87, 0x1b, 0x1b, 0xc3, 0x65, 0xf1, 0x5e}, nil,
		},
		{
			fixtures[3], []byte{0xc8, 0x33, 0x8f, 0xf9, 0x7b, 0x87, 0x1, 0x30, 0xc7, 0xc5, 0x77, 0xf4, 0x1b, 0x8a, 0x6c, 0x95}, nil,
		},
		{
			fixtures[4], []byte{0xf9, 0xdd, 0xed, 0xfc, 0xd8, 0x7e, 0xdc, 0x47, 0x31, 0x58, 0x80, 0x41, 0x34, 0x82, 0x22, 0x39}, nil,
		},
		{
			fixtures[5], []byte{0x9e, 0xd2, 0xca, 0xe, 0x35, 0xd2, 0x11, 0xc8, 0xe2, 0x2a, 0x46, 0xf1, 0x46, 0xb0, 0xf7, 0x71}, nil,
		},
		{
			fixtures[6], []byte{0x7f, 0x6a, 0xad, 0xbb, 0xec, 0x1a, 0xcc, 0x9, 0x9f, 0xb0, 0x6e, 0x9e, 0x64, 0x67, 0xf0, 0x5c}, nil,
		},
		{
			fixtures[7], []byte{0xbe, 0x68, 0x6a, 0xd5, 0x70, 0x5, 0xef, 0x12, 0xc9, 0xd6, 0x20, 0x37, 0xc2, 0x3b, 0xd1, 0x89}, nil,
		},
	}

	for i, tc := range testCases {
		hash, err := tc.event.HashStrict()
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected an error for test case %d", i)
		} else {
			require.NoError(t, err, "expected no error for test case %d", i)
			require.Equal(t, tc.expected, hash, "expected hash equality for test case %d", i)
		}
	}
}

func TestHashDatagram(t *testing.T) {
	fixtures, err := loadFixtures()
	require.NoError(t, err, "could not load event fixtures for strict hashing tests")

	testCases := []struct {
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			&EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			fixtures[0], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
		{
			fixtures[1], []byte{0x57, 0xad, 0xd1, 0x29, 0x33, 0xd1, 0xfb, 0xda, 0xec, 0x1f, 0x5d, 0x8d, 0x8, 0x22, 0x2, 0xc2}, nil,
		},
		{
			fixtures[2], []byte{0x42, 0x60, 0x33, 0x63, 0xdb, 0xff, 0x5, 0x2c, 0x82, 0x92, 0xbd, 0x6a, 0x35, 0x42, 0x7f, 0x72}, nil,
		},
		{
			fixtures[3], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
		{
			fixtures[4], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
		{
			fixtures[5], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
		{
			fixtures[6], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
		{
			fixtures[7], []byte{0x29, 0x69, 0x26, 0xd7, 0x58, 0x65, 0x52, 0xb9, 0x3f, 0xd5, 0x5c, 0x3d, 0xf1, 0x30, 0x29, 0x4e}, nil,
		},
	}

	for i, tc := range testCases {
		hash, err := tc.event.HashDatagram()
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected an error for test case %d", i)
		} else {
			require.NoError(t, err, "expected no error for test case %d", i)
			require.Equal(t, tc.expected, hash, "expected hash equality for test case %d", i)
		}
	}
}

func TestHashKeyGroup(t *testing.T) {
	fixtures, err := loadFixtures()
	require.NoError(t, err, "could not load event fixtures for strict hashing tests")

	testCases := []struct {
		keys     []string
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			[]string{"foo", "bar"}, &EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			nil, fixtures[0], nil, ErrNoKeys,
		},
		{
			[]string{"alpha"}, fixtures[0], []byte{0xb6, 0x8c, 0x41, 0x51, 0x48, 0x79, 0x3d, 0xed, 0x22, 0xd3, 0xc2, 0xd9, 0x8a, 0xba, 0x97, 0xf1}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[0], []byte{0x1c, 0x7c, 0x5f, 0x13, 0xd3, 0xe1, 0xa7, 0x3d, 0xb4, 0x73, 0xbc, 0x29, 0x6c, 0x7c, 0xba, 0xd0}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[0], []byte{0xec, 0x41, 0xdc, 0xed, 0x2, 0x6c, 0xfe, 0xd3, 0xf6, 0xbc, 0xab, 0xb3, 0x32, 0x6a, 0xff, 0x36}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[0], []byte{0xc6, 0x94, 0x6f, 0xf2, 0x9a, 0xcd, 0xbb, 0x7a, 0x7b, 0x3d, 0x66, 0x7c, 0xba, 0x6b, 0xd0, 0xf8}, nil,
		},
		{
			[]string{"alpha"}, fixtures[1], []byte{0x37, 0x13, 0x18, 0xee, 0x86, 0xcd, 0xad, 0xef, 0x44, 0xf9, 0x7d, 0xb4, 0xe2, 0x9, 0x4b, 0x65}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[1], []byte{0x85, 0xe7, 0x60, 0xa5, 0x23, 0x1f, 0x7f, 0x6b, 0x61, 0x2b, 0x43, 0x6e, 0x3c, 0x76, 0x92, 0x50}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[1], []byte{0xb1, 0xf1, 0xcd, 0xa4, 0x27, 0x5f, 0x9c, 0xfd, 0x88, 0x6b, 0x16, 0x3f, 0x3b, 0x9d, 0x75, 0xaf}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[1], []byte{0xc6, 0x7b, 0x7f, 0xcf, 0x2f, 0x8, 0x82, 0x32, 0xb0, 0xea, 0x32, 0x3f, 0x29, 0x37, 0x4a, 0x28}, nil,
		},
		{
			[]string{"alpha"}, fixtures[3], []byte{0xe0, 0xab, 0x20, 0x5c, 0xca, 0xe6, 0x97, 0xa1, 0xae, 0x1f, 0x76, 0xda, 0xb4, 0x97, 0x1b, 0x54}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[3], []byte{0x47, 0x80, 0x8, 0x9, 0xad, 0xe1, 0x7d, 0x42, 0xc, 0xd4, 0x76, 0xe0, 0x9d, 0xee, 0xeb, 0x3a}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[3], []byte{0xb9, 0x5c, 0x2, 0x4f, 0xc0, 0xda, 0xd8, 0xf4, 0xe4, 0xe6, 0x6e, 0x76, 0x4b, 0x83, 0xcb, 0x37}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[3], []byte{0xc6, 0x60, 0x2, 0xf, 0x9c, 0xa0, 0x12, 0x8b, 0x77, 0x83, 0x58, 0xb8, 0x3a, 0x1b, 0x21, 0x1b}, nil,
		},
	}

	for i, tc := range testCases {
		hash, err := tc.event.HashKeyGrouped(tc.keys)
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected an error for test case %d", i)
		} else {
			require.NoError(t, err, "expected no error for test case %d", i)
			require.Equal(t, tc.expected, hash, "expected hash equality for test case %d", i)
		}
	}
}

func TestHashUniqueKey(t *testing.T) {
	fixtures, err := loadFixtures()
	require.NoError(t, err, "could not load event fixtures for strict hashing tests")

	testCases := []struct {
		keys     []string
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			[]string{"foo", "bar"}, &EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			nil, fixtures[0], nil, ErrNoKeys,
		},
		{
			[]string{"alpha"}, fixtures[0], []byte{0x32, 0x59, 0xfa, 0xa7, 0xed, 0x90, 0xe9, 0xae, 0xab, 0xbc, 0x75, 0xac, 0xa1, 0xf7, 0xac, 0x9a}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[0], []byte{0x5a, 0xbb, 0x4a, 0x32, 0xa9, 0x3d, 0x54, 0x8e, 0xeb, 0xc8, 0xb6, 0x1e, 0xc6, 0x8, 0xef, 0xb4}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[0], []byte{0x85, 0x53, 0xb2, 0x1a, 0x1e, 0x2f, 0x63, 0x29, 0x5b, 0x1, 0x6, 0x1f, 0x64, 0x7a, 0x45, 0x87}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[0], []byte{0x4b, 0xa4, 0x42, 0x74, 0x80, 0xbf, 0x8f, 0x56, 0x19, 0x88, 0x77, 0x0, 0xa4, 0x8, 0x63, 0x9c}, nil,
		},
		{
			[]string{"alpha"}, fixtures[1], []byte{0x32, 0x59, 0xfa, 0xa7, 0xed, 0x90, 0xe9, 0xae, 0xab, 0xbc, 0x75, 0xac, 0xa1, 0xf7, 0xac, 0x9a}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[1], []byte{0x5a, 0xbb, 0x4a, 0x32, 0xa9, 0x3d, 0x54, 0x8e, 0xeb, 0xc8, 0xb6, 0x1e, 0xc6, 0x8, 0xef, 0xb4}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[1], []byte{0x85, 0x53, 0xb2, 0x1a, 0x1e, 0x2f, 0x63, 0x29, 0x5b, 0x1, 0x6, 0x1f, 0x64, 0x7a, 0x45, 0x87}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[1], []byte{0x4b, 0xa4, 0x42, 0x74, 0x80, 0xbf, 0x8f, 0x56, 0x19, 0x88, 0x77, 0x0, 0xa4, 0x8, 0x63, 0x9c}, nil,
		},
		{
			[]string{"alpha"}, fixtures[3], []byte{0xfe, 0x48, 0xa9, 0xff, 0x5b, 0x55, 0x53, 0xb0, 0x8c, 0x78, 0x8b, 0x65, 0x2f, 0x38, 0x52, 0x30}, nil,
		},
		{
			[]string{"alpha", "bravo"}, fixtures[3], []byte{0x65, 0x36, 0xf5, 0x79, 0x3f, 0x74, 0xfa, 0x3e, 0x94, 0xd, 0xe3, 0xed, 0x71, 0x95, 0x61, 0x63}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie"}, fixtures[3], []byte{0xf1, 0x80, 0x3f, 0xaa, 0x3a, 0xb7, 0xfe, 0x2b, 0xe9, 0x63, 0xdb, 0xc8, 0x15, 0x62, 0x32, 0x9}, nil,
		},
		{
			[]string{"alpha", "bravo", "charlie", "delta", "missing"}, fixtures[3], []byte{0x1c, 0xdc, 0x82, 0x74, 0x32, 0x4b, 0x84, 0x11, 0xbd, 0xf0, 0xe3, 0x9, 0x78, 0x9c, 0x95, 0x11}, nil,
		},
	}

	for i, tc := range testCases {
		hash, err := tc.event.HashUniqueKey(tc.keys)
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected an error for test case %d", i)
		} else {
			require.NoError(t, err, "expected no error for test case %d", i)
			require.Equal(t, tc.expected, hash, "expected hash equality for test case %d", i)
		}
	}
}

func TestHashUniqueField(t *testing.T) {
	fixtures, err := loadFixtures()
	require.NoError(t, err, "could not load event fixtures for strict hashing tests")

	testCases := []struct {
		fields   []string
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			[]string{"foo", "bar"}, &EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			nil, fixtures[0], nil, ErrNoFields,
		},
	}

	for i, tc := range testCases {
		hash, err := tc.event.HashUniqueField(tc.fields)
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected an error for test case %d", i)
		} else {
			require.NoError(t, err, "expected no error for test case %d", i)
			require.Equal(t, tc.expected, hash, "expected hash equality for test case %d", i)
		}
	}
}

const fixturePath = "testdata/events.json"

func loadFixtures() (_ []*EventWrapper, err error) {
	if _, err = os.Stat(fixturePath); os.IsNotExist(err) {
		if err = generateFixtures(); err != nil {
			return nil, err
		}
	}

	var f *os.File
	if f, err = os.Open(fixturePath); err != nil {
		return nil, err
	}
	defer f.Close()

	events := make([]*EventWrapper, 0)
	if err = json.NewDecoder(f).Decode(&events); err != nil {
		return nil, err
	}
	return events, nil
}

func generateFixtures() (err error) {
	block1 := randb(1024)
	block2 := randb(1024)
	block3 := randb(1892)

	key1 := rands(92)
	key2 := rands(112)
	key3 := rands(32)

	events := []*Event{
		{
			Data: block1,
			Metadata: map[string]string{
				"alpha":   key1,
				"bravo":   key2,
				"charlie": key3,
			},
			Mimetype: mimetype.ApplicationOctetStream,
			Type: &Type{
				Name:         "RandomData",
				MajorVersion: 2,
				MinorVersion: 11,
				PatchVersion: 3,
			},
			Created: timestamppb.Now(),
		},
		{
			Data: block2,
			Metadata: map[string]string{
				"alpha":   key1,
				"bravo":   key2,
				"charlie": key3,
			},
			Mimetype: mimetype.ApplicationOctetStream,
			Type: &Type{
				Name:         "RandomData",
				MajorVersion: 2,
				MinorVersion: 11,
				PatchVersion: 3,
			},
			Created: timestamppb.Now(),
		},
		{
			Data: block3,
			Metadata: map[string]string{
				"alpha":   key3,
				"bravo":   key1,
				"charlie": key2,
			},
			Mimetype: mimetype.ApplicationOctetStream,
			Type: &Type{
				Name:         "RandomData",
				MajorVersion: 2,
				MinorVersion: 11,
				PatchVersion: 3,
			},
			Created: timestamppb.Now(),
		},
		{
			Data: block1,
			Metadata: map[string]string{
				"alpha":   key2,
				"bravo":   key3,
				"charlie": key1,
			},
			Mimetype: mimetype.ApplicationOctetStream,
			Type: &Type{
				Name:         "RandomData",
				MajorVersion: 2,
				MinorVersion: 11,
				PatchVersion: 3,
			},
			Created: timestamppb.Now(),
		},
		{
			Data: block1,
			Metadata: map[string]string{
				"alpha":   key1,
				"bravo":   key2,
				"charlie": key3,
			},
			Mimetype: mimetype.ApplicationAvro,
			Type: &Type{
				Name:         "RandomData",
				MajorVersion: 2,
				MinorVersion: 11,
				PatchVersion: 3,
			},
			Created: timestamppb.Now(),
		},
		{
			Data: block1,
			Metadata: map[string]string{
				"alpha":   key1,
				"bravo":   key2,
				"charlie": key3,
			},
			Mimetype: mimetype.ApplicationOctetStream,
			Type:     nil,
			Created:  timestamppb.Now(),
		},
		{
			Data: block1,
			Metadata: map[string]string{
				"alpha":   key1,
				"bravo":   key2,
				"charlie": key3,
			},
			Mimetype: mimetype.ApplicationOctetStream,
			Type: &Type{
				Name:         "RandomData",
				MajorVersion: 2,
				MinorVersion: 12,
				PatchVersion: 0,
			},
			Created: timestamppb.Now(),
		},
		{
			Data:     block1,
			Metadata: nil,
			Mimetype: mimetype.ApplicationOctetStream,
			Type: &Type{
				Name:         "RandomData",
				MajorVersion: 2,
				MinorVersion: 11,
				PatchVersion: 3,
			},
			Created: timestamppb.Now(),
		},
	}

	var f *os.File
	if f, err = os.Create(fixturePath); err != nil {
		return err
	}
	defer f.Close()

	wraps := make([]*EventWrapper, 0, len(events))
	for _, event := range events {
		wrap := &EventWrapper{
			TopicId:   ulid.MustParse("01HBETJKP2ES10XXMK27M651GA").Bytes(),
			Committed: timestamppb.Now(),
		}
		wrap.Wrap(event)
		wraps = append(wraps, wrap)
	}

	if err = json.NewEncoder(f).Encode(wraps); err != nil {
		return err
	}
	return nil
}

func createRandomEvent(mime mimetype.MIME, etype *Type, meta map[string]string) *EventWrapper {
	event := &Event{
		Data:     randb(1024),
		Metadata: make(map[string]string),
		Mimetype: mime,
		Type:     etype,
		Created:  timestamppb.Now(),
	}

	for key, val := range meta {
		if val == "rand" {
			val = rands(96)
		}
		event.Metadata[key] = val
	}

	wrap := &EventWrapper{}
	wrap.Wrap(event)
	return wrap
}

func randb(s int) []byte {
	data := make([]byte, s)
	rand.Read(data)
	return data
}

func rands(s int) string {
	return base64.RawURLEncoding.EncodeToString(randb(s))
}

package api_test

import (
	"testing"

	. "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	// For the same event, but different policies, each hashing method should return
	// a different hash signature of the fixture.

}

func TestHashStrict(t *testing.T) {
	testCases := []struct {
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			&EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			createRandomEvent(mimetype.ApplicationAvro, nil), []byte{0xe1, 0xbc, 0xb6, 0xf5, 0x90, 0xa4, 0xf5, 0xb2, 0x77, 0x20, 0x52, 0xce, 0x4, 0x16, 0x38, 0x1a}, nil,
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
	testCases := []struct {
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			&EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			createRandomEvent(mimetype.ApplicationAvro, nil), []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, nil,
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
	testCases := []struct {
		keys     []string
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			nil, &EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			nil, createRandomEvent(mimetype.ApplicationAvro, nil), []byte{0xe1, 0xbc, 0xb6, 0xf5, 0x90, 0xa4, 0xf5, 0xb2, 0x77, 0x20, 0x52, 0xce, 0x4, 0x16, 0x38, 0x1a}, nil,
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
	testCases := []struct {
		keys     []string
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			nil, &EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			nil, createRandomEvent(mimetype.ApplicationAvro, nil), []byte{0xe1, 0xbc, 0xb6, 0xf5, 0x90, 0xa4, 0xf5, 0xb2, 0x77, 0x20, 0x52, 0xce, 0x4, 0x16, 0x38, 0x1a}, nil,
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
	testCases := []struct {
		fields   []string
		event    *EventWrapper
		expected []byte
		err      error
	}{
		{
			nil, &EventWrapper{Event: nil}, nil, ErrNoEvent,
		},
		{
			nil, createRandomEvent(mimetype.ApplicationAvro, nil), []byte{0xe1, 0xbc, 0xb6, 0xf5, 0x90, 0xa4, 0xf5, 0xb2, 0x77, 0x20, 0x52, 0xce, 0x4, 0x16, 0x38, 0x1a}, nil,
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

func createRandomEvent(mime mimetype.MIME, etype *Type) *EventWrapper {
	event := &Event{
		Mimetype: mime,
		Type:     etype,
	}

	wrap := &EventWrapper{}
	wrap.Wrap(event)
	return wrap
}

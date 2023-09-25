package api_test

import (
	"testing"

	. "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	"github.com/stretchr/testify/require"
)

func TestStrictHashing(t *testing.T) {
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

func createRandomEvent(mime mimetype.MIME, etype *Type) *EventWrapper {
	event := &Event{
		Mimetype: mime,
		Type:     etype,
	}

	wrap := &EventWrapper{}
	wrap.Wrap(event)
	return wrap
}

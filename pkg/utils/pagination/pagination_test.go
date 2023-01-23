package pagination_test

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPaginationToken(t *testing.T) {
	cursor := pagination.New("foo", "bar", 0)
	token, err := cursor.NextPageToken()
	require.NoError(t, err, "could not create next page token")
	require.Greater(t, len(token), 32, "the token created above should at least be 32 characters")
	require.Less(t, len(token), 64, "the token create above should be less than 64 characters")

	parsed, err := pagination.Parse(token)
	require.NoError(t, err, "could not parse token")
	require.True(t, proto.Equal(cursor, parsed), "parsed token should match cursor")

	cursor.Expires = nil
	_, err = cursor.NextPageToken()
	require.ErrorIs(t, err, pagination.ErrMissingExpiration, "should not be able to create a token with no expiration")

	data, err := proto.Marshal(cursor)
	require.NoError(t, err, "could not marshal protocol buffers")
	token = base64.RawURLEncoding.EncodeToString(data)
	_, err = pagination.Parse(token)
	require.ErrorIs(t, err, pagination.ErrMissingExpiration, "should not be able to parse a token with no expiration")

	_, err = pagination.Parse("")
	require.ErrorIs(t, err, pagination.ErrUnparsableToken, "should not be able to parse an invalid token")

	_, err = pagination.Parse(base64.RawStdEncoding.EncodeToString([]byte("badtokendata")))
	require.ErrorIs(t, err, pagination.ErrUnparsableToken, "should not be able to parse an invalid token")

	cursor.Expires = timestamppb.New(time.Now().Add(-5 * time.Minute))
	_, err = cursor.NextPageToken()
	require.ErrorIs(t, err, pagination.ErrCursorExpired, "should not be able to create an expired token")

	data, err = proto.Marshal(cursor)
	require.NoError(t, err, "could not marshal protocol buffers")
	token = base64.RawURLEncoding.EncodeToString(data)
	_, err = pagination.Parse(token)
	require.ErrorIs(t, err, pagination.ErrCursorExpired, "should not be able to parse an expired token")
}

func TestPaginationExpired(t *testing.T) {
	cursor := &pagination.Cursor{}
	expired, err := cursor.HasExpired()
	require.ErrorIs(t, err, pagination.ErrMissingExpiration)
	require.False(t, expired, "if err is not nil, expired should be false")

	cursor.Expires = timestamppb.New(time.Now().Add(5 * time.Minute))
	expired, err = cursor.HasExpired()
	require.NoError(t, err, "cursor should compute expiration without error")
	require.False(t, expired, "cusor should not be expired for 5 minutes")

	cursor.Expires = timestamppb.New(time.Now().Add(-5 * time.Minute))
	expired, err = cursor.HasExpired()
	require.NoError(t, err, "cursor should compute expiration without error")
	require.True(t, expired, "cusor should have expired 5 minutes ago")
}

func TestPaginationIsZero(t *testing.T) {
	cursor := &pagination.Cursor{}
	require.True(t, cursor.IsZero(), "empty cursor should be zero valued")

	cursor = pagination.New("foo", "bar", 0)
	require.False(t, cursor.IsZero(), "new cursor should not be zero valued")
}

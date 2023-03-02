package health_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/uptime/health"
	"github.com/stretchr/testify/require"
)

func TestBaseStatusKey(t *testing.T) {
	sid := uuid.MustParse("80628302-8b04-47eb-aca4-f3d32a6d661b")
	status := &health.BaseStatus{}

	// A service status is required
	_, err := status.Key()
	require.ErrorIs(t, err, health.ErrNoServiceID)

	// A timestamp is required
	status.SetServiceID(sid)
	_, err = status.Key()
	require.ErrorIs(t, err, health.ErrNoTimestamp)

	// Should be able to generate a valid key
	status.Timestamp = time.Date(2023, 04, 07, 12, 15, 26, 123456, time.UTC)

	key, err := status.Key()
	require.NoError(t, err)
	require.Len(t, key, 32)

	// The prefix should be the SID
	require.True(t, bytes.HasPrefix(key, sid[:]))

	// The suffix should be a valid ulid
	uid := &ulid.ULID{}
	err = uid.UnmarshalBinary(key[16:])
	require.NoError(t, err)

	require.Equal(t, ulid.Timestamp(status.Timestamp), uid.Time())
}

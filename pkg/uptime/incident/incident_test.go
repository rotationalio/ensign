package incident_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/uptime/incident"
	"github.com/stretchr/testify/require"
)

func TestIncidentGroupKey(t *testing.T) {
	group := &incident.Group{}

	_, err := group.Key()
	require.Error(t, err, "should not be able to create a group without a date")

	group.Date = time.Now()
	_, err = group.Key()
	require.Error(t, err, "should not be able to create a group with a full timestamp")

	group.Date = incident.Today()
	key, err := group.Key()
	require.NoError(t, err, "was unable to create a key")
	require.Len(t, key, 32)

	require.True(t, bytes.HasPrefix(key, incident.Prefix), "key does not have prefix")

	// Should be able to unmarshal the date from the end of the key
	ts := time.Time{}
	err = ts.UnmarshalBinary(key[len(incident.Prefix):])
	require.NoError(t, err, "could not unmarshal timestamp from key")
	require.True(t, ts.Equal(group.Date), "timestamp did not match group")
}

package broker_test

import (
	"testing"
	"time"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/broker"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/stretchr/testify/require"
)

func TestPublishResultAck(t *testing.T) {
	res := broker.PublishResult{
		LocalID:   rlid.Make(42).Bytes(),
		Committed: time.Date(2023, 8, 32, 14, 18, 23, 0, time.UTC),
	}

	// Test Ack from result
	rep := res.Reply()
	require.NotNil(t, rep.Embed)

	ack, ok := rep.Embed.(*api.PublisherReply_Ack)
	require.True(t, ok, "expected an ack to be returned")
	require.Equal(t, res.LocalID, ack.Ack.Id)
	require.True(t, res.Committed.Equal(ack.Ack.Committed.AsTime()))
}

func TestPublishResultNack(t *testing.T) {
	testCases := []struct {
		code     api.Nack_Code
		emsg     string
		expected string
	}{
		{api.Nack_INTERNAL, "this is a non-standard error message", "this is a non-standard error message"},
		{api.Nack_MAX_EVENT_SIZE_EXCEEDED, "", broker.NackMaxEventSizeExceeded},
		{api.Nack_TOPIC_UKNOWN, "", broker.NackTopicUnknown},
		{api.Nack_TOPIC_ARCHVIVED, "", broker.NackTopicArchived},
		{api.Nack_TOPIC_DELETED, "", broker.NackTopicDeleted},
		{api.Nack_PERMISSION_DENIED, "", broker.NackPermissionDenied},
		{api.Nack_CONSENSUS_FAILURE, "", broker.NackConsensusFailure},
		{api.Nack_SHARDING_FAILURE, "", broker.NackShardingFailure},
		{api.Nack_REDIRECT, "", broker.NackRedirect},
		{api.Nack_INTERNAL, "", broker.NackInternal},
		{api.Nack_UNPROCESSED, "", ""},
		{api.Nack_TIMEOUT, "", ""},
		{api.Nack_UNHANDLED_MIMETYPE, "", ""},
		{api.Nack_UNKNOWN_TYPE, "", ""},
		{api.Nack_DELIVER_AGAIN_ANY, "", ""},
		{api.Nack_DELIVER_AGAIN_NOT_ME, "", ""},
	}

	seq := rlid.Sequence(0)
	for i, tc := range testCases {
		res := broker.PublishResult{
			LocalID: seq.Next().Bytes(),
			Code:    tc.code,
			Error:   tc.emsg,
		}

		// Get nack from result
		rep := res.Reply()
		require.NotNil(t, rep.Embed, "test case %d failed", i)

		nack, ok := rep.Embed.(*api.PublisherReply_Nack)
		require.True(t, ok, "could not get nack from reply in test case %d", i)
		require.Equal(t, res.LocalID, nack.Nack.Id, "test case %d failed", i)
		require.Equal(t, tc.code, nack.Nack.Code, "test case %d failed", i)
		require.Equal(t, tc.expected, nack.Nack.Error, "test case %d failed", i)
	}
}

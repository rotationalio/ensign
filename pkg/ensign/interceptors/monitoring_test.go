package interceptors_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/ensign/interceptors"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	"github.com/stretchr/testify/require"
)

func TestParseMethod(t *testing.T) {
	tests := []struct {
		FullMethod string
		service    string
		rpc        string
	}{
		{mock.PublishRPC, "ensign.v1beta1.Ensign", "Publish"},
		{mock.SubscribeRPC, "ensign.v1beta1.Ensign", "Subscribe"},
		{mock.ListTopicsRPC, "ensign.v1beta1.Ensign", "ListTopics"},
		{mock.CreateTopicRPC, "ensign.v1beta1.Ensign", "CreateTopic"},
		{mock.DeleteTopicRPC, "ensign.v1beta1.Ensign", "DeleteTopic"},
		{mock.StatusRPC, "ensign.v1beta1.Ensign", "Status"},
	}

	for _, tc := range tests {
		service, rpc := interceptors.ParseMethod(tc.FullMethod)
		require.Equal(t, tc.service, service, "unexpected service parsed from %q", tc.FullMethod)
		require.Equal(t, tc.rpc, rpc, "unexpected rpc parsed from %q", tc.FullMethod)
	}
}

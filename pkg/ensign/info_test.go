package ensign_test

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	store "github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"google.golang.org/grpc/codes"
)

func (s *serverTestSuite) TestInfo() {
	require := s.Require()
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
		},
		OrgID: "01GKHJRF01YXHZ51YMMKV3RCMK",
	}

	ctx := context.Background()
	req := &api.InfoRequest{}

	// Should not be able to get project info when not authenticated
	_, err := s.client.Info(ctx, req)
	s.GRPCErrorIs(err, codes.Unauthenticated, "missing credentials")

	// Should not be able to get project info without the read topic and read metrics permissions
	token, err := s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")

	_, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// ProjectID is required in the claims
	claims.Permissions = []string{permissions.ReadTopics, permissions.ReadMetrics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for user")

	_, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// ProjectID must be valid and parsable
	claims.ProjectID = "foo"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for user")

	_, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// Claims should be valid and accepted from this point on in the test
	claims.ProjectID = "01GV6G705RV812J20S6RKJHVGE"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for user")

	// Internal error should be returned if database is not accessible
	s.store.UseError(store.ListTopics, errors.ErrIterReleased)
	_, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Internal, "unable to process project info request")

	// Empty result should be returned if there are no projects
	s.store.OnListTopics = func(ulid.ULID) iterator.TopicIterator {
		return store.NewTopicIterator(nil)
	}

	info, err := s.client.Info(ctx, req, mock.PerRPCToken(token))
	require.NoError(err, "could not fetch project info")
	require.Equal(ulid.MustParse("01GV6G705RV812J20S6RKJHVGE").Bytes(), info.ProjectId)
	require.Zero(info.NumTopics)
	require.Zero(info.NumReadonlyTopics)
	require.Zero(info.Events)
	require.Zero(info.Duplicates)
	require.Zero(info.DataSizeBytes)
	require.Zero(info.Topics)

	// Set up mock to return topics and topic infos
	err = s.store.UseFixture(store.ListTopics, "testdata/topics.json")
	require.NoError(err, "could not open topics fixture")

	s.store.OnTopicInfo, err = MockTopicInfo("testdata/topic_infos.json")
	require.NoError(err, "could not open topic infos fixture")

	// Test project info without filtering
	info, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	require.NoError(err, "could not fetch project info")
	require.Equal(ulid.MustParse("01GV6G705RV812J20S6RKJHVGE").Bytes(), info.ProjectId)
	require.Equal(uint64(5), info.NumTopics)
	require.Equal(uint64(2), info.NumReadonlyTopics)
	require.Equal(uint64(0x14df9), info.Events)
	require.Equal(uint64(0x14d), info.Duplicates)
	require.Equal(uint64(0x2451f07b), info.DataSizeBytes)
	require.Len(info.Topics, 5)

	// Test project info with filtering
	req.Topics = [][]byte{
		ulid.MustParse("01GTSN2NQV61P2R4WFYF1NF1JG").Bytes(),
		ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ").Bytes(),
		ulid.MustParse("01GTSMSX1M9G2Z45VGG4M12WC0").Bytes(),
	}

	info, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	require.NoError(err, "could not fetch project info")
	require.Equal(ulid.MustParse("01GV6G705RV812J20S6RKJHVGE").Bytes(), info.ProjectId)
	require.Equal(uint64(3), info.NumTopics) // TODO: is this the wrong number?
	require.Equal(uint64(1), info.NumReadonlyTopics)
	require.Equal(uint64(0x902), info.Events)
	require.Equal(uint64(0x12a), info.Duplicates)
	require.Equal(uint64(0xcc9386), info.DataSizeBytes)
	require.Len(info.Topics, 3)

	// Cannot filter invalid topic IDs
	req.Topics = [][]byte{[]byte("foo")}
	_, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "could not parse topic id in info request filter")
}

func (s *serverTestSuite) TestInfoSingleTopic() {
	// Should be able to get info for a single topic in a project
	// This test ensures that a Beacon requirement is fulfilled
	require := s.Require()
	ctx := context.Background()

	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
		},
		OrgID:       "01GKHJRF01YXHZ51YMMKV3RCMK",
		ProjectID:   "01GTSMZNRYXNAZQF5R8NHQ14NM",
		Permissions: []string{permissions.ReadTopics, permissions.ReadMetrics},
	}

	token, err := s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")

	err = s.store.UseFixture(store.ListTopics, "testdata/topics.json")
	require.NoError(err, "could not open topics fixture")

	s.store.OnTopicInfo, err = MockTopicInfo("testdata/topic_infos.json")
	require.NoError(err, "could not open topic infos fixture")

	projectID := ulids.MustParse("01GTSMZNRYXNAZQF5R8NHQ14NM")
	topicID := ulids.MustParse("01GTSN2NQV61P2R4WFYF1NF1JG")
	req := &api.InfoRequest{
		Topics: [][]byte{topicID[:]},
	}

	info, err := s.client.Info(ctx, req, mock.PerRPCToken(token))
	require.NoError(err, "could not execute info request")

	// Make assertions on the response based on the fixtures
	require.Equal(projectID.Bytes(), info.ProjectId, "expected topic Id to be part of the response")
	require.Equal(uint64(1), info.NumTopics)
	require.Equal(uint64(1), info.NumReadonlyTopics)
	require.Equal(uint64(1266), info.Events)
	require.Equal(uint64(298), info.Duplicates)
	require.Equal(uint64(10163651), info.DataSizeBytes)
	require.Len(info.Topics, 1, "expected only a single topic to be returned")

	topic := info.Topics[0]
	require.Equal(topicID.Bytes(), topic.TopicId)
	require.Equal(projectID.Bytes(), topic.ProjectId)
	require.Equal(info.Events, topic.Events)
	require.Equal(info.Duplicates, topic.Duplicates)
	require.Equal(info.DataSizeBytes, topic.DataSizeBytes)
	require.False(topic.Modified.AsTime().IsZero())
	require.Len(topic.Types, 3)
}

func MockTopicInfo(fixture string) (_ func(ulid.ULID) (*api.TopicInfo, error), err error) {
	var data []byte
	if data, err = os.ReadFile(fixture); err != nil {
		return nil, err
	}

	var infos map[string]*api.TopicInfo
	if infos, err = store.UnmarshalTopicInfoList(data); err != nil {
		return nil, err
	}

	return func(topicID ulid.ULID) (*api.TopicInfo, error) {
		ids := topicID.String()
		if info, ok := infos[ids]; ok {
			return info, nil
		}
		return nil, fmt.Errorf("topic id %s not found", ids)
	}, nil
}

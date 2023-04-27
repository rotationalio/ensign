package ensign_test

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	store "github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"google.golang.org/grpc/codes"
)

func (s *serverTestSuite) TestInfo() {
	require := s.Require()
	defer s.store.Reset()

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

	// Should not be able to get project info without the read topic permission
	token, err := s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for the user")

	_, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// ProjectID is required in the claims
	claims.Permissions = []string{permissions.ReadTopics}
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for user")

	_, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

	// ProjectID must be valid and parsable
	claims.ProjectID = "foo"
	token, err = s.quarterdeck.CreateAccessToken(claims)
	require.NoError(err, "could not create valid claims for user")

	_, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

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
	require.Equal("01GV6G705RV812J20S6RKJHVGE", info.ProjectId)
	require.Zero(info.Topics)
	require.Zero(info.ReadonlyTopics)
	require.Zero(info.Events)

	s.store.UseFixture(store.ListTopics, "testdata/topics.json")

	// Test project info without filtering
	info, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	require.NoError(err, "could not fetch project info")
	require.Equal("01GV6G705RV812J20S6RKJHVGE", info.ProjectId)
	require.Equal(uint64(4), info.Topics) // TODO: is this the wrong number?
	require.Equal(uint64(2), info.ReadonlyTopics)
	require.Equal(uint64(0x946), info.Events)

	// Test project info with filtering
	req.Topics = [][]byte{
		ulid.MustParse("01GTSN2NQV61P2R4WFYF1NF1JG").Bytes(),
		ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ").Bytes(),
		ulid.MustParse("01GTSMSX1M9G2Z45VGG4M12WC0").Bytes(),
	}

	info, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	require.NoError(err, "could not fetch project info")
	require.Equal("01GV6G705RV812J20S6RKJHVGE", info.ProjectId)
	require.Equal(uint64(3), info.Topics) // TODO: is this the wrong number?
	require.Equal(uint64(1), info.ReadonlyTopics)
	require.Equal(uint64(0x902), info.Events)

	// Cannot filter invalid topic IDs
	req.Topics = [][]byte{[]byte("foo")}
	_, err = s.client.Info(ctx, req, mock.PerRPCToken(token))
	s.GRPCErrorIs(err, codes.InvalidArgument, "could not parse topic id in info request filter")
}

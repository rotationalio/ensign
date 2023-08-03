package ensign_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	store "github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestStreamHandler(t *testing.T) {
	meta, err := store.Open(config.StorageConfig{ReadOnly: false, Testing: true})
	require.NoError(t, err, "could not open mock store for testing")

	stream := &mock.ServerStream{}
	handler := ensign.NewStreamHandler(ensign.UnknownStream, stream, meta)

	// Should not be able to get the ProjectID or AllowedTopics without authorization.
	_, err = handler.ProjectID()
	GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")

	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")

	// When there are no claims on the context, the handler should return unauthorized
	_, err = handler.Authorize("publisher")
	GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")

	// Add claims to the context for the remainder of the tests
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "http://localhost",
			Subject:   "01H6PGFB4T34D4WWEXQMAGJNMK",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		OrgID:       "01H6PGFG71N0AFEVTK3NJB71T9",
		ProjectID:   "01H6PGFTK2X53RGG2KMSGR2M61",
		Permissions: []string{"publisher", "subscriber"},
	}

	ctx := contexts.WithClaims(context.Background(), claims)
	stream.WithContext(ctx)

	// Should return unauthorized when the claims do not have the specific permission
	_, err = handler.Authorize("cookinthekitchen")
	GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")

	// When unauthorized, should not be able to get the ProjectID or AllowedTopics
	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")

	// Should be able to authorize with valid permissions
	actualClaims, err := handler.Authorize("publisher")
	require.NoError(t, err)
	require.Equal(t, claims, actualClaims)

	// After authorization, should be able to get the ProjectID
	projectID, err := handler.ProjectID()
	require.NoError(t, err)
	require.Equal(t, ulid.MustParse("01H6PGFTK2X53RGG2KMSGR2M61"), projectID)

	// When no topics are available, should get an error
	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.FailedPrecondition, "no topics available")

	// Internal error should be returned if topics cannot be fetched.
	meta.UseError(store.AllowedTopics, errors.New("this is a testing error"))
	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.Internal, "could not open unknown stream")

	meta.OnAllowedTopics = MockAllowedTopics
	meta.UseError(store.TopicName, errors.New("this is a testing error"))
	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.Internal, "could not open unknown stream")

	// Should be able to fetch topics
	meta.OnTopicName = MockTopicName
	group, err := handler.AllowedTopics()
	require.NoError(t, err)
	require.Equal(t, 4, group.Length())
}

func TestStreamHandlerInvalidProjectID(t *testing.T) {
	stream := &mock.ServerStream{}
	handler := ensign.NewStreamHandler(ensign.UnknownStream, stream, nil)

	testCases := []string{
		"", "notavalidulid", "00000000000000000000000000",
	}

	for _, tc := range testCases {
		claims := &tokens.Claims{ProjectID: tc, Permissions: []string{"publisher"}}
		stream.WithContext(contexts.WithClaims(context.Background(), claims))
		_, err := handler.Authorize("publisher")
		require.NoError(t, err)

		// Empty ProjectID not allowed
		projectID, err := handler.ProjectID()
		GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")
		require.True(t, ulids.IsZero(projectID))
	}
}

func MockAllowedTopics(projectID ulid.ULID) ([]ulid.ULID, error) {
	if ulids.IsZero(projectID) {
		return nil, errors.New("cannot get topics for empty ulid")
	}

	topics := make([]ulid.ULID, 0, 4)
	if projectID.Compare(ulid.MustParse("01H6PGFTK2X53RGG2KMSGR2M61")) == 0 {
		topics = append(topics,
			ulid.MustParse("01H6XTAPN0HZ1S7KEPFBF1MMPX"),
			ulid.MustParse("01H6XTAVNM21F6JXNGAJF1SJ4S"),
			ulid.MustParse("01H6XTB1780D2YKMC2MBNZ4V2X"),
			ulid.MustParse("01H6XTB5DS8YG0YZEVQ385QRTB"),
		)
	}
	return topics, nil
}

func MockTopicName(topicID ulid.ULID) (string, error) {
	switch topicID.String() {
	case "01H6XTAPN0HZ1S7KEPFBF1MMPX":
		return "example-topic-1", nil
	case "01H6XTAVNM21F6JXNGAJF1SJ4S":
		return "example-topic-2", nil
	case "01H6XTB1780D2YKMC2MBNZ4V2X":
		return "example-topic-3", nil
	case "01H6XTB5DS8YG0YZEVQ385QRTB":
		return "example-topic-4", nil
	default:
		return "", errors.New("unknown topic id")
	}
}

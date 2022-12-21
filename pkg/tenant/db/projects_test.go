package db_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/stretchr/testify/require"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProjectModel(t *testing.T) {
	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project-example",
		Created:  time.Unix(1668660681, 0).In(time.UTC),
		Modified: time.Unix(1668661302, 0).In(time.UTC),
	}

	key, err := project.Key()
	require.NoError(t, err, "could not marshal the project")
	require.Equal(t, project.TenantID[:], key[0:16], "unexpected marshaling of the tenant id half of the key")
	require.Equal(t, project.ID[:], key[16:], "unexpected marshaling of the project id half of the key")

	require.Equal(t, db.ProjectNamespace, project.Namespace(), "unexpected project namespace")

	// Test marshal and unmarshal
	data, err := project.MarshalValue()
	require.NoError(t, err, "could not marshal the project")

	other := &db.Project{}
	err = other.UnmarshalValue(data)
	require.NoError(t, err, "could not unmarshal the project")

	ProjectsEqual(t, project, other, "unmarshaled project does not match marshaled project")
}

func (s *dbTestSuite) TestCreateProject() {
	require := s.Require()
	ctx := context.Background()
	project := &db.Project{Name: "project-example"}

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.ProjectNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err := db.CreateProject(ctx, project)
	require.NoError(err, "could not create project")

	require.NotEqual("", project.ID, "expected non-zero ulid to be populated")
}

/* func (s *dbTestSuite) TestListProjects() {
	require := s.Require()
	ctx := context.Background()

	prefix := []byte("test")
	namespace := "testing"

	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back some data and terminate
		for i := 0; i < 7; i++ {
			stream.Send(&pb.KVPair{
				Key:       []byte(fmt.Sprintf("key %d", i)),
				Value:     []byte(fmt.Sprintf("value %d", i)),
				Namespace: in.Namespace,
			})
		}

		return nil
	}

	values, err := db.ListProjects(ctx, prefix, namespace)
	require.NoError(err, "error returned from list request")
	require.Len(values, 7, "unexpected number of values returned")
} */

func (s *dbTestSuite) TestRetrieveProject() {
	require := s.Require()
	ctx := context.Background()
	project := &db.Project{
		ID:   ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name: "project-example",
	}

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.ProjectNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		if !bytes.Equal(in.Key, project.ID[:]) {
			return nil, status.Error(codes.NotFound, "project not found")
		}

		// Marshal the data with msgpack
		data, err := project.MarshalValue()
		require.NoError(err, "could not marshal data")

		other := &db.Project{}
		err = other.UnmarshalValue(data)
		require.NoError(err, "could not unmarshal data")

		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "could not read fixture: %s", err)
		}

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	project, err := db.RetrieveProject(ctx, project.ID)
	require.NoError(err, "could not retrieve project")

	require.Equal(ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"), project.ID)
	require.Equal("project-example", project.Name)

	// Test NotFound path
	_, err = db.RetrieveProject(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestUpdateProject() {
	require := s.Require()
	ctx := context.Background()
	project := &db.Project{
		ID:   ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name: "project-example",
	}

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.ProjectNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		if !bytes.Equal(in.Key, project.ID[:]) {
			return nil, status.Error(codes.NotFound, "project not found")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err := db.UpdateProject(ctx, project)
	require.NoError(err, "could not update project")
}

func (s *dbTestSuite) TestDeleteProject() {
	require := s.Require()
	ctx := context.Background()
	projectID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	s.mock.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.ProjectNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Delete request")
		}

		if !bytes.Equal(in.Key, projectID[:]) {
			return nil, status.Error(codes.NotFound, "project not found")
		}

		return &pb.DeleteReply{
			Success: true,
		}, nil
	}

	err := db.DeleteProject(ctx, projectID)
	require.NoError(err, "could not delete project")

	// Test NotFound path
	err = db.DeleteProject(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}

func ProjectsEqual(t *testing.T, expected, actual *db.Project, msgAndArgs ...interface{}) {
	require.Equal(t, expected.ID, actual.ID, msgAndArgs...)
	require.Equal(t, expected.Name, actual.Name, msgAndArgs...)
	require.True(t, expected.Created.Equal(actual.Created), msgAndArgs...)
	require.True(t, expected.Modified.Equal(actual.Modified), msgAndArgs...)
}

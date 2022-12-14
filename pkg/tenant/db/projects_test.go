package db_test

import (
	"bytes"
	"context"
	"fmt"
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
		Created:  time.Unix(1670424445, 0).In(time.UTC),
		Modified: time.Unix(1670424445, 0).In(time.UTC),
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
	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project-example",
	}

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

	// Verify that below fields have been populated.
	require.NotZero(project.ID, "expected non-zero ulid to be populated")
	require.NotZero(project.Created, "expected project to have a created timestamp")
	require.Equal(project.Created, project.Modified, "expected the same created and modified timestamp")
}

func (s *dbTestSuite) TestRetrieveProject() {
	require := s.Require()
	ctx := context.Background()
	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project-example",
		Created:  time.Unix(1670424445, 0),
		Modified: time.Unix(1670424445, 0),
	}

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.ProjectNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}
		if !bytes.Equal(in.Key[16:], project.ID[:]) {
			return nil, status.Error(codes.NotFound, "project not found")
		}

		// TODO: Add msgpack fixture helpers

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

	// Verify the fields below have been populated.
	require.Equal(ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"), project.ID, "expected project id to match")
	require.Equal("project-example", project.Name, "expected project name to match")
	require.Equal(time.Unix(1670424445, 0), project.Created, "expected created timestamp to not have changed")
	require.True(time.Unix(1670424444, 0).Before(project.Modified), "expected modified timestamp to be updated")

	// Test NotFound path
	// TODO: Use crypto rand and monotonic entropy with ulid.New
	_, err = db.RetrieveProject(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestListProjects() {
	require := s.Require()
	ctx := context.Background()

	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project-example",
		Created:  time.Unix(1670424445, 0).In(time.UTC),
		Modified: time.Unix(1670424445, 0).In(time.UTC),
	}

	prefix := project.TenantID[:]
	namespace := "projects"

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

	values, err := db.List(ctx, prefix, namespace)
	require.NoError(err, "could not get project values")
	require.Len(values, 7)

	projects := make([]*db.Project, 0, len(values))
	projects = append(projects, project)
	require.Len(projects, 1)

	_, err = db.ListProjects(ctx, project.TenantID)
	require.Error(err, "could not list projects")
}

func (s *dbTestSuite) TestUpdateProject() {
	require := s.Require()
	ctx := context.Background()
	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project-example",
		Created:  time.Unix(1668660681, 0),
		Modified: time.Unix(1668660681, 0),
	}

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.ProjectNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}
		if !bytes.Equal(in.Key[16:], project.ID[:]) {
			return nil, status.Error(codes.NotFound, "project not found")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err := db.UpdateProject(ctx, project)
	require.NoError(err, "could not update project")

	// The fields below should have been populated
	require.Equal(ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"), project.ID, "project ID should not have changed")
	require.Equal(time.Unix(1668660681, 0), project.Created, "expected created timestamp to not have changed")
	require.True(time.Unix(1668660681, 0).Before(project.Modified), "expected modified timestamp to be updated")

	// Test NotFound path
	// TODO: Use crypto rand and monotonic entropy with ulid.New
	err = db.UpdateProject(ctx, &db.Project{ID: ulid.Make()})
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestDeleteProject() {
	require := s.Require()
	ctx := context.Background()
	projectID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	s.mock.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.ProjectNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Delete request")
		}
		if !bytes.Equal(in.Key[16:], projectID[:]) {
			return nil, status.Error(codes.NotFound, "project not found")
		}

		return &pb.DeleteReply{
			Success: true,
		}, nil
	}

	err := db.DeleteProject(ctx, projectID)
	require.NoError(err, "could not delete project")

	// Test NotFound path
	// TODO: Use crypto rand and monotonic entropy with ulid.New
	err = db.DeleteProject(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}

// ProjectsEqual tests assertions in the ProjectModel.
// Note: require.True compares the actual.Created and actual.Modified
// timestamps because MsgPack does not preserve time zone information.
func ProjectsEqual(t *testing.T, expected, actual *db.Project, msgAndArgs ...interface{}) {
	require.Equal(t, expected.ID, actual.ID, msgAndArgs...)
	require.Equal(t, expected.Name, actual.Name, msgAndArgs...)
	require.True(t, expected.Created.Equal(actual.Created), msgAndArgs...)
	require.True(t, expected.Modified.Equal(actual.Modified), msgAndArgs...)
}

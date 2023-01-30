package db_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/stretchr/testify/require"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProjectModel(t *testing.T) {
	project := &db.Project{
		OrgID:    ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project001",
		Created:  time.Unix(1670424445, 0).In(time.UTC),
		Modified: time.Unix(1670424445, 0).In(time.UTC),
	}

	err := project.Validate()
	require.NoError(t, err, "could not validate project data")

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
		OrgID:    ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		Name:     "project001",
	}

	err := project.Validate()
	require.NoError(err, "could not validate project data")

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.ProjectNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err = db.CreateProject(ctx, project)
	require.NoError(err, "could not create project")

	// Verify that below fields have been populated.
	require.NotEmpty(project.ID, "expected non-zero ulid to be populated")
	require.NotEmpty(project.Name, "project name is required")
	require.NotZero(project.Created, "expected project to have a created timestamp")
	require.Equal(project.Created, project.Modified, "expected the same created and modified timestamp")
}

func (s *dbTestSuite) TestRetrieveProject() {
	require := s.Require()
	ctx := context.Background()
	project := &db.Project{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project001",
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
	require.Equal("project001", project.Name, "expected project name to match")
	require.Equal(time.Unix(1670424445, 0), project.Created, "expected created timestamp to not have changed")
	require.True(time.Unix(1670424444, 0).Before(project.Modified), "expected modified timestamp to be updated")

	// Test NotFound path
	_, err = db.RetrieveProject(ctx, ulids.New())
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestListProjects() {
	require := s.Require()
	ctx := context.Background()
	tenantID := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")

	projects := []*db.Project{
		{
			TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:       ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5FA2G"),
			Name:     "project001",
			Created:  time.Unix(1670424445, 0),
			Modified: time.Unix(1670424445, 0),
		},
		{
			TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:       ulid.MustParse("01GQ38JP6CCWPNDS6KG5WDA59T"),
			Name:     "project002",
			Created:  time.Unix(1673659941, 0),
			Modified: time.Unix(1673659941, 0),
		},
		{
			TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:       ulid.MustParse("01GQ38K6YPE0ZA9ADC2BGSVWRM"),
			Name:     "project003",
			Created:  time.Unix(1674073941, 0),
			Modified: time.Unix(1674073941, 0),
		},
	}

	prefix := tenantID[:]
	namespace := "projects"

	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back some data and terminate
		for i, project := range projects {
			data, err := project.MarshalValue()
			require.NoError(err, "could not marshal data")
			stream.Send(&pb.KVPair{
				Key:       []byte(fmt.Sprintf("key %d", i)),
				Value:     data,
				Namespace: in.Namespace,
			})
		}
		return nil
	}

	values, err := db.List(ctx, prefix, namespace)
	require.NoError(err, "could not get project values")
	require.Len(values, 3, "expected 3 values")

	rep, err := db.ListProjects(ctx, tenantID)
	require.NoError(err, "could not list projects")
	require.Len(rep, 3, "expected 3 projects")

	for i := range projects {
		require.Equal(projects[i].ID, rep[i].ID, "expected project id to match")
		require.Equal(projects[i].Name, rep[i].Name, "expected project name to match")
	}
}

func (s *dbTestSuite) TestUpdateProject() {
	require := s.Require()
	ctx := context.Background()
	project := &db.Project{
		OrgID:    ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "project001",
		Created:  time.Unix(1668660681, 0),
		Modified: time.Unix(1668660681, 0),
	}

	err := project.Validate()
	require.NoError(err, "could not validate project data")

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

	err = db.UpdateProject(ctx, project)
	require.NoError(err, "could not update project")

	// The fields below should have been populated
	require.Equal(ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"), project.ID, "project ID should not have changed")
	require.Equal(time.Unix(1668660681, 0), project.Created, "expected created timestamp to not have changed")
	require.True(time.Unix(1668660681, 0).Before(project.Modified), "expected modified timestamp to be updated")

	// Test NotFound path
	err = db.UpdateProject(ctx, &db.Project{OrgID: ulids.New(), TenantID: ulids.New(), ID: ulids.New(), Name: "project002"})
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
	err = db.DeleteProject(ctx, ulids.New())
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

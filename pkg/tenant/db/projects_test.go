package db_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/stretchr/testify/require"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProjectModel(t *testing.T) {
	project := &db.Project{
		ID:   ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name: "project-example",
	}

	key, err := project.Key()
	require.NoError(t, err, "could not marshal the project")
	require.Equal(t, project.ID[:], key, "unexpected marshaling of the key")

	require.Equal(t, db.ProjectNamespace, project.Namespace(), "unexpected project namespace")

	// Test marshal and unmarshal
	data, err := project.MarshalValue()
	require.NoError(t, err, "could not marshal the project")

	other := &db.Project{}
	err = other.UnmarshalValue(data)
	require.NoError(t, err, "could not unmarshal the project")

	require.Equal(t, project, other, "unmarshaled project does not match marshaled project")

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

func (s *dbTestSuite) TestRetrieveProject() {
	require := s.Require()
	ctx := context.Background()
	projectID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.ProjectNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		if !bytes.Equal(in.Key, projectID[:]) {
			return nil, status.Error(codes.NotFound, "project not found")
		}

		// Load fixture from disk
		data, err := os.ReadFile("testdata/project.json")
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "could not read fixture: %s", err)
		}

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	project, err := db.RetrieveProject(ctx, projectID)
	require.NoError(err, "could not retrieve project")

	require.Equal(projectID, project.ID)
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

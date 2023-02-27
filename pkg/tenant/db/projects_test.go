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

func TestProjectValidate(t *testing.T) {
	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	project := &db.Project{
		OrgID: orgID,
		Name:  "Hello World",
	}

	// Test missing orgID
	project.OrgID = ulids.Null
	require.ErrorIs(t, project.Validate(), db.ErrMissingOrgID, "expected missing org id error")

	// Test missing name
	project.OrgID = orgID
	project.Name = ""
	require.ErrorIs(t, project.Validate(), db.ErrMissingProjectName, "expected missing name error")

	// Test invalid name
	project.Name = "Hello World;"
	require.ErrorIs(t, project.Validate(), db.ErrInvalidProjectName, "expected invalid name error")

	// Test valid project
	project.Name = "Hello World"
	require.NoError(t, project.Validate(), "expected valid project")
}

func TestProjectKey(t *testing.T) {
	// Test that the key can't be created when ID is missing
	id := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")
	tenantID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	project := &db.Project{
		TenantID: tenantID,
	}
	_, err := project.Key()
	require.ErrorIs(t, err, db.ErrMissingID, "expected missing project id error")

	// Test that the key can't be created when TenantID is missing
	project.ID = id
	project.TenantID = ulids.Null
	_, err = project.Key()
	require.ErrorIs(t, err, db.ErrMissingTenantID, "expected missing tenant id error")

	// Test that the key is created correctly
	project.TenantID = tenantID
	key, err := project.Key()
	require.NoError(t, err, "could not marshal the project")
	require.Equal(t, project.TenantID[:], key[0:16], "unexpected marshaling of the tenant id half of the key")
	require.Equal(t, project.ID[:], key[16:], "unexpected marshaling of the project id half of the key")
}

func TestProjectKeyModel(t *testing.T) {
	key := &db.ProjectKey{
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		TenantID: ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
	}

	// Key method should return the ID
	id, err := key.Key()
	require.NoError(t, err, "could not retrieve key from model")
	require.Equal(t, key.ID[:], id, "unexpected ID returned from model")

	// Namespace should return the correct namespace
	require.Equal(t, db.ProjectKeysNamespace, key.Namespace(), "unexpected namespace returned from model")

	// Should be able to marshal and unmarshal the model
	data, err := key.MarshalValue()
	require.NoError(t, err, "could not marshal the project key model")

	other := &db.ProjectKey{}
	require.NoError(t, other.UnmarshalValue(data), "could not unmarshal the project key model")
	require.Equal(t, key, other, "unmarshaled project key model does not match marshaled project key model")
}

func (s *dbTestSuite) TestCreateTenantProject() {
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
		switch len(in.Key) {
		case 16:
			if in.Namespace != db.ProjectKeysNamespace {
				return nil, status.Errorf(codes.InvalidArgument, "bad key for namespace %s", in.Namespace)
			}
		case 32:
			if in.Namespace != db.ProjectNamespace {
				return nil, status.Errorf(codes.InvalidArgument, "bad key for namespace %s", in.Namespace)
			}
		default:
			return nil, status.Errorf(codes.InvalidArgument, "bad key length %d", len(in.Key))
		}

		if len(in.Value) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "empty value")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err = db.CreateTenantProject(ctx, project)
	require.NoError(err, "could not create project")

	// Verify that below fields have been populated.
	require.NotEmpty(project.ID, "expected non-zero ulid to be populated")
	require.NotEmpty(project.Name, "project name is required")
	require.NotZero(project.Created, "expected project to have a created timestamp")
	require.Equal(project.Created, project.Modified, "expected the same created and modified timestamp")

	// Should error if tenant ID is not set.
	project.TenantID = ulids.Null
	require.ErrorIs(db.CreateTenantProject(ctx, project), db.ErrMissingTenantID, "expected missing tenant id error")

	// Should error if project is not valid.
	project.TenantID = ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")
	project.Name = ""
	require.ErrorIs(db.CreateTenantProject(ctx, project), db.ErrMissingProjectName, "expected missing project name error")

	// Test trtl returning not found
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return nil, status.Error(codes.NotFound, "not found")
	}
	project.Name = "project001"
	require.ErrorIs(db.CreateTenantProject(ctx, project), db.ErrNotFound, "expected not found error")
}

func (s *dbTestSuite) TestCreateProject() {
	require := s.Require()
	ctx := context.Background()
	project := &db.Project{
		TenantID: ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		Name:     "project001",
	}

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		switch len(in.Key) {
		case 16:
			if in.Namespace != db.ProjectKeysNamespace {
				return nil, status.Errorf(codes.InvalidArgument, "bad key for namespace %s", in.Namespace)
			}
		case 32:
			if in.Namespace != db.ProjectNamespace {
				return nil, status.Errorf(codes.InvalidArgument, "bad key for namespace %s", in.Namespace)
			}
		default:
			return nil, status.Errorf(codes.InvalidArgument, "bad key length %d", len(in.Key))
		}

		if len(in.Value) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "empty value")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err := db.CreateProject(ctx, project)
	require.NoError(err, "could not create project")

	// Verify that below fields have been populated.
	require.NotEmpty(project.ID, "expected non-zero ulid to be populated")
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
	key := &db.ProjectKey{
		ID:       project.ID,
		TenantID: project.TenantID,
	}

	projectData, err := project.MarshalValue()
	require.NoError(err, "could not marshal project data")

	keyData, err := key.MarshalValue()
	require.NoError(err, "could not marshal project key data")

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		// TODO: Add msgpack fixture helpers
		var data []byte
		switch len(in.Key) {
		case 16:
			if in.Namespace != db.ProjectKeysNamespace {
				return nil, status.Error(codes.InvalidArgument, "expected 16 byte key for project keys namespace")
			}

			if !bytes.Equal(in.Key, key.ID[:]) {
				return nil, status.Error(codes.NotFound, "project key not found")
			}

			data = keyData
		case 32:
			if in.Namespace != db.ProjectNamespace {
				return nil, status.Error(codes.InvalidArgument, "expected 32 byte key for project namespace")
			}

			if !bytes.Equal(in.Key[:16], key.TenantID[:]) {
				return nil, status.Error(codes.NotFound, "project not found")
			}

			if !bytes.Equal(in.Key[16:], key.ID[:]) {
				return nil, status.Error(codes.NotFound, "project not found")
			}

			data = projectData
		default:
			return nil, status.Error(codes.InvalidArgument, "invalid key")
		}

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	project, err = db.RetrieveProject(ctx, project.ID)
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
	key := &db.ProjectKey{
		ID:       project.ID,
		TenantID: project.TenantID,
	}

	err := project.Validate()
	require.NoError(err, "could not validate project data")

	projectData, err := project.MarshalValue()
	require.NoError(err, "could not marshal project data")

	keyData, err := key.MarshalValue()
	require.NoError(err, "could not marshal project key data")

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		var data []byte
		switch len(in.Key) {
		case 16:
			if in.Namespace != db.ProjectKeysNamespace {
				return nil, status.Errorf(codes.InvalidArgument, "bad key for namespace %s", in.Namespace)
			}

			data = keyData
		case 32:
			if in.Namespace != db.ProjectNamespace {
				return nil, status.Errorf(codes.InvalidArgument, "bad key for namespace %s", in.Namespace)
			}

			data = projectData
		default:
			return nil, status.Errorf(codes.InvalidArgument, "bad key length %d", len(in.Key))
		}

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		switch len(in.Key) {
		case 16:
			if in.Namespace != db.ProjectKeysNamespace {
				return nil, status.Errorf(codes.InvalidArgument, "bad key for namespace %s", in.Namespace)
			}
		case 32:
			if in.Namespace != db.ProjectNamespace {
				return nil, status.Errorf(codes.InvalidArgument, "bad key for namespace %s", in.Namespace)
			}
		default:
			return nil, status.Errorf(codes.InvalidArgument, "bad key length %d", len(in.Key))
		}

		if len(in.Value) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "empty value")
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

	// If created timestamp is missing then it should be updated
	project.Created = time.Time{}
	require.NoError(db.UpdateProject(ctx, project), "could not update project")
	require.Equal(project.Modified, project.Created, "expected created timestamp to be updated")

	// Should fail if project ID is missing
	project.ID = ulid.ULID{}
	require.ErrorIs(db.UpdateProject(ctx, project), db.ErrMissingID, "expected error for missing project ID")

	// Should fail if project is invalid
	project.ID = ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")
	project.Name = ""
	require.ErrorIs(db.UpdateProject(ctx, project), db.ErrMissingProjectName, "expected missing project name error")

	// Test NotFound path
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	err = db.UpdateProject(ctx, &db.Project{OrgID: ulids.New(), TenantID: ulids.New(), ID: ulids.New(), Name: "project002"})
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestDeleteProject() {
	require := s.Require()
	ctx := context.Background()
	key := &db.ProjectKey{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
	}

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.ProjectKeysNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}

		if in.Namespace != db.ProjectKeysNamespace {
			return nil, status.Error(codes.InvalidArgument, "expected project keys namespace")
		}

		if !bytes.Equal(in.Key, key.ID[:]) {
			return nil, status.Error(codes.NotFound, "project key not found")
		}

		data, err := key.MarshalValue()
		require.NoError(err, "could not marshal project key")

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	s.mock.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
		switch len(in.Key) {
		case 16:
			if in.Namespace != db.ProjectKeysNamespace {
				return nil, status.Errorf(codes.InvalidArgument, "bad key for namespace %s", in.Namespace)
			}

			if !bytes.Equal(in.Key, key.ID[:]) {
				return nil, status.Error(codes.NotFound, "project key not found")
			}
		case 32:
			if in.Namespace != db.ProjectNamespace {
				return nil, status.Errorf(codes.InvalidArgument, "bad key for namespace %s", in.Namespace)
			}

			if !bytes.Equal(in.Key[:16], key.TenantID[:]) {
				return nil, status.Error(codes.NotFound, "project not found")
			}

			if !bytes.Equal(in.Key[16:], key.ID[:]) {
				return nil, status.Error(codes.NotFound, "project not found")
			}
		default:
			return nil, status.Errorf(codes.InvalidArgument, "bad key length %d", len(in.Key))
		}

		return &pb.DeleteReply{
			Success: true,
		}, nil
	}

	err := db.DeleteProject(ctx, key.ID)
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

package db_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTenantModel(t *testing.T) {
	tenant := &db.Tenant{
		OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ID:              ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:            "tenant001",
		EnvironmentType: "prod",
		Created:         time.Unix(1668660681, 0).In(time.UTC),
		Modified:        time.Unix(1668661302, 0).In(time.UTC),
	}

	err := tenant.Validate()
	require.NoError(t, err, "could not validate tenant data")

	key, err := tenant.Key()
	require.NoError(t, err, "could not marshal the key")
	require.Equal(t, tenant.OrgID[:], key[0:16], "unexpected marshaling of org id half of the key")
	require.Equal(t, tenant.ID[:], key[16:], "unexpected marshaling of the tenant id half of the key")

	require.Equal(t, db.TenantNamespace, tenant.Namespace(), "unexpected tenant namespace")

	// Test marshal and unmarshal
	data, err := tenant.MarshalValue()
	require.NoError(t, err, "could not marshal the tenant")

	other := &db.Tenant{}
	err = other.UnmarshalValue(data)
	require.NoError(t, err, "could not unmarshal the tenant")

	TenantsEqual(t, tenant, other, "unmarshaled tenant does not match marshaled tenant")
}

func TestTenantValidate(t *testing.T) {
	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	tenant := &db.Tenant{
		OrgID:           orgID,
		Name:            "tenant001",
		EnvironmentType: "dev",
	}

	// Test missing orgID
	tenant.OrgID = ulid.ULID{}
	require.ErrorIs(t, tenant.Validate(), db.ErrMissingOrgID, "expected missing org id error")

	// Test missing name
	tenant.OrgID = orgID
	tenant.Name = ""
	require.ErrorIs(t, tenant.Validate(), db.ErrMissingTenantName, "expected missing name error")

	// Test missing environment type
	tenant.Name = "tenant001"
	tenant.EnvironmentType = ""
	require.ErrorIs(t, tenant.Validate(), db.ErrMissingEnvType, "expected missing environment type error")

	// Test invalid name
	tenant.EnvironmentType = "dev"
	tenant.Name = "tenant*001"
	require.ErrorIs(t, tenant.Validate(), db.ErrInvalidTenantName, "expected invalid name error")

	// Valid tenant
	tenant.Name = "tenant001"
	require.NoError(t, tenant.Validate(), "expected valid tenant")
}

func TestTenantKey(t *testing.T) {
	// Test that the key can't be created without an ID
	id := ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	tenant := &db.Tenant{
		OrgID: orgID,
	}

	_, err := tenant.Key()
	require.ErrorIs(t, err, db.ErrMissingID, "expected missing tenant id error")

	// Test that the key can't be created without an orgID
	tenant.ID = id
	tenant.OrgID = ulid.ULID{}
	_, err = tenant.Key()
	require.ErrorIs(t, err, db.ErrMissingOrgID, "expected missing org id error")

	// Test that the key is created correctly
	tenant.OrgID = orgID
	key, err := tenant.Key()
	require.NoError(t, err, "could not marshal the key")
	require.Equal(t, tenant.OrgID[:], key[0:16], "unexpected marshaling of org id half of the key")
	require.Equal(t, tenant.ID[:], key[16:], "unexpected marshaling of the tenant id half of the key")
}

func (s *dbTestSuite) TestCreateTenant() {
	require := s.Require()
	ctx := context.Background()
	tenant := &db.Tenant{
		OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ID:              ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:            "tenant001",
		EnvironmentType: "prod",
	}

	err := tenant.Validate()
	require.NoError(err, "could not validate tenant data")

	// Call mock trtl database
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (out *pb.PutReply, err error) {
		switch in.Namespace {
		case db.TenantNamespace:
			return &pb.PutReply{Success: true}, nil
		case db.OrganizationNamespace:
			return &pb.PutReply{}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	err = db.CreateTenant(ctx, tenant)
	require.NoError(err, "could not create tenant")

	// Fields should have been populated
	require.NotEmpty(tenant.ID, "expected non-zero ulid to be populated")
	require.NotEmpty(tenant.Name, "tenant name is required")
	require.NotEmpty(tenant.EnvironmentType, "tenant environment type is required")
	require.NotZero(tenant.Created, "expected tenant to have a created timestamp")
	require.Equal(tenant.Created, tenant.Modified, "expected the same created and modified timestamp")
}

func (s *dbTestSuite) TestListTenants() {
	require := s.Require()
	ctx := context.Background()
	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")

	tenants := []*db.Tenant{
		{
			OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
			ID:              ulid.MustParse("01GQ38QWNR7MYQXSQ682PJQM7T"),
			Name:            "tenant001",
			EnvironmentType: "prod",
			Created:         time.Unix(1668660681, 0),
			Modified:        time.Unix(1668661302, 0),
		},

		{
			OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
			ID:              ulid.MustParse("01GQ38QMW7FGKG7AN1TVJTGHJA"),
			Name:            "tenant002",
			EnvironmentType: "staging",
			Created:         time.Unix(1673659941, 0),
			Modified:        time.Unix(1673659941, 0),
		},

		{
			OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
			ID:              ulid.MustParse("01GQ38QBN8XYA2S0KTW8AHPXHR"),
			Name:            "tenant003",
			EnvironmentType: "dev",
			Created:         time.Unix(1674073941, 0),
			Modified:        time.Unix(1674073941, 0),
		},
	}

	prefix := orgID[:]
	namespace := "tenants"

	// Configure trtl to return the tenant records on cursor
	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back the tenant data
		for _, tenant := range tenants {
			key, err := tenant.Key()
			if err != nil {
				return status.Error(codes.FailedPrecondition, "could not marshal tenant key for trtl response")
			}
			data, err := tenant.MarshalValue()
			if err != nil {
				return status.Error(codes.FailedPrecondition, "could not marshal tenant data for trtl response")
			}
			stream.Send(&pb.KVPair{
				Key:       key,
				Value:     data,
				Namespace: in.Namespace,
			})
		}

		return nil
	}

	s.Run("Single Page", func() {
		// If all the results are on a single page then the next cursor is nil
		cursor := &pg.Cursor{
			PageSize: 100,
		}

		rep, cursor, err := db.ListTenants(ctx, orgID, cursor)
		require.NoError(err, "could not list tenants")
		require.Len(rep, 3, "expected 3 tenants")
		require.Nil(cursor, "next page cursor should not be set since there isn't a next page")

		for i := range tenants {
			require.Equal(tenants[i].ID, rep[i].ID, "expected tenant id to match")
			require.Equal(tenants[i].EnvironmentType, rep[i].EnvironmentType, "expected tenant environment type to match")
			require.Equal(tenants[i].Name, rep[i].Name, "expected tenant name to match")
		}
	})

	s.Run("Multiple Pages", func() {
		// If results are on multiple pages then the next cursor is not nil
		cursor := &pg.Cursor{
			PageSize: 2,
		}
		rep, cursor, err := db.ListTenants(ctx, orgID, cursor)
		require.NoError(err, "could not list tenants")
		require.Len(rep, 2, "expected 2 tenants on the first page")
		require.NotNil(cursor, "expected cursor to be not nil because there is a next page")

		// Ensure the new start index is correct
		startBytes, err := tenants[2].Key()
		require.NoError(err, "could not marshal tenants key")
		startKey := &db.Key{}
		require.NoError(startKey.UnmarshalValue(startBytes), "could not unmarshal start key")
		startString, err := startKey.String()
		require.NoError(err, "could not convert start key to string")
		require.Equal(startString, cursor.StartIndex, "expected cursor start index to match")
		require.Empty(cursor.EndIndex, "expected cursor end index to be empty")

		// Configure trtl to return the rest of the tenants
		tenants = tenants[2:]
		rep, cursor, err = db.ListTenants(ctx, orgID, cursor)
		require.NoError(err, "could not list tenants")
		require.Len(rep, 1, "expected 1 tenant on the second page")
		require.Nil(cursor, "expected cursor to be nil because there is no next page")
	})
}

func (s *dbTestSuite) TestRetrieveTenant() {
	require := s.Require()
	ctx := context.Background()
	tenant := &db.Tenant{
		OrgID:           ulid.MustParse("02DEF3NDEKTSV4RRFFQ69G5FAZ"),
		ID:              ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:            "example-staging",
		EnvironmentType: "prod",
		Created:         time.Unix(1668660681, 0).In(time.UTC),
		Modified:        time.Unix(1668661302, 0).In(time.UTC),
	}

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.TenantNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}

		if !bytes.Equal(in.Key[0:16], tenant.OrgID[:]) {
			return nil, status.Error(codes.NotFound, "tenant organization not found")
		}

		if !bytes.Equal(in.Key[16:], tenant.ID[:]) {
			return nil, status.Error(codes.NotFound, "tenant not found")
		}

		// Marshal the data with msgpack
		data, err := tenant.MarshalValue()
		require.NoError(err, "could not marshal the tenant")

		// Unmarshal the data with msgpack
		other := &db.Tenant{}
		err = other.UnmarshalValue(data)
		require.NoError(err, "could not unmarshal the tenant")

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	tenant, err := db.RetrieveTenant(ctx, tenant.OrgID, tenant.ID)
	require.NoError(err, "could not retrieve tenant")

	// Fields should have been populated
	require.Equal(ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"), tenant.ID)
	require.Equal("example-staging", tenant.Name)
	require.Equal("prod", tenant.EnvironmentType)
	require.Equal(time.Unix(1668660681, 0), tenant.Created, "expected created timestamp to not have changed")
	require.True(time.Unix(1668661301, 0).Before(tenant.Modified), "expected modified timestamp to be updated")

	// Test NotFound path
	_, err = db.RetrieveTenant(ctx, tenant.OrgID, ulids.New())
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestUpdateTenant() {
	require := s.Require()
	ctx := context.Background()
	tenant := &db.Tenant{
		OrgID:           ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ID:              ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:            "tenant001",
		EnvironmentType: "dev",
		Created:         time.Unix(1668574281, 0),
		Modified:        time.Unix(1668574281, 0),
	}

	err := tenant.Validate()
	require.NoError(err, "could not validate tenant data")

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.TenantNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		if !bytes.Equal(in.Key[0:16], tenant.OrgID[:]) {
			return nil, status.Error(codes.NotFound, "tenant organization not found")
		}

		if !bytes.Equal(in.Key[16:], tenant.ID[:]) {
			return nil, status.Error(codes.NotFound, "tenant not found")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err = db.UpdateTenant(ctx, tenant)
	require.NoError(err, "could not update tenant")

	// Fields should have been populated
	require.Equal(ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"), tenant.ID, "tenant ID should not have changed")
	require.NotEmpty(tenant.Name, "tenant name is required")
	require.Equal(time.Unix(1668574281, 0), tenant.Created, "expected created timestamp to not have changed")
	require.True(time.Unix(1668574281, 0).Before(tenant.Modified), "expected modified timestamp to be updated")

	// If created timestamp is missing then it should be updated
	tenant.Created = time.Time{}
	require.NoError(db.UpdateTenant(ctx, tenant), "could not update tenant")
	require.Equal(tenant.Modified, tenant.Created, "expected created timestamp to be updated")

	// Should fail if tenant ID is missing
	tenant.ID = ulid.ULID{}
	require.ErrorIs(db.UpdateTenant(ctx, tenant), db.ErrMissingID, "expected error when tenant ID is missing")

	// Should fail if tenant is invalid
	tenant.ID = ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	tenant.Name = ""
	require.ErrorIs(db.UpdateTenant(ctx, tenant), db.ErrMissingTenantName, "expected error when tenant is invalid")

	// Test NotFound path
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return nil, status.Error(codes.NotFound, "tenant not found")
	}
	err = db.UpdateTenant(ctx, &db.Tenant{OrgID: ulids.New(), ID: ulids.New(), Name: "tenant002", EnvironmentType: "dev"})
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestDeleteTenant() {
	require := s.Require()
	ctx := context.Background()
	tenant := &db.Tenant{
		OrgID:    ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ID:       ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:     "example-dev",
		Created:  time.Unix(1668574281, 0),
		Modified: time.Unix(1668574281, 0),
	}

	s.mock.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.TenantNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Delete request")
		}

		if !bytes.Equal(in.Key[0:16], tenant.OrgID[:]) {
			return nil, status.Error(codes.NotFound, "tenant not found")
		}

		if !bytes.Equal(in.Key[16:], tenant.ID[:]) {
			return nil, status.Error(codes.NotFound, "tenant not found")
		}

		return &pb.DeleteReply{
			Success: true,
		}, nil
	}

	err := db.DeleteTenant(ctx, tenant.OrgID, tenant.ID)
	require.NoError(err, "could not delete tenant")

	// Test NotFound path
	err = db.DeleteTenant(ctx, tenant.OrgID, ulids.New())
	require.ErrorIs(err, db.ErrNotFound)
}

func TenantsEqual(t *testing.T, expected, actual *db.Tenant, msgAndArgs ...interface{}) {
	require.Equal(t, expected.OrgID, actual.OrgID, msgAndArgs...)
	require.Equal(t, expected.ID, actual.ID, msgAndArgs...)
	require.Equal(t, expected.Name, actual.Name, msgAndArgs...)
	require.Equal(t, expected.EnvironmentType, actual.EnvironmentType, msgAndArgs...)
	require.True(t, expected.Created.Equal(actual.Created), msgAndArgs...)
	require.True(t, expected.Modified.Equal(actual.Modified), msgAndArgs...)
}

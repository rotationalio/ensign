package db_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
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

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.TenantNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
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

	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back some data and terminate
		for i, tenant := range tenants {
			data, err := tenant.MarshalValue()
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
	require.NoError(err, "could not get tenant values")
	require.Len(values, 3, "expected 3 values")

	rep, err := db.ListTenants(ctx, orgID)
	require.NoError(err, "could not list tenants")
	require.Len(rep, 3, "expected 3 tenants")

	// Test first tenant data has been populated.
	require.Equal(tenants[0].ID, rep[0].ID, "expected tenant id to match")
	require.Equal(tenants[0].Name, rep[0].Name, "expected tenant name to match")
	require.Equal(tenants[0].EnvironmentType, rep[0].EnvironmentType, "expected tenant environment type to match")

	// Test second tenant data has been populated.
	require.Equal(tenants[1].ID, rep[1].ID, "expected tenant id to match")
	require.Equal(tenants[1].Name, rep[1].Name, "expected tenant name to match")
	require.Equal(tenants[1].EnvironmentType, rep[1].EnvironmentType, "expected tenant environment type to match")

	// Test third tenant data has been populated.
	require.Equal(tenants[2].ID, rep[2].ID, "expected tenant id to match")
	require.Equal(tenants[2].Name, rep[2].Name, "expected tenant name to match")
	require.Equal(tenants[2].EnvironmentType, rep[2].EnvironmentType, "expected tenant environment type to match")
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

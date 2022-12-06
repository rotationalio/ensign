package db_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/stretchr/testify/require"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTenantModel(t *testing.T) {
	tenant := &db.Tenant{
		ID:              ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:            "example-dev",
		EnvironmentType: "prod",
		Created:         time.Unix(1668660681, 0).In(time.UTC),
		Modified:        time.Unix(1668661302, 0).In(time.UTC),
	}

	key, err := tenant.Key()
	require.NoError(t, err, "could not marshal the key")
	require.Equal(t, tenant.ID[:], key, "unexpected marshaling of the key")

	require.Equal(t, db.TenantNamespace, tenant.Namespace(), "unexpected tenant namespace")

	// Test marshal and unmarshal
	data, err := tenant.MarshalValue()
	require.NoError(t, err, "could not marshal the tenant")

	other := &db.Tenant{}
	err = other.UnmarshalValue(data)
	require.NoError(t, err, "could not unmarshal the tenant")

	TenantsEqual(t, tenant, other, "unmarshaled tenant does not match marshaled tenant")
}

func (s *dbTestSuite) TestCreateTenant() {
	require := s.Require()
	ctx := context.Background()
	tenant := &db.Tenant{Name: "example-dev"}

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.TenantNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err := db.CreateTenant(ctx, tenant)
	require.NoError(err, "could not create tenant")

	// Fields should have been populated
	require.NotEqual("", tenant.ID, "expected non-zero ulid to be populated")
	require.NotZero(tenant.Created, "expected tenant to have a created timestamp")
	require.Equal(tenant.Created, tenant.Modified, "expected the same created and modified timestamp")
}

func (s *dbTestSuite) TestRetrieveTenant() {
	// TODO: this test will change if how the key or data marshaling is changed.
	require := s.Require()
	ctx := context.Background()
	tenantID := ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV")

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.TenantNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}

		if !bytes.Equal(in.Key, tenantID[:]) {
			return nil, status.Error(codes.NotFound, "tenant not found")
		}

		// Load fixture from disk
		data, err := os.ReadFile("testdata/tenant.json")
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "could not read fixture: %s", err)
		}

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	tenant, err := db.RetrieveTenant(ctx, tenantID)
	require.NoError(err, "could not retrieve tenant")

	ts, err := time.Parse(time.RFC3339, "2022-11-16T16:58:07-06:00")
	require.NoError(err, "could not parse timestamp fixture")

	require.Equal(tenantID, tenant.ID)
	require.Equal("example-staging", tenant.Name)
	require.Equal(ts, tenant.Created)
	require.Equal(ts, tenant.Modified)

	// Test NotFound path
	_, err = db.RetrieveTenant(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestUpdateTenant() {
	require := s.Require()
	ctx := context.Background()
	tenant := &db.Tenant{
		ID:       ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:     "example-dev",
		Created:  time.Unix(1668574281, 0),
		Modified: time.Unix(1668574281, 0),
	}

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.TenantNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		if !bytes.Equal(in.Key, tenant.ID[:]) {
			return nil, status.Error(codes.NotFound, "tenant not found")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err := db.UpdateTenant(ctx, tenant)
	require.NoError(err, "could not update tenant")

	// Fields should have been populated
	require.Equal(ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"), tenant.ID, "tenant ID should not have changed")
	require.Equal(time.Unix(1668574281, 0), tenant.Created, "expected created timestamp to not have changed")
	require.True(time.Unix(1668574281, 0).Before(tenant.Modified), "expected modified timestamp to be updated")

	// Test NotFound path
	err = db.UpdateTenant(ctx, &db.Tenant{ID: ulid.Make()})
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestDeleteTenant() {
	require := s.Require()
	ctx := context.Background()
	tenantID := ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV")

	s.mock.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.TenantNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Delete request")
		}

		if !bytes.Equal(in.Key, tenantID[:]) {
			return nil, status.Error(codes.NotFound, "tenant not found")
		}

		return &pb.DeleteReply{
			Success: true,
		}, nil
	}

	err := db.DeleteTenant(ctx, tenantID)
	require.NoError(err, "could not delete tenant")

	// Test NotFound path
	err = db.DeleteTenant(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}

func TenantsEqual(t *testing.T, expected, actual *db.Tenant, msgAndArgs ...interface{}) {
	require.Equal(t, expected.ID, actual.ID, msgAndArgs...)
	require.Equal(t, expected.Name, actual.Name, msgAndArgs...)
	require.Equal(t, expected.EnvironmentType, actual.EnvironmentType, msgAndArgs...)
	require.True(t, expected.Created.Equal(actual.Created), msgAndArgs...)
	require.True(t, expected.Modified.Equal(actual.Modified), msgAndArgs...)
}

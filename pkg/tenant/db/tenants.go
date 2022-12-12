package db

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/vmihailenco/msgpack/v5"
)

const TenantNamespace = "tenants"

type Tenant struct {
	ID              ulid.ULID `msgpack:"id"`
	Name            string    `msgpack:"name"`
	EnvironmentType string    `msgpack:"environment_type"`
	Created         time.Time `msgpack:"created"`
	Modified        time.Time `msgpack:"modified"`
}

// Compiler time check to ensure that tenant implements the Model interface
var _ Model = &Tenant{}

func (t *Tenant) Key() ([]byte, error) {
	// TODO: do we need any other key components for listing tenants?
	return t.ID.MarshalBinary()
}

func (t *Tenant) Namespace() string {
	return TenantNamespace
}

func (t *Tenant) MarshalValue() ([]byte, error) {
	return msgpack.Marshal(t)
}

func (t *Tenant) UnmarshalValue(data []byte) error {
	return msgpack.Unmarshal(data, t)
}

// An ID passed in by the User will be used. If an ID is not passed in,
// a new ID will be created.
func CreateTenant(ctx context.Context, tenant *Tenant) (err error) {
	if tenant.ID.Compare(ulid.ULID{}) == 0 {
		// TODO: use crypto rand and monotonic entropy with ulid.New
		tenant.ID = ulid.Make()
	}

	tenant.Created = time.Now()
	tenant.Modified = tenant.Created

	if err = Put(ctx, tenant); err != nil {
		return err
	}
	return nil
}

// ListTenants retrieves all tenants assigned to an organization.
func ListTenants(ctx context.Context, prefix []byte, namespace string) (values [][]byte, err error) {

	if values, err = List(ctx, prefix, namespace); err != nil {
		return nil, err
	}
	return values, err
}

func RetrieveTenant(ctx context.Context, id ulid.ULID) (tenant *Tenant, err error) {
	// Enough information must be stored on tenant to compute the key before Get
	tenant = &Tenant{
		ID: id,
	}

	// Get will populate the rest of the tenant struct from the database
	if err = Get(ctx, tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func UpdateTenant(ctx context.Context, tenant *Tenant) (err error) {
	if tenant.ID.Compare(ulid.ULID{}) == 0 {
		return ErrMissingID
	}

	tenant.Modified = time.Now()

	if err = Put(ctx, tenant); err != nil {
		return err
	}
	return nil
}

func DeleteTenant(ctx context.Context, id ulid.ULID) (err error) {
	tenant := &Tenant{
		ID: id,
	}

	if err = Delete(ctx, tenant); err != nil {
		return err
	}
	return nil
}

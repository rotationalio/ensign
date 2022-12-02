package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
)

const TenantNamespace = "tenants"

type Tenant struct {
	ID              ulid.ULID
	Name            string
	EnvironmentType string
	Created         time.Time
	Modified        time.Time
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
	// TODO: look into bson, msgpack, etc.
	return json.Marshal(t)
}

func (t *Tenant) UnmarshalValue(data []byte) error {
	// TODO: look into bson, msgpack, etc.
	return json.Unmarshal(data, t)
}

func CreateTenant(ctx context.Context, tenant *Tenant) (err error) {
	if tenant.ID.String() == "" {
		tenant.ID = ulid.Make()
	}

	tenant.Created = time.Now()
	tenant.Modified = tenant.Created

	if err = Put(ctx, tenant); err != nil {
		return err
	}
	return nil
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
	if tenant.ID.String() == "" {
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

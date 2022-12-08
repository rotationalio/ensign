package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const TenantNamespace = "tenants"

type Tenant struct {
	ID       uuid.UUID
	Name     string
	Created  time.Time
	Modified time.Time
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
	if tenant.ID == uuid.Nil {
		tenant.ID = uuid.New()
	}

	tenant.Created = time.Now()
	tenant.Modified = tenant.Created

	if err = Put(ctx, tenant); err != nil {
		return err
	}
	return nil
}

func ListTenants(ctx context.Context, prefix []byte, namespace string) (values [][]byte, err error) {
	if _, err := List(ctx, prefix, TenantNamespace); err != nil {
		return nil, err
	}
	return values, err
}

func RetrieveTenant(ctx context.Context, id uuid.UUID) (tenant *Tenant, err error) {
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
	if tenant.ID == uuid.Nil {
		return ErrMissingID
	}

	tenant.Modified = time.Now()

	if err = Put(ctx, tenant); err != nil {
		return err
	}
	return nil
}

func DeleteTenant(ctx context.Context, id uuid.UUID) (err error) {
	tenant := &Tenant{
		ID: id,
	}

	if err = Delete(ctx, tenant); err != nil {
		return err
	}
	return nil
}

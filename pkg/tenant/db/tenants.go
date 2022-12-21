package db

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/vmihailenco/msgpack/v5"
)

const TenantNamespace = "tenants"

type Tenant struct {
	OrgID           ulid.ULID `msgpack:"org_id"`
	ID              ulid.ULID `msgpack:"id"`
	Name            string    `msgpack:"name"`
	EnvironmentType string    `msgpack:"environment_type"`
	Created         time.Time `msgpack:"created"`
	Modified        time.Time `msgpack:"modified"`
}

// Compiler time check to ensure that tenant implements the Model interface
var _ Model = &Tenant{}

// Key is a 32 byte composite key combining the org id and tenant id.
func (t *Tenant) Key() (key []byte, err error) {
	// Create a 32 byte array so that the first 16 bytes hold
	// the org id and the last 16 bytes hold the tenant id.
	key = make([]byte, 32)

	// Marshal the org id to the first 16 bytes of the key.
	if err = t.OrgID.MarshalBinaryTo(key[0:16]); err != nil {
		return nil, err
	}

	// Marshal the tenant id to the second 16 bytes of the key.
	if err = t.ID.MarshalBinaryTo(key[16:]); err != nil {
		return nil, err
	}

	return key, err
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

// / CreateTenant adds a new project to the database.
// Note: If a tenant id is not passed in by the User, a new tenant id will be generated.
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
func ListTenants(ctx context.Context, orgID ulid.ULID) (tenants []*Tenant, err error) {
	// TODO: ensure that the tenants are stored with the orgID as their prefix!
	var prefix []byte
	if orgID.Compare(ulid.ULID{}) != 0 {
		prefix = orgID[:]
	}

	// TODO: it would be better to use the cursor directly rather than have duplicate data in memory
	var values [][]byte
	if values, err = List(ctx, prefix, TenantNamespace); err != nil {
		return nil, err
	}

	// Parse the tenants from the data
	tenants = make([]*Tenant, 0, len(values))
	for _, data := range values {
		tenant := &Tenant{}
		if err = tenant.UnmarshalValue(data); err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
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

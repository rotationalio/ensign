package db

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
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
	// OrgID and TenantID are required
	if ulids.IsZero(t.OrgID) {
		return nil, ErrMissingOrgID
	}

	if ulids.IsZero(t.ID) {
		return nil, ErrMissingID
	}

	var k Key
	if k, err = CreateKey(t.OrgID, t.ID); err != nil {
		return nil, err
	}

	return k.MarshalValue()
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

func (t *Tenant) Validate() error {
	if ulids.IsZero(t.OrgID) {
		return ErrMissingOrgID
	}

	if t.Name == "" {
		return ErrMissingTenantName
	}

	if t.EnvironmentType == "" {
		return ErrMissingEnvType
	}

	if !alphaNum.MatchString(t.Name) {
		return ErrInvalidTenantName
	}

	return nil
}

// Convert the model to an API response.
func (t *Tenant) ToAPI() *api.Tenant {
	return &api.Tenant{
		ID:              t.ID.String(),
		Name:            t.Name,
		EnvironmentType: t.EnvironmentType,
		Created:         TimeToString(t.Created),
		Modified:        TimeToString(t.Modified),
	}
}

// CreateTenant adds a new project to the database.
// Note: If a tenant id is not passed in by the User, a new tenant id will be generated.
func CreateTenant(ctx context.Context, tenant *Tenant) (err error) {
	if ulids.IsZero(tenant.ID) {
		tenant.ID = ulids.New()
	}

	if err = tenant.Validate(); err != nil {
		return err
	}

	tenant.Created = time.Now()
	tenant.Modified = tenant.Created

	if err = Put(ctx, tenant); err != nil {
		return err
	}
	return nil
}

// ListTenants retrieves all tenants assigned to an organization.
func ListTenants(ctx context.Context, orgID ulid.ULID, c *pg.Cursor) (tenants []*Tenant, cursor *pg.Cursor, err error) {
	// TODO: ensure that the tenants are stored with the orgID as their prefix!
	var prefix []byte
	if orgID.Compare(ulid.ULID{}) != 0 {
		prefix = orgID[:]
	}

	// Check to see if a default cursor exists and create one if it does not.
	if c == nil {
		c = pg.New("", "", 0)
	}

	if c.PageSize <= 0 {
		return nil, nil, ErrMissingPageSize
	}

	// TODO: it would be better to use the cursor directly rather than have duplicate data in memory
	var values [][]byte
	if values, cursor, err = List(ctx, prefix, TenantNamespace, c); err != nil {
		return nil, nil, err
	}

	// Parse the members from the data
	tenants = make([]*Tenant, 0, len(values))
	for _, data := range values {
		tenant := &Tenant{}
		if err = tenant.UnmarshalValue(data); err != nil {
			return nil, nil, err
		}
		tenants = append(tenants, tenant)
	}

	if len(values) > 0 {
		cursor = pg.New(string(values[0]), string(values[len(values)-1]), c.PageSize)
	}

	return tenants, cursor, nil
}

// Retrieve a tenant from the orgID and tenantID.
func RetrieveTenant(ctx context.Context, orgID, tenantID ulid.ULID) (tenant *Tenant, err error) {
	// Enough information must be stored on tenant to compute the key before Get
	tenant = &Tenant{
		OrgID: orgID,
		ID:    tenantID,
	}

	// Get will populate the rest of the tenant struct from the database
	if err = Get(ctx, tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func UpdateTenant(ctx context.Context, tenant *Tenant) (err error) {
	if ulids.IsZero(tenant.ID) {
		return ErrMissingID
	}

	// Validate tenant data.
	if err = tenant.Validate(); err != nil {
		return err
	}

	tenant.Modified = time.Now()
	if tenant.Created.IsZero() {
		tenant.Created = tenant.Modified
	}

	if err = Put(ctx, tenant); err != nil {
		return err
	}
	return nil
}

// Delete a tenant from the orgID and tenantID.
func DeleteTenant(ctx context.Context, orgID, tenantID ulid.ULID) (err error) {
	tenant := &Tenant{
		ID:    tenantID,
		OrgID: orgID,
	}

	if err = Delete(ctx, tenant); err != nil {
		return err
	}
	return nil
}

package db

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/vmihailenco/msgpack/v5"
)

const MembersNamespace = "members"

type Member struct {
	OrgID    ulid.ULID `msgpack:"org_id"`
	TenantID ulid.ULID `msgpack:"tenant_id"`
	ID       ulid.ULID `msgpack:"id"`
	Name     string    `msgpack:"name"`
	Role     string    `msgpack:"role"`
	Created  time.Time `msgpack:"created"`
	Modified time.Time `msgpack:"modified"`
}

var _ Model = &Member{}

// Key is a 32 byte composite key combining the tennat id and member id.
func (m *Member) Key() (key []byte, err error) {
	// Create a 32 byte array so that the first 16 bytes hold
	// the tenant id and the last 16 bytes hold the member id.
	key = make([]byte, 32)

	// Marshal the tenant id to the first 16 bytes of the key.
	if err = m.TenantID.MarshalBinaryTo(key[0:16]); err != nil {
		return nil, err
	}

	// Marshal the member id to the last 16 bytes of the key.
	if err = m.ID.MarshalBinaryTo(key[16:]); err != nil {
		return nil, err
	}
	return key, err
}

func (m *Member) Namespace() string {
	return MembersNamespace
}

func (m *Member) MarshalValue() ([]byte, error) {
	return msgpack.Marshal(m)
}

func (m *Member) UnmarshalValue(data []byte) error {
	return msgpack.Unmarshal(data, m)
}

func (m *Member) Validate() error {
	if ulids.IsZero(m.OrgID) {
		return ErrMissingOrgID
	}

	if m.Name == "" {
		return ErrMissingMemberName
	}

	if !alphaNum.MatchString(m.Name) || !alphaNum.MatchString(m.Role) {
		return ValidatonError("member")
	}

	return nil
}

// CreateTenantMember adds a new Member to a tenant in the database.
// Note: If a memberID is not passed in by the User, a new member id will be generated.
func CreateTenantMember(ctx context.Context, member *Member) (err error) {
	if ulids.IsZero(member.ID) {
		member.ID = ulids.New()
	}

	// Verify tenant ID is not missing.
	// Note: Tenant ID is not included in the authentication claims, so the ULID in the URL is currently
	// parsed on the server-side for this check. CreateTenantMember and CreateMember will be combined once
	// if it is determined that the tenant ID will be added to the claims or if it is determined that only one
	// create handler is needed on the server side.
	if ulids.IsZero(member.TenantID) {
		return ErrMissingTenantID
	}

	// Validate tenant member data.
	if err = member.Validate(); err != nil {
		return err
	}

	member.Created = time.Now()
	member.Modified = member.Created

	if err = Put(ctx, member); err != nil {
		return err
	}
	return nil
}

// CreateMember adds a new Member to an organization in the database.
// Note: If a memberID is not passed in by the User, a new member id will be generated.
func CreateMember(ctx context.Context, member *Member) (err error) {
	if ulids.IsZero(member.ID) {
		member.ID = ulids.New()
	}

	// Validate tenant member data.
	if err = member.Validate(); err != nil {
		return err
	}

	member.Created = time.Now()
	member.Modified = member.Created

	if err = Put(ctx, member); err != nil {
		return err
	}
	return nil
}

// RetrieveMember gets a member from the database with a given id.
func RetrieveMember(ctx context.Context, id ulid.ULID) (member *Member, err error) {
	member = &Member{
		ID: id,
	}

	if err = Get(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

// ListMembers retrieves all members assigned to a tenant.
func ListMembers(ctx context.Context, tenantID ulid.ULID) (members []*Member, err error) {
	// Store the tenant ID as the prefix.
	var prefix []byte
	if tenantID.Compare(ulid.ULID{}) != 0 {
		prefix = tenantID[:]
	}

	// TODO: Use the cursor directly instead of having duplicate data in memory
	var values [][]byte
	if values, err = List(ctx, prefix, MembersNamespace); err != nil {
		return nil, err
	}

	// Parse the members from the data
	members = make([]*Member, 0, len(values))
	for _, data := range values {
		member := &Member{}
		if err = member.UnmarshalValue(data); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

// UpdateMember updates the record of a member by its id.
func UpdateMember(ctx context.Context, member *Member) (err error) {
	// Check if memberID exists and return a missing
	// id error response if it does not.
	if ulids.IsZero(member.ID) {
		return ErrMissingID
	}

	// Validate member data.
	if err = member.Validate(); err != nil {
		return err
	}

	member.Modified = time.Now()

	if err = Put(ctx, member); err != nil {
		return err
	}

	return nil
}

// DeleteMember deletes a member with a given id.
func DeleteMember(ctx context.Context, id ulid.ULID) (err error) {
	member := &Member{
		ID: id,
	}

	if err = Delete(ctx, member); err != nil {
		return err
	}
	return nil
}

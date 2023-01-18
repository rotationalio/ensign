package db

import (
	"context"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/vmihailenco/msgpack/v5"
)

const MembersNamespace = "members"

type Member struct {
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
	if m.TenantID.Compare(ulid.ULID{}) == 0 {
		return ErrMissingTenantID
	}

	memberName := m.Name

	if memberName == "" {
		return ErrMissingMemberName
	}

	if strings.ContainsAny(string(memberName[0]), "0123456789") {
		return ErrNumberFirstCharacter
	}

	if strings.ContainsAny(memberName, " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		return ErrSpecialCharacters
	}

	return nil
}

// CreateTenantMember adds a new Member to a tenant in the database.
// Note: If a memberID is not passed in by the User, a new member id will be generated.
func CreateTenantMember(ctx context.Context, member *Member) (err error) {
	// TODO: Use crypto rand and monotonic entropy with ulid.New

	if member.ID.Compare(ulid.ULID{}) == 0 {
		member.ID = ulid.Make()
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
	// TODO: Use crypto rand and monotonic entropy with ulid.New

	if member.ID.Compare(ulid.ULID{}) == 0 {
		member.ID = ulid.Make()
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

		// Marshal and unmarshal the data with msgPack.
		member.MarshalData()
		member.UnmarshalData(data)
		members = append(members, member)
	}

	return members, nil
}

// UpdateMember updates the record of a member by its id.
func UpdateMember(ctx context.Context, member *Member) (err error) {

	// TODO: Use crypto rand and monotonic entropy with ulid.New

	// Check if memberID exists and return a missing
	// id error response if it does not.
	if member.ID.Compare(ulid.ULID{}) == 0 {
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

// Marshals data with msgPack.
func (p *Member) MarshalData() (data []byte, err error) {
	if data, err = p.MarshalValue(); err != nil {
		return nil, err
	}

	return data, nil
}

// Unmarshals data with msgPack.
func (p *Member) UnmarshalData(data []byte) (err error) {

	if err := p.UnmarshalValue(data); err != nil {
		return err
	}
	return nil
}

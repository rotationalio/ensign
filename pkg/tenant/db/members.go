package db

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/vmihailenco/msgpack/v5"
)

const MembersNamespace = "members"

type Member struct {
	OrgID    ulid.ULID `msgpack:"org_id"`
	ID       ulid.ULID `msgpack:"id"`
	Name     string    `msgpack:"name"`
	Role     string    `msgpack:"role"`
	Created  time.Time `msgpack:"created"`
	Modified time.Time `msgpack:"modified"`
}

var _ Model = &Member{}

// Key is a 32 byte composite key combining the org id and member id.
func (m *Member) Key() (key []byte, err error) {
	// Key requires an orgID and member ID
	if ulids.IsZero(m.ID) {
		return nil, ErrMissingID
	}

	if ulids.IsZero(m.OrgID) {
		return nil, ErrMissingOrgID
	}

	var k Key
	if k, err = CreateKey(m.OrgID, m.ID); err != nil {
		return nil, err
	}

	return k.MarshalValue()
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

// Validate checks if the member is ready for storage.
func (m *Member) Validate() error {
	if ulids.IsZero(m.OrgID) {
		return ErrMissingOrgID
	}

	if m.Name == "" {
		return ErrMissingMemberName
	}

	if !alphaNum.MatchString(m.Name) {
		return ErrInvalidMemberName
	}

	if m.Role == "" {
		return ErrMissingMemberRole
	}

	if !perms.IsRole(m.Role) {
		return ErrUnknownMemberRole
	}

	return nil
}

// Convert the model to an API response
func (m *Member) ToAPI() *api.Member {
	return &api.Member{
		ID:       m.ID.String(),
		Name:     m.Name,
		Role:     m.Role,
		Created:  TimeToString(m.Created),
		Modified: TimeToString(m.Modified),
	}
}

// CreateMember adds a new Member to an organization in the database.
// Note: If a memberID is not passed in by the User, a new member id will be generated.
func CreateMember(ctx context.Context, member *Member) (err error) {
	if ulids.IsZero(member.ID) {
		member.ID = ulids.New()
	}

	// Validate member data.
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

// RetrieveMember gets a member from the database with the given orgID and member ID.
func RetrieveMember(ctx context.Context, orgID, memberID ulid.ULID) (member *Member, err error) {
	member = &Member{
		ID:    memberID,
		OrgID: orgID,
	}

	if err = Get(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

// ListMembers retrieves all members assigned to an organization.
func ListMembers(ctx context.Context, orgID ulid.ULID, c *pg.Cursor) (members []*Member, cursor *pg.Cursor, err error) {
	// Store the org ID as the prefix.
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

	// TODO: Use the cursor directly instead of having duplicate data in memory
	var values [][]byte
	if values, cursor, err = List(ctx, prefix, MembersNamespace, c); err != nil {
		return nil, nil, err
	}

	// Parse the members from the data
	members = make([]*Member, 0, len(values))
	for _, data := range values {
		member := &Member{}
		if err = member.UnmarshalValue(data); err != nil {
			return nil, nil, err
		}
		members = append(members, member)
	}

	if len(values) > 0 {
		cursor = pg.New(string(values[0]), string(values[len(values)-1]), c.PageSize)
	}

	return members, cursor, nil
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
	if member.Created.IsZero() {
		member.Created = member.Modified
	}

	if err = Put(ctx, member); err != nil {
		return err
	}

	return nil
}

// DeleteMember deletes a member with a given orgID and member ID.
func DeleteMember(ctx context.Context, orgID, memberID ulid.ULID) (err error) {
	member := &Member{
		OrgID: orgID,
		ID:    memberID,
	}

	if err = Delete(ctx, member); err != nil {
		return err
	}
	return nil
}

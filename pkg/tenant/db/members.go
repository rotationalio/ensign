package db

import (
	"context"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/vmihailenco/msgpack/v5"
)

const MembersNamespace = "members"

type Member struct {
	OrgID        ulid.ULID    `msgpack:"org_id"`
	ID           ulid.ULID    `msgpack:"id"`
	Email        string       `msgpack:"email"`
	Name         string       `msgpack:"name"`
	Role         string       `msgpack:"role"`
	status       MemberStatus `msgpack:"status"`
	Created      time.Time    `msgpack:"created"`
	Modified     time.Time    `msgpack:"modified"`
	DateAdded    time.Time    `msgpack:"date_added"`
	LastActivity time.Time    `msgpack:"last_activity"`
}

type MemberStatus string

var _ Model = &Member{}

const (
	MemberConfirmed = "Confirmed"
	MemberPending   = "Pending"
)

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

	if m.Email == "" {
		return ErrMissingMemberEmail
	}

	if strings.TrimSpace(m.Name) == "" {
		return ErrMissingMemberName
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
		ID:           m.ID.String(),
		Email:        m.Email,
		Name:         m.Name,
		Role:         m.Role,
		Status:       string(m.status),
		Created:      TimeToString(m.Created),
		Modified:     TimeToString(m.Modified),
		DateAdded:    TimeToString(m.DateAdded),
		LastActivity: TimeToString(m.LastActivity),
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
	member.DateAdded = member.Created
	member.LastActivity = member.Created

	member.status = member.Status()

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

// ListMembers retrieves a paginated list of members.
func ListMembers(ctx context.Context, orgID ulid.ULID, c *pg.Cursor) (members []*Member, cursor *pg.Cursor, err error) {
	// Store the org ID as the prefix.
	var prefix []byte
	if orgID.Compare(ulid.ULID{}) != 0 {
		prefix = orgID[:]
	}

	var seekKey []byte
	if c.EndIndex != "" {
		var start ulid.ULID
		if start, err = ulid.Parse(c.EndIndex); err != nil {
			return nil, nil, err
		}
		seekKey = start[:]
	}

	// Check to see if a default cursor exists and create one if it does not.
	if c == nil {
		c = pg.New("", "", 0)
	}

	if c.PageSize <= 0 {
		return nil, nil, ErrMissingPageSize
	}

	members = make([]*Member, 0)
	onListItem := func(item *trtl.KVPair) error {
		member := &Member{}
		if err = member.UnmarshalValue(item.Value); err != nil {
			return err
		}
		members = append(members, member)
		return nil
	}

	if cursor, err = List(ctx, prefix, seekKey, MembersNamespace, onListItem, c); err != nil {
		return nil, nil, err
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

	member.LastActivity = member.Modified

	return Put(ctx, member)
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

func (m *Member) Status() MemberStatus {
	if m.status == "" {
		switch {
		case m.LastActivity.IsZero():
			m.status = MemberPending
		default:
			m.status = MemberConfirmed
		}
	}

	return m.status
}

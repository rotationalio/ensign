package db

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/vmihailenco/msgpack/v5"
)

const MembersNamespace = "members"

type Member struct {
	ID       ulid.ULID `msgpack:"id"`
	Name     string    `msgpack:"name"`
	Role     string    `msgpack:"role"`
	Created  time.Time `msgpack:"created"`
	Modified time.Time `msgpack:"modified"`
}

var _ Model = &Member{}

func (m *Member) Key() ([]byte, error) {
	return m.ID.MarshalBinary()
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

// CreateMembers adds a new Member to the database.
// Note: If a memberID is not passed in by the User, a
// new id will be generated.
func CreateMember(ctx context.Context, member *Member) (err error) {
	// TODO: Use crypto rand and monotonic entropy with ulid.New

	// Check if a memberID exists and create a new
	// one if it does not.
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

// RetrieveMember gets a member from the data with a given id.
func RetrieveMember(ctx context.Context, id ulid.ULID) (member *Member, err error) {
	member = &Member{
		ID: id,
	}

	if err = Get(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

// UpdateMember updates the record of a member by their id.
func UpdateMember(ctx context.Context, member *Member) (err error) {
	// TODO: Use crypto rand and monotonic entropy with ulid.New

	// Check if memberID exists and return a missing
	// id error response if it does not.
	if member.ID.Compare(ulid.ULID{}) == 0 {
		return ErrMissingID
	}

	if err = Put(ctx, member); err != nil {
		return err
	}

	return nil
}

// DeleteMemeber deletes a member with a given id.
func DeleteMember(ctx context.Context, id ulid.ULID) (err error) {
	member := &Member{
		ID: id,
	}

	if err = Delete(ctx, member); err != nil {
		return err
	}
	return nil
}

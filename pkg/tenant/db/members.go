package db

import (
	"context"
	"errors"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/utils/gravatar"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	MembersNamespace       = "members"
	MembersDefaultPageSize = 100
	MaxNameLength          = 1024
)

var (
	// Must be at least 3 characters, cannot start with a number, and is alphanumeric with + _ and -
	WorkspaceNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+_-]{2,}$`)
)

type Member struct {
	OrgID             ulid.ULID          `msgpack:"org_id"`
	ID                ulid.ULID          `msgpack:"id"`
	Email             string             `msgpack:"email"`
	Name              string             `msgpack:"name"`
	Organization      string             `msgpack:"organization"`
	Workspace         string             `msgpack:"workspace"`
	ProfessionSegment ProfessionSegment  `msgpack:"profession_segment"`
	DeveloperSegment  []DeveloperSegment `msgpack:"developer_segment"`
	Role              string             `msgpack:"role"`
	Invited           bool               `msgpack:"invited"`
	JoinedAt          time.Time          `msgpack:"joined_at"`
	LastActivity      time.Time          `msgpack:"last_activity"`
	Created           time.Time          `msgpack:"created"`
	Modified          time.Time          `msgpack:"modified"`
	gravatar          string
}

type MemberStatus uint8

const (
	MemberStatusPending MemberStatus = iota
	MemberStatusOnboarding
	MemberStatusActive
)

var MemberStatusStrings = map[MemberStatus]string{
	MemberStatusPending:    "Pending",
	MemberStatusOnboarding: "Onboarding",
	MemberStatusActive:     "Active",
}

func (m MemberStatus) String() string {
	return MemberStatusStrings[m]
}

type ProfessionSegment uint8

const (
	ProfessionSegmentUnspecified ProfessionSegment = iota
	ProfessionSegmentWork
	ProfessionSegmentEducation
	ProfessionSegmentPersonal
)

var ProfessionSegmentStrings = map[ProfessionSegment]string{
	ProfessionSegmentUnspecified: "Unspecified",
	ProfessionSegmentWork:        "Work",
	ProfessionSegmentEducation:   "Education",
	ProfessionSegmentPersonal:    "Personal",
}

func (p ProfessionSegment) IsZero() bool {
	return p == ProfessionSegmentUnspecified
}

func (p ProfessionSegment) String() string {
	return ProfessionSegmentStrings[p]
}

// Parse a segment string into a ProfessionSegment, empty string is considered
// unspecified but not an error.
func ParseProfessionSegment(segment string) (ProfessionSegment, error) {
	segment = strings.ToLower(strings.TrimSpace(segment))
	switch segment {
	case "":
		return ProfessionSegmentUnspecified, nil
	case "unspecified":
		return ProfessionSegmentUnspecified, nil
	case "work":
		return ProfessionSegmentWork, nil
	case "education":
		return ProfessionSegmentEducation, nil
	case "personal":
		return ProfessionSegmentPersonal, nil
	default:
		return ProfessionSegmentUnspecified, ErrProfessionUnknown
	}
}

type DeveloperSegment uint8

const (
	DeveloperSegmentUnspecified DeveloperSegment = iota
	DeveloperSegmentSomethingElse
	DeveloperSegmentApplicationDevelopment
	DeveloperSegmentDataScience
	DeveloperSegmentDataEngineering
	DeveloperSegmentDeveloperExperience
	DeveloperSegmentCybersecurity
	DeveloperSegmentDevOps
)

var DeveloperSegmentStrings = map[DeveloperSegment]string{
	DeveloperSegmentUnspecified:            "Unspecified",
	DeveloperSegmentSomethingElse:          "Something else",
	DeveloperSegmentApplicationDevelopment: "Application development",
	DeveloperSegmentDataScience:            "Data science",
	DeveloperSegmentDataEngineering:        "Data engineering",
	DeveloperSegmentDeveloperExperience:    "Developer experience",
	DeveloperSegmentCybersecurity:          "Cybersecurity (blue or purple team)",
	DeveloperSegmentDevOps:                 "DevOps and observability",
}

func (d DeveloperSegment) String() string {
	return DeveloperSegmentStrings[d]
}

// Parse a segment string into a DeveloperSegment, empty string is considered
// unspecified but not an error.
func ParseDeveloperSegment(segment string) (DeveloperSegment, error) {
	segment = strings.ToLower(strings.TrimSpace(segment))
	switch segment {
	case "":
		return DeveloperSegmentUnspecified, nil
	case "unspecified":
		return DeveloperSegmentUnspecified, nil
	case "something else":
		return DeveloperSegmentSomethingElse, nil
	case "application development":
		return DeveloperSegmentApplicationDevelopment, nil
	case "data science":
		return DeveloperSegmentDataScience, nil
	case "data engineering":
		return DeveloperSegmentDataEngineering, nil
	case "developer experience":
		return DeveloperSegmentDeveloperExperience, nil
	case "cybersecurity (blue or purple team)":
		return DeveloperSegmentCybersecurity, nil
	case "devops and observability":
		return DeveloperSegmentDevOps, nil
	default:
		return DeveloperSegmentUnspecified, ErrDeveloperUnknown
	}
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

// Validate checks if the member is ready for storage, which can potentially return
// multiple errors in the form of a ValidationErrors.
func (m *Member) Validate() error {
	if ulids.IsZero(m.OrgID) {
		return ErrMissingOrgID
	}

	if m.Email == "" {
		return ErrMissingMemberEmail
	}

	if m.Role == "" {
		return ErrMissingMemberRole
	}

	if !perms.IsRole(m.Role) {
		return ErrUnknownMemberRole
	}

	// Validate the onboarding fields and return a ValidationErrors if there are any
	// provided fields that are invalid.
	errs := make(ValidationErrors, 0)

	m.Name = strings.TrimSpace(m.Name)
	if m.Name != "" && len(m.Name) > MaxNameLength {
		errs = append(errs, validationError("name", ErrNameTooLong))
	}

	m.Organization = strings.TrimSpace(m.Organization)
	if m.Organization != "" && len(m.Organization) > MaxNameLength {
		errs = append(errs, validationError("organization", ErrOrganizationTooLong))
	}

	m.Workspace = strings.TrimSpace(m.Workspace)
	if m.Workspace != "" {
		if len(m.Workspace) > MaxNameLength {
			errs = append(errs, validationError("workspace", ErrWorkspaceTooLong))
		} else if !WorkspaceNameRegex.MatchString(m.Workspace) {
			errs = append(errs, validationError("workspace", ErrInvalidWorkspace))
		}
	}

	if len(m.DeveloperSegment) > 0 {
		for i, segment := range m.DeveloperSegment {
			if segment == DeveloperSegmentUnspecified {
				errs = append(errs, validationError("developer_segment", ErrDeveloperUnspecified).AtIndex(i))
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// OnboardingStatus returns the current status of the member which is a derived value
// based on the information in the member record.
func (m *Member) OnboardingStatus() MemberStatus {
	switch {
	case m.Invited && m.JoinedAt.IsZero():
		return MemberStatusPending
	case !m.IsOnboarded():
		return MemberStatusOnboarding
	default:
		return MemberStatusActive
	}
}

// IsOnboarded returns true if there is enough information to consider the member fully
// onboarded into the organization.
func (m *Member) IsOnboarded() bool {
	return m.Name != "" && m.Organization != "" && m.Workspace != "" && !m.ProfessionSegment.IsZero() && len(m.DeveloperSegment) > 0
}

func (m *Member) Picture() string {
	if m.gravatar == "" {
		m.gravatar = gravatar.New(m.Email, nil)
	}

	return m.gravatar
}

// Convert the model to an API response
func (m *Member) ToAPI() *api.Member {
	ret := &api.Member{
		ID:                m.ID.String(),
		Email:             m.Email,
		Name:              m.Name,
		Organization:      m.Organization,
		Workspace:         m.Workspace,
		ProfessionSegment: m.ProfessionSegment.String(),
		Picture:           m.Picture(),
		Role:              m.Role,
		Invited:           m.Invited,
		OnboardingStatus:  m.OnboardingStatus().String(),
		Created:           TimeToString(m.Created),
		DateAdded:         TimeToString(m.JoinedAt),
		LastActivity:      TimeToString(m.LastActivity),
	}

	for _, segment := range m.DeveloperSegment {
		ret.DeveloperSegment = append(ret.DeveloperSegment, segment.String())
	}

	return ret
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

// ListMembers retrieves a paginated list of members.
func ListMembers(ctx context.Context, orgID ulid.ULID, c *pg.Cursor) (members []*Member, cursor *pg.Cursor, err error) {
	// Store the org ID as the prefix.
	var prefix []byte
	if orgID.Compare(ulid.ULID{}) != 0 {
		prefix = orgID[:]
	}

	// Check to see if a default cursor exists and create one if it does not.
	if c == nil {
		c = pg.New("", "", MembersDefaultPageSize)
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

	if cursor, err = List(ctx, prefix, MembersNamespace, onListItem, c); err != nil {
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

// Helper method that returns an error if an email address is invalid or already exists
// in the organization.
func VerifyMemberEmail(ctx context.Context, orgID ulid.ULID, email string) (err error) {
	if email == "" {
		return ErrMissingMemberEmail
	}

	var members []*Member
	if members, _, err = ListMembers(ctx, orgID, nil); err != nil {
		return err
	}

	for _, member := range members {
		if member.Email == email {
			return ErrMemberExists
		}
	}

	return nil
}

// GetMemberByEmail returns a member by the exact email address without any lowercasing validation.
func GetMemberByEmail(ctx context.Context, orgID ulid.ULID, email string) (member *Member, err error) {
	if ulids.IsZero(orgID) {
		return nil, ErrMissingOrgID
	}

	if email == "" {
		return nil, ErrMissingMemberEmail
	}

	req := &trtl.CursorRequest{
		Prefix:    orgID[:],
		Namespace: MembersNamespace,
	}

	var stream trtl.Trtl_CursorClient
	if stream, err = client.Cursor(ctx, req); err != nil {
		return nil, err
	}

	onListItem := func(item *trtl.KVPair) error {
		member = &Member{}
		if err = member.UnmarshalValue(item.Value); err != nil {
			return err
		}
		if member.Email == email {
			return ErrListBreak
		}
		return nil
	}

	for {
		var item *trtl.KVPair
		if item, err = stream.Recv(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if err = onListItem(item); err != nil {
			if errors.Is(err, ErrListBreak) {
				return member, nil
			}
			return nil, err
		}
	}
	return nil, ErrMemberEmailNotFound
}

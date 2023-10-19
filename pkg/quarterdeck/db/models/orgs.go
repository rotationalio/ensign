package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/random"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

const (
	DefaultDomainLength = 8
)

// Organization is a model that represents a row in the organizations table and provides
// database functionality for interacting with an organizations's data. It should not be
// used for API serialization.
type Organization struct {
	Base
	ID        ulid.ULID
	Name      string
	Domain    string
	projects  int
	lastLogin sql.NullString
}

// OrganizationUser is a model representing a many-to-many mapping between users and
// organizations and describes the role each user has in their organization. This model
// is primarily used by the User and Organization models and is not intended for direct
// use generally.
//
// NOTE: a user can only have one role in an organization, so roles must be defined as
// overlapping sets rather than as disjoint sets where users have multiple roles.
type OrganizationUser struct {
	Base
	OrgID              ulid.ULID
	UserID             ulid.ULID
	RoleID             int64
	DeleteConfirmToken sql.NullString
	LastLogin          sql.NullString
	user               *User
	org                *Organization
	role               *Role
}

const (
	getOrgSQL       = "SELECT name, domain, created, modified FROM organizations WHERE id=:id"
	getOrgProjects  = "SELECT COUNT(*) FROM organization_projects WHERE organization_id=:orgID"
	getOrgLastLogin = "SELECT last_login FROM organization_users WHERE organization_id=:orgID AND user_id=:userID"
)

func GetOrg(ctx context.Context, id any) (org *Organization, err error) {
	org = &Organization{}
	if org.ID, err = ulids.Parse(id); err != nil {
		return nil, err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = org.populate(tx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return org, nil
}

const (
	lookupOrgWorkspace = "SELECT id, name, domain, created, modified FROM organizations WHERE domain=:domain"
)

func LookupWorkspace(ctx context.Context, domain string) (org *Organization, err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	org = &Organization{}
	if err = tx.QueryRow(lookupOrgWorkspace, sql.Named("domain", domain)).Scan(&org.ID, &org.Name, &org.Domain, &org.Created, &org.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err = tx.QueryRow(getOrgProjects, sql.Named("orgID", org.ID)).Scan(&org.projects); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return org, nil
}

const (
	getAllOrgsSQL = "SELECT id, name, domain, created, modified FROM organizations"
)

// ListAllOrgs returns a paginated collection of all organizations in the database.
// This does not filter by user so ListUserOrgs should be used when only the
// organizations for a specific user are needed.
func ListAllOrgs(ctx context.Context, prevPage *pagination.Cursor) (organizations []*Organization, cursor *pagination.Cursor, err error) {
	if prevPage == nil {
		// Create a default cursor, e.g. the previous page was nothing
		prevPage = pagination.New("", "", 0)
	}

	if prevPage.PageSize <= 0 {
		return nil, nil, invalid(ErrMissingPageSize)
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	// Build parameterized query with WHERE clause
	// ---------------------------------------------------------------------------------------------------
	// Query construction with pageSize only:
	// SELECT id, name, domain, created, modified FROM organizations LIMIT :pageSize
	// ---------------------------------------------------------------------------------------------------
	// Query construction with pageSize and endIndex:
	// SELECT id, name, domain, created, modified FROM organizations WHERE id > :endIndex LIMIT :pageSize
	var query strings.Builder
	query.WriteString(getAllOrgsSQL)

	// Keep track of the parameters to pass with the query
	params := make([]any, 0, 3)

	if prevPage.EndIndex != "" {
		var endIndex ulid.ULID
		if endIndex, err = ulid.Parse(prevPage.EndIndex); err != nil {
			return nil, nil, invalid(ErrInvalidCursor)
		}

		// The WHERE clause ensures that we start after the end of the previous page
		params = append(params, sql.Named("endIndex", endIndex))
		query.WriteString(" WHERE id > :endIndex")
	}

	// Add the limit as the page size + 1 to perform a has next page check.
	// pageSize controls the number of results returned from the query
	params = append(params, sql.Named("pageSize", prevPage.PageSize+1))
	query.WriteString(" LIMIT :pageSize")

	// Do the query
	var rows *sql.Rows
	if rows, err = tx.Query(query.String(), params...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	nRows := int32(0)
	organizations = make([]*Organization, 0, prevPage.PageSize)
	for rows.Next() {
		// The query will request one additional message past the page size to check if
		// there is a next page. We should not process any messages after the page size.
		nRows++
		if nRows > prevPage.PageSize {
			continue
		}

		// Create organization object to append to the organizations list
		org := &Organization{}

		// Populate organization details
		if err = rows.Scan(&org.ID, &org.Name, &org.Domain, &org.Created, &org.Modified); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil, nil
			}
			return nil, nil, err
		}

		// Retrieve the number of projects associated with the organization
		if err = tx.QueryRow(getOrgProjects, sql.Named("orgID", org.ID)).Scan(&org.projects); err != nil {
			return nil, nil, err
		}

		organizations = append(organizations, org)
	}

	if err = rows.Close(); err != nil {
		return nil, nil, err
	}

	// Create the cursor to return if there is a next page of results
	if len(organizations) > 0 && nRows > prevPage.PageSize {
		cursor = pagination.New(organizations[0].ID.String(), organizations[len(organizations)-1].ID.String(), prevPage.PageSize)
	}
	return organizations, cursor, nil
}

const (
	getOrgsForUserSQL = "SELECT id, name, domain, created, modified FROM organizations WHERE id IN (SELECT organization_id FROM organization_users"
)

// ListUserOrgs returns a paginated collection of organizations filtered by the userID.
// The orgID must be a valid non-zero value of type ulid.ULID,
// a string representation of a type ulid.ULID, or a slice of bytes
// The number of organizations resturned is controlled by the prevPage cursor.
// To return the first page with a default number of results pass nil for the prevPage;
// Otherwise pass an empty page with the specified PageSize.
// If the prevPage contains an EndIndex then the next page is returned.
//
// A organizations slice with the maximum length of the page size will be returned or an
// empty (nil) slice if there are no results. If there is a next page of results, e.g.
// there is another row after the page returned, then a cursor will be returned to
// compute the next page token with.
func ListUserOrgs(ctx context.Context, userID any, prevPage *pagination.Cursor) (organizations []*Organization, cursor *pagination.Cursor, err error) {
	var user ulid.ULID
	if user, err = ulids.Parse(userID); err != nil {
		return nil, nil, err
	}

	if ulids.IsZero(user) {
		return nil, nil, ErrMissingModelID
	}

	if prevPage == nil {
		// Create a default cursor, e.g. the previous page was nothing
		prevPage = pagination.New("", "", 0)
	}

	if prevPage.PageSize <= 0 {
		return nil, nil, invalid(ErrMissingPageSize)
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	// Build parameterized query with WHERE clause
	// ---------------------------------------------------------------------------------------------------
	// Query construction with pageSize only:
	// SELECT id, name, domain, created, modified FROM organizations
	// WHERE id IN (SELECT organization_id FROM organization_users WHERE user_id=:userID) LIMIT :pageSize
	// ---------------------------------------------------------------------------------------------------
	// Query construction with pageSize and endIndex:
	// SELECT id, name, domain, created, modified FROM organizations
	// WHERE id IN (SELECT organization_id FROM organization_users WHERE user_id=:userID) AND id > :endIndex LIMIT :pageSize
	var query strings.Builder
	query.WriteString(getOrgsForUserSQL)

	// Construct the where clause
	params := make([]any, 0, 4)
	where := make([]string, 0, 3)

	params = append(params, sql.Named("userID", user))
	where = append(where, "user_id=:userID)")

	if prevPage.EndIndex != "" {
		var endIndex ulid.ULID
		if endIndex, err = ulid.Parse(prevPage.EndIndex); err != nil {
			return nil, nil, invalid(ErrInvalidCursor)
		}

		// endIndex is the id of the last user in prevPage
		// add the endIndex parameter to ensure that the next set
		// of results are greater than that id
		params = append(params, sql.Named("endIndex", endIndex))
		where = append(where, "id > :endIndex")
	}

	// Add the where clause to the query
	query.WriteString(" WHERE ")
	query.WriteString(strings.Join(where, " AND "))

	// Add the limit as the page size + 1 to perform a has next page check.
	// pageSize controls the number of results returned from the query
	params = append(params, sql.Named("pageSize", prevPage.PageSize+1))
	query.WriteString(" LIMIT :pageSize")
	// Fetch list of users associated with the orgID
	var rows *sql.Rows
	if rows, err = tx.Query(query.String(), params...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	nRows := int32(0)
	organizations = make([]*Organization, 0, prevPage.PageSize)
	for rows.Next() {
		// The query will request one additional message past the page size to check if
		// there is a next page. We should not process any messages after the page size.
		nRows++
		if nRows > prevPage.PageSize {
			continue
		}

		//create organization object to append to the organizations list
		org := &Organization{}

		// populate organization details
		if err = rows.Scan(&org.ID, &org.Name, &org.Domain, &org.Created, &org.Modified); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil, nil
			}
			return nil, nil, err
		}

		// retrieve the number of projects associated with the organization
		if err = tx.QueryRow(getOrgProjects, sql.Named("orgID", org.ID)).Scan(&org.projects); err != nil {
			return nil, nil, err
		}

		// retrieve the time the user last logged in to this organization
		if err = tx.QueryRow(getOrgLastLogin, sql.Named("orgID", org.ID), sql.Named("userID", user)).Scan(&org.lastLogin); err != nil {
			return nil, nil, err
		}

		organizations = append(organizations, org)
	}

	if err = rows.Close(); err != nil {
		return nil, nil, err
	}

	// Create the cursor to return if there is a next page of results
	if len(organizations) > 0 && nRows > prevPage.PageSize {
		cursor = pagination.New(organizations[0].ID.String(), organizations[len(organizations)-1].ID.String(), prevPage.PageSize)
	}
	return organizations, cursor, nil
}

const (
	insertOrgSQL = "INSERT INTO organizations VALUES (:id, :name, :domain, :created, :modified)"
)

// Create an organization, inserting the record into the database. If the record already
// exists or a uniqueness constraint is violated an error is returned. This method sets
// the ID, created, and modified timestamps even if the user has already set them.
//
// If the organization name or domain are empty a validation error is returned.
func (o *Organization) Create(ctx context.Context) (err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	if err = o.create(tx); err != nil {
		return err
	}

	return tx.Commit()
}

const (
	updateOrgSQL = "UPDATE organizations SET name=:name, domain=:domain, modified=:modified WHERE id=:id"
)

// Save an organization to the database, updating the name and domain fields. This
// assumes that the organization already exists in the database.
func (o *Organization) Save(ctx context.Context) (err error) {
	if err = o.Validate(); err != nil {
		return err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	o.SetModified(time.Now())
	params := make([]any, 4)
	params[0] = sql.Named("id", o.ID)
	params[1] = sql.Named("name", o.Name)
	params[2] = sql.Named("domain", o.Domain)
	params[3] = sql.Named("modified", o.Modified)

	var res sql.Result
	if res, err = tx.Exec(updateOrgSQL, params...); err != nil {
		var dberr sqlite3.Error
		if errors.As(err, &dberr) {
			if dberr.Code == sqlite3.ErrConstraint {
				return constraint(dberr)
			}
		}

		return err
	}

	// Return not found if no rows were modified
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNotFound
	}

	return tx.Commit()
}

const (
	deleteOrgSQL = "DELETE FROM organizations WHERE id=:id"
)

// Delete an organization from the database. If there are any API keys in the
// organization then this will return a constraint error.
func (o *Organization) Delete(ctx context.Context) (err error) {
	if ulids.IsZero(o.ID) {
		return ErrMissingModelID
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the organization
	if _, err = tx.Exec(deleteOrgSQL, sql.Named("id", o.ID)); err != nil {
		var dberr sqlite3.Error
		if errors.As(err, &dberr) {
			if dberr.Code == sqlite3.ErrConstraint {
				return constraint(dberr)
			}
		}
		return err
	}

	return tx.Commit()
}

// Validate that an organization is ready to be inserted or updated into the database.
func (o *Organization) Validate() error {
	if ulids.IsZero(o.ID) {
		return invalid(ErrMissingModelID)
	}

	if o.Name == "" || o.Domain == "" {
		return invalid(ErrInvalidOrganization)
	}

	return nil
}

func (o *Organization) create(tx *sql.Tx) (err error) {
	o.ID = ulids.New()

	// Create a random domain if one is not provided
	if o.Domain == "" {
		o.Domain = random.Name(DefaultDomainLength)
	}

	now := time.Now()
	o.SetCreated(now)
	o.SetModified(now)

	params := make([]any, 5)
	params[0] = sql.Named("id", o.ID)
	params[1] = sql.Named("name", o.Name)
	params[2] = sql.Named("domain", o.Domain)
	params[3] = sql.Named("created", o.Created)
	params[4] = sql.Named("modified", o.Modified)

	if _, err = tx.Exec(insertOrgSQL, params...); err != nil {
		var dberr sqlite3.Error
		if errors.As(err, &dberr) {
			if dberr.Code == sqlite3.ErrConstraint {
				return constraint(dberr)
			}
		}
		return err
	}
	return nil
}

func (o *Organization) populate(tx *sql.Tx) (err error) {
	if ulids.IsZero(o.ID) {
		return ErrMissingModelID
	}

	if err = tx.QueryRow(getOrgSQL, sql.Named("id", o.ID)).Scan(&o.Name, &o.Domain, &o.Created, &o.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if err = tx.QueryRow(getOrgProjects, sql.Named("orgID", o.ID)).Scan(&o.projects); err != nil {
		return err
	}

	return nil
}

const (
	orgExistsSQL = "SELECT EXISTS(SELECT 1 FROM organizations WHERE id=:orgID)"
)

func (o *Organization) exists(tx *sql.Tx) (ok bool, err error) {
	if ulids.IsZero(o.ID) {
		return false, nil
	}

	if err = tx.QueryRow(orgExistsSQL, sql.Named("orgID", o.ID)).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (o *Organization) ToAPI() *api.Organization {
	org := &api.Organization{
		ID:       o.ID,
		Name:     o.Name,
		Domain:   o.Domain,
		Projects: o.ProjectCount(),
	}
	org.LastLogin, _ = o.LastLogin()
	org.Created, _ = o.GetCreated()
	org.Modified, _ = o.GetModified()
	return org
}

func (o *Organization) ProjectCount() int {
	return o.projects
}

// Return the last login timestamp for the user this organization was loaded for.
func (o *Organization) LastLogin() (time.Time, error) {
	if o.lastLogin.Valid {
		return time.Parse(time.RFC3339Nano, o.lastLogin.String)
	}
	return time.Time{}, nil
}

const (
	getOrgUserSQL = "SELECT role_id, delete_confirmation_token, last_login, created, modified FROM organization_users WHERE user_id=:userID AND organization_id=:orgID"
)

func GetOrgUser(ctx context.Context, userID, orgID any) (ou *OrganizationUser, err error) {
	ou = &OrganizationUser{}
	if ou.UserID, err = ulids.Parse(userID); err != nil {
		return nil, err
	}
	if ou.OrgID, err = ulids.Parse(orgID); err != nil {
		return nil, err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getOrgUserSQL, sql.Named("userID", ou.UserID), sql.Named("orgID", ou.OrgID)).Scan(&ou.RoleID, &ou.DeleteConfirmToken, &ou.LastLogin, &ou.Created, &ou.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return ou, nil
}

// Returns the user associated with the OrganizationUser struct, ready to query with
// the given organization. This object is cached on the struct and can be refreshed.
// TODO: fetch on GetOrgUser to reduce number of raft queries.
func (o *OrganizationUser) User(ctx context.Context, refresh bool) (_ *User, err error) {
	if refresh || o.user == nil {
		if o.user, err = GetUser(ctx, o.UserID, o.OrgID); err != nil {
			return nil, err
		}
	}
	return o.user, nil
}

// Returns the organization associated with the OrganizationUser struct. The object is
// cached on the struct and can be refreshed on demand.
// TODO: fetch on GetOrgUser to reduce number of raft queries.
func (o *OrganizationUser) Organization(ctx context.Context, refresh bool) (_ *Organization, err error) {
	if refresh || o.org == nil {
		if o.org, err = GetOrg(ctx, o.OrgID); err != nil {
			return nil, err
		}
	}
	return o.org, nil
}

// Returns the role associated with the organization and user. The object is cached on
// the struct and can be refreshed on demand.
// TODO: fetch on GetOrgUser to reduce number of raft queries.
func (o *OrganizationUser) Role(ctx context.Context, refresh bool) (_ *Role, err error) {
	if refresh || o.role == nil {
		if o.role, err = GetRole(ctx, o.RoleID); err != nil {
			return nil, err
		}
	}
	return o.role, nil
}

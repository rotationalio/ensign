package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

// User is a model that represents a row in the users table and provides database
// functionality for interacting with a user's data. It should not be used for API
// serialization. Users may be retrieved from the database either via their ID (e.g.
// from the sub claim in a JWT token) or via their email address (e.g. on login). The
// user password should be stored as an argon2 hash and should be verified using the
// argon2 hashing algorithm. Care should be taken to ensure this model stays secure.
//
// Users are associated with one or more organizations. When the user model is loaded
// from the database one organization must be supplied so that permissions and role
// can be retrieved correctly. If no orgID is supplied then one of the user's
// organizations is selected from the database as the default organiztion. Use the
// SwitchOrganization method to switch the user model to a different organization to
// retrieve a different and permissions or use the UserRoles method to determine which
// organizations the user belongs to.
type User struct {
	Base
	ID           ulid.ULID
	Name         string
	Email        string
	Password     string
	AgreeToS     sql.NullBool
	AgreePrivacy sql.NullBool
	LastLogin    sql.NullString
	orgID        ulid.ULID
	orgRoles     map[ulid.ULID]string
	permissions  []string
}

const (
	getUserIDSQL    = "SELECT name, email, password, terms_agreement, privacy_agreement, last_login, created, modified FROM users WHERE id=:id"
	getUserEmailSQL = "SELECT id, name, password, terms_agreement, privacy_agreement, last_login, created, modified FROM users WHERE email=:email"
)

//===========================================================================
// Retrieve Users from Database
//===========================================================================

// GetUser by ID. The ID can be either a string, which is parsed into a ULID or it can
// be a valid ULID. The query is then executed as a read-only transaction against the
// database and the user is returned. An orgID can be specified to load the user in that
// organization. If the orgID is Null then one of the organizations the user belongs to
// is loaded (the default user organization).
func GetUser(ctx context.Context, userID, orgID any) (u *User, err error) {
	// Create the user struct and parse the ID input.
	u = &User{}
	if u.ID, err = ulids.Parse(userID); err != nil {
		return nil, err
	}

	var userOrg ulid.ULID
	if userOrg, err = ulids.Parse(orgID); err != nil {
		return nil, err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getUserIDSQL, sql.Named("id", u.ID)).Scan(&u.Name, &u.Email, &u.Password, &u.AgreeToS, &u.AgreePrivacy, &u.LastLogin, &u.Created, &u.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Load user in the specified organization or default organization if null is
	// specified; this also verifies the user is part of the organization and caches
	// the organizations and roles the user belongs to as well as the permissions of
	// the current organization.
	if err = u.loadOrganization(tx, userOrg); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}

// GetUser by Email. This query is executed as a read-only transaction. An orgID can be
// specified to load the user in that organization. If the orgID is Null then one of the
// organizations the user belongs to is loaded (the default user organization).
func GetUserEmail(ctx context.Context, email string, orgID any) (u *User, err error) {
	u = &User{Email: email}

	var userOrg ulid.ULID
	if userOrg, err = ulids.Parse(orgID); err != nil {
		return nil, err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getUserEmailSQL, sql.Named("email", u.Email)).Scan(&u.ID, &u.Name, &u.Password, &u.AgreeToS, &u.AgreePrivacy, &u.LastLogin, &u.Created, &u.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Load user in the specified organization or default organization if null is
	// specified; this also verifies the user is part of the organization and caches
	// the organizations and roles the user belongs to as well as the permissions of
	// the current organization.
	if err = u.loadOrganization(tx, userOrg); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}

//===========================================================================
// Create, Update, and Validate the User Model
//===========================================================================

const (
	insertUserSQL    = "INSERT INTO users (id, name, email, password, terms_agreement, privacy_agreement, last_login, created, modified) VALUES (:id, :name, :email, :password, :agreeTerms, :agreePrivacy, :lastLogin, :created, :modified)"
	insertUserOrgSQL = "INSERT INTO organization_users (user_id, organization_id, role_id, created, modified) VALUES (:userID, :orgID, (SELECT id FROM roles WHERE name=:role), :created, :modified)"
)

// Create a user, inserting the record in the database. If the record already exists or
// a uniqueness constraint is violated an error is returned. The user will also be
// associated with the specified organization and the specified role name. If the
// organization doesn't exist, it will be created. If the role does not exist in the
// database, an error will be returned. This method sets the user ID, created and
// modified timestamps even if they are already set on the model.
func (u *User) Create(ctx context.Context, org *Organization, role string) (err error) {
	u.ID = ulids.New()

	now := time.Now()
	u.SetCreated(now)
	u.SetModified(now)

	if err = u.Validate(); err != nil {
		return err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	params := make([]any, 9)
	params[0] = sql.Named("id", u.ID)
	params[1] = sql.Named("name", u.Name)
	params[2] = sql.Named("email", u.Email)
	params[3] = sql.Named("password", u.Password)
	params[4] = sql.Named("lastLogin", u.LastLogin)
	params[5] = sql.Named("agreeTerms", u.AgreeToS)
	params[6] = sql.Named("agreePrivacy", u.AgreePrivacy)
	params[7] = sql.Named("created", u.Created)
	params[8] = sql.Named("modified", u.Modified)

	if _, err = tx.Exec(insertUserSQL, params...); err != nil {
		var dberr sqlite3.Error
		if errors.As(err, &dberr) {
			if dberr.Code == sqlite3.ErrConstraint {
				return constraint(dberr)
			}
		}
		return err
	}

	// Check if the organization exists, if not create it
	var exists bool
	if exists, _ = org.exists(tx); !exists {
		if err = org.create(tx); err != nil {
			return err
		}
	} else {
		if err = org.populate(tx); err != nil {
			return err
		}
	}

	// Associate the user and the organization with the specified role
	orguser := make([]any, 5)
	orguser[0] = sql.Named("userID", u.ID)
	orguser[1] = sql.Named("orgID", org.ID)
	orguser[2] = sql.Named("role", role)
	orguser[3] = sql.Named("created", u.Created)
	orguser[4] = sql.Named("modified", u.Modified)

	// Associate the user and the role
	if _, err = tx.Exec(insertUserOrgSQL, orguser...); err != nil {
		var dberr sqlite3.Error
		if errors.As(err, &dberr) {
			if dberr.Code == sqlite3.ErrConstraint {
				return constraint(dberr)
			}
		}
		return err
	}

	// Load user in the specified organization or default organization if null is
	// specified; this also verifies the user is part of the organization and caches
	// the organizations and roles the user belongs to as well as the permissions of
	// the current organization.
	if err = u.loadOrganization(tx, org.ID); err != nil {
		return err
	}

	return tx.Commit()
}

const (
	updateUserSQL = "UPDATE users SET name=:name, email=:email, password=:password, terms_agreement=:agreeToS, privacy_agreement=:agreePrivacy, last_login=:lastLogin, modified=:modified WHERE id=:id"
)

// Save a user's name, email, password, agreements, and last login. The modified
// timestamp is set to the current time and neither the ID nor the created timestamp are
// modified. This query is executed as a write-transaction. The user must be fully
// populated and exist in the database for this method to execute successfully.
func (u *User) Save(ctx context.Context) (err error) {
	if err = u.Validate(); err != nil {
		return err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	u.SetModified(time.Now())
	params := make([]any, 8)
	params[0] = sql.Named("id", u.ID)
	params[1] = sql.Named("name", u.Name)
	params[2] = sql.Named("email", u.Email)
	params[3] = sql.Named("password", u.Password)
	params[4] = sql.Named("agreeToS", u.AgreeToS)
	params[5] = sql.Named("agreePrivacy", u.AgreePrivacy)
	params[6] = sql.Named("lastLogin", u.LastLogin)
	params[7] = sql.Named("modified", u.Modified)

	if _, err = tx.Exec(updateUserSQL, params...); err != nil {
		return err
	}
	return tx.Commit()
}

// Validate that the user should be inserted or updated into the database.
func (u *User) Validate() error {
	if ulids.IsZero(u.ID) {
		return invalid(ErrMissingModelID)
	}

	if u.Email == "" || u.Password == "" {
		return invalid(ErrInvalidUser)
	}

	if !passwd.IsDerivedKey(u.Password) {
		return invalid(ErrInvalidPassword)
	}
	return nil
}

const (
	updateLastLoginSQL = "UPDATE users SET last_login=:lastLogin, modified=:modified WHERE id=:id"
)

// UpdateLastLogin is a quick helper method to set the last_login and modified timestamp.
func (u *User) UpdateLastLogin(ctx context.Context) (err error) {
	now := time.Now()
	u.SetLastLogin(now)
	u.SetModified(now)

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(updateLastLoginSQL, sql.Named("id", u.ID), sql.Named("lastLogin", u.LastLogin), sql.Named("modified", u.Modified)); err != nil {
		return err
	}

	return tx.Commit()
}

const (
	getUserSQL = "SELECT id, name, email, terms_agreement, privacy_agreement, last_login, created, modified from users where id in (select user_id FROM organization_users where organization_id = $1)"
)

// ListUsers returns a paginated collection of users filtered by the orgID.
// The orgID must be a valid non-zero value of type ulid.ULID,
// a string representation of a type ulid.ULID, or a slice of bytes
// The number of users resturned is controlled by the prevPage cursor.
// To return the first page with a default number of results pass nil for the prevPage;
// Otherwise pass an empty page with the specified PageSize.
// If the prevPage contains an EndIndex then the next page is returned.
//
// A users slice with the maximum length of the page size will be returned or an
// empty (nil) slice if there are no results. If there is a next page of results, e.g.
// there is another row after the page returned, then a cursor will be returned to
// compute the next page token with.
func ListUsers(ctx context.Context, orgID any, prevPage *pagination.Cursor) (users []*User, cursor *pagination.Cursor, err error) {
	var userOrg ulid.ULID
	if userOrg, err = ulids.Parse(orgID); err != nil {
		return nil, nil, err
	}

	if ulids.IsZero(userOrg) {
		return nil, nil, ErrMissingOrgID
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

	// Fetch list of users associated with the orgID
	var rows *sql.Rows
	if rows, err = tx.Query(getUserSQL, userOrg); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	nRows := int32(0)
	users = make([]*User, 0, prevPage.PageSize)
	for rows.Next() {
		// The query will request one additional message past the page size to check if
		// there is a next page. We should not process any messages after the page size.
		nRows++
		if nRows > prevPage.PageSize {
			continue
		}

		//create user object to append to the users list and add the orgID to it
		user := &User{orgID: userOrg}

		if err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.AgreeToS, &user.AgreePrivacy, &user.LastLogin, &user.Created, &user.Modified); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil, nil
			}
			return nil, nil, err
		}

		//fetch the user's role within the organization
		var role string
		if role, err = user.UserRole(ctx, userOrg, false); err != nil {
			return nil, nil, err
		}
		user.orgRoles[userOrg] = role

		//fetch the permissions associated with the user
		if err = user.fetchPermissions(tx); err != nil {
			return nil, nil, err
		}
		users = append(users, user)
	}

	if err = rows.Close(); err != nil {
		return nil, nil, err
	}

	// Create the cursor to return if there is a next page of results
	if len(users) > 0 && nRows > prevPage.PageSize {
		cursor = pagination.New(users[0].ID.String(), users[len(users)-1].ID.String(), prevPage.PageSize)
	}
	return users, cursor, nil
}

//===========================================================================
// User Organization Management
//===========================================================================

// OrgID returns the organization id that the user was loaded for. If the model doesn't
// have an orgID then an error is returned. This method requires that the user was
// loaded using one of the fetch and catch methods such as GetUserID or that the
// SwitchOrganization method was used to load the user.
func (u *User) OrgID() (ulid.ULID, error) {
	if ulids.IsZero(u.orgID) {
		return ulids.Null, ErrMissingOrgID
	}
	return u.orgID, nil
}

// Role returns the current role for the user in the organization the user was loaded
// for. If the model does not have an orgID or the user doesn't belong to the
// organization then an error is returned. This method requires that the UserRoles have
// been fetched and cached (e.g. that the user was retrieved from the database with an
// organization or that SwitchOrganization) was used.
func (u *User) Role() (role string, _ error) {
	if ulids.IsZero(u.orgID) {
		return "", ErrMissingOrgID
	}

	var ok bool
	if role, ok = u.orgRoles[u.orgID]; !ok {
		return "", ErrUserOrganization
	}
	return role, nil
}

// SwitchOrganization loads the user role and permissions for the specified organization
// returning an error if the user is not in the specified organization.
func (u *User) SwitchOrganization(ctx context.Context, orgID any) (err error) {
	var userOrg ulid.ULID
	if userOrg, err = ulids.Parse(orgID); err != nil {
		return err
	}

	if ulids.IsZero(userOrg) {
		return ErrMissingOrgID
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return err
	}
	defer tx.Rollback()

	if err = u.loadOrganization(tx, userOrg); err != nil {
		return err
	}

	return tx.Commit()
}

// Fetches the organization roles and permissions for the specified orgID. If the orgID
// is Null then one of the user's organizations is used. If the user is not part of the
// organization with that orgID then an error is returned and the orgID on the user is
// not set or changed. This method is used by the GetUser methods as well as the
// SwitchOrganization method but with two different external contexts and validations.
func (u *User) loadOrganization(tx *sql.Tx, orgID ulid.ULID) (err error) {
	// If no orgID is specified, fetch the "default orgID"
	if ulids.IsZero(orgID) {
		if orgID, err = u.defaultOrganization(tx); err != nil {
			return err
		}
	}

	// Decache the current roles and load them again
	u.orgRoles = nil
	if err = u.fetchRoles(tx); err != nil {
		return err
	}

	// If the user is in the specified organization set the orgID, otherwise error.
	if _, ok := u.orgRoles[orgID]; !ok {
		return ErrUserOrganization
	}
	u.orgID = orgID

	// Decache the current permissions and load them again
	u.permissions = nil
	if err = u.fetchPermissions(tx); err != nil {
		return err
	}
	return nil
}

const (
	getDefaultOrgSQL = "SELECT organization_id FROM organization_users WHERE user_id=:userID LIMIT 1"
)

// Fetch the default organization for the user. This method returns at most one orgID,
// even if the user belongs to multiple organizations. It is not guaranteed that
// multiple calls to this method will return the same orgID. If the user doesn't exist
// or is not assigned to an organization an error is returned.
// TODO: right now the first organization is returned, use last logged in organization.
func (u *User) defaultOrganization(tx *sql.Tx) (orgID ulid.ULID, err error) {
	if err = tx.QueryRow(getDefaultOrgSQL, sql.Named("userID", u.ID)).Scan(&orgID); err != nil {
		return orgID, err
	}
	return orgID, nil
}

//===========================================================================
// Cacheing Database Queries
//===========================================================================

const (
	getUserRolesSQL = "SELECT ur.organization_id, r.name FROM organization_users ur JOIN roles r ON ur.role_id=r.id WHERE user_id=:userID"
)

// Returns the name of the user role associated with the user for the specified
// organization. Queries the cached information when the user is fetched unless refresh
// is true, which reloads the cached information from the database on demand.
func (u *User) UserRole(ctx context.Context, orgID ulid.ULID, refresh bool) (role string, err error) {
	if refresh || len(u.orgRoles) == 0 {
		var tx *sql.Tx
		if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
			return "", err
		}
		defer tx.Rollback()

		if err = u.fetchRoles(tx); err != nil {
			return "", err
		}
		tx.Commit()
	}

	var ok bool
	if role, ok = u.orgRoles[orgID]; !ok {
		return "", ErrUserOrganization
	}
	return role, nil
}

func (u *User) fetchRoles(tx *sql.Tx) (err error) {
	u.orgRoles = make(map[ulid.ULID]string)

	var rows *sql.Rows
	if rows, err = tx.Query(getUserRolesSQL, sql.Named("userID", u.ID)); err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			orgID ulid.ULID
			role  string
		)

		if err = rows.Scan(&orgID, &role); err != nil {
			return err
		}
		u.orgRoles[orgID] = role
	}

	return rows.Err()
}

const (
	getUserPermsSQL = "SELECT permission FROM user_permissions WHERE user_id=:userID AND organization_id=:orgID"
)

// Returns the Permissions associated with the user as a list of strings.
// The permissions are cached to prevent multiple queries; use the refresh bool to force
// a new database query to reload the permissions of the user.
func (u *User) Permissions(ctx context.Context, refresh bool) (_ []string, err error) {
	if refresh || len(u.permissions) == 0 {
		var tx *sql.Tx
		if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
			return nil, err
		}
		defer tx.Rollback()

		if err = u.fetchPermissions(tx); err != nil {
			return nil, err
		}
		tx.Commit()
	}
	return u.permissions, nil
}

func (u *User) fetchPermissions(tx *sql.Tx) (err error) {
	if ulids.IsZero(u.orgID) {
		return ErrMissingOrgID
	}

	u.permissions = make([]string, 0)

	var rows *sql.Rows
	if rows, err = tx.Query(getUserPermsSQL, sql.Named("userID", u.ID), sql.Named("orgID", u.orgID)); err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var permission string
		if err = rows.Scan(&permission); err != nil {
			return err
		}
		u.permissions = append(u.permissions, permission)
	}

	return rows.Err()
}

//===========================================================================
// Field Helper Methods
//===========================================================================

// GetLastLogin returns the parsed LastLogin timestamp if it is not null. If it is null
// then a zero-valued timestamp is returned without an error.
func (u *User) GetLastLogin() (time.Time, error) {
	if u.LastLogin.Valid {
		return time.Parse(time.RFC3339Nano, u.LastLogin.String)
	}
	return time.Time{}, nil
}

// SetLastLogin ensures the LastLogin timestamp is serialized to a string correctly.
func (u *User) SetLastLogin(ts time.Time) {
	u.LastLogin = sql.NullString{
		Valid:  true,
		String: ts.Format(time.RFC3339Nano),
	}
}

// SetAgreement marks if the user has accepted the terms of service and privacy policy.
func (u *User) SetAgreement(agreeToS, agreePrivacy bool) {
	u.AgreeToS = sql.NullBool{Valid: true, Bool: agreeToS}
	u.AgreePrivacy = sql.NullBool{Valid: true, Bool: agreePrivacy}
}

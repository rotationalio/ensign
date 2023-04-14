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
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
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
	ID                       ulid.ULID
	Name                     string
	Email                    string
	Password                 string
	AgreeToS                 sql.NullBool
	AgreePrivacy             sql.NullBool
	EmailVerified            bool
	EmailVerificationExpires sql.NullString
	EmailVerificationToken   sql.NullString
	EmailVerificationSecret  []byte
	LastLogin                sql.NullString
	orgID                    ulid.ULID
	orgRoles                 map[ulid.ULID]string
	permissions              []string
}

// User invitations are used to invite users to organizations. Users may not exist at
// the point of invitation, so users are referenced by email address rather than by
// user ID. The crucial part of the invitation is the token, which is a
// cryptographically secure random string that is used to uniquely identify pending
// invitations in the database. Once invitations are accepted then they are deleted
// from the database.
type UserInvitation struct {
	Base
	UserID    ulid.ULID
	OrgID     ulid.ULID
	Email     string
	Role      string
	Expires   string
	Token     string
	Secret    []byte
	CreatedBy ulid.ULID
	name      string
}

const (
	getUserIDSQL    = "SELECT name, email, password, terms_agreement, privacy_agreement, email_verified, email_verification_expires, email_verification_token, email_verification_secret, last_login, created, modified FROM users WHERE id=:id"
	getUserEmailSQL = "SELECT id, name, password, terms_agreement, privacy_agreement, email_verified, email_verification_expires, email_verification_token, email_verification_secret, last_login, created, modified FROM users WHERE email=:email"
	getUserTokenSQL = "SELECT id, name, email, password, terms_agreement, privacy_agreement, email_verified, email_verification_expires, email_verification_secret, last_login, created, modified FROM users WHERE email_verification_token=:token"
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

	if err = tx.QueryRow(getUserIDSQL, sql.Named("id", u.ID)).Scan(&u.Name, &u.Email, &u.Password, &u.AgreeToS, &u.AgreePrivacy, &u.EmailVerified, &u.EmailVerificationExpires, &u.EmailVerificationToken, &u.EmailVerificationSecret, &u.LastLogin, &u.Created, &u.Modified); err != nil {
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

	if err = tx.QueryRow(getUserEmailSQL, sql.Named("email", u.Email)).Scan(&u.ID, &u.Name, &u.Password, &u.AgreeToS, &u.AgreePrivacy, &u.EmailVerified, &u.EmailVerificationExpires, &u.EmailVerificationToken, &u.EmailVerificationSecret, &u.LastLogin, &u.Created, &u.Modified); err != nil {
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

// GetUser by verification token by executing a read-only transaction against the
// database.
func GetUserByToken(ctx context.Context, token string) (u *User, err error) {
	u = &User{
		EmailVerificationToken: sql.NullString{String: token, Valid: true},
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getUserTokenSQL, sql.Named("token", u.EmailVerificationToken)).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.AgreeToS, &u.AgreePrivacy, &u.EmailVerified, &u.EmailVerificationExpires, &u.EmailVerificationSecret, &u.LastLogin, &u.Created, &u.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
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
	insertUserSQL    = "INSERT INTO users (id, name, email, password, terms_agreement, privacy_agreement, email_verified, email_verification_expires, email_verification_token, email_verification_secret, last_login, created, modified) VALUES (:id, :name, :email, :password, :agreeTerms, :agreePrivacy, :emailVerified, :emailExpires, :emailToken, :emailSecret, :lastLogin, :created, :modified)"
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

	params := make([]any, 13)
	params[0] = sql.Named("id", u.ID)
	params[1] = sql.Named("name", u.Name)
	params[2] = sql.Named("email", u.Email)
	params[3] = sql.Named("password", u.Password)
	params[4] = sql.Named("lastLogin", u.LastLogin)
	params[5] = sql.Named("agreeTerms", u.AgreeToS)
	params[6] = sql.Named("agreePrivacy", u.AgreePrivacy)
	params[7] = sql.Named("emailVerified", u.EmailVerified)
	params[8] = sql.Named("emailExpires", u.EmailVerificationExpires)
	params[9] = sql.Named("emailToken", u.EmailVerificationToken)
	params[10] = sql.Named("emailSecret", u.EmailVerificationSecret)
	params[11] = sql.Named("created", u.Created)
	params[12] = sql.Named("modified", u.Modified)

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

	// Add the user to the organization
	if err = u.addOrganizationRole(tx, org, role); err != nil {
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

// Associate the user with the organization and organization role.
func (u *User) addOrganizationRole(tx *sql.Tx, org *Organization, role string) (err error) {
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

	return nil
}

const (
	updateUserSQL = "UPDATE users SET name=:name, email=:email, password=:password, terms_agreement=:agreeToS, privacy_agreement=:agreePrivacy, email_verified=:emailVerified, email_verification_expires=:emailExpires, email_verification_token=:emailToken, email_verification_secret=:emailSecret, last_login=:lastLogin, modified=:modified WHERE id=:id"
)

// Save a user's name, email, password, agreements, verification data, and last login.
// The modified timestamp is set to the current time and neither the ID nor the created
// timestamp are modified. This query is executed as a write-transaction. The user must
// be fully populated and exist in the database for this method to execute successfully.
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
	params := make([]any, 12)
	params[0] = sql.Named("id", u.ID)
	params[1] = sql.Named("name", u.Name)
	params[2] = sql.Named("email", u.Email)
	params[3] = sql.Named("password", u.Password)
	params[4] = sql.Named("agreeToS", u.AgreeToS)
	params[5] = sql.Named("agreePrivacy", u.AgreePrivacy)
	params[6] = sql.Named("emailVerified", u.EmailVerified)
	params[7] = sql.Named("emailExpires", u.EmailVerificationExpires)
	params[8] = sql.Named("emailToken", u.EmailVerificationToken)
	params[9] = sql.Named("emailSecret", u.EmailVerificationSecret)
	params[10] = sql.Named("lastLogin", u.LastLogin)
	params[11] = sql.Named("modified", u.Modified)

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
	userInviteSQL = "INSERT INTO user_invitations (user_id, organization_id, role, email, expires, token, secret, created_by, created, modified) VALUES (:userID, :orgID, :role, :email, :expires, :token, :secret, :createdBy, :created, :modified)"
)

// Create an invitation in the database from the user to the specified email address
// and return the invitation token to send to the user. This method returns an error if
// the invitee is already associated with the organization.
func (u *User) CreateInvite(ctx context.Context, email, role string) (userInvite *UserInvitation, err error) {
	var (
		invite *db.VerificationToken
		userID ulid.ULID
		user   *User
		name   string
	)

	if role == "" {
		return nil, ErrMissingRole
	}

	if !perms.IsRole(role) {
		return nil, ErrInvalidRole
	}

	// Create a token that expires in 7 days
	if invite, err = db.NewVerificationToken(email); err != nil {
		return nil, err
	}

	// Attempt to retrieve the invited user from the database
	if user, err = GetUserEmail(ctx, email, ulid.ULID{}); err != nil {
		if !errors.Is(err, ErrNotFound) {
			return nil, err
		}
	}

	if user != nil {
		// If the user already exists, make sure they are not already in the organization
		if _, err := GetOrgUser(ctx, user.ID, u.orgID); err == nil {
			return nil, ErrUserOrgExists
		} else if !errors.Is(err, ErrNotFound) {
			return nil, err
		}

		// Use the user's ID since they already exist
		userID = user.ID
		name = user.Name
	} else {
		// Create an ID if this is a new user
		userID = ulids.New()
	}

	// Create the invitation
	userInvite = &UserInvitation{
		UserID:    userID,
		OrgID:     u.orgID,
		Email:     email,
		Role:      role,
		Expires:   invite.ExpiresAt.Format(time.RFC3339Nano),
		CreatedBy: u.ID,
		name:      name,
	}

	// Sign the token to ensure we can verify it later
	var secret []byte
	if userInvite.Token, secret, err = invite.Sign(); err != nil {
		return nil, err
	}

	// Set model timestamps
	now := time.Now()
	userInvite.SetCreated(now)
	userInvite.SetModified(now)

	// Save the invitation in the database
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	params := make([]any, 10)
	params[0] = sql.Named("userID", userInvite.UserID)
	params[1] = sql.Named("orgID", userInvite.OrgID)
	params[2] = sql.Named("role", userInvite.Role)
	params[3] = sql.Named("email", userInvite.Email)
	params[4] = sql.Named("expires", userInvite.Expires)
	params[5] = sql.Named("token", userInvite.Token)
	params[6] = sql.Named("secret", secret)
	params[7] = sql.Named("createdBy", userInvite.CreatedBy)
	params[8] = sql.Named("created", userInvite.Created)
	params[9] = sql.Named("modified", userInvite.Modified)

	if _, err = tx.Exec(userInviteSQL, params...); err != nil {
		return nil, err
	}

	return userInvite, tx.Commit()
}

const (
	getUserInviteSQL = "SELECT user_id, organization_id, role, email, expires, token, secret, created_by, created, modified FROM user_invitations WHERE token=:token"
)

// Get an invitation from the database by the token.
func GetUserInvite(ctx context.Context, token string) (invite *UserInvitation, err error) {
	inv := &UserInvitation{}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getUserInviteSQL, sql.Named("token", token)).Scan(&inv.UserID, &inv.OrgID, &inv.Role, &inv.Email, &inv.Expires, &inv.Token, &inv.Secret, &inv.CreatedBy, &inv.Created, &inv.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return inv, tx.Commit()
}

// Validate the invitation against a user provided email address.
func (u *UserInvitation) Validate(email string) (err error) {
	if u.Email != email {
		return ErrInvalidEmail
	}

	// Ensure the role is a recognized role
	if !perms.IsRole(u.Role) {
		return ErrInvalidRole
	}

	token := &db.VerificationToken{
		Email: u.Email,
	}
	if token.ExpiresAt, err = time.Parse(time.RFC3339Nano, u.Expires); err != nil {
		return err
	}

	// Verify that the invite is not expired and signed by Quarterdeck
	if err = token.Verify(u.Token, u.Secret); err != nil {
		return err
	}

	return nil
}

// Name returns the name of the invited user if available.
// TODO: Should this be saved in the database?
func (u *UserInvitation) Name() string {
	return u.name
}

const (
	deleteUserInviteSQL = "DELETE FROM user_invitations WHERE token=:token"
)

// Delete an invitation from the database by the token.
func DeleteInvite(ctx context.Context, token string) (err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(deleteUserInviteSQL, sql.Named("token", token)); err != nil {
		return err
	}

	return tx.Commit()
}

const (
	getUserSQL = "SELECT id, name, email, terms_agreement, privacy_agreement, last_login, created, modified FROM users WHERE id IN (SELECT user_id FROM organization_users"
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

	// Build parameterized query with WHERE clause
	// ---------------------------------------------------------------------------------------------------
	// Query construction with pageSize only:
	// SELECT id, name, email, terms_agreement, privacy_agreement, last_login, created, modified FROM users
	// WHERE id IN (SELECT user_id FROM organization_users WHERE organization_id=:orgID) LIMIT :pageSize
	// ---------------------------------------------------------------------------------------------------
	// Query construction with pageSize and endIndex:
	// SELECT id, name, email, terms_agreement, privacy_agreement, last_login, created, modified FROM users
	// WHERE id IN (SELECT user_id FROM organization_users WHERE organization_id=:orgID) AND id > :endIndex LIMIT :pageSize
	var query strings.Builder
	query.WriteString(getUserSQL)

	// Construct the where clause
	params := make([]any, 0, 4)
	where := make([]string, 0, 3)

	params = append(params, sql.Named("orgID", orgID))
	where = append(where, "organization_id=:orgID)")

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

const (
	verifyUserOrgSQL = "SELECT EXISTS(SELECT 1 FROM organization_users where user_id=:user_id and organization_id=:organization_id)"
	userUpdateSQL    = "UPDATE users SET name=:name, modified=:modified WHERE id=:id"
)

// Update a User in the database.  The requester needs to be in the same orgID as the user.
// This check is performed by verifying that the orgID and the user_id exist in the organization_users table
// The orgID must be a valid non-zero value of type ulid.ULID,
// a string representation of a type ulid.ULID, or a slice of bytes
func (u *User) Update(ctx context.Context, orgID any) (err error) {
	//Validate the ID
	if ulids.IsZero(u.ID) {
		return invalid(ErrMissingModelID)
	}

	//Validate the orgID
	var userOrg ulid.ULID
	if userOrg, err = ulids.Parse(orgID); err != nil {
		return invalid(ErrMissingOrgID)
	}

	//Validate the Name
	if u.Name == "" {
		return invalid(ErrInvalidUser)
	}

	now := time.Now()
	u.SetModified(now)

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	//verify the user_id organization_id mapping
	var exists bool
	if err = tx.QueryRow(verifyUserOrgSQL, sql.Named("user_id", u.ID), sql.Named("organization_id", userOrg)).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}

	if _, err = tx.Exec(userUpdateSQL, sql.Named("id", u.ID), sql.Named("name", u.Name), sql.Named("modified", u.Modified)); err != nil {
		return err
	}

	if err = u.loadOrganization(tx, userOrg); err != nil {
		return err
	}

	return tx.Commit()
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

// AddOrganization adds the user to the specified organization with the specified role.
// An error is returned if the organization doesn't exist.
func (u *User) AddOrganization(ctx context.Context, org *Organization, role string) (err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	// Load the organization from the model
	var exists bool
	if exists, _ = org.exists(tx); !exists {
		return ErrNotFound
	} else {
		if err = org.populate(tx); err != nil {
			return err
		}
	}

	// Add the user to the organization
	if err = u.addOrganizationRole(tx, org, role); err != nil {
		return err
	}

	return tx.Commit()
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

const (
	deleteUserOrgSQL           = "DELETE FROM organization_users WHERE user_id=:userID AND organization_id=:orgID"
	deleteUserInviteByEmailSQL = "DELETE FROM user_invitations WHERE email=:email AND organization_id=:orgID"
)

// RemoveOrganization removes the user from the specified organization. If this results
// in the user having no organizations then the user is deleted. This also deletes all
// invitations that have been sent to the user for the organization and all api keys
// created by the user.
func (u *User) RemoveOrganization(ctx context.Context, orgID any) (err error) {
	var userOrg ulid.ULID
	if userOrg, err = ulids.Parse(orgID); err != nil {
		return err
	}

	if ulids.IsZero(userOrg) {
		return ErrMissingOrgID
	}

	// Delete the rest of the organization user resources in a transaction
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	// Revoke all the keys the user created for this organization
	if err = u.revokeKeys(tx, userOrg); err != nil {
		return err
	}

	// Delete the organization user mapping
	if _, err = tx.Exec(deleteUserOrgSQL, sql.Named("userID", u.ID), sql.Named("orgID", userOrg)); err != nil {
		return err
	}

	// Delete all invitations for the user in the organization
	if _, err = tx.Exec(deleteUserInviteByEmailSQL, sql.Named("email", u.Email), sql.Named("orgID", userOrg)); err != nil {
		return err
	}

	// If the user doesn't have any organizations then delete the user
	if _, err = u.defaultOrganization(tx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err = u.delete(tx); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return tx.Commit()
}

const (
	insertRevokedUserKeysSQL = `INSERT INTO revoked_api_keys
		SELECT k.id, k.key_id, k.name, k.organization_id, k.project_id, k.created_by, k.source, k.user_agent, k.last_used, p.perms, k.created, k.modified
		FROM api_keys k
		JOIN (
			SELECT id, api_key_id, json_group_array(name) AS perms
			FROM api_key_permissions akp
			JOIN permissions p ON p.id=akp.permission_id
			GROUP BY api_key_id
		) p ON p.api_key_id=k.id
		WHERE k.created_by=:userID AND k.organization_id=:orgID`
	deleteUserKeysSQL = "DELETE FROM api_keys WHERE created_by=:userID AND organization_id=:orgID"
)

// RevoveKeys revokes all of the keys that the user has created in the specified
// organization. Note that this does not revoke any keys that were created by other
// users and shared with this user via side channels.
func (u *User) revokeKeys(tx *sql.Tx, orgID ulid.ULID) (err error) {
	// Move the keys to the revoked table
	if _, err = tx.Exec(insertRevokedUserKeysSQL, sql.Named("userID", u.ID), sql.Named("orgID", orgID)); err != nil {
		return err
	}

	// Delete the keys from the live table
	if _, err = tx.Exec(deleteUserKeysSQL, sql.Named("userID", u.ID), sql.Named("orgID", orgID)); err != nil {
		return err
	}

	return nil
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

	// filter u.orgRoles to only the specified orgID
	var key ulid.ULID
	for key = range u.orgRoles {
		if key.Compare(orgID) != 0 {
			delete(u.orgRoles, key)
		}
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

const (
	deleteUserSQL = "DELETE FROM users WHERE id=:userID"
)

// Delete the user from the database. This is normally not done directly but as a
// result of removing the user from all their organizations.
// TODO: Preserve the email address <> user ID mapping.
func (u *User) delete(tx *sql.Tx) (err error) {
	if _, err = tx.Exec(deleteUserSQL, sql.Named("userID", u.ID)); err != nil {
		return err
	}
	return nil
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
	getUserOrgRoleSQL = "SELECT r.name FROM organization_users ur JOIN roles r ON ur.role_id=r.id WHERE user_id=:userID AND organization_id=:orgID"
	getNumOwnersSQL   = "SELECT COUNT(*) FROM organization_users WHERE organization_id=:orgID and role_id IN (SELECT id from roles where name=:ownerRole)"
	userRoleUpdateSQL = "UPDATE organization_users SET modified=:modified, role_id=(SELECT id FROM roles WHERE name=:role) WHERE user_id=:userID AND organization_id=:orgID"
)

// ChangeRole updates the role of the user in specified organization.
func (u *User) ChangeRole(ctx context.Context, orgID any, role string) (err error) {
	// Validate the orgID
	var userOrg ulid.ULID
	if userOrg, err = ulids.Parse(orgID); err != nil {
		return err
	}

	if ulids.IsZero(userOrg) {
		return ErrMissingOrgID
	}

	// Validate the role
	if !perms.IsRole(role) {
		return ErrInvalidRole
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	u.SetModified(time.Now())

	// Get the user's current role
	var currentRole string
	if err = tx.QueryRow(getUserOrgRoleSQL, sql.Named("userID", u.ID), sql.Named("orgID", userOrg)).Scan(&currentRole); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserOrganization
		}
		return err
	}

	// Make sure that the organization has at least one owner if the user is being
	// changed to a non-owner role
	if currentRole == perms.RoleOwner && role != perms.RoleOwner {
		var numOwners int
		if err = tx.QueryRow(getNumOwnersSQL, sql.Named("orgID", userOrg), sql.Named("ownerRole", perms.RoleOwner)).Scan(&numOwners); err != nil {
			return err
		}

		switch numOwners {
		case 0:
			return ErrNoOwnerRole
		case 1:
			return ErrOwnerRoleConstraint
		}
	}

	// Update the organization_users table with the new user role
	if _, err = tx.Exec(userRoleUpdateSQL, sql.Named("userID", u.ID), sql.Named("orgID", userOrg), sql.Named("role", role), sql.Named("modified", u.Modified)); err != nil {
		return err
	}

	if err = u.loadOrganization(tx, userOrg); err != nil {
		return err
	}

	return tx.Commit()
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

func (u *User) ToAPI() (user *api.User, err error) {
	user = &api.User{
		UserID: u.ID,
		Name:   u.Name,
		Email:  u.Email,
	}

	if user.Role, err = u.Role(); err != nil {
		return nil, err
	}

	if user.LastLogin, err = u.GetLastLogin(); err != nil {
		return nil, err
	}

	return user, nil
}

//===========================================================================
// Field Helper Methods
//===========================================================================

// GetVerificationToken returns the verification token for the user if it is not null.
func (u *User) GetVerificationToken() string {
	if u.EmailVerificationToken.Valid {
		return u.EmailVerificationToken.String
	}
	return ""
}

// GetVerificationExpires returns the verification token expiration time for the user
// or a zero time if the token is null.
func (u *User) GetVerificationExpires() (time.Time, error) {
	if u.EmailVerificationExpires.Valid {
		return time.Parse(time.RFC3339Nano, u.EmailVerificationExpires.String)
	}
	return time.Time{}, nil
}

// CreateVerificationToken creates a new verification token for the user, setting the
// email verification fields on the model and returning the token that should be given
// to the user.
func (u *User) CreateVerificationToken() (err error) {
	var (
		verify *db.VerificationToken
		token  string
		secret []byte
	)

	// Create a unqiue token from the user's email address
	if verify, err = db.NewVerificationToken(u.Email); err != nil {
		return err
	}

	// Sign the token to ensure that Quarterdeck can verify it later
	if token, secret, err = verify.Sign(); err != nil {
		return err
	}

	u.EmailVerificationToken = sql.NullString{Valid: true, String: token}
	u.EmailVerificationExpires = sql.NullString{Valid: true, String: verify.ExpiresAt.Format(time.RFC3339Nano)}
	u.EmailVerificationSecret = secret
	return nil
}

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

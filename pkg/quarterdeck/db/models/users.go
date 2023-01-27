package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

// User is a model that represents a row in the users table and provides database
// functionality for interacting with a user's data. It should not be used for API
// serialization. Users may be retrieved from the database either via their ID (e.g.
// from the sub claim in a JWT token) or via their email address (e.g. on login). The
// user password should be stored as an argon2 hash and should be verified using the
// argon2 hashing algorithm.
type User struct {
	Base
	ID           ulid.ULID
	Name         string
	Email        string
	Password     string
	AgreeToS     sql.NullBool
	AgreePrivacy sql.NullBool
	LastLogin    sql.NullString
	orgRoles     map[ulid.ULID]string
	permissions  []string
}

const (
	getUserIDSQL    = "SELECT name, email, password, terms_agreement, privacy_agreement, last_login, created, modified FROM users WHERE id=:id"
	getUserEmailSQL = "SELECT id, name, password, terms_agreement, privacy_agreement, last_login, created, modified FROM users WHERE email=:email"
	countUsersSQL   = "SELECT count(id) FROM users"
)

// CountUsers returns the number of users currently in the database.
func CountUsers(ctx context.Context) (count int64, err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(countUsersSQL).Scan(&count); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return count, nil
}

// GetUser by ID. The ID can be either a string, which is parsed into a ULID or it can
// be a valid ULID. The query is then executed as a read-only transaction against the
// database and the user is returned.
func GetUser(ctx context.Context, id any) (u *User, err error) {
	// Create the user struct and parse the ID input.
	u = &User{}
	switch t := id.(type) {
	case string:
		if u.ID, err = ulid.Parse(t); err != nil {
			return nil, err
		}
	case ulid.ULID:
		u.ID = t
	default:
		return nil, fmt.Errorf("unknown type %T for user id", t)
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

	// Cache permissions on the user
	if err = u.fetchPermissions(tx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}

// GetUser by Email. This query is executed as a read-only transaction.
func GetUserEmail(ctx context.Context, email string) (u *User, err error) {
	u = &User{Email: email}
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

	// Cache permissions on the user
	if err = u.fetchPermissions(tx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}

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
		return err
	}

	// Check if the organization exists, if not create it
	var exists bool
	if exists, _ = org.exists(tx); !exists {
		if err = org.create(tx); err != nil {
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
		return err
	}
	return tx.Commit()
}

const (
	updateUserSQL = "UPDATE users SET name=:name, email=:email, password=:password, terms_agreement=:agreeToS, privacy_agreement=:agreePrivacy, last_login=:lastLogin, modified=:modified WHERE id=:id"
)

// Save a user's name, email, password, and last login. The modified timestamp is set to
// the current time and neither the ID nor the created timestamp is modified. This query
// is executed as a write-transaction. The user must be fully populated and exist in
// the database for this method to execute successfully.
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
	getUserRolesSQL = "SELECT ur.organization_id, r.name FROM organization_users WHERE user_id=:userID"
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
	getUserPermsSQL = "SELECT permission FROM user_permissions WHERE user_id=:userID"
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
	u.permissions = make([]string, 0)

	var rows *sql.Rows
	if rows, err = tx.Query(getUserPermsSQL, sql.Named("userID", u.ID)); err != nil {
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

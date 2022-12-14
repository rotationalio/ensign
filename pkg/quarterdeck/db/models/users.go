package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
)

// User is a model that represents a row in the users table and provides database
// functionality for interacting with a user's data. It should not be used for API
// serialization. Users may be retrieved from the database either via their ID (e.g.
// from the sub claim in a JWT token) or via their email address (e.g. on login). The
// user password should be stored as an argon2 hash and should be verified using the
// argon2 hashing algorithm.
type User struct {
	Base
	ID        ulid.ULID
	Name      string
	Email     string
	Password  string
	LastLogin sql.NullString
}

// UserRole is a model representing a many-to-many mapping between users and rolls
// This model is primarily used by the User and Permission models and is not intended
// for direct use generally.
type UserRole struct {
	Base
	UserID ulid.ULID
	RoleID ulid.ULID
}

const (
	getUserIDSQL    = "SELECT name, email, password, last_login, created, modified FROM users WHERE id=:id"
	getUserEmailSQL = "SELECT id, name, password, last_login, created, modified FROM users WHERE email=:email"
)

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

	if err = tx.QueryRow(getUserIDSQL, sql.Named("id", u.ID)).Scan(&u.Name, &u.Email, &u.Password, &u.LastLogin, &u.Created, &u.Modified); err != nil {
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

// GetUser by Email. This query is executed as a read-only transaction.
func GetUserEmail(ctx context.Context, email string) (u *User, err error) {
	u = &User{Email: email}
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getUserEmailSQL, sql.Named("email", u.Email)).Scan(&u.ID, &u.Name, &u.Password, &u.LastLogin, &u.Created, &u.Modified); err != nil {
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

// Create a user, inserting the record in the database. If the record already exists or
// a uniqueness constraint is violated an error is returned.

const (
	updateUserSQL = "UPDATE users SET name=:name, email=:email, password=:password, last_login=:lastLogin, modified=:modified WHERE id=:id"
)

// Save a user's name, email, password, and last login. The modified timestamp is set to
// the current time and neither the ID nor the created timestamp is modified This query
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
	params := make([]any, 6)
	params[0] = sql.Named("id", u.ID)
	params[1] = sql.Named("name", u.Name)
	params[2] = sql.Named("email", u.Email)
	params[3] = sql.Named("password", u.Password)
	params[4] = sql.Named("lastLogin", u.LastLogin)
	params[5] = sql.Named("modified", u.Modified)

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
	if u.Email == "" || u.Password == "" {
		return ErrInvalidUser
	}
	return nil
}

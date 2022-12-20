package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
)

// APIKey is a model that represents a row in the api_keys table and provides database
// functionality for interacting with api key data. It should not be used for API
// serialization.
type APIKey struct {
	Base
	ID          ulid.ULID
	KeyID       string
	Secret      string
	Name        string
	ProjectID   string
	CreatedBy   sql.NullByte
	LastUsed    sql.NullString
	permissions []string
}

// APIKeyPermission is a model representing a many-to-many mapping between api keys and
// permissions. This model is primarily used by the APIKey and Permission models and is
// not intended for direct use generally.
type APIKeyPermission struct {
	Base
	RoleID       ulid.ULID
	PermissionID int64
}

const (
	getAPIKeySQL = "SELECT id, secret, name, project_id, created_by, last_used, created, modified FROM api_keys WHERE key_id=:keyID"
)

// GetAPIKey by Client ID. This query is executed as a read-only transaction.
func GetAPIKey(ctx context.Context, clientID string) (key *APIKey, err error) {
	key = &APIKey{KeyID: clientID}
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getAPIKeySQL, sql.Named("keyID", key.KeyID)).Scan(&key.ID, &key.Secret, &key.Name, &key.ProjectID, &key.CreatedBy, &key.LastUsed, &key.Created, &key.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Cache permissions on the api key
	if err = key.fetchPermissions(tx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return key, nil
}

// GetLastUsed returns the parsed LastUsed timestamp if it is not null. If it is null
// then a zero-valued timestamp is returned without an error.
func (k *APIKey) GetLastUsed() (time.Time, error) {
	if k.LastUsed.Valid {
		return time.Parse(time.RFC3339Nano, k.LastUsed.String)
	}
	return time.Time{}, nil
}

// SetLastUsed ensures the LastUsed timestamp is serialized to a string correctly.
func (k *APIKey) SetLastUsed(ts time.Time) {
	k.LastUsed = sql.NullString{
		Valid:  true,
		String: ts.Format(time.RFC3339Nano),
	}
}

const (
	updateLastUsedSQL = "UPDATE api_keys SET last_used=:lastUsed, modified=:modified WHERE id=:id"
)

// UpdateLastUsed is a quick helper to set the last_used and modified timestamp.
func (k *APIKey) UpdateLastUsed(ctx context.Context) (err error) {
	now := time.Now()
	k.SetLastUsed(now)
	k.SetModified(now)

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(updateLastUsedSQL, sql.Named("id", k.ID), sql.Named("lastUsed", k.LastUsed), sql.Named("modified", k.Modified)); err != nil {
		return err
	}

	return tx.Commit()
}

const (
	getKeyPermsSQL = "SELECT p.name FROM api_key_permissions akp JOIN permissions p ON p.id=akp.permission_id WHERE akp.api_key_id=:keyID"
)

// Returns the Permissions associated with the user as a list of strings.
// The permissions are cached to prevent multiple queries; use the refresh bool to force
// a new database query to reload the permissions of the user.
func (k *APIKey) Permissions(ctx context.Context, refresh bool) (_ []string, err error) {
	if refresh || len(k.permissions) == 0 {
		var tx *sql.Tx
		if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
			return nil, err
		}
		defer tx.Rollback()

		if err = k.fetchPermissions(tx); err != nil {
			return nil, err
		}
		tx.Commit()
	}
	return k.permissions, nil
}

func (k *APIKey) fetchPermissions(tx *sql.Tx) (err error) {
	k.permissions = make([]string, 0)

	var rows *sql.Rows
	if rows, err = tx.Query(getKeyPermsSQL, sql.Named("keyID", k.ID)); err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var permission string
		if err = rows.Scan(&permission); err != nil {
			return err
		}
		k.permissions = append(k.permissions, permission)
	}

	return rows.Err()
}

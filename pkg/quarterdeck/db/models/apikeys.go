package models

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
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
	ProjectID   ulid.ULID
	CreatedBy   sql.NullByte
	Source      sql.NullString
	UserAgent   sql.NullString
	LastUsed    sql.NullString
	permissions []string
}

// APIKeyPermission is a model representing a many-to-many mapping between api keys and
// permissions. This model is primarily used by the APIKey and Permission models and is
// not intended for direct use generally.
type APIKeyPermission struct {
	Base
	KeyID        ulid.ULID
	PermissionID int64
}

const (
	getAPIKeySQL = "SELECT id, secret, name, project_id, created_by, source, user_agent, last_used, created, modified FROM api_keys WHERE key_id=:keyID"
)

// GetAPIKey by Client ID. This query is executed as a read-only transaction.
func GetAPIKey(ctx context.Context, clientID string) (key *APIKey, err error) {
	key = &APIKey{KeyID: clientID}
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getAPIKeySQL, sql.Named("keyID", key.KeyID)).Scan(&key.ID, &key.Secret, &key.Name, &key.ProjectID, &key.CreatedBy, &key.Source, &key.UserAgent, &key.LastUsed, &key.Created, &key.Modified); err != nil {
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

const (
	insertAPIKeySQL  = "INSERT INTO api_keys (id, key_id, secret, name, project_id, created_by, source, user_agent, last_used, created, modified) VALUES (:id, :keyID, :secret, :name, :projectID, :createdBy, :source, :userAgent, :lastUsed, :created, :modified)"
	insertKeyPermSQL = "INSERT INTO api_key_permissions (api_key_id, permission_id, created, modified) VAlUES (:keyID, (SELECT id FROM permissions WHERE name=:permission), :created, :modified)"
)

// Create an APIKey, inserting the record in the database. If the record already exists
// or a uniqueness constraint is violated an error is returned. Creating the APIKey will
// also associate the permissions with the key. If a permission does not exist in the
// database, an error will be returned. This method sets the ID, created, and modified
// timestamps even if the user has already set them on the model.
func (k *APIKey) Create(ctx context.Context) (err error) {
	k.ID = ulids.New()

	now := time.Now()
	k.SetCreated(now)
	k.SetModified(now)

	if err = k.Validate(); err != nil {
		return err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	params := make([]any, 11)
	params[0] = sql.Named("id", k.ID)
	params[1] = sql.Named("keyID", k.KeyID)
	params[2] = sql.Named("secret", k.Secret)
	params[3] = sql.Named("name", k.Name)
	params[4] = sql.Named("projectID", k.ProjectID)
	params[5] = sql.Named("createdBy", k.CreatedBy)
	params[6] = sql.Named("source", k.Source)
	params[7] = sql.Named("userAgent", k.UserAgent)
	params[8] = sql.Named("lastUsed", k.LastUsed)
	params[9] = sql.Named("created", k.Created)
	params[10] = sql.Named("modified", k.Modified)

	if _, err = tx.Exec(insertAPIKeySQL, params...); err != nil {
		return err
	}

	// Associate the apikey with its permissions
	permparams := make([]any, 4)
	permparams[0] = sql.Named("keyID", k.ID)
	permparams[2] = sql.Named("created", k.Created)
	permparams[3] = sql.Named("modified", k.Modified)

	for _, permission := range k.permissions {
		permparams[1] = sql.Named("permission", permission)
		if _, err = tx.Exec(insertKeyPermSQL, permparams...); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Validate an API key is ready to be inserted into the database. Note that this
// validation does not perform database constraint validation such as if the permission
// foreign keys exist in the database, uniqueness, or not null checks.
// TODO: should we validate timestamps?
func (k *APIKey) Validate() error {
	if ulids.IsZero(k.ID) {
		return ErrMissingModelID
	}

	if k.KeyID == "" || k.Secret == "" {
		return ErrMissingKeyMaterial
	}

	if !passwd.IsDerivedKey(k.Secret) {
		return ErrInvalidSecret
	}

	if k.Name == "" {
		return ErrMissingKeyName
	}

	if ulids.IsZero(k.ProjectID) {
		return ErrMissingProjectID
	}

	if len(k.permissions) == 0 {
		return ErrNoPermissions
	}
	return nil
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

// AddPermissions to an APIKey that has not been created yet. If the APIKey has an ID
// an error is returned since APIKey permissions cannot be modified. This method will
// append permissions to the uncreated Key.
func (k *APIKey) AddPermissions(permissions ...string) error {
	if !ulids.IsZero(k.ID) {
		return ErrModifyPermissions
	}

	// Add the permissions, sort and deduplicate
	k.permissions = append(k.permissions, permissions...)
	sort.Strings(k.permissions)
	for i := len(k.permissions) - 1; i > 0; i-- {
		if k.permissions[i] == k.permissions[i-1] {
			k.permissions = append(k.permissions[:i], k.permissions[i+1:]...)
		}
	}

	return nil
}

// SetPermissions on an APIKey that has not been created yet. If the APIKey has an ID
// an error is returned since APIKey permissions cannot be modified. This method will
// overwrite any permissions already added to the APIKey.
func (k *APIKey) SetPermissions(permissions ...string) error {
	k.permissions = nil
	return k.AddPermissions(permissions...)
}

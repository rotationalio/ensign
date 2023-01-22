package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/keygen"
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
	OrgID       ulid.ULID
	ProjectID   ulid.ULID
	CreatedBy   ulid.ULID
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
	getAPIKeySQL = "SELECT id, secret, name, organization_id, project_id, created_by, source, user_agent, last_used, created, modified FROM api_keys WHERE key_id=:keyID"
	retAPIKeySQL = "SELECT key_id, name, organization_id, project_id, created_by, source, user_agent, last_used, created, modified FROM api_keys WHERE id=:id"
)

// GetAPIKey by Client ID. This query is executed as a read-only transaction. When
// fetching by Client ID we expect that an authentication is being performed, so the
// secret is also fetched.
func GetAPIKey(ctx context.Context, clientID string) (key *APIKey, err error) {
	key = &APIKey{KeyID: clientID}
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getAPIKeySQL, sql.Named("keyID", key.KeyID)).Scan(&key.ID, &key.Secret, &key.Name, &key.OrgID, &key.ProjectID, &key.CreatedBy, &key.Source, &key.UserAgent, &key.LastUsed, &key.Created, &key.Modified); err != nil {
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

// RetrieveAPIKey by ID. This query is executed as a read-only transaction. When
// retrieving a key by ID we expect that this is for informational purposes and not for
// authentication so the secret is not returned.
func RetrieveAPIKey(ctx context.Context, id ulid.ULID) (key *APIKey, err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	key = &APIKey{ID: id}
	if err = populateAPIKey(tx, key); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return key, nil
}

func populateAPIKey(tx *sql.Tx, key *APIKey) (err error) {
	if err = tx.QueryRow(retAPIKeySQL, sql.Named("id", key.ID)).Scan(&key.KeyID, &key.Name, &key.OrgID, &key.ProjectID, &key.CreatedBy, &key.Source, &key.UserAgent, &key.LastUsed, &key.Created, &key.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	// Cache permissions on the api key
	if err = key.fetchPermissions(tx); err != nil {
		return err
	}
	return nil
}

const (
	deleteAPIKeySQL = "DELETE FROM api_keys WHERE id=:id AND organization_id=:orgID"
	revokeAPIKeySQL = "INSERT INTO revoked_api_keys VALUES (:id, :keyID, :name, :orgID, :projectID, :createdBy, :source, :userAgent, :lastUsed, :permissions, :created, :modified)"
)

// DeleteAPIKey by ID restricted to the organization ID supplied. E.g. in order to
// delete an API Key both the key ID and the organization ID must match otherwise an
// ErrNotFound is returned. This method expects to only delete one row at a time and
// rolls back the operation if multiple rows are deleted.
//
// For auditing purposes the api key is deleted from the live api_keys table but then
// inserted without a secret into the revoked_api_keys table. This is because logs and
// other information like events may have API key IDs associated with them; in order to
// trace those keys back to its owners, some information must be preserved. The revoked
// table helps maintain those connections without a lot of constraints.
func DeleteAPIKey(ctx context.Context, id, orgID ulid.ULID) (err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	// Fetch key by ID from the database
	key := &APIKey{ID: id}
	if err = populateAPIKey(tx, key); err != nil {
		return err
	}

	var result sql.Result
	if result, err = tx.Exec(deleteAPIKeySQL, sql.Named("id", id), sql.Named("orgID", orgID)); err != nil {
		return err
	}

	// SQLite supports the RowsAffected Interface so there is no need to check for error.
	nRows, _ := result.RowsAffected()
	if nRows == 0 {
		return ErrNotFound
	} else if nRows > 1 {
		return fmt.Errorf("%d rows were deleted from api_keys table", nRows)
	}

	// Insert into the revoked_api_keys_table
	params := make([]any, 12)
	params[0] = sql.Named("id", key.ID)
	params[1] = sql.Named("keyID", key.KeyID)
	params[2] = sql.Named("name", key.Name)
	params[3] = sql.Named("orgID", key.OrgID)
	params[4] = sql.Named("projectID", key.ProjectID)
	params[5] = sql.Named("createdBy", key.CreatedBy)
	params[6] = sql.Named("source", key.Source)
	params[7] = sql.Named("userAgent", key.UserAgent)
	params[8] = sql.Named("lastUsed", key.LastUsed)
	params[9] = sql.Named("created", key.Created)
	params[10] = sql.Named("modified", key.Modified)

	var permissions []byte
	if permissions, err = json.Marshal(key.permissions); err != nil {
		return err
	}
	params[11] = sql.Named("permissions", string(permissions))

	if _, err = tx.Exec(revokeAPIKeySQL, params...); err != nil {
		return err
	}

	return tx.Commit()
}

const (
	insertAPIKeySQL  = "INSERT INTO api_keys (id, key_id, secret, name, organization_id, project_id, created_by, source, user_agent, last_used, created, modified) VALUES (:id, :keyID, :secret, :name, :orgID, :projectID, :createdBy, :source, :userAgent, :lastUsed, :created, :modified)"
	insertKeyPermSQL = "INSERT INTO api_key_permissions (api_key_id, permission_id, created, modified) VAlUES (:keyID, (SELECT id FROM permissions WHERE name=:permission), :created, :modified)"
)

// Create an APIKey, inserting the record in the database. If the record already exists
// or a uniqueness constraint is violated an error is returned. Creating the APIKey will
// also associate the permissions with the key. If a permission does not exist in the
// database, an error will be returned. This method sets the ID, created, and modified
// timestamps even if the user has already set them on the model. If the APIKey does not
// have a client ID and secret, they're generated before the model is created.
func (k *APIKey) Create(ctx context.Context) (err error) {
	k.ID = ulids.New()

	now := time.Now()
	k.SetCreated(now)
	k.SetModified(now)

	if k.KeyID == "" && k.Secret == "" {
		k.KeyID = keygen.KeyID()
		if k.Secret, err = passwd.CreateDerivedKey(keygen.Secret()); err != nil {
			return fmt.Errorf("could not create derived secret: %s", err)
		}
	}

	if err = k.Validate(); err != nil {
		return err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	params := make([]any, 12)
	params[0] = sql.Named("id", k.ID)
	params[1] = sql.Named("keyID", k.KeyID)
	params[2] = sql.Named("secret", k.Secret)
	params[3] = sql.Named("name", k.Name)
	params[4] = sql.Named("orgID", k.OrgID)
	params[5] = sql.Named("projectID", k.ProjectID)
	params[6] = sql.Named("createdBy", k.CreatedBy)
	params[7] = sql.Named("source", k.Source)
	params[8] = sql.Named("userAgent", k.UserAgent)
	params[9] = sql.Named("lastUsed", k.LastUsed)
	params[10] = sql.Named("created", k.Created)
	params[11] = sql.Named("modified", k.Modified)

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

const (
	updateAPIKeySQL = "UPDATE api_keys SET name=:name, modified=:modified WHERE id=:id AND organization_id=:orgID"
)

// Update an APIKey, modifying the record in the database with the key's ID and OrgID.
// After the key is updated, it is populated with the latest results from the database
// and returned to the user.
//
// NOTE: only the name field is updated so only limited validation is performed.
func (k *APIKey) Update(ctx context.Context) (err error) {
	// Lightweight validation that we can perform the update
	switch {
	case ulids.IsZero(k.ID):
		return invalid(ErrMissingModelID)
	case ulids.IsZero(k.OrgID):
		return invalid(ErrMissingOrgID)
	case k.Name == "":
		return invalid(ErrMissingKeyName)
	}

	// Update the modified timestamp
	k.SetModified(time.Now())

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return err
	}
	defer tx.Rollback()

	params := make([]any, 4)
	params[0] = sql.Named("id", k.ID)
	params[1] = sql.Named("name", k.Name)
	params[2] = sql.Named("orgID", k.OrgID)
	params[3] = sql.Named("modified", k.Modified)

	var result sql.Result
	if result, err = tx.Exec(updateAPIKeySQL, params...); err != nil {
		return err
	}

	// SQLite supports the RowsAffected Interface so there is no need to check for error.
	nRows, _ := result.RowsAffected()
	if nRows == 0 {
		return ErrNotFound
	} else if nRows > 1 {
		return fmt.Errorf("%d rows were updated from api_keys table", nRows)
	}

	// Update the model from the database
	if err = populateAPIKey(tx, k); err != nil {
		return err
	}
	return tx.Commit()
}

// Validate an API key is ready to be inserted into the database. Note that this
// validation does not perform database constraint validation such as if the permission
// foreign keys exist in the database, uniqueness, or not null checks.
// TODO: should we validate timestamps?
func (k *APIKey) Validate() error {
	if ulids.IsZero(k.ID) {
		return invalid(ErrMissingModelID)
	}

	if k.KeyID == "" || k.Secret == "" {
		return invalid(ErrMissingKeyMaterial)
	}

	if !passwd.IsDerivedKey(k.Secret) {
		return invalid(ErrInvalidSecret)
	}

	if k.Name == "" {
		return invalid(ErrMissingKeyName)
	}

	if ulids.IsZero(k.OrgID) {
		return invalid(ErrMissingOrgID)
	}

	if ulids.IsZero(k.ProjectID) {
		return invalid(ErrMissingProjectID)
	}

	if len(k.permissions) == 0 {
		return invalid(ErrNoPermissions)
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

// ToAPI creates a Quarterdeck API response from the model, populating all fields
// except for the ClientSecret since this is not returned in most API requests.
func (k *APIKey) ToAPI(ctx context.Context) *api.APIKey {
	key := &api.APIKey{
		ID:        k.ID,
		ClientID:  k.KeyID,
		Name:      k.Name,
		OrgID:     k.OrgID,
		ProjectID: k.ProjectID,
		CreatedBy: k.CreatedBy,
		Source:    k.Source.String,
		UserAgent: k.UserAgent.String,
	}
	key.Permissions, _ = k.Permissions(ctx, false)
	key.LastUsed, _ = k.GetLastUsed()
	key.Created, _ = k.GetCreated()
	key.Modified, _ = k.GetModified()
	return key
}

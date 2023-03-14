package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/quarterdeck/keygen"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
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
	Partial     bool
	status      APIKeyStatus
	revoked     bool
	permissions []string
}

type APIKeyStatus string

const (
	APIKeyStatusUnknown APIKeyStatus = ""
	APIKeyStatusUnused  APIKeyStatus = "unused"
	APIKeyStatusActive  APIKeyStatus = "active"
	APIKeyStatusStale   APIKeyStatus = "stale"
	APIKeyStatusRevoked APIKeyStatus = "revoked"
)

const APIKeyStalenessThreshold = 90 * 24 * time.Hour

// APIKeyPermission is a model representing a many-to-many mapping between api keys and
// permissions. This model is primarily used by the APIKey and Permission models and is
// not intended for direct use generally.
type APIKeyPermission struct {
	Base
	KeyID        ulid.ULID
	PermissionID int64
}

const (
	listAPIKeySQL = "SELECT id, key_id, name, organization_id, project_id, last_used, created, modified, partial FROM api_keys"
)

// ListAPIKeys returns a paginated collection of APIKeys from the database filtered by
// the orgID and optionally by the projectID. To fetch all keys for an organization,
// pass a zero-valued ULID as the projectID. The number of results returned is
// controlled by the prevPage cursor. To return the first page with a default number of
// results pass nil for the prevPage; otherwise pass an empty page with the specified
// PageSize. If the prevPage contains an EndIndex then the next page is returned.
//
// An apikeys slice with the maximum length of the page size will be returned or an
// empty (nil) slice if there are no results. If there is a next page of results, e.g.
// there is another row after the page returned, then a cursor will be returned to
// compute the next page token with.
func ListAPIKeys(ctx context.Context, orgID, projectID ulid.ULID, prevPage *pagination.Cursor) (keys []*APIKey, cursor *pagination.Cursor, err error) {
	if ulids.IsZero(orgID) {
		return nil, nil, invalid(ErrMissingOrgID)
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
	var query strings.Builder
	query.WriteString(listAPIKeySQL)

	// Construct the where clause
	params := make([]any, 0, 4)
	where := make([]string, 0, 3)

	params = append(params, sql.Named("orgID", orgID))
	where = append(where, "organization_id=:orgID")

	if !ulids.IsZero(projectID) {
		params = append(params, sql.Named("projectID", projectID))
		where = append(where, "project_id=:projectID")
	}

	if prevPage.EndIndex != "" {
		var endIndex ulid.ULID
		if endIndex, err = ulid.Parse(prevPage.EndIndex); err != nil {
			return nil, nil, invalid(ErrInvalidCursor)
		}

		params = append(params, sql.Named("endIndex", endIndex))
		where = append(where, "id > :endIndex")
	}

	// Add the where clause to the query
	query.WriteString(" WHERE ")
	query.WriteString(strings.Join(where, " AND "))

	// Add the limit as the page size + 1 to perform a has next page check.
	params = append(params, sql.Named("pageSize", prevPage.PageSize+1))
	query.WriteString(" LIMIT :pageSize")

	// Fetch rows
	var rows *sql.Rows
	if rows, err = tx.Query(query.String(), params...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	// Process rows into a result page
	nRows := int32(0)
	keys = make([]*APIKey, 0, prevPage.PageSize)
	for rows.Next() {
		// The query will request one additional message past the page size to check if
		// there is a next page. We should not process any messages after the page size.
		nRows++
		if nRows > prevPage.PageSize {
			continue
		}

		apikey := &APIKey{}
		if err = rows.Scan(&apikey.ID, &apikey.KeyID, &apikey.Name, &apikey.OrgID, &apikey.ProjectID, &apikey.LastUsed, &apikey.Created, &apikey.Modified, &apikey.Partial); err != nil {
			return nil, nil, err
		}
		keys = append(keys, apikey)
	}

	if err = rows.Close(); err != nil {
		return nil, nil, err
	}

	// Create the cursor to return if there is a next page of results
	if len(keys) > 0 && nRows > prevPage.PageSize {
		cursor = pagination.New(keys[0].ID.String(), keys[len(keys)-1].ID.String(), prevPage.PageSize)
	}

	// TODO: Iterate over the keys in the revoked table if the user specified all keys.
	return keys, cursor, nil
}

const (
	getAPIKeySQL = "SELECT id, secret, name, organization_id, project_id, created_by, source, user_agent, partial, last_used, created, modified FROM api_keys WHERE key_id=:keyID"
	retAPIKeySQL = "SELECT key_id, name, organization_id, project_id, created_by, source, user_agent, partial, last_used, created, modified FROM api_keys WHERE id=:id"
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

	if err = tx.QueryRow(getAPIKeySQL, sql.Named("keyID", key.KeyID)).Scan(&key.ID, &key.Secret, &key.Name, &key.OrgID, &key.ProjectID, &key.CreatedBy, &key.Source, &key.UserAgent, &key.Partial, &key.LastUsed, &key.Created, &key.Modified); err != nil {
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
	if err = tx.QueryRow(retAPIKeySQL, sql.Named("id", key.ID)).Scan(&key.KeyID, &key.Name, &key.OrgID, &key.ProjectID, &key.CreatedBy, &key.Source, &key.UserAgent, &key.Partial, &key.LastUsed, &key.Created, &key.Modified); err != nil {
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
	insertAPIKeySQL  = "INSERT INTO api_keys (id, key_id, secret, name, organization_id, project_id, created_by, source, user_agent, partial, last_used, created, modified) VALUES (:id, :keyID, :secret, :name, :orgID, :projectID, :createdBy, :source, :userAgent, :partial, :lastUsed, :created, :modified)"
	insertKeyPermSQL = "INSERT INTO api_key_permissions (api_key_id, permission_id, created, modified) VAlUES (:keyID, (SELECT id FROM permissions WHERE name=:permission AND allow_api_keys=true), :created, :modified)"
	updatePartialSQL = "UPDATE api_keys SET partial=(SELECT EXISTS (SELECT p.id FROM permissions p WHERE p.allow_api_keys=true EXCEPT SELECT kp.permission_id FROM api_key_permissions kp WHERE kp.api_key_id=:keyID))"
	queryPartialSQL  = "SELECT partial FROM api_keys WHERE id=:keyID"
)

// Create an APIKey, inserting the record in the database. If the record already exists
// or a uniqueness constraint is violated an error is returned. Creating the APIKey will
// also associate the permissions with the key. If a permission does not exist in the
// database, an error will be returned. This method sets the ID, partial, created, and
// modified timestamps even if the user has already set them on the model. If the APIKey
// does not have a client ID and secret, they're generated before the model is created.
//
// NOTE: the OrgID and ProjectID on the APIKey must be associated in the Quarterdeck
// database otherwise an ErrInvalidProjectID error is returned. Callers should populate
// the OrgID from the claims of the user and NOT from user submitted input.
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

	// Ensure the ProjectID is associated with the OrgID
	op := &OrganizationProject{
		OrgID:     k.OrgID,
		ProjectID: k.ProjectID,
	}
	if projectExists, operr := op.exists(tx); operr != nil || !projectExists {
		if operr != nil {
			return operr
		}
		return invalid(ErrInvalidProjectID)
	}

	params := make([]any, 13)
	params[0] = sql.Named("id", k.ID)
	params[1] = sql.Named("keyID", k.KeyID)
	params[2] = sql.Named("secret", k.Secret)
	params[3] = sql.Named("name", k.Name)
	params[4] = sql.Named("orgID", k.OrgID)
	params[5] = sql.Named("projectID", k.ProjectID)
	params[6] = sql.Named("createdBy", k.CreatedBy)
	params[7] = sql.Named("source", k.Source)
	params[8] = sql.Named("userAgent", k.UserAgent)
	params[9] = sql.Named("partial", k.Partial)
	params[10] = sql.Named("lastUsed", k.LastUsed)
	params[11] = sql.Named("created", k.Created)
	params[12] = sql.Named("modified", k.Modified)

	if _, err = tx.Exec(insertAPIKeySQL, params...); err != nil {
		var dberr sqlite3.Error
		if errors.As(err, &dberr) {
			if dberr.Code == sqlite3.ErrConstraint {
				return constraint(dberr)
			}
		}
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
			var dberr sqlite3.Error
			if errors.As(err, &dberr) {
				if dberr.Code == sqlite3.ErrConstraint {
					return invalid(ErrInvalidPermission)
				}
			}
			return err
		}
	}

	// Update the query with its partial status
	// Set the partial flag by query on the struct
	if _, err = tx.Exec(updatePartialSQL, sql.Named("keyID", k.ID)); err != nil {
		var dberr sqlite3.Error
		if errors.As(err, &dberr) {
			if dberr.Code == sqlite3.ErrConstraint {
				return constraint(dberr)
			}
		}
		return err
	}

	// Fetch the partial status from database
	// NOTE: this is done this way to allow for logicless transactions in Raft replicated SQLite queries
	if err = tx.QueryRow(queryPartialSQL, sql.Named("keyID", k.ID)).Scan(&k.Partial); err != nil {
		return err
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

const (
	APIKeyPermissionsSQL = "SELECT name FROM permissions WHERE allow_api_keys=true"
)

// Fetch all eligible API key permissions from the database as a map for quick checks.
func GetAPIKeyPermissions(ctx context.Context) (permissions map[string]struct{}, err error) {
	permissions = make(map[string]struct{})

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var rows *sql.Rows
	if rows, err = tx.Query(APIKeyPermissionsSQL); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission string
		if err = rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions[permission] = struct{}{}
	}

	return permissions, tx.Commit()
}

// Validate an API key is ready to be inserted into the database. Note that this
// validation does not perform database constraint validation such as if the permission
// foreign keys exist in the database, uniqueness, or not null checks.
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

	if ulids.IsZero(k.CreatedBy) {
		return invalid(ErrMissingCreatedBy)
	}

	if len(k.permissions) == 0 {
		return invalid(ErrNoPermissions)
	}
	return nil
}

// Status of the APIKey based on the LastUsed timestamp if the api keys have not been
// revoked. If the keys have never been used the unused status is returned; if they have
// not been used in 90 days then the stale status is returned; otherwise the apikey is
// considered active unless it has been revoked.
func (k *APIKey) Status() APIKeyStatus {
	if k.status == APIKeyStatusUnknown {
		lastUsed, _ := k.GetLastUsed()
		switch {
		case k.revoked:
			k.status = APIKeyStatusRevoked
		case lastUsed.IsZero():
			k.status = APIKeyStatusUnused
		case time.Since(lastUsed) > APIKeyStalenessThreshold:
			k.status = APIKeyStatusStale
		default:
			k.status = APIKeyStatusActive
		}
	}
	return k.status
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

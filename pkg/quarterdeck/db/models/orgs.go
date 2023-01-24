package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
)

// Organization is a model that represents a row in the organizations table and provides
// database functionality for interacting with an organizations's data. It should not be
// used for API serialization.
type Organization struct {
	Base
	ID     ulid.ULID
	Name   string
	Domain string
}

// OrganizationUser is a model representing a many-to-many mapping between users and
// organizations. This model is primarily used by the User and Organization models and
// is not intended for direct use generally.
type OrganizationUser struct {
	Base
	OrgID  ulid.ULID
	UserID ulid.ULID
}

// OrganizationProject is a model representing the many-to-one mapping between projects
// and organizations. The project model is not stored in the database (but rather in
// the tenant database) so only the projectID is stored, but it must be unique to
// prevent a security hole where a user issues APIKeys to a project in an organization
// that they do not belong to. Before issuing APIKeys with a projectID, Quarterdeck
// checks to ensure that the project actually belongs to the organization via a lookup
// in this table. Otherwise, all information about the project is stored in Tenant.
type OrganizationProject struct {
	Base
	OrgID     ulid.ULID
	ProjectID ulid.ULID
}

const (
	getOrgSQL = "SELECT name, domain, created, modified FROM organizations WHERE id=:id"
)

func GetOrg(ctx context.Context, orgID ulid.ULID) (org *Organization, err error) {
	org = &Organization{ID: orgID}
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getOrgSQL, sql.Named("id", org.ID)).Scan(&org.Name, &org.Domain, &org.Created, &org.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return nil, nil
}

// Exists checks if an organization project mapping exists in order to verify that a
// project is allowed to be associated with an APIKey or other claims resource for the
// user with the specified OrgID claims. Only the OrgID and ProjectID are used for this
// so no preliminary fetch from the database is required to execute the query.
func (op *OrganizationProject) Exists(ctx context.Context) (ok bool, err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return false, err
	}
	defer tx.Rollback()

	if ok, err = op.exists(tx); err != nil {
		return ok, err
	}
	tx.Commit()
	return ok, err
}

const (
	orgProjExistsSQL = "SELECT EXISTS(SELECT 1 FROM organization_projects WHERE organization_id=:orgID AND project_id=:projectID);"
)

// Check if a project is associated with an organization inside of a transaction.
func (op *OrganizationProject) exists(tx *sql.Tx) (ok bool, err error) {
	if err = tx.QueryRow(orgProjExistsSQL, sql.Named("orgID", op.OrgID), sql.Named("projectID", op.ProjectID)).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

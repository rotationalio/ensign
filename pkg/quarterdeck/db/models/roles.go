package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
)

// Role is a model that represents a row in the roles table and provides database
// functionality for interacting with role data. It should not be used for API
// serialization.
type Role struct {
	Base
	ID          int64
	Name        string
	Description sql.NullString
	permissions []*Permission
}

// Permission is a model that represents a row in the permissions table and provides
// database functionality for interacting with permission data. It should not be used
// for API serialization.
type Permission struct {
	Base
	ID           int64
	Name         string
	Description  sql.NullString
	AllowAPIKeys bool
	AllowRoles   bool
}

// RolePermission is a model representing a many-to-many mapping between roles and
// permissions. This model is primarily used by the Role and Permission models and is
// not intended for direct use generally.
type RolePermission struct {
	Base
	RoleID       int64
	PermissionID int64
}

const (
	getRoleSQL      = "SELECT name, description, created, modified FROM roles WHERE id=:roleID"
	getRolePermsSQL = "SELECT p.id, p.name, p.description, p.allow_api_keys, p.allow_roles, p.created, p.modified FROM role_permissions rp JOIN permissions p ON rp.permission_id=p.id WHERE rp.role_id=:roleID"
)

func GetRole(ctx context.Context, roleID int64) (role *Role, err error) {
	role = &Role{
		ID: roleID,
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(getRoleSQL, sql.Named("roleID", role.ID)).Scan(&role.Name, &role.Description, &role.Created, &role.Modified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err = role.fetchPermissions(tx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return role, nil
}

func (r *Role) Permissions(ctx context.Context, refresh bool) (_ []*Permission, err error) {
	if refresh || len(r.permissions) == 0 {
		var tx *sql.Tx
		if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
			return nil, err
		}
		defer tx.Rollback()

		if err = r.fetchPermissions(tx); err != nil {
			return nil, err
		}

		tx.Commit()
	}
	return r.permissions, nil
}

func (r *Role) fetchPermissions(tx *sql.Tx) (err error) {
	if r.ID == 0 {
		return ErrMissingModelID
	}

	r.permissions = make([]*Permission, 0)

	var rows *sql.Rows
	if rows, err = tx.Query(getRolePermsSQL, sql.Named("roleID", r.ID)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	defer rows.Close()

	for rows.Next() {
		p := &Permission{}
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.AllowAPIKeys, &p.AllowRoles, &p.Created, &p.Modified); err != nil {
			return err
		}
		r.permissions = append(r.permissions, p)
	}

	return nil
}

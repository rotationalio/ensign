package models

import (
	"context"
	"database/sql"

	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
)

const (
	countUsersSQL = "SELECT count(id) FROM users"
	countOrgsSQL  = "SELECT count(id) FROM organizations"
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

// CountOrganizations returns the number of organizations currently in the database.
func CountOrganizations(ctx context.Context) (count int64, err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(countOrgsSQL).Scan(&count); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return count, nil
}

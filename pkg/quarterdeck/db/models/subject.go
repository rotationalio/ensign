package models

import (
	"context"
	"database/sql"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

type SubjectType uint8

const (
	UnknownSubject SubjectType = iota
	UserSubject
	APIKeySubject
)

const identifySubjectSQL = `WITH user_exists AS (
	SELECT EXISTS(
		SELECT user_id FROM organization_users WHERE organization_id=:orgID AND user_id=:subjectID
	)
), apikey_exists AS (
	SELECT EXISTS(
		SELECT id FROM api_keys WHERE organization_id=:orgID AND id=:subjectID
	)
) SELECT * FROM user_exists, apikey_exists;`

// IdentifySubject is a helper tool to check the database to determine if the specified
// subject from authentication claims refers to a user or to an apikey. This method
// relies on the fact that the database IDs are ULIDs, meaning that it is highly
// unlikely that a user id and an apikey id are the same (technically possible, but with
// the same collision probability as a UUID). This is made further unlikely by requiring
// an organization ID is specified to return the subject type.
func IdentifySubject(ctx context.Context, subjectID, orgID any) (_ SubjectType, err error) {
	var sub ulid.ULID
	if sub, err = ulids.Parse(subjectID); err != nil {
		return UnknownSubject, err
	}

	var org ulid.ULID
	if org, err = ulids.Parse(orgID); err != nil {
		return UnknownSubject, err
	}

	if ulids.IsZero(sub) || ulids.IsZero(org) {
		return UnknownSubject, ErrNotFound
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return UnknownSubject, err
	}
	defer tx.Rollback()

	var isUser, isAPIKey bool
	if err = tx.QueryRow(identifySubjectSQL, sql.Named("subjectID", sub), sql.Named("orgID", org)).Scan(&isUser, &isAPIKey); err != nil {
		return UnknownSubject, err
	}

	switch {
	case isUser && !isAPIKey:
		return UserSubject, nil
	case isAPIKey && !isUser:
		return APIKeySubject, nil
	case isUser && isAPIKey:
		return UnknownSubject, ErrSubjectCollision
	default:
		return UnknownSubject, ErrNotFound
	}
}

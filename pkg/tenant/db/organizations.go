package db

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

const OrganizationNamespace = "organizations"

// Use the resourceID to retrieve the orgID from the database.
func GetOrgIndex(ctx context.Context, resourceID ulid.ULID) (orgID ulid.ULID, err error) {
	if ulids.IsZero(resourceID) {
		return ulid.ULID{}, ErrMissingID
	}

	if err = orgID.UnmarshalBinary(resourceID[:]); err != nil {
		return ulid.ULID{}, err
	}

	return orgID, nil
}

// Store the resourceID as a key and the orgID as a value in the database.
func PutOrgIndex(ctx context.Context, resourceID, orgID ulid.ULID) error {
	if ulids.IsZero(resourceID) {
		return ErrMissingID
	}

	if ulids.IsZero(orgID) {
		return ErrMissingOrgID
	}

	if err := putRequest(ctx, OrganizationNamespace, resourceID[:], orgID[:]); err != nil {
		return err
	}

	return nil
}

// VerifyOrg will check that resources are allocated to the correct organization.
// The method will take in an orgID from the claims and will return true if the orgID
// from the database is the same and an error if it is not.
func VerifyOrg(ctx context.Context, claimsOrgID, resourceID ulid.ULID) (bool, error) {
	if ulids.IsZero(claimsOrgID) {
		return false, ErrMissingOrgID
	}

	if ulids.IsZero(resourceID) {
		return false, ErrMissingID
	}

	orgID, err := GetOrgIndex(ctx, resourceID)
	if err != nil {
		return false, err
	}

	err = PutOrgIndex(ctx, resourceID, orgID)
	if err != nil {
		return false, err
	}
	if orgID.Compare(claimsOrgID) == 0 {
		return true, nil
	} else {
		return false, ErrOrgNotVerified
	}
}

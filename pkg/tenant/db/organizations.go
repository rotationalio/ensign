package db

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

const OrganizationNamespace = "organizations"

// Use the resourceID to retrieve the orgID from the database.
func GetOrgIndex(ctx context.Context, resourceID ulid.ULID) (orgID ulid.ULID, err error) {
	if ulids.IsZero(resourceID) {
		return ulid.ULID{}, ErrMissingID
	}

	var data []byte

	if data, err = getRequest(ctx, OrganizationNamespace, resourceID[:]); err != nil {
		return orgID, err
	}

	if err = orgID.UnmarshalBinary(data); err != nil {
		return orgID, err
	}

	return orgID, nil
}

// Store the resourceID as a key and the orgID as a value.
func PutOrgIndex(ctx context.Context, resourceID, orgID ulid.ULID) error {
	if err := putRequest(ctx, OrganizationNamespace, resourceID[:], orgID[:]); err != nil {
		return err
	}

	return nil
}

// VerifyOrg will check that resources are allocated to the correct organization.
// The method will take in an orgID and will return true if the orgID of a resource
// (tenant, member, project, topic, api key) is the same and an error if it is not.
func VerifyOrg(ctx context.Context, claimsOrgID, resourceID ulid.ULID) (bool, error) {
	if ulids.IsZero(claimsOrgID) {
		return false, ErrMissingOrgID
	}

	if ulids.IsZero(resourceID) {
		return false, ErrMissingID
	}

	orgID, err := GetOrgIndex(ctx, resourceID)
	if err != nil {
		fmt.Println(orgID)
		return false, err
	}

	err = PutOrgIndex(ctx, resourceID, orgID)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	fmt.Println(err)
	if orgID.Compare(claimsOrgID) == 0 {
		return true, nil
	} else {
		return false, ErrOrgNotVerified
	}
}

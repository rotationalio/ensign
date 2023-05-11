package db

import (
	"context"
	"time"

	"github.com/gosimple/slug"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
)

// CreateUserResources creates all the necessary database objects for a new user given
// a partially constructed member model. This method should be called after a new user
// has been successfully registered with Quarterdeck in order to allow the user to
// access default resources such as the tenant and user profile info when they login.
func CreateUserResources(ctx context.Context, orgName string, member *Member) (err error) {
	// Ensure the user data is valid before creating anything
	if err = member.Validate(); err != nil {
		return err
	}

	// New user should have a tenant
	tenant := &Tenant{
		OrgID:           member.OrgID,
		Name:            slug.Make(orgName),
		EnvironmentType: "development",
	}
	if err = CreateTenant(ctx, tenant); err != nil {
		return err
	}

	// Create the member record for the user
	if err = CreateMember(ctx, member); err != nil {
		return err
	}

	return nil
}

// UpdateLastLogin is a helper method that updates the last login time for a member
// given an access token. This should normally be called in a background task to
// prevent blocking the user from logging in.
func UpdateLastLogin(ctx context.Context, accessToken string, login time.Time) (err error) {
	// Parse the claims from the access token
	var claims *tokens.Claims
	if claims, err = tokens.ParseUnverifiedTokenClaims(accessToken); err != nil {
		return err
	}

	// Retrieve the member record to update
	// TODO: This should be in a trtl transaction to prevent updates being stomped
	var member *Member
	if member, err = RetrieveMember(ctx, claims.ParseOrgID(), claims.ParseUserID()); err != nil {
		return err
	}

	// Update the last login time
	member.LastActivity = login
	if err = UpdateMember(ctx, member); err != nil {
		return err
	}

	return nil
}

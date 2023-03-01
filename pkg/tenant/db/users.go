package db

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
)

// CreateUserResources creates all the necessary database objects for a new user given
// a partially constructed member model. This method should be called after a new user
// has been successfully registered with Quarterdeck in order to allow the user to
// access default resources such as the tenant and project when they login.
func CreateUserResources(ctx context.Context, projectID ulid.ULID, member *Member) (err error) {
	// Ensure the user data is valid before creating anything
	if err = member.Validate(); err != nil {
		return err
	}

	// New user should have a tenant
	tenant := &Tenant{
		OrgID:           member.OrgID,
		Name:            fmt.Sprintf("%s's Tenant", member.Name),
		EnvironmentType: "development",
	}
	if err = CreateTenant(ctx, tenant); err != nil {
		return err
	}

	// Create the member record for the user
	if err = CreateTenantMember(ctx, member); err != nil {
		return err
	}

	// New user should have a default project
	project := &Project{
		ID:       projectID,
		OrgID:    member.OrgID,
		TenantID: tenant.ID,
		Name:     fmt.Sprintf("%s's Project", member.Name),
	}
	if err = CreateTenantProject(ctx, project); err != nil {
		return err
	}

	return nil
}

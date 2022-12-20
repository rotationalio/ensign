package models

import "github.com/oklog/ulid/v2"

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

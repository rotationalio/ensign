package permissions

import "strings"

// These permissions are used to authorize user requests and should match the defined
// permissions in the quarterdeck database.
//
// NOTE: if adding or removing permissions from this list, they also need to be updated
// in a database migration. Please also ensure that the AllPermissions variable is also
// updated to ensure that the tests pass.
const (
	// Organizations management
	EditOrganizations   = "organizations:edit"
	DeleteOrganizations = "organizations:delete"
	ReadOrganizations   = "organizations:read"

	// Organization collaborators management
	AddCollaborators    = "collaborators:add"
	RemoveCollaborators = "collaborators:remove"
	EditCollaborators   = "collaborators:edit"
	ReadCollaborators   = "collaborators:read"

	// Organization projects management
	EditProjects   = "projects:edit"
	DeleteProjects = "projects:delete"
	ReadProjects   = "projects:read"

	// Project API Keys management
	EditAPIKeys   = "apikeys:edit"
	DeleteAPIKeys = "apikeys:delete"
	ReadAPIKeys   = "apikeys:read"

	// Project topics management
	CreateTopics  = "topics:create"
	EditTopics    = "topics:edit"
	DestroyTopics = "topics:destroy"
	ReadTopics    = "topics:read"

	// Eventing permissions
	ReadMetrics = "metrics:read"
	Publisher   = "publisher"
	Subscriber  = "subscriber"
)

// Prefixes allow for easy checking of permission groups
const (
	PrefixOrganizations = "organizations:"
	PrefixCollaborators = "collaborators:"
	PrefixProjects      = "projects"
	PrefixAPIKeys       = "apikeys:"
	PrefixTopics        = "topics:"
	PrefixMetrics       = "metrics:"
)

// Roles define collections of permissions; these constants are the roles defined in
// the Quarterdeck datbase and should be kept up to date with the database schema.
const (
	RoleOwner    = "Owner"
	RoleAdmin    = "Admin"
	RoleMember   = "Member"
	RoleObserver = "Observer"
)

// AllPermissions contains the list of all available permissions and is primarily used
// for testing or determining if a string is a valid permission without doing a database
// query. It maps the permission string to the primary key of the permission, helping
// with database migration generation.
var AllPermissions = map[string]uint8{
	EditOrganizations:   1,
	DeleteOrganizations: 2,
	ReadOrganizations:   3,
	AddCollaborators:    4,
	RemoveCollaborators: 5,
	EditCollaborators:   6,
	ReadCollaborators:   7,
	EditProjects:        8,
	DeleteProjects:      9,
	ReadProjects:        10,
	EditAPIKeys:         11,
	DeleteAPIKeys:       12,
	ReadAPIKeys:         13,
	CreateTopics:        14,
	EditTopics:          15,
	DestroyTopics:       16,
	ReadTopics:          17,
	ReadMetrics:         18,
	Publisher:           19,
	Subscriber:          20,
}

// InGroup is a quick test to check if a permission belongs to the specified group.
// E.g. if the "topics:read" permission is part of the "topics" group based on the prefix.
func InGroup(permission, group string) bool {
	return strings.HasPrefix(permission, group)
}

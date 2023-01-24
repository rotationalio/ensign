package permissions

// These permissions are used to authorize user requests and should match the defined
// permissions in the quarterdeck database.
//
// NOTE: if adding or removing permissions from this list, they also need to be updated
// in a database migration. Please also ensure that the AllPermissions variable is also
// updated to ensure that the tests pass.
const (
	// Organizations management
	CreateOrganizations = "organizations:create"
	DeleteOrganizations = "organizations:delete"
	ListOrganizations   = "organizations:list"
	EditOrganizations   = "organizations:edit"
	DetailOrganizations = "organizations:detail"

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
	CreateOrganizations: 1,
	DeleteOrganizations: 2,
	ListOrganizations:   3,
	EditOrganizations:   4,
	DetailOrganizations: 5,
	AddCollaborators:    6,
	RemoveCollaborators: 7,
	EditCollaborators:   8,
	ReadCollaborators:   9,
	EditProjects:        10,
	DeleteProjects:      11,
	ReadProjects:        12,
	EditAPIKeys:         13,
	DeleteAPIKeys:       14,
	ReadAPIKeys:         15,
	CreateTopics:        16,
	EditTopics:          17,
	DestroyTopics:       18,
	ReadTopics:          19,
	ReadMetrics:         20,
	Publisher:           21,
	Subscriber:          22,
}

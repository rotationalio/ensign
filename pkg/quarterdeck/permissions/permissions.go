package permissions

// These permissions are used to authorize user requests and should match the defined
// permissions in the quarterdeck database.
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

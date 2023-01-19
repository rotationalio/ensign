package tenant

const (
	// Tenant management
	ReadTenantPermission   = "read:tenant"
	WriteTenantPermission  = "write:tenant"
	DeleteTenantPermission = "delete:tenant"

	// Members management
	ReadMemberPermission   = "read:member"
	WriteMemberPermission  = "write:member"
	DeleteMemberPermission = "delete:member"

	// Projects management
	ReadProjectPermission   = "read:project"
	WriteProjectPermission  = "write:project"
	DeleteProjectPermission = "delete:project"

	// Topics management
	ReadTopicPermission   = "read:topic"
	WriteTopicPermission  = "write:topic"
	DeleteTopicPermission = "delete:topic"

	// API Keys management
	ReadAPIKey   = "read:key"
	WriteAPIKey  = "write:key"
	DeleteAPIKey = "delete:key"
)

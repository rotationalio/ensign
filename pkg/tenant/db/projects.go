package db

import (
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/net/context"
)

const (
	ProjectNamespace     = "projects"
	ProjectKeysNamespace = "project_keys"
)

type Project struct {
	OrgID    ulid.ULID `msgpack:"org_id"`
	TenantID ulid.ULID `msgpack:"tenant_id"`
	ID       ulid.ULID `msgpack:"id"`
	Name     string    `msgpack:"name"`
	Created  time.Time `msgpack:"created"`
	Modified time.Time `msgpack:"modified"`
}

// ProjectKey stores the components of the project key to enable direct lookup from the
// project ID.
type ProjectKey struct {
	TenantID ulid.ULID `msgpack:"tenant_id"`
	ID       ulid.ULID `msgpack:"id"`
}

var _ Model = &Project{}
var _ Model = &ProjectKey{}

// Key is a 32 composite key combining the tenant id and the project id.
func (p *Project) Key() (key []byte, err error) {
	// Key requires a tenant id and project id.
	if ulids.IsZero(p.ID) {
		return nil, ErrMissingID
	}

	if ulids.IsZero(p.TenantID) {
		return nil, ErrMissingTenantID
	}

	// Create a 32 byte array so that the first 16 bytes hold the tenant id
	// and the last 16 bytes hold the project id.
	key = make([]byte, 32)

	// Marshal the tenant id to the first 16 bytes of the key.
	if err = p.TenantID.MarshalBinaryTo(key[0:16]); err != nil {
		return nil, err
	}

	// Marshal the project id to the last 16 bytes of the key.
	if err = p.ID.MarshalBinaryTo(key[16:]); err != nil {
		return nil, err
	}
	return key, err
}

func (p *Project) Namespace() string {
	return ProjectNamespace
}

func (p *Project) MarshalValue() ([]byte, error) {
	return msgpack.Marshal(p)
}

func (p *Project) UnmarshalValue(data []byte) error {
	return msgpack.Unmarshal(data, p)
}

func (p *Project) Validate() error {
	if ulids.IsZero(p.OrgID) {
		return ErrMissingOrgID
	}

	if p.Name == "" {
		return ErrMissingProjectName
	}

	if !alphaNum.MatchString(p.Name) {
		return ErrInvalidProjectName
	}

	return nil
}

// Convert the model to an API response.
func (p *Project) ToAPI() *api.Project {
	return &api.Project{
		ID:       p.ID.String(),
		Name:     p.Name,
		Created:  TimeToString(p.Created),
		Modified: TimeToString(p.Modified),
	}
}

func (p *ProjectKey) Key() (key []byte, err error) {
	if ulids.IsZero(p.ID) {
		return nil, ErrMissingID
	}

	return p.ID[:], nil
}

func (p *ProjectKey) Namespace() string {
	return ProjectKeysNamespace
}

func (p *ProjectKey) MarshalValue() ([]byte, error) {
	return msgpack.Marshal(p)
}

func (p *ProjectKey) UnmarshalValue(data []byte) error {
	return msgpack.Unmarshal(data, p)
}

// CreateTenantProject adds a new project to a tenant in the database.
// Note: If a project id is not passed in by the User, a new project id will be generated.
func CreateTenantProject(ctx context.Context, project *Project) (err error) {
	if ulids.IsZero(project.ID) {
		project.ID = ulids.New()
	}

	if ulids.IsZero(project.TenantID) {
		return ErrMissingTenantID
	}

	// Validate project data.
	if err = project.Validate(); err != nil {
		return err
	}

	project.Created = time.Now()
	project.Modified = project.Created

	if err = Put(ctx, project); err != nil {
		return err
	}

	// Store the project key in the database to allow direct lookups by project id.
	key := &ProjectKey{
		TenantID: project.TenantID,
		ID:       project.ID,
	}
	if err = Put(ctx, key); err != nil {
		return err
	}
	return nil
}

// CreateProject adds a new project to an organization in the database.
// Note: If a project id is not passed in by the User, a new project id will be generated.
func CreateProject(ctx context.Context, project *Project) (err error) {
	if ulids.IsZero(project.ID) {
		project.ID = ulids.New()
	}

	project.Created = time.Now()
	project.Modified = project.Created

	if err = Put(ctx, project); err != nil {
		return err
	}

	// Store the project key in the database to allow direct lookups by project id.
	key := &ProjectKey{
		TenantID: project.TenantID,
		ID:       project.ID,
	}
	if err = Put(ctx, key); err != nil {
		return err
	}
	return nil
}

// RetrieveProject gets a project from the database with the given project id.
func RetrieveProject(ctx context.Context, projectID ulid.ULID) (project *Project, err error) {
	// Lookup the project key in the database
	key := &ProjectKey{
		ID: projectID,
	}

	if err = Get(ctx, key); err != nil {
		return nil, err
	}

	// Use the key to lookup the project
	project = &Project{
		TenantID: key.TenantID,
		ID:       key.ID,
	}
	if err = Get(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// ListProjects retrieves all projects assigned to a tenant.
func ListProjects(ctx context.Context, tenantID ulid.ULID) (projects []*Project, err error) {
	// Store the tenant ID as the prefix.
	var prefix []byte
	if tenantID.Compare(ulid.ULID{}) != 0 {
		prefix = tenantID[:]
	}

	// TODO: Use the cursor directly instead of having duplicate data in memory.
	var values [][]byte
	if values, err = List(ctx, prefix, ProjectNamespace); err != nil {
		return nil, err
	}

	// Parse the projects from the data
	projects = make([]*Project, 0, len(values))
	for _, data := range values {
		project := &Project{}
		if err = project.UnmarshalValue(data); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

// UpdateProject updates the record of a project by its id.
func UpdateProject(ctx context.Context, project *Project) (err error) {
	if ulids.IsZero(project.ID) {
		return ErrMissingID
	}

	// Validate project data.
	if err = project.Validate(); err != nil {
		return err
	}

	// Retrieve the project key to update the project.
	// Note: There is a possible concurrency issue if the project is deleted between
	// Get and Put.
	key := &ProjectKey{
		ID: project.ID,
	}
	if err = Get(ctx, key); err != nil {
		return err
	}

	project.Modified = time.Now()
	if project.Created.IsZero() {
		project.Created = project.Modified
	}

	project.TenantID = key.TenantID
	if err = Put(ctx, project); err != nil {
		return err
	}
	return nil
}

// DeleteProject deletes a project with the given project id.
func DeleteProject(ctx context.Context, projectID ulid.ULID) (err error) {
	project := &Project{
		ID: projectID,
	}

	// Retrieve the project key to delete the project.
	key := &ProjectKey{
		ID: projectID,
	}
	if err = Get(ctx, key); err != nil {
		return err
	}

	// Delete the project and its key from the database.
	project.TenantID = key.TenantID
	if err = Delete(ctx, project); err != nil {
		return err
	}

	if err = Delete(ctx, key); err != nil {
		return err
	}
	return nil
}

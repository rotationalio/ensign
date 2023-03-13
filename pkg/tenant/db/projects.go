package db

import (
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/net/context"
)

const ProjectNamespace = "projects"

type Project struct {
	OrgID    ulid.ULID `msgpack:"org_id"`
	TenantID ulid.ULID `msgpack:"tenant_id"`
	ID       ulid.ULID `msgpack:"id"`
	Name     string    `msgpack:"name"`
	Created  time.Time `msgpack:"created"`
	Modified time.Time `msgpack:"modified"`
}

var _ Model = &Project{}

// Key is a 32 composite key combining the tenant id and the project id.
func (p *Project) Key() (key []byte, err error) {
	// Key requires a tenant id and project id.
	if ulids.IsZero(p.ID) {
		return nil, ErrMissingID
	}

	if ulids.IsZero(p.TenantID) {
		return nil, ErrMissingTenantID
	}

	var k Key
	if k, err = CreateKey(p.TenantID, p.ID); err != nil {
		return nil, err
	}

	return k.MarshalValue()
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

	if strings.TrimSpace(p.Name) == "" {
		return ErrMissingProjectName
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
	if err = PutObjectKey(ctx, project); err != nil {
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
	if err = PutObjectKey(ctx, project); err != nil {
		return err
	}
	return nil
}

// RetrieveProject gets a project from the database with the given project id.
func RetrieveProject(ctx context.Context, projectID ulid.ULID) (project *Project, err error) {
	// Lookup the project key in the database
	var key []byte
	if key, err = GetObjectKey(ctx, projectID); err != nil {
		return nil, err
	}

	// Use the key to lookup the project
	var data []byte
	if data, err = getRequest(ctx, ProjectNamespace, key); err != nil {
		return nil, err
	}

	// Unmarshal the data into the project
	project = &Project{}
	if err = project.UnmarshalValue(data); err != nil {
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

// UpdateProject updates the record of a project from its database model.
func UpdateProject(ctx context.Context, project *Project) (err error) {
	if ulids.IsZero(project.ID) {
		return ErrMissingID
	}

	// Validate project data.
	if err = project.Validate(); err != nil {
		return err
	}

	project.Modified = time.Now()
	if project.Created.IsZero() {
		project.Created = project.Modified
	}

	return Put(ctx, project)
}

// DeleteProject deletes a project with the given project id.
func DeleteProject(ctx context.Context, projectID ulid.ULID) (err error) {
	// Retrieve the project key to delete the project.
	var key []byte
	if key, err = GetObjectKey(ctx, projectID); err != nil {
		return err
	}

	// Delete the project and its key from the database.
	if err = deleteRequest(ctx, ProjectNamespace, key); err != nil {
		return err
	}

	if err = DeleteObjectKey(ctx, key); err != nil {
		return err
	}
	return nil
}

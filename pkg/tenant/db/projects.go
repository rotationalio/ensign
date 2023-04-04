package db

import (
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
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

	// Store the project tenant ID as a key and project org ID as a value in the database for org verification.
	if err = PutOrgIndex(ctx, project.TenantID, project.OrgID); err != nil {
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

	// Store the project ID as a key and project org ID as a value in the database for org verification.
	if err = PutOrgIndex(ctx, project.ID, project.OrgID); err != nil {
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

// ListProjects retrieves a paginated list of projects.
func ListProjects(ctx context.Context, tenantID ulid.ULID, c *pg.Cursor) (projects []*Project, cursor *pg.Cursor, err error) {
	// Store the tenant ID as the prefix.
	var prefix []byte
	if tenantID.Compare(ulid.ULID{}) != 0 {
		prefix = tenantID[:]
	}

	var seekKey []byte
	if c.EndIndex != "" {
		var start ulid.ULID
		if start, err = ulid.Parse(c.EndIndex); err != nil {
			return nil, nil, err
		}
		seekKey = start[:]
	}

	// Check to see if a default cursor exists and create one if it does not.
	if c == nil {
		c = pg.New("", "", 0)
	}

	if c.PageSize <= 0 {
		return nil, nil, ErrMissingPageSize
	}

	projects = make([]*Project, 0)
	onListItem := func(item *trtl.KVPair) error {
		project := &Project{}
		if err = project.UnmarshalValue(item.Value); err != nil {
			return err
		}
		projects = append(projects, project)
		return nil
	}

	if cursor, err = List(ctx, prefix, seekKey, ProjectNamespace, onListItem, c); err != nil {
		return nil, nil, err
	}

	return projects, cursor, nil
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

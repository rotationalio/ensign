package db

import (
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/net/context"
)

const ProjectNamespace = "projects"

type Project struct {
	TenantID ulid.ULID `msgpack:"tenant_id"`
	ID       ulid.ULID `msgpack:"id"`
	Name     string    `msgpack:"name"`
	Created  time.Time `msgpack:"created"`
	Modified time.Time `msgpack:"modified"`
}

var _ Model = &Project{}

// Key is a 32 composite key combining the tenant id and the project id.
func (p *Project) Key() (key []byte, err error) {
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

// CreateProject adds a new project to the database.
// Note: If a project id is not passed in by the User, a new id will be generated.
func CreateProject(ctx context.Context, project *Project) (err error) {
	if project.ID.Compare(ulid.ULID{}) == 0 {
		project.ID = ulid.Make()
	}

	project.Created = time.Now()
	project.Modified = project.Created

	if err = Put(ctx, project); err != nil {
		return err
	}
	return nil
}

// RetrieveProject gets a project from the database with a given id.
func RetrieveProject(ctx context.Context, id ulid.ULID) (project *Project, err error) {
	project = &Project{
		ID: id,
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
	if project.ID.Compare(ulid.ULID{}) == 0 {
		return ErrMissingID
	}

	project.Modified = time.Now()

	if err = Put(ctx, project); err != nil {
		return err
	}
	return nil
}

// DeleteProject deletes a project with a given id.
func DeleteProject(ctx context.Context, id ulid.ULID) (err error) {
	project := &Project{
		ID: id,
	}

	if err = Delete(ctx, project); err != nil {
		return err
	}
	return nil
}

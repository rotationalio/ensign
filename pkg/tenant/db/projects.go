package db

import (
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/net/context"
)

const ProjectNamespace = "projects"

type Project struct {
	ID       ulid.ULID
	Name     string
	Created  time.Time
	Modified time.Time
}

var (
	UlidNil ulid.ULID
)

var _ Model = &Project{}

func (p *Project) Key() ([]byte, error) {
	return p.ID.MarshalBinary()
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

func CreateProject(ctx context.Context, project *Project) (err error) {
	if project.ID == UlidNil {
		project.ID = ulid.Make()
	}

	project.Created = time.Now()
	project.Modified = project.Created

	if err = Put(ctx, project); err != nil {
		return err
	}
	return nil
}

// ListProjects retrieves all projects assigned to a tenant.
func ListProjects(ctx context.Context, prefix []byte, namespace string) (values [][]byte, err error) {

	if values, err = List(ctx, prefix, namespace); err != nil {
		return nil, err
	}
	return values, err
}

func RetrieveProject(ctx context.Context, id ulid.ULID) (project *Project, err error) {
	project = &Project{
		ID: id,
	}

	if err = Get(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func UpdateProject(ctx context.Context, project *Project) (err error) {
	if project.ID == UlidNil {
		return ErrMissingID
	}

	project.Modified = time.Now()

	if err = Put(ctx, project); err != nil {
		return err
	}
	return nil
}

func DeleteProject(ctx context.Context, id ulid.ULID) (err error) {
	project := &Project{
		ID: id,
	}

	if err = Delete(ctx, project); err != nil {
		return err
	}
	return nil
}

// TODO: Add

package db

import (
	"encoding/json"

	"github.com/oklog/ulid/v2"
	"golang.org/x/net/context"
)

const ProjectNamespace = "projects"

type Project struct {
	ID   ulid.ULID
	Name string
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
	return json.Marshal(p)
}

func (p *Project) UnmarshalValue(data []byte) error {
	return json.Unmarshal(data, p)
}

func CreateProject(ctx context.Context, project *Project) (err error) {
	if project.ID == UlidNil {
		project.ID = ulid.Make()
	}

	if err = Put(ctx, project); err != nil {
		return err
	}
	return nil
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

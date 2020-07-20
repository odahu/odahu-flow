package stubclients

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/filter"
)

type TIStubClient struct {
	db map[string]training.ToolchainIntegration
}

func NewTIStubClient() TIStubClient{
	return TIStubClient{
		db: make(map[string]training.ToolchainIntegration),
	}
}

func (t TIStubClient) GetToolchainIntegration(name string) (*training.ToolchainIntegration, error) {
	entity, ok := t.db[name]
	if !ok {
		return nil, errors.NotFoundError{Entity: name}
	}
	return &training.ToolchainIntegration{
		ID:     name,
		Spec:   entity.Spec,
		Status: entity.Status,
	}, nil
}

func (t TIStubClient) GetToolchainIntegrationList(
	options ...filter.ListOption) ([]training.ToolchainIntegration, error) {
	panic("implement me")
}

func (t TIStubClient) DeleteToolchainIntegration(name string) error {
	if _, err := t.GetToolchainIntegration(name); err != nil {
		return err
	}
	delete(t.db, name)
	return nil
}

func (t TIStubClient) UpdateToolchainIntegration(md *training.ToolchainIntegration) error {
	panic("implement me")
}

func (t TIStubClient) CreateToolchainIntegration(md *training.ToolchainIntegration) error {
	t.db[md.ID] = training.ToolchainIntegration{
		ID:    md.ID,
		Spec:  md.Spec,
		Status: md.Status,
	}
	return nil
}


type PIStubClient struct {
	db map[string]packaging.PackagingIntegration
}

func NewPIStubClient() PIStubClient{
	return PIStubClient{
		db: make(map[string]packaging.PackagingIntegration),
	}
}

func (p PIStubClient) GetPackagingIntegration(name string) (*packaging.PackagingIntegration, error) {
	entity, ok := p.db[name]
	if !ok {
		return nil, errors.NotFoundError{Entity: name}
	}
	return &packaging.PackagingIntegration{
		ID:     name,
		Spec:   entity.Spec,
		Status: entity.Status,
	}, nil
}

func (p PIStubClient) GetPackagingIntegrationList(
	options ...filter.ListOption) ([]packaging.PackagingIntegration, error) {

	panic("implement me")
}

func (p PIStubClient) DeletePackagingIntegration(name string) error {
	if _, err := p.GetPackagingIntegration(name); err != nil {
		return err
	}
	delete(p.db, name)
	return nil
}

func (p PIStubClient) UpdatePackagingIntegration(md *packaging.PackagingIntegration) error {
	panic("implement me")
}

func (p PIStubClient) CreatePackagingIntegration(md *packaging.PackagingIntegration) error {
	p.db[md.ID] = packaging.PackagingIntegration{
		ID:    md.ID,
		Spec:  md.Spec,
		Status: md.Status,
	}
	return nil
}

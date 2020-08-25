package stubclients

import (
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
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

func (t TIStubClient) GetModelTraining(id string) (*training.ModelTraining, error) {
	panic("implement me")
}

func (t TIStubClient) SaveModelTrainingResult(id string, result *v1alpha1.TrainingResult) error {
	panic("implement me")
}

func (t TIStubClient) GetToolchainIntegration(name string) (*training.ToolchainIntegration, error) {
	entity, ok := t.db[name]
	if !ok {
		return nil, errors.NotFoundError{Entity: name}
	}
	return &entity, nil
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

func (p PIStubClient) GetModelPackaging(id string) (*packaging.ModelPackaging, error) {
	panic("implement me")
}

func (p PIStubClient) SaveModelPackagingResult(id string, result []v1alpha1.ModelPackagingResult) error {
	panic("implement me")
}

func (p PIStubClient) GetPackagingIntegration(name string) (*packaging.PackagingIntegration, error) {
	entity, ok := p.db[name]
	if !ok {
		return nil, errors.NotFoundError{Entity: name}
	}
	return &entity, nil
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

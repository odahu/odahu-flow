package deploymentclient

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
)

type Client interface {
	GetModelDeployment(name string) (*deployment.ModelDeployment, error)
	GetModelDeploymentList(options ...filter.ListOption) ([]deployment.ModelDeployment, error)
	DeleteModelDeployment(name string) error
	UpdateModelDeployment(md *deployment.ModelDeployment) error
	CreateModelDeployment(md *deployment.ModelDeployment) error
	GetModelRoute(name string) (*deployment.ModelRoute, error)
	GetModelRouteList(options ...filter.ListOption) ([]deployment.ModelRoute, error)
	DeleteModelRoute(name string) error
	UpdateModelRoute(md *deployment.ModelRoute) error
	CreateModelRoute(md *deployment.ModelRoute) error
}

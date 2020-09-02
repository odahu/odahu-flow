package packaging

import (
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
)

type Client interface {
	GetModelPackaging(id string) (*packaging.ModelPackaging, error)
	SaveModelPackagingResult(id string, result []odahuflowv1alpha1.ModelPackagingResult) error
	GetPackagingIntegration(id string) (*packaging.PackagingIntegration, error)
}

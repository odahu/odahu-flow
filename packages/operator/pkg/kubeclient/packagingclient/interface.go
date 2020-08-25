package packagingclient

import (
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
)

type Client interface {
	GetModelPackaging(id string) (*packaging.ModelPackaging, error)
	GetModelPackagingList(options ...filter.ListOption) ([]packaging.ModelPackaging, error)
	DeleteModelPackaging(id string) error
	UpdateModelPackaging(md *packaging.ModelPackaging) error
	CreateModelPackaging(md *packaging.ModelPackaging) error
	GetModelPackagingLogs(id string, writer utils.Writer, follow bool) error
	SaveModelPackagingResult(id string, result []odahuflowv1alpha1.ModelPackagingResult) error
	GetModelPackagingResult(id string) ([]odahuflowv1alpha1.ModelPackagingResult, error)
}

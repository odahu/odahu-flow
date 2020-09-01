package trainingclient

import (
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
)

type Client interface {
	GetModelTraining(id string) (*training.ModelTraining, error)
	GetModelTrainingList(options ...filter.ListOption) ([]training.ModelTraining, error)
	DeleteModelTraining(id string) error
	UpdateModelTraining(md *training.ModelTraining) error
	CreateModelTraining(md *training.ModelTraining) error
	GetModelTrainingLogs(id string, writer utils.Writer, follow bool) error
	SaveModelTrainingResult(id string, result *odahuflowv1alpha1.TrainingResult) error
	GetModelTrainingResult(id string) (*odahuflowv1alpha1.TrainingResult, error)
}

package training

import (
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
)

type Client interface {
	GetModelTraining(id string) (*training.ModelTraining, error)
	SaveModelTrainingResult(id string, result *v1alpha1.TrainingResult) error
	GetTrainingIntegration(name string) (*training.TrainingIntegration, error)
}

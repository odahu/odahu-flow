package controllers_test

import (
	"github.com/google/uuid"
	kubetypes "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/controllers/mocks"
	apitypes "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

const (
	TestImage = "TestImage"
)

var tests = []struct{
	jobSpec kubetypes.BatchInferenceJobSpec
	serviceSpec apitypes.InferenceServiceSpec
	taskSpec tektonv1beta1.TaskSpec
}{
	{
		jobSpec:     kubetypes.BatchInferenceJobSpec{
			BatchInferenceServiceID: "test-batch-job",
			Command:              []string{"predict"},
			InputPath:               "/data/input",
			OutputPath:              "data/output",
			NodeSelector:            nil,
		},
		serviceSpec: apitypes.InferenceServiceSpec{
			Image:            "",
			InputConnection:  "",
			OutputConnection: "",
			ModelConnection:  "",
			ModelPath:        "",
			Triggers:         apitypes.InferenceServiceTriggers{},
		},
		taskSpec:    tektonv1beta1.TaskSpec{},
	},
}

func TestBatchJobToTaskRun(t *testing.T) {
	uid := uuid.New().String()[:5]
	name := "testJob" + uid
	for _, tt := range tests {
		_ = kubetypes.BatchInferenceJob{
			ObjectMeta: metav1.ObjectMeta{
				Name:                       name,
				Namespace:                  testNamespace,
			},
			Spec:      tt.jobSpec,
		}
		service := apitypes.InferenceService{
			ID:     "test-batch-job",
			Spec:   apitypes.InferenceServiceSpec{
				Image:            "",
				InputConnection:  "inputConn",
				OutputConnection: "outputConn",
				ModelConnection:  "modelConn",
				ModelPath:        "/model/1",
				Triggers:         apitypes.InferenceServiceTriggers{},
			},
		}
		connAPI := mocks.ConnectionAPI{}
		batchServAPI := mocks.BatchInferenceServiceAPI{}
		batchServAPI.On("Get", "test-batch-job").Return(service, nil)
		connAPI.On("GetConnection", "")

	}

}

//
//    Copyright 2019 EPAM Systems
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package training_test

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	train_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	conn_k8s_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/kubernetes"
	mt_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	mt_post_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"go.uber.org/multierr"
	"testing"
)

// TODO: add multiple test error

const (
	gpuResourceName     = "nvidia"
	mtvOutputConnection = "training-validation-output-connection"
	invalidK8sLabel     = "invalid label!"
)

var (
	defaultTrainingResource = config.NewDefaultModelTrainingConfig().DefaultResources
	cpuNodeSelector         = map[string]string{"mode": "training"}
	gpuNodeSelector         = map[string]string{"mode": "gpuTraining"}
	validTraining           = training.ModelTraining{
		ID: "model-id",
		Spec: v1alpha1.ModelTrainingSpec{
			Model: v1alpha1.ModelIdentity{
				Name:    "model-name",
				Version: "1",
			},
			OutputConnection: mtvOutputConnection,
			Toolchain:        testToolchainIntegrationID,
			AlgorithmSource: v1alpha1.AlgorithmSource{
				VCS: v1alpha1.VCS{
					ConnName:  testMtVCSID,
					Reference: testVcsReference,
				},
			},
			NodeSelector: cpuNodeSelector,
		},
	}
)

type ModelTrainingValidationSuite struct {
	suite.Suite
	g              *GomegaWithT
	validator      *train_route.MtValidator
	mtRepository   mt_repository.ToolchainRepository
	connRepository conn_repository.Repository
}

func (s *ModelTrainingValidationSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func (s *ModelTrainingValidationSuite) SetupSuite() {

	s.mtRepository = mt_post_repository.ToolchainRepo{DB: db}
	s.connRepository = conn_k8s_repository.NewRepository(testNamespace, kubeClient)

	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.NodePools = append(trainingConfig.NodePools, config.NodePool{
		NodeSelector: cpuNodeSelector})
	trainingConfig.GPUNodePools = append(trainingConfig.GPUNodePools, config.NodePool{
		NodeSelector: gpuNodeSelector,
	})

	s.validator = train_route.NewMtValidator(
		s.mtRepository,
		s.connRepository,
		trainingConfig,
		gpuResourceName,
	)

	// Create the connection that will be used as the vcs param for a training.
	if err := s.connRepository.CreateConnection(&connection.Connection{
		ID: testMtVCSID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			Reference: testVcsReference,
		},
	}); err != nil {
		// If we get a panic that we have a test configuration problem
		panic(err)
	}

	// Create the connection that will be used as the outputConnection param for a training.
	if err := s.connRepository.CreateConnection(&connection.Connection{
		ID: mtvOutputConnection,
		Spec: v1alpha1.ConnectionSpec{
			Type: connection.GcsType,
		},
	}); err != nil {
		// If we get a panic that we have a test configuration problem
		panic(err)
	}

	// Create the toolchain integration that will be used for a training.
	if err := s.mtRepository.CreateToolchainIntegration(&training.ToolchainIntegration{
		ID: testToolchainIntegrationID,
		Spec: v1alpha1.ToolchainIntegrationSpec{
			DefaultImage: testToolchainMtImage,
		},
	}); err != nil {
		// If we get a panic that we have a test configuration problem
		panic(err)
	}
}

func (s *ModelTrainingValidationSuite) TearDownSuite() {
	if err := s.mtRepository.DeleteToolchainIntegration(testToolchainIntegrationID); err != nil {
		panic(err)
	}

	if err := s.connRepository.DeleteConnection(testMtVCSID); err != nil {
		panic(err)
	}
}

func TestModelTrainingValidationSuite(t *testing.T) {
	suite.Run(t, new(ModelTrainingValidationSuite))
}

func (s *ModelTrainingValidationSuite) TestMtDefaultResource() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{},
	}

	_ = s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(mt.Spec.Resources).ShouldNot(BeNil())
	s.g.Expect(*mt.Spec.Resources).Should(Equal(defaultTrainingResource))
}

func (s *ModelTrainingValidationSuite) TestMtVcsReference() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			AlgorithmSource: v1alpha1.AlgorithmSource{
				VCS: v1alpha1.VCS{
					ConnName: testMtVCSID,
				},
			},
		},
	}

	_ = s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(mt.Spec.AlgorithmSource.VCS.Reference).To(Equal(testVcsReference))
}

func (s *ModelTrainingValidationSuite) TestMtMtImage() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			Toolchain: testToolchainIntegrationID,
		},
	}

	_ = s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(mt.Spec.Image).To(Equal(testToolchainMtImage))
}

func (s *ModelTrainingValidationSuite) TestMtMtImageExplicitly() {
	image := "image-test"
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			Toolchain: testToolchainIntegrationID,
			Image:     image,
		},
	}

	_ = s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(mt.Spec.Image).To(Equal(image))
}

func (s *ModelTrainingValidationSuite) TestMtExplicitMTReference() {
	mtExplicitReference := "test-ref"
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			AlgorithmSource: v1alpha1.AlgorithmSource{
				VCS: v1alpha1.VCS{
					ConnName:  testMtVCSID,
					Reference: mtExplicitReference,
				},
			},
		},
	}

	_ = s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(mt.Spec.AlgorithmSource.VCS.Reference).To(Equal(mtExplicitReference))
}

func (s *ModelTrainingValidationSuite) TestMtNotExplicitMTReference() {
	conn := &connection.Connection{
		ID: "vcs",
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			Reference: "",
		},
	}

	err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).Should(BeNil())
	defer s.connRepository.DeleteConnection(conn.ID)

	mt := validTraining
	mt.Spec.AlgorithmSource.VCS.Reference = ""

	err = s.validator.ValidatesAndSetDefaults(&mt)
	s.Assertions.NoError(err)
}

func (s *ModelTrainingValidationSuite) TestMtEmptyAlgorithmSourceName() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).ShouldNot(BeNil())
	s.g.Expect(err.Error()).To(ContainSubstring(train_route.EmptyAlgorithmSourceNameMessageError))
}

func (s *ModelTrainingValidationSuite) TestMtWrongVcsConnectionType() {
	conn := &connection.Connection{
		ID: "wrong-type",
		Spec: v1alpha1.ConnectionSpec{
			Type: connection.S3Type,
		},
	}

	err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).Should(BeNil())
	defer s.connRepository.DeleteConnection(conn.ID)

	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			AlgorithmSource: v1alpha1.AlgorithmSource{
				VCS: v1alpha1.VCS{
					ConnName: conn.ID,
				},
			},
		},
	}

	err = s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).ShouldNot(BeNil())
	s.g.Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(train_route.WrongVcsTypeErrorMessage, conn.Spec.Type)))
}

func (s *ModelTrainingValidationSuite) TestMtWrongObjectStorageConnectionType() {
	conn := &connection.Connection{
		ID: "wrong-type",
		Spec: v1alpha1.ConnectionSpec{
			Type: connection.GITType,
		},
	}

	err := s.connRepository.CreateConnection(conn)
	s.g.Expect(err).Should(BeNil())
	defer s.connRepository.DeleteConnection(conn.ID)

	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			AlgorithmSource: v1alpha1.AlgorithmSource{
				ObjectStorage: v1alpha1.ObjectStorage{
					ConnName: conn.ID,
				},
			},
		},
	}

	err = s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).ShouldNot(BeNil())
	s.g.Expect(err.Error()).To(ContainSubstring("data binding has wrong data type"))
}

func (s *ModelTrainingValidationSuite) TestMtToolchainType() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			Toolchain: "not-exists",
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(
		"entity \"not-exists\" is not found"))
}

func (s *ModelTrainingValidationSuite) TestMtVcsNotExists() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			AlgorithmSource: v1alpha1.AlgorithmSource{
				VCS: v1alpha1.VCS{
					ConnName: "not-exists",
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(
		"entity \"not-exists\" is not found"))
}

func (s *ModelTrainingValidationSuite) TestMtObjectStorageNotExists() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			AlgorithmSource: v1alpha1.AlgorithmSource{
				ObjectStorage: v1alpha1.ObjectStorage{
					ConnName: "not-exists",
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(
		"entity \"not-exists\" is not found"))
}

func (s *ModelTrainingValidationSuite) TestMtAlgorithmSourceName() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			AlgorithmSource: v1alpha1.AlgorithmSource{
				VCS: v1alpha1.VCS{
					ConnName: "",
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(train_route.EmptyAlgorithmSourceNameMessageError))
}

func (s *ModelTrainingValidationSuite) TestMtMultipleAlgorithmSourceName() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			AlgorithmSource: v1alpha1.AlgorithmSource{
				VCS: v1alpha1.VCS{
					ConnName: "not-empty",
				},
				ObjectStorage: v1alpha1.ObjectStorage{
					ConnName: "not-empty",
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(train_route.MultipleAlgorithmSourceMessageError))
}

func (s *ModelTrainingValidationSuite) TestMtToolchainEmptyName() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			Toolchain: "",
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(train_route.ToolchainEmptyErrorMessage))
}

func (s *ModelTrainingValidationSuite) TestMtToolchain_invalid() {
	mt := validTraining
	mt.Spec.Toolchain = invalidK8sLabel

	err := s.validator.ValidatesAndSetDefaults(&mt)
	s.Assertions.Error(err)
	s.Assertions.Contains(err.Error(), mt.Spec.Toolchain)
}

func (s *ModelTrainingValidationSuite) TestMtEmptyModelName() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(train_route.EmptyModelNameErrorMessage))
}

func (s *ModelTrainingValidationSuite) TestMtName_invalid() {
	mt := validTraining
	mt.Spec.Model.Name = invalidK8sLabel

	err := s.validator.ValidatesAndSetDefaults(&mt)
	s.Assertions.Error(err)
	s.Assertions.Contains(err.Error(), mt.Spec.Model.Name)
}

func (s *ModelTrainingValidationSuite) TestMtEmptyModelVersion() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(train_route.EmptyModelVersionErrorMessage))
}

func (s *ModelTrainingValidationSuite) TestMtVersion_invalid() {
	mt := validTraining
	mt.Spec.Model.Version = invalidK8sLabel

	err := s.validator.ValidatesAndSetDefaults(&mt)
	s.Assertions.Error(err)
	s.Assertions.Contains(err.Error(), mt.Spec.Model.Version)
}

func (s *ModelTrainingValidationSuite) TestMtGenerationOutputArtifactName() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(mt.Spec.Model.ArtifactNameTemplate).ShouldNot(BeEmpty())
}

func (s *ModelTrainingValidationSuite) TestMtWrongDataType() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			Data: []v1alpha1.DataBindingDir{
				{
					Connection: testMtVCSID,
					LocalPath:  testMtDataPath,
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(
		"odahu-flow-test data binding has wrong data type. " +
			"Currently supported the following types of connections for data bindings:"))
}

func (s *ModelTrainingValidationSuite) TestMtEmptyDataName() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			Data: []v1alpha1.DataBindingDir{
				{
					Connection: "",
					LocalPath:  testMtDataPath,
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(fmt.Sprintf(
		train_route.EmptyDataBindingNameErrorMessage, 0)))
}

func (s *ModelTrainingValidationSuite) TestMtEmptyDataPath() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			Data: []v1alpha1.DataBindingDir{
				{
					Connection: testMtVCSID,
					LocalPath:  "",
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(fmt.Sprintf(
		train_route.EmptyDataBindingPathErrorMessage, 0)))
}

func (s *ModelTrainingValidationSuite) TestMtNotFoundData() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			Data: []v1alpha1.DataBindingDir{
				{
					Connection: "not-present",
					LocalPath:  testMtDataPath,
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(
		"entity \"not-present\" is not found"))
}

func (s *ModelTrainingValidationSuite) TestMtResourcesValidation() {
	wrongResourceValue := "wrong res"
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{
			Resources: &v1alpha1.ResourceRequirements{
				Limits: &v1alpha1.ResourceList{
					Memory: &wrongResourceValue,
					GPU:    &wrongResourceValue,
					CPU:    &wrongResourceValue,
				},
				Requests: &v1alpha1.ResourceList{
					Memory: &wrongResourceValue,
					GPU:    &wrongResourceValue,
					CPU:    &wrongResourceValue,
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).Should(HaveOccurred())

	errorMessage := err.Error()
	s.g.Expect(errorMessage).Should(ContainSubstring(
		"validation of memory request is failed: quantities must match the regular expression"))
	s.g.Expect(errorMessage).Should(ContainSubstring(
		"validation of cpu request is failed: quantities must match the regular expression"))
	s.g.Expect(errorMessage).Should(ContainSubstring(
		"validation of memory limit is failed: quantities must match the regular expression"))
	s.g.Expect(errorMessage).Should(ContainSubstring(
		"validation of cpu limit is failed: quantities must match the regular expression"))
	s.g.Expect(errorMessage).Should(ContainSubstring(
		"validation of gpu limit is failed: quantities must match the regular expression"))
}

// If default output connection is not set in config and user doesn't provide it, validation fails
func (s *ModelTrainingValidationSuite) TestOutputConnection_NoDefault_NoParam() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{},
	}
	testConfig := config.ModelTrainingConfig{
		DefaultResources:   defaultTrainingResource,
		OutputConnectionID: "",
	}
	err := train_route.NewMtValidator(
		s.mtRepository,
		s.connRepository,
		testConfig,
		gpuResourceName,
	).ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(validation.EmptyValueStringError, "OutputConnection")))
}

// If default output connection is set in config and user doesn't pass output connection, use default
func (s *ModelTrainingValidationSuite) TestOutputConnection_Default_NoParam() {
	mt := &training.ModelTraining{}

	testConfig := config.ModelTrainingConfig{OutputConnectionID: testMtOutConnDefault}
	_ = train_route.NewMtValidator(
		s.mtRepository,
		s.connRepository,
		testConfig,
		gpuResourceName,
	).ValidatesAndSetDefaults(mt)
	s.g.Expect(mt.Spec.OutputConnection).Should(Equal(testMtOutConnDefault))
}

// If output connection is set in both config and user request, use one from user
func (s *ModelTrainingValidationSuite) TestOutputConnection_Both_Default_Param() {
	mt := &training.ModelTraining{}
	mt.Spec.OutputConnection = mtvOutputConnection
	_ = s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(mt.Spec.OutputConnection).Should(Equal(mtvOutputConnection))
}

// If connection repository doesn't contain connection with passed ID validation must raise NotFoundError
func (s *ModelTrainingValidationSuite) TestOutputConnection_ConnectionNotFound() {
	mt := &training.ModelTraining{
		Spec: v1alpha1.ModelTrainingSpec{OutputConnection: testMpOutConnNotFound},
	}
	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("entity %q is not found", testMpOutConnNotFound))

}

func (s *ModelTrainingValidationSuite) TestValidateID() {
	mt := &training.ModelTraining{
		ID: "not-VALID-id-",
	}

	err := s.validator.ValidatesAndSetDefaults(mt)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(validation.ErrIDValidation.Error()))
}

// Tests that nil node selector is considered valid
func (s *ModelTrainingValidationSuite) TestValidateNodeSelector_nil() {
	mt := validTraining
	mt.Spec.NodeSelector = nil
	err := s.validator.ValidatesAndSetDefaults(&mt)
	s.Assertions.Nil(err)
}

// training object has valid node selector that exists in config
func (s *ModelTrainingValidationSuite) TestValidateNodeSelector_Valid() {
	mt := validTraining
	err := s.validator.ValidatesAndSetDefaults(&mt)
	s.Assertions.Nil(err)
}

// training object has invalid node selector that does not exist in config
// Expect validator to return exactly one error
func (s *ModelTrainingValidationSuite) TestValidateNodeSelector_Invalid() {
	mt := validTraining
	mt.Spec.NodeSelector = map[string]string{"mode": "invalid"}
	err := s.validator.ValidatesAndSetDefaults(&mt)
	s.Assertions.NotNil(err)
	s.Assertions.Len(multierr.Errors(err), 1)
}

// GPU training has valid node selector that exists in config
func (s *ModelTrainingValidationSuite) TestValidateNodeSelector_GPU_Valid() {
	mt := validTraining
	gpuRequest := "1"
	mt.Spec.Resources = &v1alpha1.ResourceRequirements{Requests: &v1alpha1.ResourceList{GPU: &gpuRequest}}
	mt.Spec.NodeSelector = gpuNodeSelector

	err := s.validator.ValidatesAndSetDefaults(&mt)
	s.Assertions.Nil(err)
}

// GPU training has invalid node selector that does not exist in config
// Expect validator to return exactly one error
func (s *ModelTrainingValidationSuite) TestValidateNodeSelector_GPU_Invalid() {
	mt := validTraining
	gpuRequest := "1"
	mt.Spec.Resources = &v1alpha1.ResourceRequirements{Requests: &v1alpha1.ResourceList{GPU: &gpuRequest}}
	mt.Spec.NodeSelector = map[string]string{"mode": "gpu-invalid"}

	err := s.validator.ValidatesAndSetDefaults(&mt)
	s.Assertions.NotNil(err)
	s.Assertions.Len(multierr.Errors(err), 1)
}

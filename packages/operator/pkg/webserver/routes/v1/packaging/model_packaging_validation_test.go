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

package packaging_test

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	"path/filepath"
	"testing"

	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/odahuflow/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	conn_k8s_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	odahuflow_apis "github.com/odahu/odahu-flow/packages/operator/pkg/apis"
	mp_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	mp_k8s_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/kubernetes"
	pack_route "github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/packaging"
)

var (
	piIDMpValid                      = "pi-id"
	piEntrypointMpValid              = "/usr/bin/test"
	piImageMpValid                   = "test:image"
	piArtifactNameMpValid            = "some-artifact-name.zip"
	connDockerTypeMpValid            = "docker-conn"
	connS3TypeMpValid                = "s3-conn"
	defaultTargetArgument1Connection = "default-conn-id"
	piArgumentsMpValid               = packaging.JsonSchema{
		Properties: []packaging.Property{
			{
				Name: "argument-1",
				Parameters: []packaging.Parameter{
					{
						Name:  "minimum",
						Value: float64(5),
					},
					{
						Name:  "type",
						Value: "number",
					},
				},
			},
			{
				Name: "argument-2",
				Parameters: []packaging.Parameter{
					{
						Name:  "type",
						Value: "string",
					},
				},
			},
		},
		Required: []string{"argument-1"},
	}
	piTargetsMpValid = []v1alpha1.TargetSchema{
		{
			Name: "target-1",
			ConnectionTypes: []string{
				string(connection.S3Type),
				string(connection.GcsType),
				string(connection.AzureBlobType),
			},
			Default:  defaultTargetArgument1Connection,
			Required: false,
		},
		{
			Name: "target-2",
			ConnectionTypes: []string{
				string(connection.DockerType),
			},
			Required: true,
		},
	}
)

type ModelPackagingValidationSuite struct {
	suite.Suite
	g              *GomegaWithT
	mpRepository   mp_repository.Repository
	connRepository conn_repository.Repository
	k8sEnvironment *envtest.Environment
	validator      *pack_route.MpValidator
}

func (s *ModelPackagingValidationSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func (s *ModelPackagingValidationSuite) SetupSuite() {
	var cfg *rest.Config

	s.k8sEnvironment = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "..", "..", "config", "crds")},
	}

	err := odahuflow_apis.AddToScheme(scheme.Scheme)
	if err != nil {
		s.T().Fatalf("Cannot setup the odahuflow schema: %v", err)
	}

	cfg, err = s.k8sEnvironment.Start()
	if err != nil {
		s.T().Fatalf("Cannot setup the test k8s api: %v", err)
	}

	mgr, err := manager.New(cfg, manager.Options{NewClient: utils.NewClient})
	if err != nil {
		panic(err)
	}

	s.mpRepository = mp_k8s_repository.NewRepository(testNamespace, testNamespace, mgr.GetClient(), nil)
	s.connRepository = conn_k8s_repository.NewRepository(testNamespace, mgr.GetClient())
	s.validator = pack_route.NewMpValidator(
		s.mpRepository, s.connRepository, "conn-id", config.NvidiaResourceName,
	)

	err = s.mpRepository.CreatePackagingIntegration(&packaging.PackagingIntegration{
		ID: piIDMpValid,
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   piEntrypointMpValid,
			DefaultImage: piImageMpValid,
			Schema: packaging.Schema{
				Targets:   piTargetsMpValid,
				Arguments: piArgumentsMpValid,
			},
		},
	})
	if err != nil {
		s.T().Fatal(err)
	}

	err = s.connRepository.CreateConnection(&connection.Connection{
		ID: connDockerTypeMpValid,
		Spec: v1alpha1.ConnectionSpec{
			Type: connection.DockerType,
		},
		Status: v1alpha1.ConnectionStatus{},
	})
	if err != nil {
		s.T().Fatal(err)
	}

	err = s.connRepository.CreateConnection(&connection.Connection{
		ID: connS3TypeMpValid,
		Spec: v1alpha1.ConnectionSpec{
			Type: connection.S3Type,
		},
		Status: v1alpha1.ConnectionStatus{},
	})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ModelPackagingValidationSuite) TearDownSuite() {
	if err := s.connRepository.DeleteConnection(connS3TypeMpValid); err != nil {
		s.T().Fatal(err)
	}

	if err := s.connRepository.DeleteConnection(connDockerTypeMpValid); err != nil {
		s.T().Fatal(err)
	}

	if err := s.mpRepository.DeletePackagingIntegration(piIDMpValid); err != nil {
		s.T().Fatal(err)
	}

	if err := s.k8sEnvironment.Stop(); err != nil {
		s.T().Fatal("Cannot stop the test k8s api")
	}
}

func TestModelPackagingValidationSuite(t *testing.T) {
	suite.Run(t, new(ModelPackagingValidationSuite))
}

func (s *ModelPackagingValidationSuite) TestMpIDGeneration() {
	ti := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{},
	}

	_ = s.validator.ValidateAndSetDefaults(ti)
	s.g.Expect(ti.ID).ShouldNot(BeEmpty())
}

func (s *ModelPackagingValidationSuite) TestMpIDExplicitly() {
	id := "some-id"
	mp := &packaging.ModelPackaging{
		ID:   id,
		Spec: packaging.ModelPackagingSpec{},
	}

	_ = s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(mp.ID).Should(Equal(id))
}

func (s *ModelPackagingValidationSuite) TestMpImage() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
		},
	}

	_ = s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(mp.Spec.Image).ShouldNot(BeEmpty())
	s.g.Expect(mp.Spec.Image).Should(Equal(piImageMpValid))
}

func (s *ModelPackagingValidationSuite) TestMpImageExplicitly() {
	image := "some-image"
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Image:           image,
		},
	}

	_ = s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(mp.Spec.Image).Should(Equal(image))
}

func (s *ModelPackagingValidationSuite) TestMpImageNotFound() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: "not found",
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(
		"packagingintegrations.odahuflow.odahu.org \"not found\" not found"))
}

func (s *ModelPackagingValidationSuite) TestMpArtifactName() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			ArtifactName: piArtifactNameMpValid,
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).ShouldNot(ContainSubstring(pack_route.TrainingIDOrArtifactNameErrorMessage))
}

func (s *ModelPackagingValidationSuite) TestMpArtifactNameMissed() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(pack_route.TrainingIDOrArtifactNameErrorMessage))
}

func (s *ModelPackagingValidationSuite) TestMpIntegrationNameEmpty() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(pack_route.EmptyIntegrationNameErrorMessage))
}

func (s *ModelPackagingValidationSuite) TestMpIntegrationNotFound() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: "some-packaging-name",
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(
		"packagingintegrations.odahuflow.odahu.org \"some-packaging-name\" not found"))
}

func (s *ModelPackagingValidationSuite) TestMpNotValidArgumentsSchema() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Arguments: map[string]interface{}{
				"argument-1": 4,
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("argument-1: Must be greater than or equal to 5"))
}

func (s *ModelPackagingValidationSuite) TestMpAdditionalArguments() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Arguments: map[string]interface{}{
				"argument-3": "value",
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("Additional property argument-3 is not allowed"))
}

func (s *ModelPackagingValidationSuite) TestMpRequiredArguments() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Arguments: map[string]interface{}{
				"argument-2": "value",
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("argument-1 is required"))
}

func (s *ModelPackagingValidationSuite) TestMpRequiredTargets() {
	ti := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Targets: []v1alpha1.Target{
				{
					Name:           "target-1",
					ConnectionName: connS3TypeMpValid,
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(ti)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("[target-2] are required targets"))
}

func (s *ModelPackagingValidationSuite) TestMpDefaultTargets() {
	ti := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName:  piIDMpValid,
			ArtifactName:     "test",
			OutputConnection: connS3TypeMpValid,
			Arguments: map[string]interface{}{
				"argument-1": 5,
			},
			Targets: []v1alpha1.Target{
				{
					Name:           "target-2",
					ConnectionName: connDockerTypeMpValid,
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(ti)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(ti.Spec.Targets).Should(HaveLen(2))
	s.g.Expect(ti.Spec.Targets[0].Name).Should(Equal("target-2"))
	s.g.Expect(ti.Spec.Targets[0].ConnectionName).Should(Equal(connDockerTypeMpValid))
	s.g.Expect(ti.Spec.Targets[1].Name).Should(Equal("target-1"))
	s.g.Expect(ti.Spec.Targets[1].ConnectionName).Should(Equal(defaultTargetArgument1Connection))
}

func (s *ModelPackagingValidationSuite) TestMpNotFoundTargets() {
	targetNotFoundName := "target-not-found"
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Targets: []v1alpha1.Target{
				{
					Name:           targetNotFoundName,
					ConnectionName: connS3TypeMpValid,
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(fmt.Sprintf(
		pack_route.TargetNotFoundErrorMessage, targetNotFoundName, piIDMpValid,
	)))
}

func (s *ModelPackagingValidationSuite) TestMpTargetConnNotFound() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Targets: []v1alpha1.Target{
				{
					Name:           "target-1",
					ConnectionName: "conn-not-found",
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("not found"))
}

func (s *ModelPackagingValidationSuite) TestMpTargetConnWrongType() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
			Targets: []v1alpha1.Target{
				{
					Name:           "target-1",
					ConnectionName: connDockerTypeMpValid,
				},
			},
		},
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(fmt.Sprintf(
		pack_route.NotValidConnTypeErrorMessage, "target-1", connection.DockerType, piIDMpValid,
	)))
}

func (s *ModelPackagingValidationSuite) TestMpGenerateDefaultResources() {
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
		},
	}

	_ = s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(mp.Spec.Resources).ShouldNot(BeNil())
	s.g.Expect(mp.Spec.Resources).Should(Equal(pack_route.DefaultPackagingResources))
}

func (s *ModelPackagingValidationSuite) TestMpResourcesValidation() {
	wrongResourceValue := "wrong res"
	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{
			IntegrationName: piIDMpValid,
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

	err := s.validator.ValidateAndSetDefaults(mp)
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

func (s *ModelPackagingValidationSuite) TestOutputConnection() {
	testMpOutConnDefault := testOutConnDefault
	testMpOutConn := testOutConn
	testMpOutConnNotFound := testOutConnNotFound

	mp := &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{},
	}

	// If configuration output connection is not set then user must specify it as ModelTraining parameter
	err := pack_route.NewMpValidator(
		s.mpRepository,
		s.connRepository,
		"",
		config.NvidiaResourceName,
	).ValidateAndSetDefaults(mp)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(validation.EmptyValueStringError, "OutputConnection")))

	// If configuration output connection is set and user has not passed output connection as training
	// parameter then output connection value from configuration will be used as default
	_ = pack_route.NewMpValidator(
		s.mpRepository,
		s.connRepository,
		testMpOutConnDefault,
		config.NvidiaResourceName,
	).ValidateAndSetDefaults(mp)
	s.g.Expect(mp.Spec.OutputConnection).Should(Equal(testMpOutConnDefault))

	// If configuration output connection is set but user also has passed output connection as training
	// parameter then user value is used
	mp.Spec.OutputConnection = testMpOutConn
	_ = pack_route.NewMpValidator(
		s.mpRepository,
		s.connRepository,
		testMpOutConn,
		config.NvidiaResourceName,
	).ValidateAndSetDefaults(mp)
	s.g.Expect(mp.Spec.OutputConnection).Should(Equal(testMpOutConn))

	// If connection repository doesn't contain connection with passed ID validation must raise NotFoundError
	mp = &packaging.ModelPackaging{
		Spec: packaging.ModelPackagingSpec{OutputConnection: testMpOutConnNotFound},
	}
	err = pack_route.NewMpValidator(
		s.mpRepository,
		s.connRepository,
		testMpOutConn,
		config.NvidiaResourceName,
	).ValidateAndSetDefaults(mp)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("entity %q is not found", testMpOutConnNotFound))
}

func (s *ModelPackagingValidationSuite) TestValidateID() {
	mp := &packaging.ModelPackaging{
		ID: "not-VALID-id-",
	}

	err := s.validator.ValidateAndSetDefaults(mp)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(validation.ErrIDValidation.Error()))
}

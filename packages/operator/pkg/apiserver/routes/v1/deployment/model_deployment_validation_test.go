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

package deployment_test

import (
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	md_routes "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	"github.com/stretchr/testify/suite"
	"go.uber.org/multierr"
	"testing"

	. "github.com/onsi/gomega"
)

var (
	mdRoleName        = "test-tole"
	validNodeSelector = map[string]string{"mode": "valid"} // should be injected to config in SetupTest
	validDeployment   = deployment.ModelDeployment{
		ID: "valid-id",
		Spec: v1alpha1.ModelDeploymentSpec{
			Image:        "image",
			NodeSelector: validNodeSelector,
		},
	}
)

type ModelDeploymentValidationSuite struct {
	suite.Suite
	g                     *GomegaWithT
	defaultModelValidator *md_routes.ModelDeploymentValidator
}

func (s *ModelDeploymentValidationSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
	deployConfig := config.NewDefaultModelDeploymentConfig()
	deployConfig.NodePools = append(deployConfig.NodePools, config.NodePool{NodeSelector: validNodeSelector})
	s.defaultModelValidator = md_routes.NewModelDeploymentValidator(
		deployConfig,
		config.NvidiaResourceName,
	)
}

func TestModelDeploymentValidationSuite(t *testing.T) {
	suite.Run(t, new(ModelDeploymentValidationSuite))
}

func (s *ModelDeploymentValidationSuite) TestMDMinReplicasDefaultValue() {
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{},
	}
	_ = s.defaultModelValidator.ValidatesMDAndSetDefaults(md)

	s.g.Expect(*md.Spec.MinReplicas).To(Equal(md_routes.MdDefaultMinimumReplicas))
}

func (s *ModelDeploymentValidationSuite) TestMDMaxReplicasDefaultValue() {
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{},
	}
	_ = s.defaultModelValidator.ValidatesMDAndSetDefaults(md)

	s.g.Expect(*md.Spec.MaxReplicas).To(Equal(md_routes.MdDefaultMaximumReplicas))
}

func (s *ModelDeploymentValidationSuite) TestMDResourcesDefaultValues() {
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{},
	}
	_ = s.defaultModelValidator.ValidatesMDAndSetDefaults(md)

	s.g.Expect(*md.Spec.Resources).To(Equal(config.NewDefaultModelTrainingConfig().DefaultResources))
}

func (s *ModelDeploymentValidationSuite) TestMDReadinessProbeDefaultValue() {
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{},
	}
	_ = s.defaultModelValidator.ValidatesMDAndSetDefaults(md)

	s.g.Expect(*md.Spec.ReadinessProbeInitialDelay).To(Equal(md_routes.MdDefaultReadinessProbeInitialDelay))
}

func (s *ModelDeploymentValidationSuite) TestMDLivenessProbeDefaultValue() {
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{},
	}
	_ = s.defaultModelValidator.ValidatesMDAndSetDefaults(md)

	s.g.Expect(*md.Spec.LivenessProbeInitialDelay).To(Equal(md_routes.MdDefaultLivenessProbeInitialDelay))
}

func (s *ModelDeploymentValidationSuite) TestValidateNegativeMinReplicas() {
	minReplicas := int32(-1)
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{
			MinReplicas: &minReplicas,
		},
	}

	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(md_routes.NegativeMinReplicasErrorMessage))
}

func (s *ModelDeploymentValidationSuite) TestValidateNegativeMaxReplicas() {
	maxReplicas := int32(-1)
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{
			MaxReplicas: &maxReplicas,
		},
	}

	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(md_routes.NegativeMaxReplicasErrorMessage))
}

func (s *ModelDeploymentValidationSuite) TestValidateMaximumReplicas() {
	maxReplicas := int32(-1)
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{
			MaxReplicas: &maxReplicas,
		},
	}

	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(md_routes.NegativeMaxReplicasErrorMessage))
}

func (s *ModelDeploymentValidationSuite) TestValidateMinLessMaxReplicas() {
	minReplicas := int32(2)
	maxReplicas := int32(1)
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{
			MinReplicas: &minReplicas,
			MaxReplicas: &maxReplicas,
		},
	}

	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(md_routes.MaxMoreThanMinReplicasErrorMessage))
}

func (s *ModelDeploymentValidationSuite) TestValidateMinModelThanDefaultMax() {
	minReplicas := int32(3)
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{
			MinReplicas: &minReplicas,
		},
	}

	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).ToNot(ContainSubstring(md_routes.MaxMoreThanMinReplicasErrorMessage))
	s.g.Expect(*md.Spec.MinReplicas).To(Equal(minReplicas))
	s.g.Expect(*md.Spec.MaxReplicas).To(Equal(minReplicas))
}

func (s *ModelDeploymentValidationSuite) TestValidateImage() {
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{},
	}

	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(md_routes.EmptyImageErrorMessage))
}

func (s *ModelDeploymentValidationSuite) TestValidateReadinessProbe() {
	readinessProbe := int32(-1)
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{
			ReadinessProbeInitialDelay: &readinessProbe,
		},
	}

	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(md_routes.ReadinessProbeErrorMessage))
}

func (s *ModelDeploymentValidationSuite) TestValidateLivenessProbe() {
	livenessProbe := int32(-1)
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{
			LivenessProbeInitialDelay: &livenessProbe,
		},
	}

	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(md_routes.LivenessProbeErrorMessage))
}

func (s *ModelDeploymentValidationSuite) TestMdResourcesValidation() {
	wrongResourceValue := "wrong res"
	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{
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

	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
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

func (s *ModelDeploymentValidationSuite) TestValidateDefaultDockerPullConnectionName() {
	newDefaultDockerPullConnectionName := "default-docker-pull-conn"
	mdConfig := config.NewDefaultModelDeploymentConfig()
	mdConfig.DefaultDockerPullConnName = newDefaultDockerPullConnectionName

	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{},
	}

	_ = md_routes.NewModelDeploymentValidator(mdConfig, config.NvidiaResourceName).ValidatesMDAndSetDefaults(md)
	s.g.Expect(md.Spec.ImagePullConnectionID).ShouldNot(BeNil())
	s.g.Expect(*md.Spec.ImagePullConnectionID).Should(Equal(newDefaultDockerPullConnectionName))
}

func (s *ModelDeploymentValidationSuite) TestValidateDockerPullConnectionName() {
	dockerPullConnectionName := "default-docker-pull-conn"

	md := &deployment.ModelDeployment{
		Spec: v1alpha1.ModelDeploymentSpec{
			ImagePullConnectionID: &dockerPullConnectionName,
		},
	}

	_ = s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
	s.g.Expect(md.Spec.ImagePullConnectionID).ShouldNot(BeNil())
	s.g.Expect(*md.Spec.ImagePullConnectionID).Should(Equal(dockerPullConnectionName))
}

func (s *ModelDeploymentValidationSuite) TestValidateID() {
	md := &deployment.ModelDeployment{
		ID: "not-VALID-id-",
	}

	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(md)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(validation.ErrIDValidation.Error()))
}

// Tests that nil node selector is considered valid
func (s *ModelDeploymentValidationSuite) TestValidateNodeSelector_nil() {
	md := validDeployment
	md.Spec.NodeSelector = nil
	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(&md)
	s.Assertions.NoError(err)
}

// Deployment object has valid node selector that exists in config
func (s *ModelDeploymentValidationSuite) TestValidateNodeSelector_Valid() {
	md := validDeployment
	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(&md)
	s.Assertions.NoError(err)
}

// Deployment object has invalid node selector that does not exist in config
// Expect validator to return exactly one error
func (s *ModelDeploymentValidationSuite) TestValidateNodeSelector_Invalid() {
	mt := validDeployment
	mt.Spec.NodeSelector = map[string]string{"mode": "invalid"}
	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(&mt)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(UnknownNodeSelector, md.Spec.NodeSelector))
}

// Deployment has no Role Name; expect a model ID as suffix of default role
func (s *ModelDeploymentValidationSuite) TestValidateRoleName_Empty() {
	mt := validDeployment
	role := ""
	mt.Spec.RoleName = &role
	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(&mt)
	s.Assertions.NoError(err)
	s.Assertions.Equal(*mt.Spec.RoleName, md_routes.DefaultRolePrefix + mt.ID)
}

// Deployment has no Role Name; expect a model ID as suffix of default role
func (s *ModelDeploymentValidationSuite) TestValidateRoleName_Nil() {
	mt := validDeployment
	var role *string = nil
	mt.Spec.RoleName = role
	err := s.defaultModelValidator.ValidatesMDAndSetDefaults(&mt)
	s.Assertions.NoError(err)
	s.Assertions.Equal(*mt.Spec.RoleName, md_routes.DefaultRolePrefix + mt.ID)
}

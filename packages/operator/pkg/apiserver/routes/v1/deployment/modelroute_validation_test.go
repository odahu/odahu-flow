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
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	md_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	dep_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/deployment"
	"github.com/stretchr/testify/suite"
	repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	"testing"

	. "github.com/onsi/gomega"
)

type ModelRouteValidationSuite struct {
	suite.Suite
	g            *GomegaWithT
	validator    *dep_route.MrValidator
}

func (s *ModelRouteValidationSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func (s *ModelRouteValidationSuite) SetupSuite() {

	s.validator = dep_route.NewMrValidator(md_service.NewService(repo.DeploymentRepo{DB: db}, nil))
}

func TestModelPackagingValidationSuite(t *testing.T) {
	suite.Run(t, new(ModelRouteValidationSuite))
}

func (s *ModelRouteValidationSuite) TestDefaultValues() {
	mr := &deployment.ModelRoute{
		Spec: v1alpha1.ModelRouteSpec{
			ModelDeploymentTargets: []v1alpha1.ModelDeploymentTarget{
				{
					Name: mdID1,
				},
			},
		},
	}

	_ = s.validator.ValidatesAndSetDefaults(mr)

	s.g.Expect(*mr.Spec.ModelDeploymentTargets[0].Weight).To(Equal(dep_route.MaxWeight))
}

func (s *ModelRouteValidationSuite) TestEmptyURLPrefix() {
	mr := &deployment.ModelRoute{
		Spec: v1alpha1.ModelRouteSpec{},
	}

	err := s.validator.ValidatesAndSetDefaults(mr)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(dep_route.URLPrefixEmptyErrorMessage))
}

func (s *ModelRouteValidationSuite) TestNotExistsMirrorMD() {
	mirrorMD := "not-exists"
	mr := &deployment.ModelRoute{
		Spec: v1alpha1.ModelRouteSpec{
			Mirror: &mirrorMD,
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mr)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("entity \"not-exists\" is not found"))
}

func (s *ModelRouteValidationSuite) TestMissingMDTargets() {
	mr := &deployment.ModelRoute{
		Spec: v1alpha1.ModelRouteSpec{},
	}

	err := s.validator.ValidatesAndSetDefaults(mr)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(dep_route.EmptyTargetErrorMessage))
}

func (s *ModelRouteValidationSuite) TestOneTargetWrongWeight() {
	weight := int32(77)
	mr := &deployment.ModelRoute{
		Spec: v1alpha1.ModelRouteSpec{
			ModelDeploymentTargets: []v1alpha1.ModelDeploymentTarget{
				{
					Name:   mdID1,
					Weight: &weight,
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mr)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(dep_route.OneTargetErrorMessage))
}

func (s *ModelRouteValidationSuite) TestOneTargetNotExist() {
	mr := &deployment.ModelRoute{
		Spec: v1alpha1.ModelRouteSpec{
			ModelDeploymentTargets: []v1alpha1.ModelDeploymentTarget{
				{
					Name: "not-exists",
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mr)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("entity \"not-exists\" is not found"))
}

func (s *ModelRouteValidationSuite) TestMultipleTargetsWrongWeight() {
	weight1 := int32(11)
	weight2 := int32(22)
	mr := &deployment.ModelRoute{
		Spec: v1alpha1.ModelRouteSpec{
			ModelDeploymentTargets: []v1alpha1.ModelDeploymentTarget{
				{
					Name:   mdID1,
					Weight: &weight1,
				},
				{
					Name:   mdID2,
					Weight: &weight2,
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mr)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(dep_route.TotalWeightErrorMessage))
}

func (s *ModelRouteValidationSuite) TestMultipleTargetsMissingWeight() {
	weight2 := int32(22)
	mr := &deployment.ModelRoute{
		Spec: v1alpha1.ModelRouteSpec{
			ModelDeploymentTargets: []v1alpha1.ModelDeploymentTarget{
				{
					Name: mdID1,
				},
				{
					Name:   mdID2,
					Weight: &weight2,
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mr)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(dep_route.MissedWeightErrorMessage))
}

func (s *ModelRouteValidationSuite) TestMultipleTargetsNotExist() {
	weight1 := int32(11)
	weight2 := int32(50)
	mr := &deployment.ModelRoute{
		Spec: v1alpha1.ModelRouteSpec{
			ModelDeploymentTargets: []v1alpha1.ModelDeploymentTarget{
				{
					Name:   "not-exists",
					Weight: &weight1,
				},
				{
					Name:   mdID2,
					Weight: &weight2,
				},
			},
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mr)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("entity \"not-exists\" is not found"))
}

func (s *ModelRouteValidationSuite) TestURLStartWithSlash() {
	mr := &deployment.ModelRoute{
		Spec: v1alpha1.ModelRouteSpec{
			URLPrefix: "test/test",
		},
	}

	err := s.validator.ValidatesAndSetDefaults(mr)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(dep_route.URLPrefixSlashErrorMessage))
}

func (s *ModelRouteValidationSuite) TestForbiddenPrefixes() {
	for _, prefix := range dep_route.ForbiddenPrefixes {
		mr := &deployment.ModelRoute{
			Spec: v1alpha1.ModelRouteSpec{
				URLPrefix: fmt.Sprintf("%s/test/test", prefix),
			},
		}

		err := s.validator.ValidatesAndSetDefaults(mr)
		s.g.Expect(err).To(HaveOccurred())
		s.g.Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(dep_route.ForbiddenPrefix, prefix)))
	}
}

func (s *ModelRouteValidationSuite) TestAllowForbiddenPrefixes() {
	for _, prefix := range dep_route.ForbiddenPrefixes {
		mr := &deployment.ModelRoute{
			ID: mrID,
			Spec: v1alpha1.ModelRouteSpec{
				URLPrefix: fmt.Sprintf("%s/test/test", prefix),
			},
		}

		err := s.validator.ValidatesAndSetDefaults(mr)
		s.g.Expect(err).To(HaveOccurred())
		s.g.Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(dep_route.ForbiddenPrefix, prefix)))
	}
}

func (s *ModelRouteValidationSuite) TestValidateID() {
	mr := &deployment.ModelRoute{
		ID: "not-VALID-id-",
	}

	err := s.validator.ValidatesAndSetDefaults(mr)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(validation.ErrIDValidation.Error()))
}

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
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	"testing"

	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/odahuflow/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	pack_route "github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/packaging"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

type modelPackagingIntegrationSuite struct {
	suite.Suite
	g *GomegaWithT
}

func (s *modelPackagingIntegrationSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func TestMPIValidateSuite(t *testing.T) {
	suite.Run(t, new(modelPackagingIntegrationSuite))
}

func (s *modelPackagingIntegrationSuite) TestPiIDValidation() {
	pi := &packaging.PackagingIntegration{
		Spec: packaging.PackagingIntegrationSpec{},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(pack_route.EmptyIDErrorMessage))
}

func (s *modelPackagingIntegrationSuite) TestPiEntrypointValidation() {
	pi := &packaging.PackagingIntegration{
		Spec: packaging.PackagingIntegrationSpec{},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(pack_route.EmptyEntrypointErrorMessage))
}

func (s *modelPackagingIntegrationSuite) TestPiDefaultImageValidation() {
	pi := &packaging.PackagingIntegration{
		Spec: packaging.PackagingIntegrationSpec{},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(pack_route.EmptyDefaultImageErrorMessage))
}

func (s *modelPackagingIntegrationSuite) TestPiEmptyTargetsValidation() {
	pi := &packaging.PackagingIntegration{
		ID: "some-id",
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   "some-entrypoint",
			DefaultImage: "some/image:tag",
			Privileged:   false,
		},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).ShouldNot(HaveOccurred())
}

func (s *modelPackagingIntegrationSuite) TestPiEmptyArgumentsValidation() {
	pi := &packaging.PackagingIntegration{
		ID: "some-id",
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   "some-entrypoint",
			DefaultImage: "some/image:tag",
			Privileged:   false,
		},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).ShouldNot(HaveOccurred())
}

func (s *modelPackagingIntegrationSuite) TestPiEmptyTargetName() {
	pi := &packaging.PackagingIntegration{
		ID: "some-id",
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   "some-entrypoint",
			DefaultImage: "some/image:tag",
			Privileged:   false,
			Schema: packaging.Schema{
				Targets: []v1alpha1.TargetSchema{
					{
						Name:            "",
						ConnectionTypes: nil,
						Required:        false,
					},
				},
			},
		},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(pack_route.TargetEmptyNameErrorMessage))
}

func (s *modelPackagingIntegrationSuite) TestPiEmptyConnectionType() {
	targetName := "some-name"
	pi := &packaging.PackagingIntegration{
		ID: "some-id",
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   "some-entrypoint",
			DefaultImage: "some/image:tag",
			Privileged:   false,
			Schema: packaging.Schema{
				Targets: []v1alpha1.TargetSchema{
					{
						Name:            targetName,
						ConnectionTypes: nil,
						Required:        false,
					},
				},
			},
		},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(fmt.Sprintf(
		pack_route.TargetEmptyConnectionTypesErrorMessage, targetName,
	)))
}

func (s *modelPackagingIntegrationSuite) TestPiUnknownConnectionType() {
	unknownConnType := "some-type"
	targetName := "some-name"
	pi := &packaging.PackagingIntegration{
		ID: "some-id",
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   "some-entrypoint",
			DefaultImage: "some/image:tag",
			Privileged:   false,
			Schema: packaging.Schema{
				Targets: []v1alpha1.TargetSchema{
					{
						Name:            targetName,
						ConnectionTypes: []string{unknownConnType},
						Required:        false,
					},
				},
			},
		},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(fmt.Sprintf(
		pack_route.TargetUnknownConnTypeErrorMessage, targetName, unknownConnType,
	)))
}

func (s *modelPackagingIntegrationSuite) TestPiNotValidJsonSchema() {
	pi := &packaging.PackagingIntegration{
		ID: "some-id",
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   "some-entrypoint",
			DefaultImage: "some/image:tag",
			Privileged:   false,
			Schema: packaging.Schema{
				Arguments: packaging.JsonSchema{
					Properties: []packaging.Property{
						{
							Name: "some-var",
							Parameters: []packaging.Parameter{
								{
									Name:  "type",
									Value: "bool",
								},
							},
						},
					},
				},
			},
		},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("given: /bool/ Expected valid values are"))
}

func (s *modelPackagingIntegrationSuite) TestValidateID() {
	pi := &packaging.PackagingIntegration{
		ID: "not-VALID-id-",
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(validation.ErrIDValidation.Error()))
}

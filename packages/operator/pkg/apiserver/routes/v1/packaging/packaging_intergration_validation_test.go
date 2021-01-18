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

	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	pack_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/packaging"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

type PackagingIntegrationValidationSuite struct {
	suite.Suite
	g *GomegaWithT
}

func (s *PackagingIntegrationValidationSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func TestMPIValidateSuite(t *testing.T) {
	suite.Run(t, new(PackagingIntegrationValidationSuite))
}

func (s *PackagingIntegrationValidationSuite) TestPiIDValidation() {
	pi := &packaging.PackagingIntegration{
		Spec: packaging.PackagingIntegrationSpec{},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(validation.EmptyValueStringError, "ID"))
}

func (s *PackagingIntegrationValidationSuite) TestPiEntrypointValidation() {
	pi := &packaging.PackagingIntegration{
		Spec: packaging.PackagingIntegrationSpec{},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(pack_route.EmptyEntrypointErrorMessage))
}

func (s *PackagingIntegrationValidationSuite) TestPiDefaultImageValidation() {
	pi := &packaging.PackagingIntegration{
		Spec: packaging.PackagingIntegrationSpec{},
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(pack_route.EmptyDefaultImageErrorMessage))
}

func (s *PackagingIntegrationValidationSuite) TestPiEmptyTargetsValidation() {
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

func (s *PackagingIntegrationValidationSuite) TestPiEmptyArgumentsValidation() {
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

func (s *PackagingIntegrationValidationSuite) TestPiEmptyTargetName() {
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

func (s *PackagingIntegrationValidationSuite) TestPiEmptyConnectionType() {
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

func (s *PackagingIntegrationValidationSuite) TestPiUnknownConnectionType() {
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

func (s *PackagingIntegrationValidationSuite) TestPiNotValidJsonSchema() {
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

func (s *PackagingIntegrationValidationSuite) TestValidateID() {
	pi := &packaging.PackagingIntegration{
		ID: "not-VALID-id-",
	}

	err := pack_route.NewPiValidator().ValidateAndSetDefaults(pi)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(validation.ErrIDValidation.Error()))
}

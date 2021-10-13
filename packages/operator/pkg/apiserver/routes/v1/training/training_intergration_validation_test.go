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
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	train_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TrainingIntegrationValidationSuite struct {
	suite.Suite
	g         *GomegaWithT
	validator *train_route.TiValidator
}

func (s *TrainingIntegrationValidationSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func (s *TrainingIntegrationValidationSuite) SetupSuite() {
	s.validator = train_route.NewTiValidator()
}

func TestTrainingIntegrationValidationSuite(t *testing.T) {
	suite.Run(t, new(TrainingIntegrationValidationSuite))
}

func (s *TrainingIntegrationValidationSuite) TestTiEntrypointEmpty() {
	ti := &training.TrainingIntegration{
		Spec: v1alpha1.TrainingIntegrationSpec{},
	}

	err := s.validator.ValidatesAndSetDefaults(ti)
	s.g.Expect(err).ShouldNot(BeNil())
	s.g.Expect(err.Error()).Should(ContainSubstring(train_route.EmptyEntrypointErrorMessage))
}

func (s *TrainingIntegrationValidationSuite) TestTiDefaultImageEmpty() {
	ti := &training.TrainingIntegration{
		Spec: v1alpha1.TrainingIntegrationSpec{},
	}

	err := s.validator.ValidatesAndSetDefaults(ti)
	s.g.Expect(err).ShouldNot(BeNil())
	s.g.Expect(err.Error()).Should(ContainSubstring(train_route.EmptyDefaultImageErrorMessage))
}

func (s *TrainingIntegrationValidationSuite) TestValidateID() {
	ti := &training.TrainingIntegration{
		ID: "not-VALID-id-",
	}

	err := s.validator.ValidatesAndSetDefaults(ti)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(validation.ErrIDValidation.Error()))
}

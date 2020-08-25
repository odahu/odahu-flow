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

package postgres_test

import (
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	mp_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	postgres_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	mpImage    = "test:new_image"
	mpNewImage = "test:new_image"
	mpType     = "test-type"
	mpID       = "mp1"
)

var (
	mpArtifactName = "someArtifactName"
	mpArguments    = map[string]interface{}{
		"key-1": "value-1",
		"key-2": float64(5),
		"key-3": true,
	}
	mpTargets = []odahuflowv1alpha1.Target{
		{
			Name:           "test",
			ConnectionName: "test-conn",
		},
	}
)

type MPRepositorySuite struct {
	suite.Suite
	rep       mp_repository.Repository
}

func generateMP() *packaging.ModelPackaging {
	return &packaging.ModelPackaging{
		ID: mpID,
		Spec: packaging.ModelPackagingSpec{
			ArtifactName:    mpArtifactName,
			IntegrationName: mpType,
			Image:           mpImage,
			Arguments:       mpArguments,
			Targets:         mpTargets,
		},
	}
}

func (s *MPRepositorySuite) SetupSuite() {

	s.rep = postgres_repo.PackagingRepo{DB: db}
}

func (s *MPRepositorySuite) TearDownTest() {
	if err := s.rep.DeleteModelPackaging(mpID); err != nil && !odahuErrors.IsNotFoundError(err) {
		// If we get the panic that we have a test configuration problem
		panic(err)
	}
}

func TestSuiteMP(t *testing.T) {
	suite.Run(t, new(MPRepositorySuite))
}

func (s *MPRepositorySuite) TestModelPackagingRepository() {
	created := generateMP()

	assert.NoError(s.T(), s.rep.CreateModelPackaging(created))

	fetched, err := s.rep.GetModelPackaging(mpID)
	assert.NoError(s.T(), err)
	assert.Exactly(s.T(), fetched.ID, created.ID)
	assert.Exactly(s.T(), fetched.Spec, created.Spec)

	updated := fetched
	updated.Spec.Image = mpNewImage
	assert.NoError(s.T(), s.rep.UpdateModelPackaging(updated))

	fetched, err = s.rep.GetModelPackaging(mpID)
	assert.NoError(s.T(), err)
	assert.Exactly(s.T(), fetched.ID, updated.ID)
	assert.Exactly(s.T(), fetched.Spec, updated.Spec)
	assert.Exactly(s.T(), fetched.Spec.Image, mpNewImage)

	assert.False(s.T(), fetched.DeletionMark)
	assert.NoError(s.T(), s.rep.SetDeletionMark(mpID, true))
	fetched, err = s.rep.GetModelPackaging(mpID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), fetched.DeletionMark)


	assert.NoError(s.T(), s.rep.DeleteModelPackaging(mpID))
	_, err = s.rep.GetModelPackaging(mpID)
	assert.Error(s.T(), err)
	assert.Exactly(s.T(), err, odahuErrors.NotFoundError{Entity: mpID})
}


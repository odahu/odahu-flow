/*
 * Copyright 2021 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package training_integration //nolint

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	training2 "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"time"
)

func NewService(repo training2.TrainingIntegrationRepository) *Service {
	return &Service{repo: repo}
}

type Service struct {
	repo training2.TrainingIntegrationRepository
}

func (tis *Service) GetTrainingIntegration(id string) (*training.TrainingIntegration, error) {
	return tis.repo.GetTrainingIntegration(id)
}

func (tis *Service) GetTrainingIntegrationList(options ...filter.ListOption) (
	[]training.TrainingIntegration, error,
) {
	return tis.repo.GetTrainingIntegrationList(options...)
}

func (tis *Service) CreateTrainingIntegration(md *training.TrainingIntegration) error {
	md.CreatedAt = time.Now()
	md.UpdatedAt = md.CreatedAt
	return tis.repo.SaveTrainingIntegration(md)
}

func (tis *Service) UpdateTrainingIntegration(md *training.TrainingIntegration) error {
	md.UpdatedAt = time.Now()
	oldMd, err := tis.repo.GetTrainingIntegration(md.ID)
	if err != nil {
		return err
	}
	md.CreatedAt = oldMd.CreatedAt
	return tis.repo.UpdateTrainingIntegration(md)
}

func (tis *Service) DeleteTrainingIntegration(id string) error {
	return tis.repo.DeleteTrainingIntegration(id)
}

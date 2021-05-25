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

package toolchain

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	training2 "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"time"
)

func NewService(repo training2.ToolchainRepository) *Service {
	return &Service{repo: repo}
}

type Service struct {
	repo training2.ToolchainRepository
}

func (tis *Service) GetToolchainIntegration(id string) (*training.ToolchainIntegration, error) {
	return tis.repo.GetToolchainIntegration(id)
}

func (tis *Service) GetToolchainIntegrationList(options ...filter.ListOption) (
	[]training.ToolchainIntegration, error,
) {
	return tis.repo.GetToolchainIntegrationList(options...)
}

func (tis *Service) CreateToolchainIntegration(md *training.ToolchainIntegration) error {
	md.CreatedAt = time.Now()
	md.UpdatedAt = md.CreatedAt
	return tis.repo.SaveToolchainIntegration(md)
}

func (tis *Service) UpdateToolchainIntegration(md *training.ToolchainIntegration) error {
	md.UpdatedAt = time.Now()

	// First try to check that row exists otherwise raise exception to fit interface
	oldTi, err := tis.GetToolchainIntegration(md.ID)
	if err != nil {
		return err
	}

	md.Status = oldTi.Status
	md.CreatedAt = oldTi.CreatedAt

	return tis.repo.UpdateToolchainIntegration(md)
}

func (tis *Service) DeleteToolchainIntegration(id string) error {
	// First try to check that row exists otherwise raise exception to fit interface
	_, err := tis.GetToolchainIntegration(id)
	if err != nil {
		return err
	}

	return tis.repo.DeleteToolchainIntegration(id)
}

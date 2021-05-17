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

type Service interface {
	GetToolchainIntegration(name string) (*training.ToolchainIntegration, error)
	GetToolchainIntegrationList(options ...filter.ListOption) ([]training.ToolchainIntegration, error)
	CreateToolchainIntegration(md *training.ToolchainIntegration) error
	UpdateToolchainIntegration(md *training.ToolchainIntegration) error
	DeleteToolchainIntegration(name string) error
}

func NewService(repo training2.ToolchainRepository) *toolchainService {
	return &toolchainService{repo: repo}
}

type toolchainService struct {
	repo training2.ToolchainRepository
}

func (tis *toolchainService) GetToolchainIntegration(id string) (*training.ToolchainIntegration, error) {
	return tis.repo.GetToolchainIntegration(id)
}

func (tis *toolchainService) GetToolchainIntegrationList(options ...filter.ListOption) (
	[]training.ToolchainIntegration, error,
) {
	return tis.repo.GetToolchainIntegrationList(options...)
}

func (tis *toolchainService) CreateToolchainIntegration(md *training.ToolchainIntegration) error {
	md.CreatedAt = time.Now()
	md.UpdatedAt = md.CreatedAt
	return tis.repo.SaveToolchainIntegration(md)
}

func (tis *toolchainService) UpdateToolchainIntegration(md *training.ToolchainIntegration) error {
	md.UpdatedAt = time.Now()
	return tis.repo.UpdateToolchainIntegration(md)
}

func (tis *toolchainService) DeleteToolchainIntegration(id string) error {
	return tis.repo.DeleteToolchainIntegration(id)
}

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

package packaging_integration //nolint

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	packaging_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"time"
)

func NewService(repo packaging_repo.PackagingIntegrationRepository) *PackagingIntegrationService {
	return &PackagingIntegrationService{repo: repo}
}

type PackagingIntegrationService struct {
	repo packaging_repo.PackagingIntegrationRepository
}

func (pis *PackagingIntegrationService) GetPackagingIntegration(id string) (*packaging.PackagingIntegration, error) {
	return pis.repo.GetPackagingIntegration(id)
}

func (pis *PackagingIntegrationService) GetPackagingIntegrationList(options ...filter.ListOption) (
	[]packaging.PackagingIntegration, error,
) {
	return pis.repo.GetPackagingIntegrationList(options...)
}

func (pis *PackagingIntegrationService) CreatePackagingIntegration(md *packaging.PackagingIntegration) error {
	md.CreatedAt = time.Now()
	md.UpdatedAt = md.CreatedAt
	return pis.repo.SavePackagingIntegration(md)
}

func (pis *PackagingIntegrationService) UpdatePackagingIntegration(md *packaging.PackagingIntegration) error {
	md.UpdatedAt = time.Now()
	return pis.repo.UpdatePackagingIntegration(md)
}

func (pis *PackagingIntegrationService) DeletePackagingIntegration(id string) error {
	return pis.repo.DeletePackagingIntegration(id)
}

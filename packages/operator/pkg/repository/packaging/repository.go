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

package packaging

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	utils "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/postgres"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
)

const (
	TagKey = "name"
)

type Repository interface {
	GetModelPackaging(ctx context.Context, qrr utils.Querier, id string) (*packaging.ModelPackaging, error)
	GetModelPackagingList(
		ctx context.Context, qrr utils.Querier, options ...filter.ListOption) ([]packaging.ModelPackaging, error)
	DeleteModelPackaging(ctx context.Context, qrr utils.Querier, id string) error
	SetDeletionMark(ctx context.Context, qrr utils.Querier, id string, value bool) error
	UpdateModelPackaging(ctx context.Context, qrr utils.Querier, mp *packaging.ModelPackaging) error
	UpdateModelPackagingStatus(ctx context.Context, qrr utils.Querier, id string, s v1alpha1.ModelPackagingStatus) error
	CreateModelPackaging(ctx context.Context, qrr utils.Querier, mp *packaging.ModelPackaging) error
}

type PackagingIntegrationRepository interface {
	GetPackagingIntegration(id string) (*packaging.PackagingIntegration, error)
	GetPackagingIntegrationList(options ...filter.ListOption) ([]packaging.PackagingIntegration, error)
	DeletePackagingIntegration(id string) error
	UpdatePackagingIntegration(md *packaging.PackagingIntegration) error
	CreatePackagingIntegration(md *packaging.PackagingIntegration) error
}

type MPFilter struct {
	Type []string `name:"type" postgres:"spec->>'type'"`
}

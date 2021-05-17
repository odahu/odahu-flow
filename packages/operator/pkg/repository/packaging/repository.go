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
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
)

const (
	TagKey = "name"
)

type Repository interface {
	GetModelPackaging(ctx context.Context, tx *sql.Tx, id string) (*packaging.ModelPackaging, error)
	GetModelPackagingList(
		ctx context.Context, tx *sql.Tx, options ...filter.ListOption) ([]packaging.ModelPackaging, error)
	DeleteModelPackaging(ctx context.Context, tx *sql.Tx, id string) error
	SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error
	UpdateModelPackaging(ctx context.Context, tx *sql.Tx, mp *packaging.ModelPackaging) error
	UpdateModelPackagingStatus(ctx context.Context, tx *sql.Tx, id string, s v1alpha1.ModelPackagingStatus) error
	SaveModelPackaging(ctx context.Context, tx *sql.Tx, mp *packaging.ModelPackaging) error
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
}

type PackagingIntegrationRepository interface {
	GetPackagingIntegration(id string) (*packaging.PackagingIntegration, error)
	GetPackagingIntegrationList(options ...filter.ListOption) ([]packaging.PackagingIntegration, error)
	DeletePackagingIntegration(id string) error
	UpdatePackagingIntegration(md *packaging.PackagingIntegration) error
	SavePackagingIntegration(md *packaging.PackagingIntegration) error
}

type MPFilter struct {
	Type []string `name:"type" postgres:"spec->>'type'"`
}

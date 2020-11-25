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

package deployment

import (
	"context"
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
)

const (
	TagKey = "name"
)

type Repository interface {
	GetModelDeployment(ctx context.Context, tx *sql.Tx, id string) (*deployment.ModelDeployment, error)
	GetModelDeploymentList(
		ctx context.Context, tx *sql.Tx, options ...filter.ListOption) ([]deployment.ModelDeployment, error)
	DeleteModelDeployment(ctx context.Context, tx *sql.Tx, id string) error
	UpdateModelDeployment(ctx context.Context, tx *sql.Tx, md *deployment.ModelDeployment) error
	UpdateModelDeploymentStatus(ctx context.Context, tx *sql.Tx, id string, s v1alpha1.ModelDeploymentStatus) error
	CreateModelDeployment(ctx context.Context, tx *sql.Tx, md *deployment.ModelDeployment) error
	SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
}

type MdFilter struct {
	RoleName []string `name:"roleName" postgres:"spec->>'roleName'"`
}

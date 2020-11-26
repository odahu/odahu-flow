//
//    Copyright 2020 EPAM Systems
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

package route

import (
	"context"
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
)

type Repository interface {
	GetModelRoute(ctx context.Context, tx *sql.Tx, name string) (*deployment.ModelRoute, error)
	GetModelRouteList(ctx context.Context, tx *sql.Tx, options ...filter.ListOption) ([]deployment.ModelRoute, error)
	DeleteModelRoute(ctx context.Context, tx *sql.Tx, name string) error
	UpdateModelRoute(ctx context.Context, tx *sql.Tx, md *deployment.ModelRoute) error
	CreateModelRoute(ctx context.Context, tx *sql.Tx, r *deployment.ModelRoute) error
	UpdateModelRouteStatus(ctx context.Context, tx *sql.Tx, id string, s v1alpha1.ModelRouteStatus) error
	SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error
	BeginTransaction(ctx context.Context) (*sql.Tx, error)

	DefaultExists(ctx context.Context, mdID string, qrr *sql.Tx) (bool, error)
	IsDefault(ctx context.Context, id string, tx *sql.Tx) (bool, error)
}


type Filter struct {
	Default	     bool `name:"is_default" postgres:"is_default"`
	MdID	     []string `name:"mdId" postgres:"spec->'modelDeployments'->0->>'mdName'"`
}
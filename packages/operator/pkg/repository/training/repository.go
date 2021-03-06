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

package training

import (
	"context"
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
)

const (
	TagKey = "name"
)

type Repository interface {
	GetModelTraining(ctx context.Context, tx *sql.Tx, id string) (*training.ModelTraining, error)
	GetModelTrainingList(
		ctx context.Context, tx *sql.Tx, options ...filter.ListOption) ([]training.ModelTraining, error)
	DeleteModelTraining(ctx context.Context, tx *sql.Tx, id string) error
	SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error
	UpdateModelTraining(ctx context.Context, tx *sql.Tx, mt *training.ModelTraining) error
	UpdateModelTrainingStatus(ctx context.Context, tx *sql.Tx, id string, s v1alpha1.ModelTrainingStatus) error
	SaveModelTraining(ctx context.Context, tx *sql.Tx, mt *training.ModelTraining) error
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
}

type ToolchainRepository interface {
	GetToolchainIntegration(name string) (*training.ToolchainIntegration, error)
	GetToolchainIntegrationList(options ...filter.ListOption) ([]training.ToolchainIntegration, error)
	DeleteToolchainIntegration(name string) error
	UpdateToolchainIntegration(md *training.ToolchainIntegration) error
	SaveToolchainIntegration(md *training.ToolchainIntegration) error
}

type MTFilter struct {
	Toolchain    []string `name:"toolchain" postgres:"spec->>'toolchain'"`
	ModelName    []string `name:"model_name" postgres:"spec->'model'->>'name'"`
	ModelVersion []string `name:"model_version" postgres:"spec->'model'->>'version'"`
}

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

package connection

import (
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
)

const (
	TagKey = "name"
)

type Repository interface {
	SaveModelPackagingResult(id string, result []odahuflowv1alpha1.ModelPackagingResult) error
	GetModelPackagingResult(id string) ([]odahuflowv1alpha1.ModelPackagingResult, error)
	GetModelPackaging(id string) (*packaging.ModelPackaging, error)
	GetModelPackagingList(options ...kubernetes.ListOption) ([]packaging.ModelPackaging, error)
	DeleteModelPackaging(id string) error
	GetModelPackagingLogs(id string, writer utils.Writer, follow bool) error
	UpdateModelPackaging(md *packaging.ModelPackaging) error
	CreateModelPackaging(md *packaging.ModelPackaging) error
	PackagingIntegrationRepository
}

type PackagingIntegrationRepository interface {
	GetPackagingIntegration(id string) (*packaging.PackagingIntegration, error)
	GetPackagingIntegrationList(options ...kubernetes.ListOption) ([]packaging.PackagingIntegration, error)
	DeletePackagingIntegration(id string) error
	UpdatePackagingIntegration(md *packaging.PackagingIntegration) error
	CreatePackagingIntegration(md *packaging.PackagingIntegration) error
}

type MPFilter struct {
	Type []string `name:"type"`
}

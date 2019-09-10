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
	. "github.com/legion-platform/legion/legion/operator/pkg/apis/packaging"
	"github.com/legion-platform/legion/legion/operator/pkg/storage/kubernetes"
	"io"
	"net/http"
)

const (
	TagKey = "name"
)

type Storage interface {
	SaveModelPackagingResult(id string, result map[string]string) error
	GetModelPackaging(id string) (*ModelPackaging, error)
	GetModelPackagingList(options ...kubernetes.ListOption) ([]ModelPackaging, error)
	DeleteModelPackaging(id string) error
	GetModelPackagingLogs(id string, writer Writer, follow bool) error
	UpdateModelPackaging(md *ModelPackaging) error
	CreateModelPackaging(md *ModelPackaging) error
	GetPackagingIntegration(id string) (*PackagingIntegration, error)
	GetPackagingIntegrationList(options ...kubernetes.ListOption) ([]PackagingIntegration, error)
	DeletePackagingIntegration(id string) error
	UpdatePackagingIntegration(md *PackagingIntegration) error
	CreatePackagingIntegration(md *PackagingIntegration) error
}

type MPFilter struct {
	Type []string `name:"type"`
}

type Writer interface {
	http.Flusher
	http.CloseNotifier
	io.Writer
}

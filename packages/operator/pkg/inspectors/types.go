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

package inspectors

import (
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	"go.uber.org/zap"
	"net/http"
)

// ModelServerInspector gets URL prefix for a deployed model
// It return Metadata and Swagger2 if passed prefix is served
// by MLServer that is known for ModelServerInspector otherwise
// returns error
type ModelServerInspector interface {
	// Discover ML Server endpoints under the prefix
	// and return model Metadata and Swagger if it is possible
	// return error if Web API behind prefix is not compatible with
	// MLServer
	Inspect(prefix string, hostHeader string, log *zap.SugaredLogger) (model_types.ServedModel, error)
}

type temporaryErr struct {
	error
}

func (t temporaryErr) Temporary() bool {
	return true
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

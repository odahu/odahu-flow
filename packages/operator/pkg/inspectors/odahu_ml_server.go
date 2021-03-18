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
	"encoding/json"
	"fmt"
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type OdahuMLServerInspector struct {
	EdgeURL    url.URL
	HTTPClient httpClient
}

// Part of swagger spec responsible for model name/version
type SwaggerMetadata struct {
	Info struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`
}

// HTTP codes that may be a temporary error and inspection has to be retried later
var temporaryErrorCodes = []int{
	http.StatusRequestTimeout,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,

	http.StatusNotFound,
	http.StatusInternalServerError,

	// Client can be used in reactive workloads (background workers) that suppose backoff retries
	// So workload can retry attempt to get data after error on server will be
	// fixed
	http.StatusInternalServerError,
}

func (o OdahuMLServerInspector) Inspect(
	prefix string, hostHeader string, log *zap.SugaredLogger) (model model_types.ServedModel, err error) {

	modelRequest := o.generateModelRequest(prefix, hostHeader)
	log.Infow("metadata inspect request", "path", modelRequest.URL.Path,
		"hostHeader", modelRequest.Host)

	var response *http.Response
	response, err = o.HTTPClient.Do(modelRequest)
	if err != nil {
		log.Error(err, "Can not get swagger response for prefix")
		return model, err
	}

	if response.StatusCode >= 400 {
		var body []byte
		_, _ = response.Body.Read(body)
		errorStr := fmt.Sprintf("Request to %s returned status code: %d. Body: %s",
			modelRequest.URL, response.StatusCode, body)

		for _, tempCode := range temporaryErrorCodes {
			if tempCode == response.StatusCode {
				return model, temporaryErr{
					fmt.Errorf(errorStr + "; may be temporary"),
				}
			}
		}
		return model, fmt.Errorf(errorStr)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Errorw("Unable to close response body", zap.Error(err))
		}
	}()

	rawBody, err := ioutil.ReadAll(response.Body)
	log.Debugw("Get response from model", "content", string(rawBody))

	swaggerMeta := &SwaggerMetadata{}
	err = json.Unmarshal(rawBody, swaggerMeta)
	if err != nil {
		return model_types.ServedModel{}, err
	}

	model = model_types.ServedModel{
		Swagger: model_types.Swagger2{Raw: rawBody},
		Metadata: model_types.Metadata{
			ModelName:    swaggerMeta.Info.Title,
			ModelVersion: swaggerMeta.Info.Version,
		},
	}

	return model, nil
}

func (o OdahuMLServerInspector) generateModelRequest(prefix string, hostHeader string) *http.Request {

	MlServerURL := url.URL{
		Scheme: o.EdgeURL.Scheme,
		Host:   o.EdgeURL.Host,
		Path:   path.Join(o.EdgeURL.Path, prefix, "/api/model/info"),
	}

	return &http.Request{
		Method: http.MethodGet,
		URL:    &MlServerURL,
		Host:   hostHeader,
	}
}

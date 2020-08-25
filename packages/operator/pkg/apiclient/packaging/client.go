/*
 * Copyright 2019 EPAM Systems
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

package packaging

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	http_util "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/http"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("model_packaging_api_client")

type packagingAPIClient struct {
	http_util.BaseAPIClient
}

func NewClient(apiURL string, token string, clientID string,
	clientSecret string, tokenURL string) Client {
	return &packagingAPIClient{
		BaseAPIClient: http_util.NewBaseAPIClient(
			apiURL,
			token,
			clientID,
			clientSecret,
			tokenURL,
			"api/v1",
		),
	}
}

func wrapMpLogger(id string) logr.Logger {
	return log.WithValues("mp_id", id)
}

func (c *packagingAPIClient) GetModelPackaging(id string) (mp *packaging.ModelPackaging, err error) {
	mpLogger := wrapMpLogger(id)

	response, err := c.DoRequest(
		http.MethodGet,
		strings.Replace("/model/packaging/:id", ":id", id, 1),
		nil,
	)
	if err != nil {
		mpLogger.Error(err, "Retrieving of the model packaging from API failed")

		return nil, err
	}

	mp = &packaging.ModelPackaging{}
	mpBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		mpLogger.Error(err, "Read all data from API response")

		return nil, err
	}
	defer func() {
		bodyCloseError := response.Body.Close()
		if bodyCloseError != nil {
			mpLogger.Error(err, "Closing model packaging response body")
		}
	}()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("error occures: %s", string(mpBytes))
	}

	err = json.Unmarshal(mpBytes, mp)
	if err != nil {
		mpLogger.Error(err, "Unmarshal the model packaging")

		return nil, err
	}

	return mp, nil
}

func (c *packagingAPIClient) SaveModelPackagingResult(
	id string,
	result []odahuflowv1alpha1.ModelPackagingResult,
) error {
	mpLogger := wrapMpLogger(id)

	response, err := c.DoRequest(
		http.MethodPut,
		strings.Replace("/model/packaging/:id/result", ":id", id, 1),
		result,
	)
	if err != nil {
		mpLogger.Error(err, "Saving of the model packaging result in API failed")

		return err
	}

	if response.StatusCode >= 400 {
		mpBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			mpLogger.Error(err, "Read all data from API response")

			return err
		}
		defer func() {
			bodyCloseError := response.Body.Close()
			if bodyCloseError != nil {
				mpLogger.Error(err, "Closing model packaging response body")
			}
			if err != nil {
				err = bodyCloseError
			}
		}()

		return fmt.Errorf("error occures: %s", string(mpBytes))
	}

	return nil
}

func (c *packagingAPIClient) GetPackagingIntegration(id string) (pi *packaging.PackagingIntegration, err error) {
	piLogger := wrapMpLogger(id)

	response, err := c.DoRequest(
		http.MethodGet,
		strings.Replace("/packaging/integration/:id", ":id", id, 1),
		nil,
	)
	if err != nil {
		piLogger.Error(err, "Retrieving of the packaging integration from API failed")

		return nil, err
	}

	pi = &packaging.PackagingIntegration{}
	piBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		piLogger.Error(err, "Read all data from API response")

		return nil, err
	}
	defer func() {
		bodyCloseError := response.Body.Close()
		if bodyCloseError != nil {
			piLogger.Error(err, "Closing packaging integration response body")
		}
	}()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("error occures: %s", string(piBytes))
	}

	err = json.Unmarshal(piBytes, pi)
	if err != nil {
		piLogger.Error(err, "Unmarshal the packaging integration")

		return nil, err
	}

	return pi, nil
}



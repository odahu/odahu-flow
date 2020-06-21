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

package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	packaging_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	http_util "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/http"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("model_packaging_http_repository")

type httpPackagingRepository struct {
	http_util.BaseAPIClient
}

// todo: doc
func NewRepository(apiURL string, token string, clientID string,
	clientSecret string, tokenURL string) packaging_repository.Repository {
	return &httpPackagingRepository{
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

func (htr *httpPackagingRepository) GetModelPackaging(id string) (mp *packaging.ModelPackaging, err error) {
	mpLogger := wrapMpLogger(id)

	response, err := htr.DoRequest(
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

func (htr *httpPackagingRepository) GetModelPackagingList(options ...kubernetes.ListOption) (
	[]packaging.ModelPackaging, error,
) {
	panic("not implemented")
}

func (htr *httpPackagingRepository) DeleteModelPackaging(id string) error {
	panic("not implemented")
}

func (htr *httpPackagingRepository) UpdateModelPackaging(mp *packaging.ModelPackaging) error {
	panic("not implemented")

}

func (htr *httpPackagingRepository) CreateModelPackaging(mp *packaging.ModelPackaging) error {
	panic("not implemented")
}

func (htr *httpPackagingRepository) GetModelPackagingLogs(
	id string, writer utils.Writer, follow bool,
) error {
	panic("not implemented")
}

func (htr *httpPackagingRepository) SaveModelPackagingResult(
	id string,
	result []odahuflowv1alpha1.ModelPackagingResult,
) error {
	mpLogger := wrapMpLogger(id)

	response, err := htr.DoRequest(
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

func (htr *httpPackagingRepository) GetModelPackagingResult(id string) (
	[]odahuflowv1alpha1.ModelPackagingResult, error,
) {
	panic("not implemented")
}

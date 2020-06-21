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
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	training_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	http_util "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/http"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("model_training_http_repository")

type httpTrainingRepository struct {
	http_util.BaseAPIClient
}

// todo: doc
func NewRepository(
	apiURL string, token string, clientID string,
	clientSecret string, tokenURL string) training_repository.Repository {
	return &httpTrainingRepository{
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

func wrapMtLogger(id string) logr.Logger {
	return log.WithValues(odahuflow.ModelTrainingIDLogPrefix, id)
}

func (htr *httpTrainingRepository) GetModelTraining(id string) (mt *training.ModelTraining, err error) {
	mtLogger := wrapMtLogger(id)

	response, err := htr.DoRequest(
		http.MethodGet,
		strings.Replace("/model/training/:id", ":id", id, 1),
		nil,
	)
	if err != nil {
		mtLogger.Error(err, "Retrieving of the training from API failed")

		return nil, err
	}

	mt = &training.ModelTraining{}
	mtBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		mtLogger.Error(err, "Read all data from API response")

		return nil, err
	}
	defer func() {
		bodyCloseError := response.Body.Close()
		if bodyCloseError != nil {
			mtLogger.Error(err, "Closing training response body")
		}
	}()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("error occures: %s", string(mtBytes))
	}

	err = json.Unmarshal(mtBytes, mt)
	if err != nil {
		mtLogger.Error(err, "Unmarshal the training")

		return nil, err
	}

	return mt, nil
}

func (htr *httpTrainingRepository) GetModelTrainingList(options ...kubernetes.ListOption) (
	[]training.ModelTraining, error,
) {
	panic("not implemented")
}

func (htr *httpTrainingRepository) DeleteModelTraining(id string) error {
	panic("not implemented")
}

func (htr *httpTrainingRepository) UpdateModelTraining(mt *training.ModelTraining) error {
	panic("not implemented")
}

func (htr *httpTrainingRepository) CreateModelTraining(mt *training.ModelTraining) error {
	panic("not implemented")
}

func (htr *httpTrainingRepository) GetModelTrainingLogs(
	id string, writer utils.Writer, follow bool,
) error {
	panic("not implemented")
}

func (htr *httpTrainingRepository) SaveModelTrainingResult(
	id string, result *v1alpha1.TrainingResult,
) error {
	mtLogger := wrapMtLogger(id)

	response, err := htr.DoRequest(
		http.MethodPut,
		strings.Replace("/model/training/:id/result", ":id", id, 1),
		result,
	)
	if err != nil {
		mtLogger.Error(err, "Saving of the model training result in API failed")

		return err
	}

	if response.StatusCode >= 400 {
		mpBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			mtLogger.Error(err, "Read all data from API response")

			return err
		}
		defer func() {
			bodyCloseError := response.Body.Close()
			if bodyCloseError != nil {
				mtLogger.Error(err, "Closing model training response body")
			}
			if err != nil {
				err = bodyCloseError
			}
		}()

		return fmt.Errorf("error occures: %s", string(mpBytes))
	}

	return nil
}

func (htr *httpTrainingRepository) GetModelTrainingResult(id string) (*v1alpha1.TrainingResult, error) {
	panic("not implemented")
}

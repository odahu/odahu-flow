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

package deployment

import (
	"encoding/json"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	http_util "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/http"
	"io/ioutil"
	"net/http"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("model_deployment_api_client")

type apiClient struct {
	http_util.BaseAPIClient
}

func NewClient(authConfig config.AuthConfig) apiClient {
	return apiClient{http_util.NewBaseAPIClient(
		authConfig.APIURL,
		authConfig.APIToken,
		authConfig.ClientID,
		authConfig.ClientSecret,
		authConfig.OAuthOIDCTokenEndpoint,
		"api/v1",
	)}
}

func (c *apiClient) GetModelDeployment(id string) (*deployment.ModelDeployment, error) {
	mpLogger := log.WithValues("mp_id", id)

	response, err := c.DoRequest(
		http.MethodGet,
		"/model/deployment/"+id,
		nil,
	)
	if err != nil {
		mpLogger.Error(err, "Retrieving of the model deployment from API failed")
		return nil, err
	}

	md := &deployment.ModelDeployment{}
	mdBytes, err := ioutil.ReadAll(response.Body)
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

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("error occures: %s", string(mdBytes))
	}

	err = json.Unmarshal(mdBytes, md)
	if err != nil {
		mpLogger.Error(err, "Unmarshal the model packaging")
		return nil, err
	}

	return md, nil
}

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
	"net/url"
	"strings"

	"github.com/go-logr/logr"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	http_util "github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/http"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("connection-http-repository")

type httpConnectionRepository struct {
	http_util.BaseAPIClient
}

func NewRepository(
	apiURL string, token string, clientID string,
	clientSecret string, tokenURL string) conn_repository.Repository {
	return &httpConnectionRepository{
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

func wrapConnLogger(id string) logr.Logger {
	return log.WithValues("conn_id", id)
}

func (hcr *httpConnectionRepository) GetConnection(id string) (conn *connection.Connection, err error) {
	connLogger := wrapConnLogger(id)

	return hcr.getConnectionFromAPI(connLogger, &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Replace("/connection/:id/decrypted", ":id", id, 1),
		},
	})
}

func (hcr *httpConnectionRepository) getConnectionFromAPI(
	connLogger logr.Logger, req *http.Request,
) (conn *connection.Connection, err error) {
	response, err := hcr.Do(req)
	if err != nil {
		connLogger.Error(err, "Retrieving of the connection from API failed")

		return nil, err
	}

	conn = &connection.Connection{}
	connBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		connLogger.Error(err, "Read all data from API response")

		return nil, err
	}
	defer func() {
		bodyCloseError := response.Body.Close()
		if bodyCloseError != nil {
			connLogger.Error(err, "Closing connection response body")
		}
	}()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("error occures: %s", string(connBytes))
	}

	err = json.Unmarshal(connBytes, conn)
	if err != nil {
		connLogger.Error(err, "Unmarshal the connection")

		return nil, err
	}

	return conn, nil
}

func (hcr *httpConnectionRepository) GetConnectionList(options ...conn_repository.ListOption) (
	[]connection.Connection, error,
) {
	panic("not implemented")
}

func (hcr *httpConnectionRepository) DeleteConnection(id string) error {
	panic("not implemented")
}

func (hcr *httpConnectionRepository) UpdateConnection(connection *connection.Connection) error {
	panic("not implemented")
}

func (hcr *httpConnectionRepository) CreateConnection(connection *connection.Connection) error {
	panic("not implemented")
}

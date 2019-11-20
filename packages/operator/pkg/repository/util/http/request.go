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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	defaultAPIRequestTimeout = 10 * time.Second
	authorizationHeaderName  = "Authorization"
	authorizationHeaderValue = "Bearer %s"
)

var log = logf.Log.WithName("connection-controller")

type BaseAPIClient struct {
	// todo: doc
	apiURL string
	// todo: doc
	token string
	// todo: doc
	apiVersion string
}

func NewBaseAPIClient(apiURL string, token string, apiVersion string) BaseAPIClient {
	return BaseAPIClient{
		apiURL:     apiURL,
		token:      token,
		apiVersion: apiVersion,
	}
}

func (bec *BaseAPIClient) Do(req *http.Request) (*http.Response, error) {
	if len(req.URL.Host) == 0 {
		apiURLStr := fmt.Sprintf("%s/%s%s", bec.apiURL, bec.apiVersion, req.URL.Path)
		apiURL, err := url.Parse(apiURLStr)
		if err != nil {
			log.Error(err, "Can not parse API URL. Most likely, it is a problem with configuration.",
				"api_url", apiURLStr)

			return nil, err
		}

		apiURL.RawQuery = req.URL.RawQuery
		req.URL = apiURL
	}

	if req.Header == nil {
		req.Header = make(map[string][]string, 1)
	}

	req.Header[authorizationHeaderName] = []string{
		fmt.Sprintf(authorizationHeaderValue, bec.token),
	}

	apiHTTPClient := http.Client{
		Timeout: defaultAPIRequestTimeout,
	}

	return apiHTTPClient.Do(req)
}

func (bec *BaseAPIClient) DoRequest(httpMethod, path string, body interface{}) (*http.Response, error) {
	var bodyStream io.ReadCloser

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		bodyStream = ioutil.NopCloser(bytes.NewReader(data))
	}

	return bec.Do(&http.Request{
		Method: httpMethod,
		URL: &url.URL{
			Path: path,
		},
		Body: bodyStream,
	})
}

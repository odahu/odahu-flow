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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"strings"
	"time"
)

const (
	defaultAPIRequestTimeout = 10 * time.Second
	authorizationHeaderName  = "Authorization"
	authorizationHeaderValue = "Bearer %s"
	serviceAccountScopes     = "openid profile offline_access groups"
	clientCredentialsGrant   = "client_credentials"
)

var log = logf.Log.WithName("connection-controller")

type BaseAPIClient struct {
	// todo: doc
	apiURL string
	// todo: doc
	token string
	// todo: doc
	apiVersion string
	// tokenURL refers to OpenID Provider Token URL
	tokenURL string
	// OpenID client id
	clientID string
	// OpenID client secret
	clientSecret string
}

type OAuthTokenResponse struct {
	IDToken string `json:"id_token"`
}

func NewBaseAPIClient(
	apiURL string, token string, clientID string,
	clientSecret string, tokenURL string, apiVersion string) BaseAPIClient {
	return BaseAPIClient{
		apiURL:       apiURL,
		token:        token,
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenURL:     tokenURL,
		apiVersion:   apiVersion,
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

	// We need store body bytes for retry in case of login
	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	resp, err := apiHTTPClient.Do(req)

	// First attempt could finished by 401 response
	if resp != nil && loginRequired(resp) {

		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// If login required (401) let's login
		loginErr := bec.Login()
		if loginErr != nil {
			log.Error(loginErr, "Login attempt is failed")
			return resp, loginErr
		}

		// Update authorization header
		req.Header[authorizationHeaderName] = []string{
			fmt.Sprintf(authorizationHeaderValue, bec.token),
		}

		// Try again
		return apiHTTPClient.Do(req)
	}

	return resp, err
}

func (bec *BaseAPIClient) DoContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return bec.Do(req)
}

func loginRequired(response *http.Response) bool {
	return response.StatusCode == http.StatusUnauthorized
}

func (bec *BaseAPIClient) Login() error {

	data := url.Values{}
	data.Set("grant_type", clientCredentialsGrant)
	data.Set("client_id", bec.clientID)
	data.Set("client_secret", bec.clientSecret)
	data.Set("scope", serviceAccountScopes)

	body := strings.NewReader(data.Encode())

	request, err := http.NewRequest(http.MethodPost, bec.tokenURL, body)

	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	apiHTTPClient := http.Client{
		Timeout: defaultAPIRequestTimeout,
	}
	resp, err := apiHTTPClient.Do(request)
	if err != nil {
		return err
	}

	bodyBytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return errRead
	}
	var OAuthResp OAuthTokenResponse
	errJSONParse := json.Unmarshal(bodyBytes, &OAuthResp)
	if errJSONParse != nil {
		return errJSONParse
	}

	bec.token = OAuthResp.IDToken

	defer func() {
		bodyCloseError := resp.Body.Close()
		if bodyCloseError != nil {
			log.Error(err, "Closing model packaging response body")
		}
	}()

	return nil

}

func isTemporaryStatusCode(code int) bool {
	for tempCode := range []int{
		http.StatusRequestTimeout,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,

		// Client can be used in reactive workloads (background workers) that suppose backoff retries
		// So workload can retry attempt to get data after error on server will be
		// fixed
		http.StatusInternalServerError,
	} {
		if code == tempCode {
			return true
		}
	}

	return false
}

type temporaryErr struct {
	error
}

func (t temporaryErr) Temporary() bool {
	return true
}


// DoRequest makes HTTP request with some pre- and post- processing
// body will be encoded to json if not nil
// response will be checked on status code > 2xx. If so error will be returned with
// func Temporary() bool interface in mind
// So that caller can assert mentioned interface to make decision to retry later
func (bec *BaseAPIClient) DoRequest(httpMethod, path string, body interface{}) (*http.Response, error) {
	var bodyStream io.ReadCloser

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		bodyStream = ioutil.NopCloser(bytes.NewReader(data))
	}

	var response *http.Response
	var err error
	response, err = bec.Do(&http.Request{
		Method: httpMethod,
		URL: &url.URL{
			Path: path,
		},
		Body: bodyStream,
	})

	if err != nil {
		log.Error(err, "Retrieving of the ModelRoute events is failed")
		return response, err
	}

	var buf []byte
	buf, err =  ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err, "Unable to read response body")
		if err := response.Body.Close(); err != nil {
			log.Error(err,"Unable to close connection")
		}
		return response, err
	}
	if err := response.Body.Close(); err != nil {
		log.Error(err,"Unable to close connection")
		return response, err
	}

	// Caller should be able to read response body by himself
	response.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	if response.StatusCode >= 400 {
		if isTemporaryStatusCode(response.StatusCode) {
			return response, temporaryErr{
				error: fmt.Errorf(
					"not correct status code: %d. Maybe temporary. Response body: %s",
					response.StatusCode, buf,
				),
			}
		}
		return response, fmt.Errorf(
			"not correct status code: %d. Response body: %s",
			response.StatusCode, buf,
		)
	}
	return response, err
}

func (bec *BaseAPIClient) DoRequestGetBody(httpMethod, path string, body interface{}) ([]byte, error) {
	var buf []byte
	response, err := bec.DoRequest(httpMethod, path, body)
	if err != nil {
		return buf, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

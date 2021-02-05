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

package inspectors_test

import (
	"encoding/json"
	"errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/servicecatalog/inspectors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const (
	someURLPrefix = "some/prefix"
)

var (
	logger *zap.Logger
)

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = l
}

type tritonInspectorSuite struct {
	suite.Suite
	inspector inspectors.TritonInspector
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(tritonInspectorSuite))
}

func (s *tritonInspectorSuite) SetupSuite() {
	s.inspector = inspectors.TritonInspector{
		EdgeHost: "host",
		EdgeURL:  url.URL{},
	}
}

func (s *tritonInspectorSuite) TestInspect() {
	// Simulates that Triton has 1 ready model
	s.inspector.HTTPClient = &httpClient{
		f: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(strings.NewReader(
					`[{"name": "my-model", "version": "11", "ready": true}]`)),
			}, nil
		},
	}

	model, err := s.inspector.Inspect(someURLPrefix, logger.Sugar())
	if err != nil {
		panic(err)
	}

	swaggerSpec := string(model.Swagger.Raw)
	s.Assertions.Contains(swaggerSpec, "/v2/models/my-model/ready")
	s.Assertions.Contains(swaggerSpec, "/v2/models/my-model/infer")
	s.Assertions.Contains(swaggerSpec, `"title": "my-model"`)

	// Check that result spec is a valid JSON
	err = json.Unmarshal(model.Swagger.Raw, &map[string]interface{}{})
	if err != nil {
		panic(err)
	}
}

func (s *tritonInspectorSuite) TestInspect_NoModels() {
	// Simulates that Triton has no models; expect an error
	s.inspector.HTTPClient = &httpClient{
		f: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(strings.NewReader(
					`[]`)),
			}, nil
		},
	}

	_, err := s.inspector.Inspect(someURLPrefix, logger.Sugar())

	s.Assertions.Error(err)
	s.Assertions.Contains(err.Error(), "server serves 0 models")
}

func (s *tritonInspectorSuite) TestInspect_MoreThanOneModel() {
	// Simulates that Triton has 2 models; we expect inspector to generate a spec for the first one only
	s.inspector.HTTPClient = &httpClient{
		f: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(strings.NewReader(
					`[{"name": "my-model", "version": "11", "ready": true}, 
                        {"name": "another-model", "version": "22", "ready": true}]`)),
			}, nil
		},
	}

	model, err := s.inspector.Inspect(someURLPrefix, logger.Sugar())
	if err != nil {
		panic(err)
	}

	swaggerSpec := string(model.Swagger.Raw)
	s.Assertions.Contains(swaggerSpec, "/v2/models/my-model/ready")
	s.Assertions.Contains(swaggerSpec, "/v2/models/my-model/infer")
	s.Assertions.Contains(swaggerSpec, `"title": "my-model"`)

	// Check that result spec is a valid JSON
	err = json.Unmarshal(model.Swagger.Raw, &map[string]interface{}{})
	if err != nil {
		panic(err)
	}
}

func (s *tritonInspectorSuite) TestInspect_NoResponse() {
	// Simulates that Triton server does not respond
	s.inspector.HTTPClient = &httpClient{
		f: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("failed to connect")
		},
	}

	_, err := s.inspector.Inspect(someURLPrefix, logger.Sugar())

	s.Assertions.Error(err)
	s.Assertions.Contains(err.Error(), "failed to fetch model repository")
}

type httpClient struct {
	f httpClientFunc
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	return c.f(req)
}

type httpClientFunc func(req *http.Request) (*http.Response, error)

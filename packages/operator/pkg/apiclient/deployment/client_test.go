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

package deployment_test

import (
	"encoding/json"
	"fmt"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/deployment"
	apis "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	md = apis.ModelDeployment{
		ID: "test-md",
		Spec: odahuflowv1alpha1.ModelDeploymentSpec{
			Image:     "image-name:tag",
			Predictor: odahuflow.OdahuMLServer.ID,
		},
	}
)

type mdSuite struct {
	suite.Suite
	testServer *httptest.Server
	client     deployment.Client
}

func (s *mdSuite) SetupSuite() {
	s.testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/model/deployment/test-md":
			if r.Method != http.MethodGet {
				notFound(w, r)
				return
			}
			w.WriteHeader(http.StatusOK)
			mdBytes, err := json.Marshal(md)
			if err != nil {
				// Must not be occurred
				panic(err)
			}

			_, err = w.Write(mdBytes)
			if err != nil {
				// Must not be occurred
				panic(err)
			}
		// Mock endpoint that returns some HTML response (e.g. simulate Nginx error)
		case "/api/v1/model/deployment/get-html-response":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("<html>some error page</html>"))
			if err != nil {
				// Must not be occurred
				panic(err)
			}
		// Mock endpoint that does not response
		case "/api/v1/model/deployment/no-response":
			panic(http.ErrAbortHandler)
		default:
			notFound(w, r)
		}
	}))

	s.client = deployment.NewClient(config.AuthConfig{APIURL: s.testServer.URL})
}

func (s *mdSuite) TearDownSuite() {
	s.testServer.Close()
}

func TestMdSuite(t *testing.T) {
	suite.Run(t, new(mdSuite))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprintf(w, "%s url not found", r.URL.Path)
	if err != nil {
		// Must not be occurred
		panic(err)
	}
}

func (s *mdSuite) TestGetDeployment() {
	mdFromClient, err := s.client.GetModelDeployment(md.ID)

	s.Assertions.NoError(err)
	s.Assertions.Equal(md, *mdFromClient)
}

func (s *mdSuite) TestGetDeployment_NotFound() {
	mdFromClient, err := s.client.GetModelDeployment("nonexistent-deployment")

	s.Assertions.Nil(mdFromClient)
	s.Assertions.Error(err)
	s.Assertions.Contains(err.Error(), "not found")
}

func (s *mdSuite) TestGetDeployment_CannotUnmarshal() {
	mdFromClient, err := s.client.GetModelDeployment("get-html-response")

	s.Assertions.Nil(mdFromClient)
	s.Assertions.Error(err)
	s.Assertions.Contains(err.Error(), "invalid character")
}

func (s *mdSuite) TestGetDeployment_NoResponse() {
	mdFromClient, err := s.client.GetModelDeployment("no-response")

	s.Assertions.Nil(mdFromClient)
	s.Assertions.Error(err)
	s.Assertions.Contains(err.Error(), "EOF")
}

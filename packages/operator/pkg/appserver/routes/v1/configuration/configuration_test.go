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

package configuration_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	conf_route "github.com/odahu/odahu-flow/packages/operator/pkg/appserver/routes/v1/configuration"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	clientSecret = "client-secret"
)

type ConfigurationRouteSuite struct {
	suite.Suite
	g      *GomegaWithT
	server *gin.Engine
	config *config.Config
}

func (s *ConfigurationRouteSuite) SetupSuite() {
	var err error
	s.config, err = config.LoadConfig()
	if err != nil {
		s.T().Fatalf("Cannot initialize config: %s", err.Error())
	}
	s.config.Operator.Auth.ClientSecret = clientSecret
	s.config.Trainer.Auth.ClientSecret = clientSecret
	s.config.Packager.Auth.ClientSecret = clientSecret

	s.server = gin.Default()
	routeGroup := s.server.Group("")
	conf_route.ConfigureRoutes(routeGroup, *s.config)
}

func (s *ConfigurationRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func TestConnectionRouteSuite(t *testing.T) {
	suite.Run(t, new(ConfigurationRouteSuite))
}

func (s *ConfigurationRouteSuite) TestGetConfiguration() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		conf_route.GetConfigurationURL,
		nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result config.Config
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(*s.config).ShouldNot(Equal(result))
	// Check that we cleanup the auth struct
	s.g.Expect(result.Operator.Auth).Should(Equal(config.AuthConfig{}))
	s.g.Expect(result.Trainer.Auth).Should(Equal(config.AuthConfig{}))
	s.g.Expect(result.Packager.Auth).Should(Equal(config.AuthConfig{}))
}

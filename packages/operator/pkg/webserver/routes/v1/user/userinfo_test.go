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

package userinfo_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/user"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	user_route "github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes/v1/user"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	userName = "John Doe"
	email    = "test@email.org"
	token    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZW1haW" +
		"wiOiJ0ZXN0QGVtYWlsLm9yZyIsImlhdCI6MTUxNjIzOTAyMn0.mDHHgcPKVidgB7VFNfSHS-K08a4a4kRHQF94waNbpzg"
)

type UserInfoRouteSuite struct {
	suite.Suite
	g      *GomegaWithT
	server *gin.Engine
}

func (s *UserInfoRouteSuite) SetupSuite() {
	s.server = gin.Default()
	v1Group := s.server.Group("")
	user_route.ConfigureRoutes(v1Group, config.NewDefaultUserConfig().Claims)
}

func (s *UserInfoRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func TestConnectionRouteSuite(t *testing.T) {
	suite.Run(t, new(UserInfoRouteSuite))
}

func (s *UserInfoRouteSuite) TestGetUserInfoExtractionFromHeaders() {
	userInfoURL, err := url.Parse(user_route.GetUserInfoURL)
	s.g.Expect(err).ShouldNot(HaveOccurred())

	w := httptest.NewRecorder()
	req := &http.Request{
		Method: http.MethodGet,
		URL:    userInfoURL,
		Header: map[string][]string{
			"Authorization": {token},
		},
	}
	s.server.ServeHTTP(w, req)

	var result user.UserInfo
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(result).Should(Equal(user.UserInfo{
		Email:    email,
		Username: userName,
	}))
}

func (s *UserInfoRouteSuite) TestGetUserInfoNonHeaders() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, user_route.GetUserInfoURL, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	var result user.UserInfo
	err = json.Unmarshal(w.Body.Bytes(), &result)
	s.g.Expect(err).NotTo(HaveOccurred())

	s.g.Expect(result).Should(Equal(user.AnonymousUser))
}

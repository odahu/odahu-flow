/*
 * Copyright 2020 EPAM Systems
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

package utils_test

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JWTTestSuite struct {
	suite.Suite
	g *GomegaWithT
}

func (s *JWTTestSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func TestJWTTestSuite(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}

func (s *JWTTestSuite) TestExtractUserInfoFromToken() {
	// {
	//  "sub": "1234567890",
	//  "name": "John Doe",
	//  "email": "test@email.org",
	//  "iat": 1516239022
	// }
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZW1haWwiOiJ0Z" +
		"XN0QGVtYWlsLm9yZyIsImlhdCI6MTUxNjIzOTAyMn0.mDHHgcPKVidgB7VFNfSHS-K08a4a4kRHQF94waNbpzg"

	userInfo, err := utils.ExtractUserInfoFromToken(token)
	s.g.Expect(err).Should(BeNil())
	s.g.Expect(userInfo.Username).Should(Equal("John Doe"))
	s.g.Expect(userInfo.Email).Should(Equal("test@email.org"))
}

func (s *JWTTestSuite) TestExtractUserInfoFromTokenMissingName() {
	// {
	//  "sub": "1234567890",
	//  "email": "test@email.org",
	//  "iat": 1516239022
	// }
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZW1haWwiOiJ0ZXN0QGVtYWlsLm9yZyIs" +
		"ImlhdCI6MTUxNjIzOTAyMn0.ezWMqjto3LXFSzR7dUbusmgUmdaxdRHM9UA6qHj0k-A"

	_, err := utils.ExtractUserInfoFromToken(token)
	s.g.Expect(err).ShouldNot(BeNil())
	s.g.Expect(err.Error()).Should(ContainSubstring("claim is missing"))
}

func (s *JWTTestSuite) TestExtractUserInfoFromTokenMissingEmail() {
	// {
	//  "sub": "1234567890",
	//  "name": "John Doe",
	//  "iat": 1516239022
	// }
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ij" +
		"oxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	_, err := utils.ExtractUserInfoFromToken(token)
	s.g.Expect(err).ShouldNot(BeNil())
	s.g.Expect(err.Error()).Should(ContainSubstring("claim is missing"))
}

func (s *JWTTestSuite) TestExtractUserInfoFromTokenEmailNotStringType() {
	// {
	//  "sub": "1234567890",
	//  "name": "John Doe",
	//  "email": ["test@email.org"],
	//  "iat": 1516239022
	// }
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZW1haWwiOl" +
		"sidGVzdEBlbWFpbC5vcmciXSwiaWF0IjoxNTE2MjM5MDIyfQ._MsFbVFXKIl5JAkiC-pY2KxEqriOtKBpqxLBoBIAvNE"

	_, err := utils.ExtractUserInfoFromToken(token)
	s.g.Expect(err).ShouldNot(BeNil())
	s.g.Expect(err.Error()).Should(ContainSubstring("claim is not the string type"))
}

func (s *JWTTestSuite) TestExtractUserInfoFromTokenNameNotStringType() {
	// {
	//  "sub": "1234567890",
	//  "name": ["John Doe"],
	//  "email": "test@email.org",
	//  "iat": 1516239022
	// }
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6WyJKb2huIERvZSJdLCJlbWF" +
		"pbCI6InRlc3RAZW1haWwub3JnIiwiaWF0IjoxNTE2MjM5MDIyfQ.9Wchix3yJEIDxZ1crkqSXE2xs_uxt8vswQPzL0-Q33E"

	_, err := utils.ExtractUserInfoFromToken(token)
	s.g.Expect(err).ShouldNot(BeNil())
	s.g.Expect(err.Error()).Should(ContainSubstring("claim is not the string type"))
}

func (s *JWTTestSuite) TestExtractUserInfoFromTokenClaimsMissing() {
	//nolint
	token := "malformed_jwt"

	_, err := utils.ExtractUserInfoFromToken(token)
	s.g.Expect(err).Should(Equal(utils.ErrMalformedJWT))
}

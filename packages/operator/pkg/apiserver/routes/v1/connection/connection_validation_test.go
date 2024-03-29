//
//    Copyright 2019 EPAM Systems
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package connection_test

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	conn_route "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes/v1/connection"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	ecrValidUri                   = "7777777777.dkr.ecr.eu-central-1.amazonaws.com/cluster-name/counter-3.3:886a4942-138d-460c-8f3c-eb3b3d2e1d98" //nolint
	failedMockKeyEvaluatorMessage = "failed mock key evaluator"
)

type ConnectionValidationSuite struct {
	suite.Suite
	g *GomegaWithT
	v *conn_route.ConnValidator
}

func (s *ConnectionValidationSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
	s.v = conn_route.NewConnValidator(successMockKeyEvaluator)
}

func successMockKeyEvaluator(uri string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte("success")), nil
}

func failedMockKeyEvaluator(uri string) (string, error) {
	return "", errors.New(failedMockKeyEvaluatorMessage)
}

func TestConnectionValidationSuite(t *testing.T) {
	suite.Run(t, new(ConnectionValidationSuite))
}

func (s *ConnectionValidationSuite) TestEmptyURL() {
	conn := &connection.Connection{
		Spec: v1alpha1.ConnectionSpec{
			Type:      "not-existed",
			Reference: connReference,
			KeySecret: creds,
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(conn_route.EmptyURIErrorMessage))
}


func (s *ConnectionValidationSuite) TestEmptyURLEcr() {
	conn := &connection.Connection{
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.EcrType,
			Reference: connReference,
			KeySecret: creds,
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(conn_route.EmptyURIErrorMessage))
	// If URI is empty then not validate ECR specific format of URI
	s.g.Expect(err.Error()).Should(Not(ContainSubstring(conn_route.ECRTypeNotValidURI)))
}

func (s *ConnectionValidationSuite) TestUnknownTypeType() {
	connType := v1alpha1.ConnectionType("not-existed")
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connType,
			URI:       connURI,
			Reference: connReference,
			KeySecret: creds,
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(
		conn_route.UnknownTypeErrorMessage, connType, connection.AllConnectionTypes,
	)))
}

func (s *ConnectionValidationSuite) TestGitTypePublicKeyBase64() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       connURI,
			Reference: connReference,
			PublicKey: "not-base64-encoding",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("must be base64-encoded"))
}

func (s *ConnectionValidationSuite) TestDockerTypeUsername() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:     connection.DockerType,
			URI:      connURI,
			Password: "password",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(conn_route.DockerTypeUsernameErrorMessage))
}

func (s *ConnectionValidationSuite) TestDockerTypePassword() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:     connection.DockerType,
			URI:      connURI,
			Username: "username",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(conn_route.DockerTypePasswordErrorMessage))
}

func (s *ConnectionValidationSuite) TestValidDockerType() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:     connection.DockerType,
			URI:      connURI,
			Username: "username",
			Password: "password",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).ShouldNot(HaveOccurred())
}

func (s *ConnectionValidationSuite) TestGcsTypeRegion() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GcsType,
			URI:       connURI,
			Role:      "role",
			KeySecret: "a2V5LXNlY3JldA==",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(conn_route.GcsTypeRegionErrorMessage))
}

func (s *ConnectionValidationSuite) TestGcsTypeSecretMissed() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:   connection.GcsType,
			URI:    connURI,
			Region: "region",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(conn_route.GcsTypeKeySecretEmptyErrorMessage))
}

func (s *ConnectionValidationSuite) TestGcsTypeRoleNotSupported() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:   connection.GcsType,
			URI:    connURI,
			Region: "region",
			Role:   "role",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(conn_route.GcsTypeRoleNotSupportedErrorMessage))
}

func (s *ConnectionValidationSuite) TestAzureBlobTypeSecretMissed() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type: connection.AzureBlobType,
			URI:  connURI,
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(conn_route.AzureBlobTypeKeySecretEmptyErrorMessage))
}

func (s *ConnectionValidationSuite) TestS3TypeRegion() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.S3Type,
			URI:       connURI,
			Role:      "role",
			KeySecret: "a2V5LXNlY3JldA==",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(conn_route.S3TypeRegionErrorMessage))
}

func (s *ConnectionValidationSuite) TestS3TypeRoleAndSecretMissed() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:   connection.S3Type,
			URI:    connURI,
			Region: "username",
		},
	}

	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(conn_route.S3TypeKeySecretEmptyErrorMessage))
}

func (s *ConnectionValidationSuite) TestS3ValidRoleParameterNotSupported() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type: connection.S3Type,
			URI:  connURI,
			Role: "role",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(conn_route.S3TypeRoleNotSupportedErrorMessage))
}

func (s *ConnectionValidationSuite) TestGitTypeValid() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       connURI,
			KeySecret: "a2V5LXNlY3JldA==",
			Reference: "branch",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).ShouldNot(HaveOccurred())
}

func (s *ConnectionValidationSuite) TestGitTypeWithoutSecretIsValid() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       connURI,
			Reference: "branch",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).ShouldNot(HaveOccurred())
}

func (s *ConnectionValidationSuite) TestGitTypeInvalidPublicKey() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       connURI,
			PublicKey: "not-base64-format",
			Reference: "branch",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("must be base64-encoded"))
}

func (s *ConnectionValidationSuite) TestGitTypeGeneratePublicKey() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       connURI,
			Reference: "branch",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).ShouldNot(HaveOccurred())
}

func (s *ConnectionValidationSuite) TestGitTypeGeneratePublicKeyError() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       connURI,
			Reference: "branch",
		},
	}

	err := conn_route.NewConnValidator(failedMockKeyEvaluator).ValidatesAndSetDefaults(conn)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(
		fmt.Sprintf(conn_route.GitTypePublicKeyExtractionErrorMessage, failedMockKeyEvaluatorMessage),
	))
}

func (s *ConnectionValidationSuite) TestS3ValidSecretParameter() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.S3Type,
			URI:       connURI,
			Region:    "region",
			KeySecret: "a2V5LXNlY3JldA==",
			KeyID:     "a2V5LWlk",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).ShouldNot(HaveOccurred())
}

func (s *ConnectionValidationSuite) TestECRTypeValidationUrl() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type: connection.EcrType,
			URI:  "not-valid-url",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring(fmt.Sprintf(conn_route.ECRTypeNotValidURI, "")))
}

func (s *ConnectionValidationSuite) TestECRValidateKeySecret() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:   connection.EcrType,
			URI:    ecrValidUri,
			Region: "region",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(conn_route.ECRTypeKeySecretEmptyErrorMessage))
}

func (s *ConnectionValidationSuite) TestECRTypeRegionFromUrl() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.EcrType,
			URI:       ecrValidUri,
			KeySecret: "a2V5LXNlY3JldA==",
			KeyID:     "a2V5LWlk",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(conn.Spec.Region).Should(Equal("eu-central-1"))
}

func (s *ConnectionValidationSuite) TestECRTypeValidParameters() {
	conn := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.EcrType,
			URI:       ecrValidUri,
			KeySecret: "a2V5LXNlY3JldA==",
			KeyID:     "a2V5LWlk",
			Region:    "region",
		},
	}
	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).ShouldNot(HaveOccurred())
}

func (s *ConnectionValidationSuite) TestValidateID() {
	conn := &connection.Connection{
		ID: "not-VALID-id-",
	}

	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring(validation.ErrIDValidation.Error()))
}

func (s *ConnectionValidationSuite) TestValidateBase64Secrets_allValid() {
	conn := &connection.Connection{
		ID: "valid-id",
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       "URI",
			KeySecret: "c2VjcmV0",
			PublicKey: "c2VjcmV0",
			KeyID:     "c2VjcmV0",
			Password:  "c2VjcmV0",
		},
	}

	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).NotTo(HaveOccurred())
}

func (s *ConnectionValidationSuite) TestValidateBase64Secrets_keySecretOnly() {
	conn := &connection.Connection{
		ID: "valid-id",
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       "URI",
			KeySecret: "c2VjcmV0",
		},
	}

	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).NotTo(HaveOccurred())
}

func (s *ConnectionValidationSuite) TestValidateBase64Secrets_invalidKeySecret() {
	conn := &connection.Connection{
		ID: "valid-id",
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       "URI",
			KeySecret: "my secret",
		},
	}

	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("must be base64-encoded"))
}

func (s *ConnectionValidationSuite) TestValidateBase64Secrets_invalidKeyID() {
	conn := &connection.Connection{
		ID: "valid-id",
		Spec: v1alpha1.ConnectionSpec{
			Type:  connection.GITType,
			URI:   "URI",
			KeyID: "my secret",
		},
	}

	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("must be base64-encoded"))
}

func (s *ConnectionValidationSuite) TestValidateBase64Secrets_invalidPassword() {
	conn := &connection.Connection{
		ID: "valid-id",
		Spec: v1alpha1.ConnectionSpec{
			Type:     connection.GITType,
			URI:      "URI",
			Password: "my secret",
		},
	}

	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("must be base64-encoded"))
}

func (s *ConnectionValidationSuite) TestValidateBase64Secrets_invalidPublicKey() {
	conn := &connection.Connection{
		ID: "valid-id",
		Spec: v1alpha1.ConnectionSpec{
			Type:      connection.GITType,
			URI:       "URI",
			PublicKey: "public key",
		},
	}

	err := s.v.ValidatesAndSetDefaults(conn)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(err.Error()).To(ContainSubstring("must be base64-encoded"))
}

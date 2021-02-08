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

package connection

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	"go.uber.org/multierr"
)

const (
	ErrorMessageTemplate                   = "%s: %s"
	EmptyURIErrorMessage                   = "empty uri"
	ValidationConnErrorMessage             = "Validation of connection is failed"
	UnknownTypeErrorMessage                = "unknown type: %s. Supported types: %s"
	PasswordDecodeErrorMessage             = "password must be base64-encoded, error: %s"
	KeyIDDecodeErrorMessage                = "key id must be base64-encoded, error: %s"
	KeySecretDecodeErrorMessage            = "key secret must be base64-encoded, error: %s"
	PublicKeyDecodeErrorMessage            = "public key must be base64-encoded, error: %s"
	DockerTypePasswordErrorMessage         = "docker type requires the password parameter" //nolint
	DockerTypeUsernameErrorMessage         = "docker type requires the username parameter"
	GitTypePublicKeyExtractionErrorMessage = "can not extract the public SSH host key from URI: %s"
	GcsTypeRegionErrorMessage              = "gcs type requires that region must be non-empty"
	GcsTypeKeySecretEmptyErrorMessage      = "gcs type requires that keySecret parameter" +
		" must be non-empty"
	GcsTypeRoleNotSupportedErrorMessage = "gcs type does not support role parameter yet" +
		" must be non-empty"
	AzureBlobTypeKeySecretEmptyErrorMessage = "azureblob type requires that keySecret parameter contains" +
		" HTTP endpoint with SAS Token"
	S3TypeRegionErrorMessage         = "s3 type requires that region must be non-empty"
	S3TypeKeySecretEmptyErrorMessage = "s3 type requires that keyID and keySecret parameters" +
		" must be non-empty"
	S3TypeRoleNotSupportedErrorMessage = "s3 type does not support role parameter yet"
	ECRTypeKeySecretEmptyErrorMessage  = "ecr type requires that keyID and keySecret parameters" +
		" must be non-empty"
	ECRTypeNotValidURI = "not valid uri for ecr type: %s"
)

type PublicKeyEvaluator func(string) (string, error)

type ConnValidator struct {
	keyEvaluator PublicKeyEvaluator
}

// Currently validator does not need any arguments
func NewConnValidator(keyEvaluator PublicKeyEvaluator) *ConnValidator {
	return &ConnValidator{keyEvaluator: keyEvaluator}
}

func (cv *ConnValidator) ValidatesAndSetDefaults(conn *connection.Connection) (err error) {
	err = multierr.Append(validation.ValidateID(conn.ID), err)
	err = multierr.Append(cv.validateBase64Fields(conn), err)

	if len(conn.Spec.URI) == 0 {
		err = multierr.Append(err, errors.New(EmptyURIErrorMessage))
	}

	switch conn.Spec.Type {
	case connection.GITType:
		err = multierr.Append(err, cv.validateGitType(conn))
	case connection.S3Type:
		err = multierr.Append(err, cv.validateS3Type(conn))
	case connection.GcsType:
		err = multierr.Append(err, cv.validateGcsType(conn))
	case connection.AzureBlobType:
		err = multierr.Append(err, cv.validateAzureBlobType(conn))
	case connection.DockerType:
		err = multierr.Append(err, cv.validateDockerType(conn))
	case connection.EcrType:
		err = multierr.Append(err, cv.validateEcrType(conn))
	default:
		err = multierr.Append(err, fmt.Errorf(UnknownTypeErrorMessage, conn.Spec.Type, connection.AllConnectionTypes))
	}

	if err != nil {
		return fmt.Errorf(ErrorMessageTemplate, ValidationConnErrorMessage, err.Error())
	}

	return err
}

// If any of secret fields provided, they must be base64-encoded
func (cv *ConnValidator) validateBase64Fields(conn *connection.Connection) error {
	var err error

	_, decodeErr := base64.StdEncoding.DecodeString(conn.Spec.Password)
	if decodeErr != nil {
		err = multierr.Append(err, fmt.Errorf(PasswordDecodeErrorMessage, decodeErr.Error()))
	}

	_, decodeErr = base64.StdEncoding.DecodeString(conn.Spec.KeySecret)
	if decodeErr != nil {
		err = multierr.Append(err, fmt.Errorf(KeySecretDecodeErrorMessage, decodeErr.Error()))
	}

	_, decodeErr = base64.StdEncoding.DecodeString(conn.Spec.KeyID)
	if decodeErr != nil {
		err = multierr.Append(err, fmt.Errorf(KeyIDDecodeErrorMessage, decodeErr.Error()))
	}

	_, decodeErr = base64.StdEncoding.DecodeString(conn.Spec.PublicKey)
	if decodeErr != nil {
		err = multierr.Append(err, fmt.Errorf(PublicKeyDecodeErrorMessage, decodeErr.Error()))
	}

	return err
}

func (cv *ConnValidator) validateEcrType(conn *connection.Connection) (err error) {
	// Try to parse URI
	if conn.Spec.URI != "" {  // If URI is empty then skip extracting Registry and further logic
		registry, urlParsingErr := api.ExtractRegistry(conn.Spec.URI)
		if urlParsingErr != nil {
			err = multierr.Append(err, fmt.Errorf(ECRTypeNotValidURI, urlParsingErr.Error()))
		}

		if len(conn.Spec.Region) == 0 && registry != nil {
			conn.Spec.Region = registry.Region
			logC.Info("Connection region is empty. Set region from url", "id", conn.ID, "region", registry.Region)
		}
	}

	if len(conn.Spec.KeySecret) == 0 || len(conn.Spec.KeyID) == 0 {
		err = multierr.Append(err, errors.New(ECRTypeKeySecretEmptyErrorMessage))
	}

	return
}

func (cv *ConnValidator) validateS3Type(conn *connection.Connection) (err error) {
	if len(conn.Spec.Region) == 0 {
		err = multierr.Append(err, errors.New(S3TypeRegionErrorMessage))
	}

	if len(conn.Spec.Role) != 0 {
		err = multierr.Append(err, errors.New(S3TypeRoleNotSupportedErrorMessage))
	}

	if len(conn.Spec.KeySecret) == 0 || len(conn.Spec.KeyID) == 0 {
		err = multierr.Append(err, errors.New(S3TypeKeySecretEmptyErrorMessage))
	}

	return
}

func (cv *ConnValidator) validateGcsType(conn *connection.Connection) (err error) {
	if len(conn.Spec.Region) == 0 {
		err = multierr.Append(err, errors.New(GcsTypeRegionErrorMessage))
	}

	if len(conn.Spec.Role) != 0 {
		err = multierr.Append(err, errors.New(GcsTypeRoleNotSupportedErrorMessage))
	}

	if len(conn.Spec.KeySecret) == 0 {
		err = multierr.Append(err, errors.New(GcsTypeKeySecretEmptyErrorMessage))
	}

	return
}

func (cv *ConnValidator) validateAzureBlobType(conn *connection.Connection) (err error) {
	if len(conn.Spec.KeySecret) == 0 {
		err = multierr.Append(err, errors.New(AzureBlobTypeKeySecretEmptyErrorMessage))
	}

	return
}

func (cv *ConnValidator) validateDockerType(conn *connection.Connection) (err error) {
	if len(conn.Spec.Password) == 0 {
		err = multierr.Append(err, errors.New(DockerTypePasswordErrorMessage))
	}

	if len(conn.Spec.Username) == 0 {
		err = multierr.Append(err, errors.New(DockerTypeUsernameErrorMessage))
	}

	return
}

func (cv *ConnValidator) validateGitType(conn *connection.Connection) (err error) {
	if len(conn.Spec.PublicKey) == 0 {

		if conn.Spec.URI == "" {
			logC.Info(".Spec.URI is empty. Skip extracting Git Public Key from URI")
			return nil
		}

		publicKey, keyError := cv.keyEvaluator(conn.Spec.URI)
		if keyError != nil {
			err = multierr.Append(keyError, fmt.Errorf(GitTypePublicKeyExtractionErrorMessage, keyError.Error()))
		} else {
			conn.Spec.PublicKey = base64.StdEncoding.EncodeToString([]byte(publicKey))
		}
	}

	return err
}

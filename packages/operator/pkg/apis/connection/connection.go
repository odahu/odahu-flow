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
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"go.uber.org/multierr"
)

const (
	S3Type            = v1alpha1.ConnectionType("s3")
	GcsType           = v1alpha1.ConnectionType("gcs")
	AzureBlobType     = v1alpha1.ConnectionType("azureblob")
	GITType           = v1alpha1.ConnectionType("git")
	DockerType        = v1alpha1.ConnectionType("docker")
	EcrType           = v1alpha1.ConnectionType("ecr")
	DecryptedDataMask = "*****"
)

var (
	AllConnectionTypes = []v1alpha1.ConnectionType{
		S3Type, GcsType, AzureBlobType, GITType, DockerType, EcrType,
	}
	AllConnectionTypesSet = map[v1alpha1.ConnectionType]interface{}{}
)

func init() {
	for _, conn := range AllConnectionTypes {
		AllConnectionTypesSet[conn] = nil
	}
}

type Connection struct {
	// Connection id
	ID string `json:"id"`
	// Connection specification
	Spec v1alpha1.ConnectionSpec `json:"spec"`
	// Connection status
	Status v1alpha1.ConnectionStatus `json:"status,omitempty"`
}

// Replace sensitive data with mask in the connection
func (c *Connection) DeleteSensitiveData() *Connection {
	if len(c.Spec.Password) != 0 {
		c.Spec.Password = DecryptedDataMask
	}

	if len(c.Spec.KeySecret) != 0 {
		c.Spec.KeySecret = DecryptedDataMask
	}

	if len(c.Spec.KeyID) != 0 {
		c.Spec.KeyID = DecryptedDataMask
	}

	return c
}

// Decodes sensitive data from base64
func (c *Connection) DecodeBase64Secrets() error {
	var err error

	decoded, decodeErr := base64.StdEncoding.DecodeString(c.Spec.Password)
	if decodeErr != nil {
		err = multierr.Append(err, decodeErr)
	}
	c.Spec.Password = string(decoded)

	decoded, decodeErr = base64.StdEncoding.DecodeString(c.Spec.KeySecret)
	if decodeErr != nil {
		err = multierr.Append(err, decodeErr)
	}
	c.Spec.KeySecret = string(decoded)

	decoded, decodeErr = base64.StdEncoding.DecodeString(c.Spec.KeyID)
	if decodeErr != nil {
		err = multierr.Append(err, decodeErr)
	}
	c.Spec.KeyID = string(decoded)

	return err
}

// Encodes sensitive data to base64
func (c *Connection) EncodeBase64Secrets() {
	c.Spec.Password = base64.StdEncoding.EncodeToString([]byte(c.Spec.Password))
	c.Spec.KeySecret = base64.StdEncoding.EncodeToString([]byte(c.Spec.KeySecret))
	c.Spec.KeyID = base64.StdEncoding.EncodeToString([]byte(c.Spec.KeyID))
}

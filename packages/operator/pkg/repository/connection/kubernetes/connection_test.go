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

package kubernetes_test

import (
	"testing"

	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	. "github.com/onsi/gomega"
)

const (
	connID   = "test-id"
	connType = connection.DockerType
)

// TODO: add more tests
func TestConnectionRepository(t *testing.T) {
	g := NewGomegaWithT(t)

	created := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type: connType,
		},
	}

	err := c.SaveConnection(created)
	g.Expect(err).NotTo(HaveOccurred())

	fetched, err := c.GetConnection(connID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.Spec))

	fetched, err = c.GetConnection(connID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.DeleteSensitiveData().Spec))

	newConnType := connection.GcsType
	updated := &connection.Connection{
		ID: connID,
		Spec: v1alpha1.ConnectionSpec{
			Type: newConnType,
		},
	}
	err = c.UpdateConnection(updated)
	g.Expect(err).NotTo(HaveOccurred())

	fetched, err = c.GetConnection(connID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(updated.ID))
	g.Expect(fetched.Spec).To(Equal(updated.DeleteSensitiveData().Spec))
	g.Expect(fetched.Spec.Type).To(Equal(newConnType))

	g.Expect(c.DeleteConnection(connID)).NotTo(HaveOccurred())
	_, err = c.GetConnection(connID)
	g.Expect(err).To(HaveOccurred())
}

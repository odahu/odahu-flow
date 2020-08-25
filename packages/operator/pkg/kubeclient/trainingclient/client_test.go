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

package trainingclient_test

import (
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	"log"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

const (
	mtID            = "foo"
	modelName       = "test-model-name"
	modelVersion    = "test-model-version"
	newModelName    = "new-test-model-name"
	newModelVersion = "new-test-model-version"
	testNamespace = "default"
)

var (
	c trainingclient.Client
)

func TestMain(m *testing.M) {
	os.Exit(WrapperTestMain(m))
}

func WrapperTestMain(m *testing.M) int {
	k8sClient, _, closeF, _, err := testenvs.SetupTestKube(filepath.Join("..", "..", "..", "config", "crds"))
	if err != nil {
		log.Println("Unable to setup kubernetes")
		return -1
	}
	defer func() {
		if err := closeF(); err != nil {
			log.Println("Unable to stop k8s test environment")
		}
	}()
	c = trainingclient.NewClient(testNamespace, testNamespace, k8sClient, nil)
	return m.Run()
}


func TestTrainingKubeClient(t *testing.T) {
	g := NewGomegaWithT(t)

	created := &training.ModelTraining{
		ID: mtID,
		Spec: v1alpha1.ModelTrainingSpec{
			Model: v1alpha1.ModelIdentity{
				Name:    modelName,
				Version: modelVersion,
			},
		},
	}

	g.Expect(c.CreateModelTraining(created)).NotTo(HaveOccurred())

	fetched, err := c.GetModelTraining(mtID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.Spec))

	updated := &training.ModelTraining{
		ID: mtID,
		Spec: v1alpha1.ModelTrainingSpec{
			Model: v1alpha1.ModelIdentity{
				Name:    newModelName,
				Version: newModelVersion,
			},
		},
	}
	g.Expect(c.UpdateModelTraining(updated)).NotTo(HaveOccurred())

	fetched, err = c.GetModelTraining(mtID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(updated.ID))
	g.Expect(fetched.Spec).To(Equal(updated.Spec))
	g.Expect(fetched.Spec.Model.Name).To(Equal(newModelName))
	g.Expect(fetched.Spec.Model.Version).To(Equal(newModelVersion))

	g.Expect(c.DeleteModelTraining(mtID)).NotTo(HaveOccurred())
	_, err = c.GetModelTraining(mtID)
	g.Expect(err).To(HaveOccurred())
}

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

package packagingclient_test

import (
	"context"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/packagingclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

const (
	mpImage    = "test:new_image"
	mpNewImage = "test:new_image"
	mpType     = "test-type"
	mpID       = "mp1"
	testNamespace = "default"
)

var (
	mpArtifactName = "someArtifactName"
	mpArguments    = map[string]interface{}{
		"key-1": "value-1",
		"key-2": float64(5),
		"key-3": true,
	}
	mpTargets = []odahuflowv1alpha1.Target{
		{
			Name:           "test",
			ConnectionName: "test-conn",
		},
	}
	cfg *rest.Config
)

func TestMain(m *testing.M) {
	os.Exit(WrapperTestMain(m))
}

func WrapperTestMain(m *testing.M) int {
	_, cfgLocal, closeF, _, err := testenvs.SetupTestKube(filepath.Join("..", "..", "..", "config", "crds"))
	if err != nil {
		log.Println("Unable to setup kubernetes")
		return -1
	}
	defer func() {
		if err := closeF(); err != nil {
			log.Println("Unable to stop k8s test environment")
		}
	}()
	cfg = cfgLocal
	return m.Run()
}

type MPRepositorySuite struct {
	suite.Suite
	g             *GomegaWithT
	k8sClient     client.Client
	packK8SCLient packagingclient.Client
}

func generateMPResultCM() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      odahuflow.GeneratePackageResultCMName(mpID),
			Namespace: testNamespace,
		},
	}
}

func generateMP() *packaging.ModelPackaging {
	return &packaging.ModelPackaging{
		ID: mpID,
		Spec: packaging.ModelPackagingSpec{
			ArtifactName:    mpArtifactName,
			IntegrationName: mpType,
			Image:           mpImage,
			Arguments:       mpArguments,
			Targets:         mpTargets,
		},
	}
}

func (s *MPRepositorySuite) SetupSuite() {
	var err error
	s.k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		// If we get a panic that we have a test configuration problem
		panic(err)
	}

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		// If we get the panic that we have a test configuration problem
		panic(err)
	}

	// k8sConfig is nil because we use this client only for getting logs
	// we do not test this functionality in unit tests
	s.packK8SCLient = packagingclient.NewClient(testNamespace, testNamespace, k8sClient, nil)
}

func (s *MPRepositorySuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func (s *MPRepositorySuite) TearDownTest() {
	if err := s.packK8SCLient.DeleteModelPackaging(mpID); err != nil && !errors.IsNotFoundError(err) {
		// If we get the panic that we have a test configuration problem
		panic(err)
	}
}

func TestSuiteMP(t *testing.T) {
	suite.Run(t, new(MPRepositorySuite))
}

func (s *MPRepositorySuite) TestPackagingKubeClient() {
	created := generateMP()

	s.g.Expect(s.packK8SCLient.CreateModelPackaging(created)).NotTo(HaveOccurred())

	fetched, err := s.packK8SCLient.GetModelPackaging(mpID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(fetched.ID).To(Equal(created.ID))
	s.g.Expect(fetched.Spec).To(Equal(created.Spec))

	updated := fetched
	updated.Spec.Image = mpNewImage
	s.g.Expect(s.packK8SCLient.UpdateModelPackaging(updated)).NotTo(HaveOccurred())

	fetched, err = s.packK8SCLient.GetModelPackaging(mpID)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.g.Expect(fetched.ID).To(Equal(updated.ID))
	s.g.Expect(fetched.Spec).To(Equal(updated.Spec))
	s.g.Expect(fetched.Spec.Image).To(Equal(mpNewImage))

	s.g.Expect(s.packK8SCLient.DeleteModelPackaging(mpID)).NotTo(HaveOccurred())
	_, err = s.packK8SCLient.GetModelPackaging(mpID)
	s.g.Expect(err).To(HaveOccurred())
	s.g.Expect(errors.IsNotFoundError(err)).Should(BeTrue())
}

func (s *MPRepositorySuite) TestModelPackagingResult() {
	created := generateMP()
	s.g.Expect(s.packK8SCLient.CreateModelPackaging(created)).NotTo(HaveOccurred())

	resultConfigMap := generateMPResultCM()
	err := s.k8sClient.Create(context.TODO(), resultConfigMap)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), resultConfigMap)

	expectedPackagingResult := []odahuflowv1alpha1.ModelPackagingResult{
		{
			Name:  "test-name-1",
			Value: "test-value-1",
		},
		{
			Name:  "test-name-2",
			Value: "test-value-2",
		},
	}
	err = s.packK8SCLient.SaveModelPackagingResult(mpID, expectedPackagingResult)
	s.g.Expect(err).ShouldNot(HaveOccurred())

	results, err := s.packK8SCLient.GetModelPackagingResult(mpID)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(expectedPackagingResult).Should(Equal(results))
}

func (s *MPRepositorySuite) TestMPResultConfigMapNotFound() {
	created := generateMP()
	s.g.Expect(s.packK8SCLient.CreateModelPackaging(created)).NotTo(HaveOccurred())

	expectedPackagingResult := []odahuflowv1alpha1.ModelPackagingResult{
		{
			Name:  "test-name-1",
			Value: "test-value-1",
		},
		{
			Name:  "test-name-2",
			Value: "test-value-2",
		},
	}
	err := s.packK8SCLient.SaveModelPackagingResult(mpID, expectedPackagingResult)
	s.g.Expect(err).Should(HaveOccurred())
	s.g.Expect(err.Error()).Should(ContainSubstring("not found"))
}

func (s *MPRepositorySuite) TestEmptyModelPackagingResult() {
	created := generateMP()
	s.g.Expect(s.packK8SCLient.CreateModelPackaging(created)).NotTo(HaveOccurred())

	resultConfigMap := generateMPResultCM()
	err := s.k8sClient.Create(context.TODO(), resultConfigMap)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), resultConfigMap)

	expectedPackagingResult := []odahuflowv1alpha1.ModelPackagingResult{}
	err = s.packK8SCLient.SaveModelPackagingResult(mpID, expectedPackagingResult)
	s.g.Expect(err).ShouldNot(HaveOccurred())

	results, err := s.packK8SCLient.GetModelPackagingResult(mpID)
	s.g.Expect(err).ShouldNot(HaveOccurred())
	s.g.Expect(expectedPackagingResult).Should(Equal(results))
}

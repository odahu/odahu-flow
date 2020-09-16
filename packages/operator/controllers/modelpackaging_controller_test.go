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
package controllers_test

import (
	"context"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/controllers"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/stubclients"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sync"
	"testing"
	"time"
)

const (
	mpName                     = "test-mp"
	testPackagingIntegrationID = "ti"
	modelPackageImage          = "model-image:packager"
	integrationImage           = "model-image:packaging"
)

var (
	mpNamespacedName  = types.NamespacedName{Name: mpName, Namespace: testNamespace}
	mpExpectedRequest = reconcile.Request{
		NamespacedName: types.NamespacedName{Name: mpName, Namespace: testNamespace},
	}
)

type ModelPackagingControllerSuite struct {
	suite.Suite
	g *GomegaWithT

	k8sClient    client.Client
	k8sManager   manager.Manager
	stubPIClient stubclients.PIStubClient
	stopMgr      chan struct{}
	mgrStopped   *sync.WaitGroup
	requests     chan reconcile.Request
}

func (s *ModelPackagingControllerSuite) createPackagingIntegration() *packaging.PackagingIntegration {
	testPackagingIntegration := &packaging.PackagingIntegration{
		ID: testPackagingIntegrationID,
		Spec: packaging.PackagingIntegrationSpec{
			DefaultImage: integrationImage,
			Schema: packaging.Schema{
				Targets: []odahuflowv1alpha1.TargetSchema{
					{
						Name: "target-1",
						ConnectionTypes: []string{
							string(connection.S3Type),
							string(connection.GcsType),
							string(connection.AzureBlobType),
						},
						Required: false,
					},
				},
				Arguments: packaging.JsonSchema{
					Properties: []packaging.Property{
						{
							Name: "argument-1",
							Parameters: []packaging.Parameter{
								{
									Name:  "type",
									Value: "string",
								},
							},
						},
					},
					Required: []string{"argument-1"},
				},
			},
		},
	}

	return testPackagingIntegration
}

func (s *ModelPackagingControllerSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())

	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
	s.g.Expect(err).NotTo(HaveOccurred())
	s.k8sClient = mgr.GetClient()
	s.k8sManager = mgr
	s.stubPIClient = stubclients.NewPIStubClient()

	s.requests = make(chan reconcile.Request)

	s.stopMgr = make(chan struct{})
	s.mgrStopped = &sync.WaitGroup{}
	s.mgrStopped.Add(1)
	go func() {
		defer s.mgrStopped.Done()
		s.g.Expect(mgr.Start(s.stopMgr)).NotTo(HaveOccurred())
	}()

	if err := s.stubPIClient.CreatePackagingIntegration(s.createPackagingIntegration()); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ModelPackagingControllerSuite) initReconciler(packagingConfig config.ModelPackagingConfig) {
	packagingConfig.PackagingIntegrationNamespace = testNamespace
	packagingConfig.Namespace = testNamespace

	cfg := config.NewDefaultConfig()
	cfg.Packaging = packagingConfig

	reconciler := controllers.NewModelPackagingReconciler(s.k8sManager, *cfg, s.stubPIClient)
	rw := NewReconcilerWrapper(reconciler, s.requests)
	s.g.Expect(rw.SetupWithManager(s.k8sManager)).NotTo(HaveOccurred())
}

func (s *ModelPackagingControllerSuite) TearDownTest() {

	if err := s.stubPIClient.DeletePackagingIntegration(s.createPackagingIntegration().ID); err != nil {
		s.T().Fatal(err)
	}

	close(s.stopMgr)
	s.mgrStopped.Wait()
}

func TestModelPackagingControllerSuite(t *testing.T) {
	suite.Run(t, new(ModelPackagingControllerSuite))
}

// Node pool provided in packaging request, use it for tekton task
func (s *ModelPackagingControllerSuite) TestNodePool_Provided() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.NodePools = []config.NodePool{{NodeSelector: nodeSelector}}
	s.initReconciler(packagingConfig)

	mp := buildPackagingWithNodeSelector(nodeSelector)

	cleanF := s.createPackaging(mp)
	defer cleanF()
	tektonTask := s.getTektonPackagingTask(mp.Name)

	s.Assertions.Nil(tektonTask.Spec.PodTemplate.Affinity)
	s.Assertions.Equal(nodeSelector, tektonTask.Spec.PodTemplate.NodeSelector)
}

// Node pool not provided, build affinity for all CPU pools from config
func (s *ModelPackagingControllerSuite) TestNodePool_NotProvided_UseAffinity() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	nodeSelector1 := map[string]string{"node-key": "node-value", "another": "another-value"}
	nodeSelector2 := map[string]string{"node-key2": "node-value2"}
	packagingConfig.NodePools = []config.NodePool{{NodeSelector: nodeSelector1}, {NodeSelector: nodeSelector2}}
	s.initReconciler(packagingConfig)

	expectedNodeSelectorRequirement1 := []v1.NodeSelectorRequirement{
		{
			Key:      "node-key",
			Operator: v1.NodeSelectorOpIn,
			Values:   []string{"node-value"},
		},
		{
			Key:      "another",
			Operator: v1.NodeSelectorOpIn,
			Values:   []string{"another-value"},
		},
	}
	expectedNodeSelectorRequirement2 := []v1.NodeSelectorRequirement{{
		Key:      "node-key2",
		Operator: v1.NodeSelectorOpIn,
		Values:   []string{"node-value2"},
	}}

	mp := buildPackagingWithNodeSelector(nil)
	cleanF := s.createPackaging(mp)
	defer cleanF()
	tektonTask := s.getTektonPackagingTask(mp.Name)

	actualAffinity := tektonTask.Spec.PodTemplate.Affinity
	actualNodeSelectorTerms := actualAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	s.Assertions.Len(actualNodeSelectorTerms, 2)

	for i, expectedNodeSelectorRequirement := range [][]v1.NodeSelectorRequirement{
		expectedNodeSelectorRequirement1, expectedNodeSelectorRequirement2} {
		s.Assertions.ElementsMatch(expectedNodeSelectorRequirement, actualNodeSelectorTerms[i].MatchExpressions)
	}

	s.Assertions.Nil(tektonTask.Spec.PodTemplate.NodeSelector)
}

// Template tolerations tests
func (s *ModelPackagingControllerSuite) templateTestTolerations(input []v1.Toleration) {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Tolerations = input
	s.initReconciler(packagingConfig)

	mp := &odahuflowv1alpha1.ModelPackaging{
		ObjectMeta: metav1.ObjectMeta{Name: mpName, Namespace: testNamespace},
		Spec: odahuflowv1alpha1.ModelPackagingSpec{
			Type: testPackagingIntegrationID,
			Resources: &odahuflowv1alpha1.ResourceRequirements{
				Limits: &odahuflowv1alpha1.ResourceList{CPU: &testResValue},
			},
		},
	}

	cleanF := s.createPackaging(mp)
	defer cleanF()
	tektonTask := s.getTektonPackagingTask(mp.Name)
	s.Assertions.Equal(input, tektonTask.Spec.PodTemplate.Tolerations)
}

// Toleration is nil in config, expect nil toleration in tekton task
func (s *ModelPackagingControllerSuite) TestToleration_nil() {
	s.templateTestTolerations(nil)
}

// Single toleration in config, expect it in tekton task
func (s *ModelPackagingControllerSuite) TestToleration_Single() {
	s.templateTestTolerations([]v1.Toleration{{Key: "taint-key", Operator: v1.TolerationOpEqual, Value: "taint-val"}})
}

// Multiple tolerations in config, expect them in tekton task
func (s *ModelPackagingControllerSuite) TestToleration_Multiple() {
	s.templateTestTolerations([]v1.Toleration{
		{Key: "taint-key", Operator: v1.TolerationOpEqual, Value: "taint-val"},
		{Key: "taint-key", Operator: v1.TolerationOpEqual, Value: "taint-val"},
	})
}

func (s *ModelPackagingControllerSuite) TestPackagingStepConfiguration() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.ModelPackagerImage = modelPackageImage
	s.initReconciler(packagingConfig)

	packResources := &odahuflowv1alpha1.ResourceRequirements{
		Limits: &odahuflowv1alpha1.ResourceList{
			CPU: &testResValue,
			GPU: &testResValue,
		},
		Requests: &odahuflowv1alpha1.ResourceList{
			CPU: &testResValue,
			GPU: &testResValue,
		},
	}
	k8sPackagingResources, err := kubernetes.ConvertOdahuflowResourcesToK8s(packResources, config.NvidiaResourceName)
	s.g.Expect(err).Should(BeNil())

	mp := &odahuflowv1alpha1.ModelPackaging{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mpName,
			Namespace: testNamespace,
		},
		Spec: odahuflowv1alpha1.ModelPackagingSpec{
			Image:     integrationImage,
			Resources: packResources,
			Type:      testPackagingIntegrationID,
		},
	}

	err = s.k8sClient.Create(context.TODO(), mp)
	s.g.Expect(err).NotTo(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), mp)

	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(mpExpectedRequest)))
	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(mpExpectedRequest)))

	s.g.Expect(s.k8sClient.Get(context.TODO(), mpNamespacedName, mp)).ToNot(HaveOccurred())

	tr := &tektonv1beta1.TaskRun{}
	trKey := types.NamespacedName{Name: mp.Name, Namespace: testNamespace}
	s.g.Expect(s.k8sClient.Get(context.TODO(), trKey, tr)).ToNot(HaveOccurred())

	expectedHelperContainerResources := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    *k8sPackagingResources.Limits.Cpu(),
			v1.ResourceMemory: *utils.DefaultHelperLimits.Memory(),
		},
		Requests: v1.ResourceList{
			v1.ResourceMemory: resource.MustParse("0"),
			v1.ResourceCPU:    resource.MustParse("0"),
		},
	}

	for _, step := range tr.Spec.TaskSpec.Steps {
		switch step.Name {
		case odahuflow.PackagerSetupStep:
			s.g.Expect(step.Image).Should(Equal(modelPackageImage))
			s.g.Expect(step.Resources).Should(Equal(expectedHelperContainerResources))
		case odahuflow.PackagerPackageStep:
			s.g.Expect(step.Image).Should(Equal(integrationImage))
			s.g.Expect(step.Resources).Should(Equal(k8sPackagingResources))
		case odahuflow.PackagerResultStep:
			s.g.Expect(step.Image).Should(Equal(modelPackageImage))
			s.g.Expect(step.Resources).Should(Equal(expectedHelperContainerResources))
		default:
			s.T().Errorf("Unexpected step name: %s", step.Name)
		}
	}
}

func (s *ModelPackagingControllerSuite) TestPackagingTimeout() {
	packagingConfig := config.NewDefaultModelPackagingConfig()
	packagingConfig.Timeout = 3 * time.Hour
	s.initReconciler(packagingConfig)

	mp := &odahuflowv1alpha1.ModelPackaging{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mpName,
			Namespace: testNamespace,
		},
		Spec: odahuflowv1alpha1.ModelPackagingSpec{
			Type: testPackagingIntegrationID,
		},
	}

	err := s.k8sClient.Create(context.TODO(), mp)
	s.g.Expect(err).NotTo(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), mp)

	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(mpExpectedRequest)))
	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(mpExpectedRequest)))

	s.g.Expect(s.k8sClient.Get(context.TODO(), mpNamespacedName, mp)).ToNot(HaveOccurred())

	tr := &tektonv1beta1.TaskRun{}
	trKey := types.NamespacedName{Name: mp.Name, Namespace: testNamespace}
	s.g.Expect(s.k8sClient.Get(context.TODO(), trKey, tr)).ToNot(HaveOccurred())

	s.g.Expect(tr.Spec.Timeout.Duration).Should(Equal(time.Hour * 3))
}

func (s *ModelPackagingControllerSuite) createPackaging(mp *odahuflowv1alpha1.ModelPackaging) (
	cleanF func()) {
	err := s.k8sClient.Create(context.TODO(), mp)
	s.Assertions.Nil(err, "Failed to create packaging in K8s")

	s.Assertions.Eventually(
		func() bool {
			select {
			case r := <-s.requests:
				return r == mpExpectedRequest
			default:
				return false
			}
		},
		5*time.Second,
		10*time.Millisecond)

	s.Assertions.Nil(s.k8sClient.Get(context.TODO(), mpNamespacedName, mp))
	return func() { s.k8sClient.Delete(context.TODO(), mp) }
}

func (s *ModelPackagingControllerSuite) getTektonPackagingTask(packagingName string) *tektonv1beta1.TaskRun {
	tr := &tektonv1beta1.TaskRun{}
	trKey := types.NamespacedName{Name: packagingName, Namespace: testNamespace}
	err := s.k8sClient.Get(context.TODO(), trKey, tr)
	s.Assertions.Nil(err, "Tekton task retrieval failed")
	return tr
}

func buildPackagingWithNodeSelector(nodeSelector map[string]string) *odahuflowv1alpha1.ModelPackaging {
	return &odahuflowv1alpha1.ModelPackaging{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mpName,
			Namespace: testNamespace,
		},
		Spec: odahuflowv1alpha1.ModelPackagingSpec{
			Type:         testPackagingIntegrationID,
			NodeSelector: nodeSelector,
		},
	}
}

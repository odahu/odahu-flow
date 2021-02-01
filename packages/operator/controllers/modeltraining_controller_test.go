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
	"fmt"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/controllers"
	training_apis "github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
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
	"math/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sync"
	"testing"
	"time"
)

const (
	mtName                     = "test-mt"
	testToolchainIntegrationID = "ti"

	modelBuildImage = "model-image:builder"
	toolchainImage  = "model-image:toolchain"
)

var (
	gpuNodeSelector = map[string]string{"gpu-key": "gpu-value"}
	nodeSelector    = map[string]string{"node-key": "node-value"}

	testResValue             = "5"
	testToolchainIntegration = &training_apis.ToolchainIntegration{
		ID: testToolchainIntegrationID,
		Spec: odahuflowv1alpha1.ToolchainIntegrationSpec{
			DefaultImage: toolchainImage,
		},
	}
	validTraining = odahuflowv1alpha1.ModelTraining{
		ObjectMeta: metav1.ObjectMeta{Name: mtName, Namespace: testNamespace},
		Spec: odahuflowv1alpha1.ModelTrainingSpec{
			NodeSelector: nodeSelector,
			Toolchain:    testToolchainIntegrationID,
			Resources: &odahuflowv1alpha1.ResourceRequirements{
				Limits: &odahuflowv1alpha1.ResourceList{
					CPU: &testResValue,
				},
			},
		},
	}
)

type ModelTrainingControllerSuite struct {
	suite.Suite
	g            *GomegaWithT
	k8sClient    client.Client
	k8sManager   manager.Manager
	stubTIClient stubclients.TIStubClient
	stopMgr      chan struct{}
	mgrStopped   *sync.WaitGroup
	requests     chan reconcile.Request
}

func (s *ModelTrainingControllerSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())

	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
	s.g.Expect(err).NotTo(HaveOccurred())
	s.k8sClient = mgr.GetClient()
	s.k8sManager = mgr
	s.stubTIClient = stubclients.NewTIStubClient()

	s.requests = make(chan reconcile.Request, 1000)

	s.stopMgr = make(chan struct{})
	s.mgrStopped = &sync.WaitGroup{}
	s.mgrStopped.Add(1)
	go func() {
		defer s.mgrStopped.Done()
		s.g.Expect(mgr.Start(s.stopMgr)).NotTo(HaveOccurred())
	}()

	// Create the toolchain integration that will be used for a training.
	if err := s.stubTIClient.CreateToolchainIntegration(testToolchainIntegration); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ModelTrainingControllerSuite) initReconciler(trainingConfig config.ModelTrainingConfig) {
	trainingConfig.ToolchainIntegrationNamespace = testNamespace
	trainingConfig.Namespace = testNamespace

	cfg := config.NewDefaultConfig()
	cfg.Training = trainingConfig

	reconciler := controllers.NewModelTrainingReconciler(s.k8sManager, *cfg, s.stubTIClient)
	rw := NewReconcilerWrapper(reconciler, s.requests)
	s.g.Expect(rw.SetupWithManager(s.k8sManager)).NotTo(HaveOccurred())

}

func (s *ModelTrainingControllerSuite) TearDownTest() {

	if err := s.stubTIClient.DeleteToolchainIntegration(testToolchainIntegrationID); err != nil {
		s.T().Fatal(err)
	}

	close(s.stopMgr)
	s.mgrStopped.Wait()
}

func TestModelTrainingControllerSuite(t *testing.T) {
	suite.Run(t, new(ModelTrainingControllerSuite))
}

// Node pool provided in training request, use it for tekton task
func (s *ModelTrainingControllerSuite) TestTrainingReconcile_NodePoolProvided() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.NodePools = []config.NodePool{{NodeSelector: nodeSelector}}
	s.initReconciler(trainingConfig)

	mt := newValidTraining()
	mt.Spec.NodeSelector = nodeSelector

	cleanF := s.createTraining(mt)
	defer cleanF()
	tektonTask := s.getTektonTrainingTask(mt)

	s.Assertions.Nil(tektonTask.Spec.PodTemplate.Affinity)
	s.Assertions.Equal(nodeSelector, tektonTask.Spec.PodTemplate.NodeSelector)
}

// Node pool not provided, build affinity for all CPU pools from config
func (s *ModelTrainingControllerSuite) TestTrainingReconcile_NodePoolNotProvided_UseAffinity() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	nodeSelector1 := map[string]string{"node-key": "node-value", "another": "another-value"}
	nodeSelector2 := map[string]string{"node-key2": "node-value2"}
	trainingConfig.NodePools = []config.NodePool{{NodeSelector: nodeSelector1}, {NodeSelector: nodeSelector2}}
	s.initReconciler(trainingConfig)

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

	mt := newValidTraining()
	mt.Spec.NodeSelector = nil
	cleanF := s.createTraining(mt)
	defer cleanF()
	tektonTask := s.getTektonTrainingTask(mt)

	actualAffinity := tektonTask.Spec.PodTemplate.Affinity
	actualNodeSelectorTerms := actualAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	s.Assertions.Len(actualNodeSelectorTerms, 2)

	for i, expectedNodeSelectorRequirement := range [][]v1.NodeSelectorRequirement{
		expectedNodeSelectorRequirement1, expectedNodeSelectorRequirement2} {
		s.Assertions.ElementsMatch(expectedNodeSelectorRequirement, actualNodeSelectorTerms[i].MatchExpressions)
	}

	s.Assertions.Nil(tektonTask.Spec.PodTemplate.NodeSelector)
}

// Node pool provided in GPU training request, use it for tekton task
func (s *ModelTrainingControllerSuite) TestTrainingReconcile_GPU_NodePoolProvided() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.GPUNodePools = []config.NodePool{{NodeSelector: gpuNodeSelector}}
	s.initReconciler(trainingConfig)

	mt := newValidTraining()
	mt.Spec.Resources.Limits = &odahuflowv1alpha1.ResourceList{GPU: &testResValue}
	mt.Spec.NodeSelector = gpuNodeSelector

	cleanF := s.createTraining(mt)
	defer cleanF()
	tektonTask := s.getTektonTrainingTask(mt)

	s.Assertions.Nil(tektonTask.Spec.PodTemplate.Affinity)
	s.Assertions.Equal(gpuNodeSelector, tektonTask.Spec.PodTemplate.NodeSelector)
}

// Node pool not provided for GPU training, build affinity for all GPU pools from config
func (s *ModelTrainingControllerSuite) TestTrainingReconcile_GPU_NodePoolNotProvided_UseAffinity() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	nodeSelector1 := map[string]string{"gpu-node-key": "gpu-node-value", "gpu-another": "gpu-another-value"}
	nodeSelector2 := map[string]string{"gpu-node-key2": "gpu-node-value2"}
	trainingConfig.GPUNodePools = []config.NodePool{{NodeSelector: nodeSelector1}, {NodeSelector: nodeSelector2}}
	s.initReconciler(trainingConfig)

	expectedNodeSelectorRequirement1 := []v1.NodeSelectorRequirement{
		{
			Key:      "gpu-node-key",
			Operator: v1.NodeSelectorOpIn,
			Values:   []string{"gpu-node-value"},
		},
		{
			Key:      "gpu-another",
			Operator: v1.NodeSelectorOpIn,
			Values:   []string{"gpu-another-value"},
		},
	}
	expectedNodeSelectorRequirement2 := []v1.NodeSelectorRequirement{{
		Key:      "gpu-node-key2",
		Operator: v1.NodeSelectorOpIn,
		Values:   []string{"gpu-node-value2"},
	}}

	mt := newValidTraining()
	mt.Spec.Resources.Limits = &odahuflowv1alpha1.ResourceList{GPU: &testResValue}
	mt.Spec.NodeSelector = nil

	cleanF := s.createTraining(mt)
	defer cleanF()
	tektonTask := s.getTektonTrainingTask(mt)

	actualAffinity := tektonTask.Spec.PodTemplate.Affinity
	actualNodeSelectorTerms := actualAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	s.Assertions.Len(actualNodeSelectorTerms, 2)

	for i, expectedNodeSelectorRequirement := range [][]v1.NodeSelectorRequirement{
		expectedNodeSelectorRequirement1, expectedNodeSelectorRequirement2} {
		s.Assertions.ElementsMatch(expectedNodeSelectorRequirement, actualNodeSelectorTerms[i].MatchExpressions)
	}

	s.Assertions.Nil(tektonTask.Spec.PodTemplate.NodeSelector)
}

// Template tolerations tests (CPU)
func (s *ModelTrainingControllerSuite) templateTestTolerations(input []v1.Toleration) {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Tolerations = input
	s.initReconciler(trainingConfig)

	mt := newValidTraining()
	cleanF := s.createTraining(mt)
	defer cleanF()
	tektonTask := s.getTektonTrainingTask(mt)
	s.Assertions.Equal(input, tektonTask.Spec.PodTemplate.Tolerations)
}

// Toleration is nil in config, expect nil toleration in tekton task
func (s *ModelTrainingControllerSuite) TestTrainingReconcile_Tolerations_nil() {
	s.templateTestTolerations(nil)
}

// Single toleration in config, expect it in tekton task
func (s *ModelTrainingControllerSuite) TestTrainingReconcile_Tolerations_Single() {
	s.templateTestTolerations([]v1.Toleration{{Key: "taint-key", Operator: v1.TolerationOpEqual, Value: "taint-val"}})
}

// Multiple tolerations in config, expect them in tekton task
func (s *ModelTrainingControllerSuite) TestTrainingReconcile_Tolerations_Multiple() {
	s.templateTestTolerations([]v1.Toleration{
		{Key: "taint-key", Operator: v1.TolerationOpEqual, Value: "taint-val"},
		{Key: "taint-key", Operator: v1.TolerationOpEqual, Value: "taint-val"},
	})
}

// Template tolerations tests (GPU)
func (s *ModelTrainingControllerSuite) templateTestGPUTolerations(input []v1.Toleration) {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.GPUTolerations = input
	s.initReconciler(trainingConfig)

	mt := newValidTraining()
	mt.Spec.Resources.Limits = &odahuflowv1alpha1.ResourceList{GPU: &testResValue}

	cleanF := s.createTraining(mt)
	defer cleanF()
	tektonTask := s.getTektonTrainingTask(mt)
	s.Assertions.Equal(input, tektonTask.Spec.PodTemplate.Tolerations)
}

// GPU Tolerations is nil in config, expect nil toleration in tekton task
func (s *ModelTrainingControllerSuite) TestTrainingReconcile_GPU_Tolerations_nil() {
	s.templateTestGPUTolerations(nil)
}

// Single GPU toleration in config, expect it in tekton task
func (s *ModelTrainingControllerSuite) TestTrainingReconcile_GPU_Tolerations_Single() {
	s.templateTestGPUTolerations([]v1.Toleration{{Key: "taint-key", Operator: v1.TolerationOpEqual, Value: "taint-val"}})
}

// Multiple GPU tolerations in config, expect them in tekton task
func (s *ModelTrainingControllerSuite) TestTrainingReconcile_GPU_Tolerations_Multiple() {
	s.templateTestGPUTolerations([]v1.Toleration{
		{Key: "taint-key", Operator: v1.TolerationOpEqual, Value: "taint-val"},
		{Key: "taint-key", Operator: v1.TolerationOpEqual, Value: "taint-val"},
	})
}

func (s *ModelTrainingControllerSuite) TestTrainingStepConfiguration() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.NodePools = []config.NodePool{{NodeSelector: nodeSelector}}
	trainingConfig.GPUNodePools = []config.NodePool{{NodeSelector: gpuNodeSelector}}
	trainingConfig.ModelTrainerImage = modelBuildImage
	s.initReconciler(trainingConfig)

	trainResources := &odahuflowv1alpha1.ResourceRequirements{
		Limits: &odahuflowv1alpha1.ResourceList{
			CPU: &testResValue,
			GPU: &testResValue,
		},
		Requests: &odahuflowv1alpha1.ResourceList{
			CPU: &testResValue,
			GPU: &testResValue,
		},
	}
	k8sTrainerResources, err := kubernetes.ConvertOdahuflowResourcesToK8s(trainResources, config.NvidiaResourceName)
	s.g.Expect(err).Should(BeNil())

	mt := newValidTraining()
	mt.Spec = odahuflowv1alpha1.ModelTrainingSpec{
		Image:     toolchainImage,
		Resources: trainResources,
		Toolchain: testToolchainIntegrationID,
	}

	err = s.k8sClient.Create(context.TODO(), mt)
	s.g.Expect(err).NotTo(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), mt)

	expectedTrainingRequest := reconcile.Request{NamespacedName: types.NamespacedName{Name: mt.Name, Namespace: mt.Namespace},
	}
	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedTrainingRequest)))

	mtNamespacedName := types.NamespacedName{Name: mt.Name, Namespace: mt.Namespace}
	s.g.Expect(s.k8sClient.Get(context.TODO(), mtNamespacedName, mt)).ToNot(HaveOccurred())

	tr := &tektonv1beta1.TaskRun{}
	trKey := types.NamespacedName{Name: mt.Name, Namespace: testNamespace}
	s.g.Expect(s.k8sClient.Get(context.TODO(), trKey, tr)).ToNot(HaveOccurred())

	expectedHelperContainerResources := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    *k8sTrainerResources.Limits.Cpu(),
			v1.ResourceMemory: *utils.DefaultHelperLimits.Memory(),
		},
		Requests: v1.ResourceList{
			v1.ResourceMemory: resource.MustParse("0"),
			v1.ResourceCPU:    resource.MustParse("0"),
		},
	}

	for _, step := range tr.Spec.TaskSpec.Steps {
		switch step.Name {
		case odahuflow.TrainerSetupStep:
			s.g.Expect(step.Image).Should(Equal(modelBuildImage))
			s.g.Expect(step.Resources).Should(Equal(expectedHelperContainerResources))
		case odahuflow.TrainerTrainStep:
			s.g.Expect(step.Image).Should(Equal(toolchainImage))
			s.g.Expect(step.Resources).Should(Equal(k8sTrainerResources))
		case odahuflow.TrainerValidationStep:
			s.g.Expect(step.Image).Should(Equal(toolchainImage))
			s.g.Expect(step.Resources).Should(Equal(expectedHelperContainerResources))
		case odahuflow.TrainerResultStep:
			s.g.Expect(step.Image).Should(Equal(modelBuildImage))
			s.g.Expect(step.Resources).Should(Equal(expectedHelperContainerResources))
		default:
			s.T().Errorf("Unexpected step name: %s", step.Name)
		}
	}
}

func (s *ModelTrainingControllerSuite) TestTrainingTimeout() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.Timeout = 3 * time.Hour
	s.initReconciler(trainingConfig)

	mt := &odahuflowv1alpha1.ModelTraining{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mtName,
			Namespace: testNamespace,
		},
		Spec: odahuflowv1alpha1.ModelTrainingSpec{
			Toolchain: testToolchainIntegrationID,
		},
	}

	err := s.k8sClient.Create(context.TODO(), mt)
	s.g.Expect(err).NotTo(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), mt)

	mtNamespacedName := types.NamespacedName{Name: mt.Name, Namespace: mt.Namespace}
	expectedTrainingRequest := reconcile.Request{NamespacedName: mtNamespacedName}
	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedTrainingRequest)))

	s.g.Expect(s.k8sClient.Get(context.TODO(), mtNamespacedName, mt)).ToNot(HaveOccurred())

	tr := &tektonv1beta1.TaskRun{}
	trKey := types.NamespacedName{Name: mt.Name, Namespace: testNamespace}
	s.g.Expect(s.k8sClient.Get(context.TODO(), trKey, tr)).ToNot(HaveOccurred())

	s.g.Expect(tr.Spec.Timeout.Duration).Should(Equal(time.Hour * 3))
}

func (s *ModelTrainingControllerSuite) TestTrainingEnvs() {
	const (
		trainingEnvKey    = "training-env-key"
		trainingEnvValue  = "training-env-value"
		toolchainEnvKey   = "toolchain-env-key"
		toolchainEnvValue = "toolchain-env-value"
	)

	s.initReconciler(config.NewDefaultModelTrainingConfig())

	ti := &training_apis.ToolchainIntegration{
		ID: testToolchainIntegrationID,
		Spec: odahuflowv1alpha1.ToolchainIntegrationSpec{
			DefaultImage: toolchainImage,
			AdditionalEnvironments: map[string]string{
				toolchainEnvKey: toolchainEnvValue,
			},
		},
	}

	mt := &odahuflowv1alpha1.ModelTraining{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mtName,
			Namespace: testNamespace,
		},
		Spec: odahuflowv1alpha1.ModelTrainingSpec{
			Toolchain: testToolchainIntegrationID,
			CustomEnvs: []odahuflowv1alpha1.EnvironmentVariable{
				{Name: trainingEnvKey, Value: trainingEnvValue},
			},
		},
	}

	// Recreate the toolchain integration
	err := s.stubTIClient.DeleteToolchainIntegration(testToolchainIntegrationID)
	s.g.Expect(err).NotTo(HaveOccurred())
	err = s.stubTIClient.CreateToolchainIntegration(ti)
	s.g.Expect(err).NotTo(HaveOccurred())

	err = s.k8sClient.Create(context.TODO(), mt)
	s.g.Expect(err).NotTo(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), mt)

	mtNamespacedName := types.NamespacedName{Name: mt.Name, Namespace: mt.Namespace}
	expectedTrainingRequest := reconcile.Request{NamespacedName: mtNamespacedName}
	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedTrainingRequest)))

	s.g.Expect(s.k8sClient.Get(context.TODO(), mtNamespacedName, mt)).ToNot(HaveOccurred())

	tr := &tektonv1beta1.TaskRun{}
	trKey := types.NamespacedName{Name: mt.Name, Namespace: testNamespace}
	s.g.Eventually(func() error { return s.k8sClient.Get(context.TODO(), trKey, tr) }).ShouldNot(HaveOccurred())

	// second container is the training container
	envs := tr.Spec.TaskSpec.Steps[1].Env
	// first envs must be toolchains envs, then training envs
	s.g.Expect(envs).Should(Equal([]v1.EnvVar{
		{Name: toolchainEnvKey, Value: toolchainEnvValue},
		{Name: trainingEnvKey, Value: trainingEnvValue},
	}))
}

// Test utilities

func (s *ModelTrainingControllerSuite) createTraining(training *odahuflowv1alpha1.ModelTraining) (cleanF func()) {
	err := s.k8sClient.Create(context.TODO(), training)
	s.Assertions.Nil(err, "Failed to create training in K8s")

	mtNamespacedName := types.NamespacedName{Name: training.Name, Namespace: training.Namespace}
	s.Assertions.Eventually(
		func() bool { return s.k8sClient.Get(context.TODO(), mtNamespacedName, training) == nil },
		10*time.Second,
		10*time.Millisecond)

	return func() { s.k8sClient.Delete(context.TODO(), training) }
}

func (s *ModelTrainingControllerSuite) getTektonTrainingTask(mt *odahuflowv1alpha1.ModelTraining) *tektonv1beta1.TaskRun {
	tr := &tektonv1beta1.TaskRun{}
	trKey := types.NamespacedName{Name: mt.Name, Namespace: mt.Namespace}
	s.Assertions.Eventually(
		func() bool { return s.k8sClient.Get(context.TODO(), trKey, tr) == nil },
		10*time.Second,
		10*time.Millisecond,
		"Task run not found!")
	return tr
}

// Returns validTraining with random Name to avoid collisions when running in parallel
func newValidTraining() *odahuflowv1alpha1.ModelTraining {
	mt := validTraining.DeepCopy()
	mt.Name = fmt.Sprintf("training-%d", rand.Int()) //nolint
	return mt
}

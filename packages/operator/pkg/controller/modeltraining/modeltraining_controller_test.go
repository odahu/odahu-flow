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
package modeltraining

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/pkg/apis/odahuflow/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	tektonv1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tektonschema "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/scheme"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sync"
	"testing"
	"time"
)

const (
	mtName                     = "test-mt"
	timeout                    = time.Second * 5
	testNamespace              = "default"
	testToolchainIntegrationID = "ti"

	gpuTolerationKey      = "gpu-key"
	gpuTolerationValue    = "gpu-value"
	gpuTolerationOperator = "gpu-operator"
	gpuTolerationEffect   = "gpu-effect"

	tolerationKey      = "key"
	tolerationValue    = "value"
	tolerationOperator = "operator"
	tolerationEffect   = "effect"

	modelBuildImage = "model-image:builder"
	toolchainImage  = "model-image:toolchain"
)

var (
	gpuNodeSelector = map[string]string{"gpu-key": "gpu-value"}
	nodeSelector    = map[string]string{"node-key": "node-value"}

	expectedRequest = reconcile.Request{
		NamespacedName: types.NamespacedName{Name: mtName, Namespace: testNamespace},
	}
	mtNamespacedName         = types.NamespacedName{Name: mtName, Namespace: testNamespace}
	testResValue             = "5"
	emptyValue               = ""
	testToolchainIntegration = &odahuflowv1alpha1.ToolchainIntegration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testToolchainIntegrationID,
			Namespace: testNamespace,
		},
		Spec: odahuflowv1alpha1.ToolchainIntegrationSpec{
			DefaultImage: toolchainImage,
		},
	}
	toleration = map[string]string{
		config.TolerationKey:      tolerationKey,
		config.TolerationValue:    tolerationValue,
		config.TolerationOperator: tolerationOperator,
		config.TolerationEffect:   tolerationEffect,
	}
	gpuToleration = map[string]string{
		config.TolerationKey:      gpuTolerationKey,
		config.TolerationValue:    gpuTolerationValue,
		config.TolerationOperator: gpuTolerationOperator,
		config.TolerationEffect:   gpuTolerationEffect,
	}
)

type ModelTrainingControllerSuite struct {
	suite.Suite
	g               *GomegaWithT
	testEnvironment *envtest.Environment
	cfg             *rest.Config

	k8sClient  client.Client
	k8sManager manager.Manager
	stopMgr    chan struct{}
	mgrStopped *sync.WaitGroup
	requests   chan reconcile.Request
}

func (s *ModelTrainingControllerSuite) SetupSuite() {
	s.testEnvironment = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "config", "crds"),
			filepath.Join("..", "..", "..", "hack", "tests", "thirdparty_crds"),
		},
	}
	err := apis.AddToScheme(scheme.Scheme)
	if err != nil {
		s.T().Fatal(err)
	}

	if err := tektonschema.AddToScheme(scheme.Scheme); err != nil {
		s.T().Fatal(err)
	}

	if s.cfg, err = s.testEnvironment.Start(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ModelTrainingControllerSuite) TearDownSuite() {
	if err := s.testEnvironment.Stop(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ModelTrainingControllerSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())

	mgr, err := manager.New(s.cfg, manager.Options{})
	s.g.Expect(err).NotTo(HaveOccurred())
	s.k8sClient = mgr.GetClient()
	s.k8sManager = mgr

	s.requests = make(chan reconcile.Request)

	s.stopMgr = make(chan struct{})
	s.mgrStopped = &sync.WaitGroup{}
	s.mgrStopped.Add(1)
	go func() {
		defer s.mgrStopped.Done()
		s.g.Expect(mgr.Start(s.stopMgr)).NotTo(HaveOccurred())
	}()

	// Create the toolchain integration that will be used for a training.
	if err := s.k8sClient.Create(context.TODO(), testToolchainIntegration); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ModelTrainingControllerSuite) initReconciler(trainingConfig config.ModelTrainingConfig) {
	trainingConfig.ToolchainIntegrationNamespace = testNamespace
	trainingConfig.Namespace = testNamespace

	mpReconciler := newReconciler(s.k8sManager, trainingConfig, config.NewDefaultOperatorConfig(), config.NewDefaultCommonConfig(), config.NvidiaResourceName)
	recFn := reconcile.Func(func(req reconcile.Request) (reconcile.Result, error) {
		result, err := mpReconciler.Reconcile(req)
		s.requests <- req
		return result, err
	})

	s.g.Expect(add(s.k8sManager, recFn)).NotTo(HaveOccurred())
}

func (s *ModelTrainingControllerSuite) TearDownTest() {
	if err := s.k8sClient.Delete(context.TODO(), testToolchainIntegration); err != nil {
		s.T().Fatal(err)
	}

	testToolchainIntegration.ResourceVersion = ""

	close(s.stopMgr)
	s.mgrStopped.Wait()
}

func TestModelTrainingControllerSuite(t *testing.T) {
	suite.Run(t, new(ModelTrainingControllerSuite))
}

func (s *ModelTrainingControllerSuite) templateNodeSelectorTest(
	mtResources *odahuflowv1alpha1.ResourceRequirements,
	expectedNodeSelector map[string]string,
	expectedToleration []v1.Toleration,
) {
	mt := &odahuflowv1alpha1.ModelTraining{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mtName,
			Namespace: testNamespace,
		},
		Spec: odahuflowv1alpha1.ModelTrainingSpec{
			Resources: mtResources,
			Toolchain: testToolchainIntegrationID,
		},
	}

	err := s.k8sClient.Create(context.TODO(), mt)
	s.g.Expect(err).NotTo(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), mt)

	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedRequest)))
	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedRequest)))

	s.g.Expect(s.k8sClient.Get(context.TODO(), mtNamespacedName, mt)).ToNot(HaveOccurred())

	tr := &tektonv1alpha1.TaskRun{}
	trKey := types.NamespacedName{Name: mt.Name, Namespace: testNamespace}
	s.g.Expect(s.k8sClient.Get(context.TODO(), trKey, tr)).ToNot(HaveOccurred())

	s.g.Expect(tr.Spec.PodTemplate.NodeSelector).Should(Equal(expectedNodeSelector))
	s.g.Expect(tr.Spec.PodTemplate.Tolerations).Should(Equal(expectedToleration))
}

func (s *ModelTrainingControllerSuite) TestEmptyGPUNodePools() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	s.initReconciler(trainingConfig)

	s.templateNodeSelectorTest(
		&odahuflowv1alpha1.ResourceRequirements{
			Limits: &odahuflowv1alpha1.ResourceList{
				GPU: &testResValue,
			},
		},
		nil,
		nil,
	)
}

func (s *ModelTrainingControllerSuite) TestEmptyGPUValue() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.GPUNodeSelector = map[string]string{}
	trainingConfig.GPUToleration = map[string]string{}
	s.initReconciler(trainingConfig)

	s.templateNodeSelectorTest(
		&odahuflowv1alpha1.ResourceRequirements{
			Limits: &odahuflowv1alpha1.ResourceList{
				GPU: &emptyValue,
			},
		},
		nil,
		nil,
	)
}

func (s *ModelTrainingControllerSuite) TestGPUNodePools() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.GPUNodeSelector = gpuNodeSelector
	trainingConfig.GPUToleration = gpuToleration
	s.initReconciler(trainingConfig)

	s.templateNodeSelectorTest(
		&odahuflowv1alpha1.ResourceRequirements{
			Limits: &odahuflowv1alpha1.ResourceList{
				GPU: &testResValue,
			},
		},
		gpuNodeSelector,
		[]v1.Toleration{{
			Key:      gpuTolerationKey,
			Operator: gpuTolerationOperator,
			Value:    gpuTolerationValue,
			Effect:   gpuTolerationEffect,
		}},
	)
}

func (s *ModelTrainingControllerSuite) TestOnlyGPUNodeSelector() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.GPUNodeSelector = gpuNodeSelector
	s.initReconciler(trainingConfig)

	s.templateNodeSelectorTest(
		&odahuflowv1alpha1.ResourceRequirements{
			Limits: &odahuflowv1alpha1.ResourceList{
				GPU: &testResValue,
			},
		},
		gpuNodeSelector,
		nil,
	)
}

func (s *ModelTrainingControllerSuite) TestEmptyNodePools() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	s.initReconciler(trainingConfig)

	s.templateNodeSelectorTest(
		&odahuflowv1alpha1.ResourceRequirements{
			Limits: &odahuflowv1alpha1.ResourceList{
				CPU: &testResValue,
			},
		},
		nil,
		nil,
	)
}

func (s *ModelTrainingControllerSuite) TestNodePools() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.NodeSelector = nodeSelector
	trainingConfig.Toleration = toleration
	s.initReconciler(trainingConfig)

	s.templateNodeSelectorTest(
		&odahuflowv1alpha1.ResourceRequirements{
			Limits: &odahuflowv1alpha1.ResourceList{
				CPU: &testResValue,
			},
		},
		nodeSelector,
		[]v1.Toleration{{
			Key:      tolerationKey,
			Operator: tolerationOperator,
			Value:    tolerationValue,
			Effect:   tolerationEffect,
		}},
	)
}

func (s *ModelTrainingControllerSuite) TestOnlyNodeSelector() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.NodeSelector = nodeSelector
	s.initReconciler(trainingConfig)

	s.templateNodeSelectorTest(
		&odahuflowv1alpha1.ResourceRequirements{
			Limits: &odahuflowv1alpha1.ResourceList{
				CPU: &testResValue,
			},
		},
		nodeSelector,
		nil,
	)
}

// If GPU and CPU resources setup on a one training, GPU node selector must be used
func (s *ModelTrainingControllerSuite) TestGPUandCPUNodePools() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.NodeSelector = nodeSelector
	trainingConfig.GPUNodeSelector = gpuNodeSelector
	trainingConfig.Toleration = toleration
	trainingConfig.GPUToleration = gpuToleration
	s.initReconciler(trainingConfig)

	s.templateNodeSelectorTest(
		&odahuflowv1alpha1.ResourceRequirements{
			Limits: &odahuflowv1alpha1.ResourceList{
				GPU: &testResValue,
				CPU: &testResValue,
			},
		},
		gpuNodeSelector,
		[]v1.Toleration{{
			Key:      gpuTolerationKey,
			Operator: gpuTolerationOperator,
			Value:    gpuTolerationValue,
			Effect:   gpuTolerationEffect,
		}},
	)
}

func (s *ModelTrainingControllerSuite) TestTrainingStepConfiguration() {
	trainingConfig := config.NewDefaultModelTrainingConfig()
	trainingConfig.NodeSelector = nodeSelector
	trainingConfig.GPUNodeSelector = gpuNodeSelector
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

	mt := &odahuflowv1alpha1.ModelTraining{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mtName,
			Namespace: testNamespace,
		},
		Spec: odahuflowv1alpha1.ModelTrainingSpec{
			Image:     toolchainImage,
			Resources: trainResources,
			Toolchain: testToolchainIntegrationID,
		},
	}

	err = s.k8sClient.Create(context.TODO(), mt)
	s.g.Expect(err).NotTo(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), mt)

	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedRequest)))
	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedRequest)))

	s.g.Expect(s.k8sClient.Get(context.TODO(), mtNamespacedName, mt)).ToNot(HaveOccurred())

	tr := &tektonv1alpha1.TaskRun{}
	trKey := types.NamespacedName{Name: mt.Name, Namespace: testNamespace}
	s.g.Expect(s.k8sClient.Get(context.TODO(), trKey, tr)).ToNot(HaveOccurred())

	expectedHelperContainerResources := corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    *k8sTrainerResources.Limits.Cpu(),
			corev1.ResourceMemory: *utils.DefaultHelperLimits.Memory(),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("0"),
			corev1.ResourceCPU:    resource.MustParse("0"),
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
			s.T().Errorf("Unexpected spep name: %s", step.Name)
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

	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedRequest)))
	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedRequest)))

	s.g.Expect(s.k8sClient.Get(context.TODO(), mtNamespacedName, mt)).ToNot(HaveOccurred())

	tr := &tektonv1alpha1.TaskRun{}
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

	ti := &odahuflowv1alpha1.ToolchainIntegration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testToolchainIntegrationID,
			Namespace: testNamespace,
		},
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
	err := s.k8sClient.Delete(context.TODO(), ti)
	s.g.Expect(err).NotTo(HaveOccurred())
	err = s.k8sClient.Create(context.TODO(), ti)
	s.g.Expect(err).NotTo(HaveOccurred())

	err = s.k8sClient.Create(context.TODO(), mt)
	s.g.Expect(err).NotTo(HaveOccurred())
	defer s.k8sClient.Delete(context.TODO(), mt)

	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedRequest)))
	s.g.Eventually(s.requests, timeout).Should(Receive(Equal(expectedRequest)))

	s.g.Expect(s.k8sClient.Get(context.TODO(), mtNamespacedName, mt)).ToNot(HaveOccurred())

	tr := &tektonv1alpha1.TaskRun{}
	trKey := types.NamespacedName{Name: mt.Name, Namespace: testNamespace}
	s.g.Expect(s.k8sClient.Get(context.TODO(), trKey, tr)).ToNot(HaveOccurred())

	// second container is the training container
	envs := tr.Spec.TaskSpec.Steps[1].Env
	// first envs must be toolchains envs, then training envs
	s.g.Expect(envs).Should(Equal([]corev1.EnvVar{
		{Name: toolchainEnvKey, Value: toolchainEnvValue},
		{Name: trainingEnvKey, Value: trainingEnvValue},
	}))
}

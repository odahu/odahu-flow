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
	. "github.com/odahu/odahu-flow/packages/operator/controllers"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	"math/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	image                   = "test/image:123"
	mdName                  = "test-md"
	mdMinReplicas           = int32(1)
	mdMaxReplicas           = int32(2)
	mdReadinessDelay        = int32(33)
	mdLivenessDelay         = int32(44)
	mdNamespace             = "default"
	mdImagePullConnectionID = ""
	validDeployment         = odahuflowv1alpha1.ModelDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: mdName, Namespace: mdNamespace},
		Spec: odahuflowv1alpha1.ModelDeploymentSpec{
			Image:                      image,
			Predictor:                  odahuflow.OdahuMLServer.ID,
			MinReplicas:                &mdMinReplicas,
			MaxReplicas:                &mdMaxReplicas,
			ReadinessProbeInitialDelay: &mdReadinessDelay,
			LivenessProbeInitialDelay:  &mdLivenessDelay,
			Resources:                  mdResources,
			ImagePullConnectionID:      &mdImagePullConnectionID,
			NodeSelector:               nodeSelector,
		},
	}
	reqMem      = "128Mi"
	reqCPU      = "125m"
	limMem      = "256Mi"
	mdResources = &odahuflowv1alpha1.ResourceRequirements{
		Limits: &odahuflowv1alpha1.ResourceList{
			CPU:    nil,
			Memory: &limMem,
		},
		Requests: &odahuflowv1alpha1.ResourceList{
			CPU:    &reqCPU,
			Memory: &reqMem,
		},
	}
)

type ModelDeploymentControllerSuite struct {
	suite.Suite
	k8sClient  client.Client
	k8sManager manager.Manager
	requests   chan reconcile.Request
	stopMgr    chan struct{}
	mgrStopped *sync.WaitGroup
}

func TestModelDeploymentControllerSuite(t *testing.T) {
	suite.Run(t, new(ModelDeploymentControllerSuite))
}

func (s *ModelDeploymentControllerSuite) SetupTest() {
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
	s.Assertions.NoError(err)

	s.k8sClient = mgr.GetClient()
	s.k8sManager = mgr

	s.requests = make(chan reconcile.Request, 1000)

	s.stopMgr = make(chan struct{})
	s.mgrStopped = &sync.WaitGroup{}
	s.mgrStopped.Add(1)
	go func() {
		defer s.mgrStopped.Done()
		s.Assertions.NoError(mgr.Start(s.stopMgr))
	}()
}

func (s *ModelDeploymentControllerSuite) initReconciler(deploymentConfig config.ModelDeploymentConfig) {
	cfg := config.NewDefaultConfig()
	cfg.Deployment = deploymentConfig

	reconciler := NewModelDeploymentReconciler(s.k8sManager, *cfg)
	rw := NewReconcilerWrapper(reconciler, s.requests)
	s.Assertions.NoError(rw.SetupWithManager(s.k8sManager))
}

func (s *ModelDeploymentControllerSuite) TearDownTest() {
	close(s.stopMgr)
	s.mgrStopped.Wait()
}

// TODO: break one super-test into separate tests
func (s *ModelDeploymentControllerSuite) TestReconcile() {
	s.initReconciler(config.NewDefaultModelDeploymentConfig())

	md := &odahuflowv1alpha1.ModelDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: mdName, Namespace: mdNamespace},
		Spec: odahuflowv1alpha1.ModelDeploymentSpec{
			Image:                      image,
			Predictor:                  odahuflow.OdahuMLServer.ID,
			MinReplicas:                &mdMinReplicas,
			MaxReplicas:                &mdMaxReplicas,
			ReadinessProbeInitialDelay: &mdReadinessDelay,
			LivenessProbeInitialDelay:  &mdLivenessDelay,
			Resources:                  mdResources,
			ImagePullConnectionID:      &mdImagePullConnectionID,
		},
	}

	cleanF := s.createDeployment(md)
	defer cleanF()
	knativeConfiguration := s.getKnativeConfiguration(md)

	configurationAnnotations := knativeConfiguration.Spec.Template.ObjectMeta.Annotations
	s.Assertions.Len(configurationAnnotations, 6)

	s.Assertions.Contains(configurationAnnotations, KnativeMinReplicasKey)
	s.Assertions.Equal(configurationAnnotations[KnativeMinReplicasKey], strconv.Itoa(int(mdMinReplicas)))

	s.Assertions.Contains(configurationAnnotations, KnativeMaxReplicasKey)
	s.Assertions.Equal(configurationAnnotations[KnativeMaxReplicasKey], strconv.Itoa(int(mdMaxReplicas)))

	s.Assertions.Contains(configurationAnnotations, KnativeAutoscalingTargetKey)
	s.Assertions.Equal(configurationAnnotations[KnativeAutoscalingTargetKey], KnativeAutoscalingTargetDefaultValue)

	s.Assertions.Contains(configurationAnnotations, KnativeAutoscalingClass)
	s.Assertions.Equal(configurationAnnotations[KnativeAutoscalingClass], DefaultKnativeAutoscalingClass)

	s.Assertions.Contains(configurationAnnotations, KnativeAutoscalingMetric)
	s.Assertions.Equal(configurationAnnotations[KnativeAutoscalingMetric], DefaultKnativeAutoscalingMetric)

	configurationLabels := knativeConfiguration.Spec.Template.ObjectMeta.Labels
	s.Assertions.Len(configurationLabels, 5)
	s.Assertions.Contains(configurationLabels, ModelNameAnnotationKey)
	s.Assertions.Equal(md.Name, configurationLabels[ModelNameAnnotationKey])

	podSpec := knativeConfiguration.Spec.Template.Spec
	s.Assertions.Len(podSpec.Containers, 1)
	s.Assertions.Equal(DefaultTerminationPeriod, *podSpec.TimeoutSeconds)

	containerSpec := podSpec.Containers[0]
	mdResources, err := kubernetes.ConvertOdahuflowResourcesToK8s(md.Spec.Resources, config.NvidiaResourceName)
	s.Assertions.NoError(err)
	s.Assertions.Equal(mdResources, containerSpec.Resources)

	s.Assertions.Equal(image, containerSpec.Image)
	s.Assertions.Len(containerSpec.Ports, 1)
	s.Assertions.Equal(DefaultPortName, containerSpec.Ports[0].Name)
	s.Assertions.Equal(DefaultModelPort, containerSpec.Ports[0].ContainerPort)
	s.Assertions.NotNil(containerSpec.LivenessProbe)
	s.Assertions.Equal(mdLivenessDelay, containerSpec.LivenessProbe.InitialDelaySeconds)
	s.Assertions.NotNil(containerSpec.ReadinessProbe)
	s.Assertions.Equal(mdReadinessDelay, containerSpec.ReadinessProbe.InitialDelaySeconds)
}

// Node pool provided in packaging request, use it for knative configuration
func (s *ModelDeploymentControllerSuite) TestDeploymentReconcile_NodePoolProvided() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	someNodeSelector := map[string]string{"mode": "deployment"}
	deploymentConfig.NodePools = []config.NodePool{{NodeSelector: someNodeSelector}}
	s.initReconciler(deploymentConfig)

	md := newValidDeployment()
	md.Spec.NodeSelector = someNodeSelector

	cleanF := s.createDeployment(md)
	defer cleanF()
	knativeConfiguration := s.getKnativeConfiguration(md)

	s.Assertions.Nil(knativeConfiguration.Spec.Template.Spec.Affinity)
	s.Assertions.Equal(someNodeSelector, knativeConfiguration.Spec.Template.Spec.NodeSelector)
}

// Node pool not provided, build affinity for all deployment pools from config
func (s *ModelDeploymentControllerSuite) TestDeploymentReconcile_NodePoolNotProvided_UseAffinity() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	nodeSelector1 := map[string]string{"mode": "deployment"}
	nodeSelector2 := map[string]string{"mode": "deployment2", "label": "value"}
	deploymentConfig.NodePools = []config.NodePool{
		{NodeSelector: nodeSelector1},
		{NodeSelector: nodeSelector2},
	}
	s.initReconciler(deploymentConfig)

	md := newValidDeployment()
	md.Spec.NodeSelector = nil
	cleanF := s.createDeployment(md)
	defer cleanF()

	knativeConfiguration := s.getKnativeConfiguration(md)

	expectedNodeSelectorRequirement1 := []v1.NodeSelectorRequirement{{
		Key:      "mode",
		Operator: v1.NodeSelectorOpIn,
		Values:   []string{"deployment"},
	}}
	expectedNodeSelectorRequirement2 := []v1.NodeSelectorRequirement{
		{
			Key:      "mode",
			Operator: v1.NodeSelectorOpIn,
			Values:   []string{"deployment2"},
		},
		{
			Key:      "label",
			Operator: v1.NodeSelectorOpIn,
			Values:   []string{"value"},
		},
	}

	actualAffinity := knativeConfiguration.Spec.Template.Spec.Affinity
	actualNodeSelectorTerms := actualAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	s.Assertions.Len(actualNodeSelectorTerms, 2)

	for i, expectedNodeSelectorRequirement := range [][]v1.NodeSelectorRequirement{
		expectedNodeSelectorRequirement1, expectedNodeSelectorRequirement2,
	} {
		s.Assertions.ElementsMatch(expectedNodeSelectorRequirement, actualNodeSelectorTerms[i].MatchExpressions)
	}

	s.Assertions.Nil(knativeConfiguration.Spec.Template.Spec.NodeSelector)
}

func (s *ModelDeploymentControllerSuite) TestDeploymentReconcile_Tolerations() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	tolerations := []v1.Toleration{{Key: "dedicated", Operator: v1.TolerationOpEqual, Value: "deploy"}}
	deploymentConfig.Tolerations = tolerations
	s.initReconciler(deploymentConfig)

	md := newValidDeployment()

	cleanF := s.createDeployment(md)
	defer cleanF()
	knativeConfiguration := s.getKnativeConfiguration(md)

	s.Assertions.Equal(tolerations, knativeConfiguration.Spec.Template.Spec.Tolerations)
}

// Test utilities

func (s *ModelDeploymentControllerSuite) createDeployment(md *odahuflowv1alpha1.ModelDeployment) (
	cleanF func(),
) {
	err := s.k8sClient.Create(context.TODO(), md)
	s.Assertions.NoError(err)

	namespacedName := types.NamespacedName{Name: md.Name, Namespace: md.Namespace}
	s.Assertions.Eventually(
		func() bool { return s.k8sClient.Get(context.TODO(), namespacedName, md) == nil },
		10*time.Second,
		10*time.Millisecond)
	return func() { s.k8sClient.Delete(context.TODO(), md) }
}

func (s *ModelDeploymentControllerSuite) getKnativeConfiguration(md *odahuflowv1alpha1.ModelDeployment,
) *knservingv1.Configuration {
	configuration := &knservingv1.Configuration{}
	configurationKey := types.NamespacedName{Name: KnativeConfigurationName(md), Namespace: md.Namespace}
	s.Assertions.Eventually(
		func() bool { return s.k8sClient.Get(context.TODO(), configurationKey, configuration) == nil },
		10*time.Second,
		10*time.Millisecond,
		"Knative configuration not found!")
	return configuration
}

// Returns validDeployment with random Name to avoid collisions when running in parallel
func newValidDeployment() *odahuflowv1alpha1.ModelDeployment {
	md := validDeployment.DeepCopy()
	md.Name = fmt.Sprintf("deployment-%d", rand.Int()) //nolint
	return md
}

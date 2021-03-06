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
	"github.com/odahu/odahu-flow/packages/operator/pkg/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	"math/rand"
	"odahu-commons/predictors"
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
			Predictor:                  predictors.OdahuMLServer.ID,
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
	logger, _ = zap.NewDevelopment()
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

	reconciler := NewModelDeploymentReconciler(s.k8sManager, *cfg, logger)
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
			Predictor:                  predictors.OdahuMLServer.ID,
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
	knativeService := s.getKnativeService(md)

	podAnnonations := knativeService.Spec.Template.ObjectMeta.Annotations
	s.Assertions.Len(podAnnonations, 7)

	s.Assertions.Contains(podAnnonations, KnativeMinReplicasKey)
	s.Assertions.Equal(podAnnonations[KnativeMinReplicasKey], strconv.Itoa(int(mdMinReplicas)))

	s.Assertions.Contains(podAnnonations, KnativeMaxReplicasKey)
	s.Assertions.Equal(podAnnonations[KnativeMaxReplicasKey], strconv.Itoa(int(mdMaxReplicas)))

	s.Assertions.Contains(podAnnonations, KnativeAutoscalingTargetKey)
	s.Assertions.Equal(podAnnonations[KnativeAutoscalingTargetKey], KnativeAutoscalingTargetDefaultValue)

	s.Assertions.Contains(podAnnonations, KnativeAutoscalingClass)
	s.Assertions.Equal(podAnnonations[KnativeAutoscalingClass], DefaultKnativeAutoscalingClass)

	s.Assertions.Contains(podAnnonations, KnativeAutoscalingMetric)
	s.Assertions.Equal(podAnnonations[KnativeAutoscalingMetric], DefaultKnativeAutoscalingMetric)

	configurationLabels := knativeService.Spec.Template.ObjectMeta.Labels
	s.Assertions.Len(configurationLabels, 4)
	s.Assertions.Contains(configurationLabels, ModelNameAnnotationKey)
	s.Assertions.Equal(md.Name, configurationLabels[ModelNameAnnotationKey])

	podSpec := knativeService.Spec.Template.Spec
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
	knativeConfiguration := s.getKnativeService(md)

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

	knativeService := s.getKnativeService(md)

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

	actualAffinity := knativeService.Spec.Template.Spec.Affinity
	actualNodeSelectorTerms := actualAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	s.Assertions.Len(actualNodeSelectorTerms, 2)

	for i, expectedNodeSelectorRequirement := range [][]v1.NodeSelectorRequirement{
		expectedNodeSelectorRequirement1, expectedNodeSelectorRequirement2,
	} {
		s.Assertions.ElementsMatch(expectedNodeSelectorRequirement, actualNodeSelectorTerms[i].MatchExpressions)
	}

	s.Assertions.Nil(knativeService.Spec.Template.Spec.NodeSelector)
}

func (s *ModelDeploymentControllerSuite) TestDeploymentReconcile_Tolerations() {
	deploymentConfig := config.NewDefaultModelDeploymentConfig()
	tolerations := []v1.Toleration{{Key: "dedicated", Operator: v1.TolerationOpEqual, Value: "deploy"}}
	deploymentConfig.Tolerations = tolerations
	s.initReconciler(deploymentConfig)

	md := newValidDeployment()

	cleanF := s.createDeployment(md)
	defer cleanF()
	knativeConfiguration := s.getKnativeService(md)

	s.Assertions.Equal(tolerations, knativeConfiguration.Spec.Template.Spec.Tolerations)
}

func (s *ModelDeploymentControllerSuite) TestPolicyCM() {
	conf := config.NewDefaultModelDeploymentConfig()
	conf.Namespace = mdNamespace
	s.initReconciler(conf)

	firstRole := "first-role"

	md := &odahuflowv1alpha1.ModelDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: mdName, Namespace: mdNamespace},
		Spec: odahuflowv1alpha1.ModelDeploymentSpec{
			RoleName: &firstRole,
			Image:                      image,
			Predictor:                  predictors.OdahuMLServer.ID,
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

	cm := v1.ConfigMap{}

	cmName := GetCMPolicyName(md)

	predictor, ok := predictors.Predictors[md.Spec.Predictor]
	s.Require().True(ok)
	expectedCMData, err := deployment.ReadDefaultPoliciesAndRender(firstRole, predictor.OpaPolicyFilename)
	s.Require().NoError(err)

	s.Assertions.Eventually(func() bool {
		if err := s.k8sClient.Get(
			context.TODO(), types.NamespacedName{Name: cmName, Namespace: mdNamespace}, &cm); err != nil {
			return false
		}
		return assert.ObjectsAreEqual(expectedCMData, cm.Data)
	}, 10*time.Second, time.Second)

	// Let's update deployment by changing roleName. We should expect new ConfigMap
	secondRole := "second-role"
	md.Spec.RoleName = &secondRole
	err = s.k8sClient.Update(context.TODO(), md)
	s.Require().NoError(err)

	expectedCMData, err = deployment.ReadDefaultPoliciesAndRender(secondRole, predictor.OpaPolicyFilename)
	s.Require().NoError(err)

	s.Assertions.Eventually(func() bool {
		if err := s.k8sClient.Get(
			context.TODO(), types.NamespacedName{Name: cmName, Namespace: mdNamespace}, &cm); err != nil {
			return false
		}
		return assert.ObjectsAreEqual(expectedCMData, cm.Data)
	}, 10*time.Second, time.Second, "ConfigMap should be changed after roleName modification")


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

func (s *ModelDeploymentControllerSuite) getKnativeService(md *odahuflowv1alpha1.ModelDeployment,
) *knservingv1.Service {
	service := &knservingv1.Service{}
	serviceKey := types.NamespacedName{Name: KnativeServiceName(md), Namespace: md.Namespace}
	s.Assertions.Eventually(
		func() bool { return s.k8sClient.Get(context.TODO(), serviceKey, service) == nil },
		10*time.Second,
		10*time.Millisecond,
		"Knative service not found!")
	err := s.k8sClient.Get(context.TODO(), serviceKey, service)
	fmt.Println(err)
	return service
}

// Returns validDeployment with random Name to avoid collisions when running in parallel
func newValidDeployment() *odahuflowv1alpha1.ModelDeployment {
	md := validDeployment.DeepCopy()
	md.Name = fmt.Sprintf("deployment-%d", rand.Int()) //nolint
	return md
}

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
	"github.com/odahu/odahu-flow/packages/operator/controllers/mocks"
	apitypes "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	connapitypes "github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis/duck/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sync"
	"testing"
	"time"
)

const (
	requestID = "requestID"
	serviceID = "test-batch-service"
)


func getJob(ID string) odahuflowv1alpha1.BatchInferenceJob{
	return odahuflowv1alpha1.BatchInferenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: ID, Namespace: testNamespace},
		Spec:       odahuflowv1alpha1.BatchInferenceJobSpec{
			BatchInferenceServiceID: serviceID,
			Command:                 []string{"python"},
			Args:                    []string{"/opt/app/src.py", "--forecast"},
			InputPath:               "/input",
			OutputPath:              "/output",
			BatchRequestID:          requestID,
		},
		Status:     odahuflowv1alpha1.BatchInferenceJobStatus{},
	}
}

var serviceTpl = apitypes.InferenceService{
	ID:     serviceID,
	Spec:   apitypes.InferenceServiceSpec{
		Image:            "odahu-tools:latest",
		InputConnection:  "connection",
		OutputConnection: "connection",
		ModelConnection:  "connection",
	},
	Status: apitypes.InferenceServiceStatus{},
}

var conn = connapitypes.Connection{
	ID:     "connection",
	Spec:   odahuflowv1alpha1.ConnectionSpec{
		Type:        connapitypes.GcsType,
		URI:         "gs://bucket/output",
		Username:    "username",
		Password:    "pswd",
		KeyID:       "key",
		KeySecret:   "secret",
	},
	Status: odahuflowv1alpha1.ConnectionStatus{},
}



type BatchInferenceSuite struct {
	suite.Suite
	k8sClient  client.Client
	k8sManager manager.Manager
	stopMgr    chan struct{}
	mgrStopped *sync.WaitGroup
}

func TestBatchInferenceSuite(t *testing.T) {
	suite.Run(t, new(BatchInferenceSuite))
}


var testCases = []struct{
	trStatus tektonv1beta1.TaskRunStatus
	podStatus corev1.PodStatus

	expectJobStatus odahuflowv1alpha1.BatchInferenceJobStatus
	expectLastSubmittedRun apitypes.InferenceJobRun
} {
	{
		trStatus:         tektonv1beta1.TaskRunStatus{Status: v1beta1.Status{
			Conditions: v1beta1.Conditions{{
				Status: corev1.ConditionUnknown,
			}},
		}},
		podStatus:        corev1.PodStatus{},
		expectJobStatus:        odahuflowv1alpha1.BatchInferenceJobStatus{
			State:   odahuflowv1alpha1.BatchScheduling,
		},
		expectLastSubmittedRun: apitypes.InferenceJobRun{},
	},
	{
		trStatus: tektonv1beta1.TaskRunStatus{
			Status: v1beta1.Status{
				Conditions: v1beta1.Conditions{{
					Status: corev1.ConditionTrue,
				}},
			},
			TaskRunStatusFields: tektonv1beta1.TaskRunStatusFields{
				PodName: "podName",
			},
		},
		podStatus:        corev1.PodStatus{},
		expectJobStatus:        odahuflowv1alpha1.BatchInferenceJobStatus{
			State:   odahuflowv1alpha1.BatchSucceeded,
			PodName: "podName",
		},
		expectLastSubmittedRun: apitypes.InferenceJobRun{},
	},
}

// Tests what job status will be generated depending on the TaskRun and Pod statuses
func (s *BatchInferenceSuite) TestJobStatus() {
	for i, test := range testCases {
		s.runController()
		s.T().Run(fmt.Sprintf("jobStatus test#%d", i), func(t *testing.T) {
			as := require.New(t)

			connGetter := mocks.ConnGetter{}
			connGetter.On("GetConnection", "connection").Return(&conn, nil)
			podGetter := mocks.PodGetter{}
			podGetter.On("GetPod").Return(corev1.Pod{Status: test.podStatus}, nil)
			serviceAPI := mocks.BatchInferenceServiceAPI{}
			serviceAPI.On("Get", serviceID).Return(serviceTpl, nil)

			s.initReconciler(BatchInferenceJobReconcilerOptions{
				Mgr:               s.k8sManager,
				BatchInferenceAPI: &serviceAPI,
				ConnGetter:        &connGetter,
				PodGetter:         &podGetter,
				Cfg:               config.Config{},
			})

			// First of all let's create InferenceJob and wait until tekton TaskRun will be
			// created
			jobID := fmt.Sprintf("test-job-status-job-%d", i)
			job := getJob(jobID)

			as.NoError(s.k8sClient.Create(context.TODO(), &job))
			tr := &tektonv1beta1.TaskRun{}
			trKey := types.NamespacedName{Name: job.Name, Namespace: job.Namespace}
			as.Eventually(
				func() bool {
					if err := s.k8sClient.Get(context.TODO(), trKey, tr); err != nil {
						return false
					}
					return true
				},
				10*time.Second,
				10*time.Millisecond)

			// Then modify status of TaskRun to emulate TektonCD controller
			tr.Status = test.trStatus
			as.NoError(s.k8sClient.Status().Update(context.TODO(), tr))
			// And check that we have expected status of BatchInferenceJob depending on this test case
			//updatedJob := &odahuflowv1alpha1.BatchInferenceJob{}
			jobKey := types.NamespacedName{Name: job.Name, Namespace: job.Namespace}
			as.Eventually(
				func() bool {
					updatedJob := &odahuflowv1alpha1.BatchInferenceJob{}
					if err := s.k8sClient.Get(context.TODO(), jobKey, updatedJob); err != nil {
						return false
					}
					return updatedJob.Status == test.expectJobStatus
				},
				10*time.Second,
				10*time.Millisecond)

		})

		s.stopController()
	}
}

func (s *BatchInferenceSuite) runController() {
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
	s.Assertions.NoError(err)

	s.k8sClient = mgr.GetClient()
	s.k8sManager = mgr

	s.stopMgr = make(chan struct{})
	s.mgrStopped = &sync.WaitGroup{}
	s.mgrStopped.Add(1)
	go func() {
		defer s.mgrStopped.Done()
		s.Assertions.NoError(mgr.Start(s.stopMgr))
	}()
}

func (s *BatchInferenceSuite) initReconciler(opts BatchInferenceJobReconcilerOptions) {

	cfg := config.NewDefaultConfig()
	cfg.Batch.Namespace = testNamespace

	opts.Mgr = s.k8sManager
	opts.Cfg = *cfg

	reconciler := NewBatchInferenceJobReconciler(opts)
	s.Assertions.NoError(reconciler.SetupWithManager(s.k8sManager))
}

func (s *BatchInferenceSuite) stopController() {
	close(s.stopMgr)
	s.mgrStopped.Wait()
}

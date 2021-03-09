//
//    Copyright 2021 EPAM Systems
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
	connapitypes "github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis/duck/v1beta1"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sync"
	"testing"
	"time"
)

const (
	requestID = "requestID"
)

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

func getJob(id string) odahuflowv1alpha1.BatchInferenceJob{
	return odahuflowv1alpha1.BatchInferenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: id, Namespace: testNamespace},
		Spec:       odahuflowv1alpha1.BatchInferenceJobSpec{
			Command:                 []string{"python"},
			Args:                    []string{"/opt/app/src.py", "--forecast"},
			InputPath:               "/input",
			OutputPath:              "/output",
			BatchRequestID:          requestID,
			Image:            "odahu-tools:latest",
			InputConnection:  "connection",
			OutputConnection: "connection",
			ModelConnection:  "connection",
		},
		Status:     odahuflowv1alpha1.BatchInferenceJobStatus{},
	}
}

// job from getJob(ID string) will be created for each subtest
var testCases = []struct{
	// Appearing tektoncd TaskRun with trStatus will be emulated inside subtest to check
	// reaction of BatchInferenceJob controller
	trStatus tektonv1beta1.TaskRunStatus
	// You can configure PodStatus of Pod that correspond to created TaskRun
	podStatus corev1.PodStatus

	// expectJobStatus is the status that we expect BatchInferenceJob controller set
	// as reaction to TaskRun with trStatus
	expectJobStatus odahuflowv1alpha1.BatchInferenceJobStatus
} {
	{
		trStatus: tektonv1beta1.TaskRunStatus{
			Status: v1beta1.Status{
				Conditions: v1beta1.Conditions{{
					Status: corev1.ConditionTrue,
					Reason: "exit code 0",
					Message: "Emulated success",
				}},
			},
		},
		podStatus:        corev1.PodStatus{},
		expectJobStatus:        odahuflowv1alpha1.BatchInferenceJobStatus{
			State:   odahuflowv1alpha1.BatchSucceeded,
			Reason: "exit code 0",
			Message: "Emulated success",
		},
	},
	{
		trStatus: tektonv1beta1.TaskRunStatus{
			Status: v1beta1.Status{
				Conditions: v1beta1.Conditions{{
					Status: corev1.ConditionFalse,
					Reason: "exit code 1",
					Message: "Emulated fail",
				}},
			},
		},
		podStatus:        corev1.PodStatus{},
		expectJobStatus:        odahuflowv1alpha1.BatchInferenceJobStatus{
			State:   odahuflowv1alpha1.BatchFailed,
			Reason: "exit code 1",
			Message: "Emulated fail",
		},
	},
	{
		trStatus: tektonv1beta1.TaskRunStatus{
			Status: v1beta1.Status{
				Conditions: v1beta1.Conditions{{
					Status: corev1.ConditionUnknown,
				}},
			},
			TaskRunStatusFields: tektonv1beta1.TaskRunStatusFields{
				PodName: "",
			},
		},
		podStatus:        corev1.PodStatus{},
		expectJobStatus:        odahuflowv1alpha1.BatchInferenceJobStatus{
			State:   odahuflowv1alpha1.BatchScheduling,
		},
	},
	{
		trStatus:         tektonv1beta1.TaskRunStatus{
			Status: v1beta1.Status{
				Conditions: v1beta1.Conditions{{
					Status: corev1.ConditionUnknown,
				}},
			},
			TaskRunStatusFields: tektonv1beta1.TaskRunStatusFields{
				PodName: "pod",
			},
		},
		podStatus:        corev1.PodStatus{
			Reason: "Evicted", Message: "Pod was evicted",
		},
		expectJobStatus:        odahuflowv1alpha1.BatchInferenceJobStatus{
			State:   odahuflowv1alpha1.BatchFailed,
			Reason: "Evicted",
			Message: "Pod was evicted",
			PodName: "pod",
		},
	},
	{
		trStatus:         tektonv1beta1.TaskRunStatus{
			Status: v1beta1.Status{
				Conditions: v1beta1.Conditions{{
					Status: corev1.ConditionUnknown,
				}},
			},
			TaskRunStatusFields: tektonv1beta1.TaskRunStatusFields{
				PodName: "pod",
			},
		},
		podStatus:        corev1.PodStatus{
			Phase: corev1.PodPending,
		},
		expectJobStatus:        odahuflowv1alpha1.BatchInferenceJobStatus{
			State:   odahuflowv1alpha1.BatchScheduling,
			PodName: "pod",
		},
	},
	{
		trStatus:         tektonv1beta1.TaskRunStatus{
			Status: v1beta1.Status{
				Conditions: v1beta1.Conditions{{
					Status: corev1.ConditionUnknown,
				}},
			},
			TaskRunStatusFields: tektonv1beta1.TaskRunStatusFields{
				PodName: "pod",
			},
		},
		podStatus:        corev1.PodStatus{
			Phase: corev1.PodUnknown,
		},
		expectJobStatus:        odahuflowv1alpha1.BatchInferenceJobStatus{
			State:   odahuflowv1alpha1.BatchScheduling,
			PodName: "pod",
		},
	},
	{
		trStatus:         tektonv1beta1.TaskRunStatus{
			Status: v1beta1.Status{
				Conditions: v1beta1.Conditions{{
					Status: corev1.ConditionUnknown,
				}},
			},
			TaskRunStatusFields: tektonv1beta1.TaskRunStatusFields{
				PodName: "pod",
			},
		},
		podStatus:        corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
		expectJobStatus:        odahuflowv1alpha1.BatchInferenceJobStatus{
			State:   odahuflowv1alpha1.BatchRunning,
			PodName: "pod",
		},
	},
}

// Tests what job status will be generated depending on the TaskRun and Pod statuses
func TestJobStatus(t *testing.T) {

	for i, test := range testCases {
		i := i
		test := test
		t.Run(fmt.Sprintf("jobStatus test#%d", i), func(t *testing.T) {
			as := require.New(t)

			connGetter := mocks.ConnGetter{}
			connGetter.On("GetConnection", "connection").Return(&conn, nil)
			podGetter := mocks.PodGetter{}
			podGetter.On(
				"GetPod", mock.Anything, "pod", testNamespace,
				).Return(corev1.Pod{Status: test.podStatus}, nil)

			// Create isolated controller to avoid conflicts with different mocks
			// in parallel tests
			mgr, closeController := runBatchController(t, &connGetter, &podGetter)
			defer closeController()
			kubeClient := mgr.GetClient()

			// First of all let's create InferenceJob and wait until tekton TaskRun will be
			// created
			jobID := fmt.Sprintf("test-job-status-job-%d", i)
			job := getJob(jobID)

			as.NoError(kubeClient.Create(context.TODO(), &job))
			tr := &tektonv1beta1.TaskRun{}
			trKey := types.NamespacedName{Name: job.Name, Namespace: job.Namespace}
			as.Eventually(
				func() bool {
					if err := kubeClient.Get(context.TODO(), trKey, tr); err != nil {
						return false
					}
					log.Print("TaskRun is detected")
					return true
				},
				time.Minute,
				2*time.Second)

			// Then modify status of TaskRun to emulate TektonCD controller
			tr.Status = test.trStatus
			as.NoError(kubeClient.Status().Update(context.TODO(), tr))
			// And check that we have expected status of BatchInferenceJob depending on this test case
			//updatedJob := &odahuflowv1alpha1.BatchInferenceJob{}
			jobKey := types.NamespacedName{Name: job.Name, Namespace: job.Namespace}
			as.Eventually(
				func() bool {
					updatedJob := &odahuflowv1alpha1.BatchInferenceJob{}
					if err := kubeClient.Get(context.TODO(), jobKey, updatedJob); err != nil {
						log.Printf("BatchInferenceJob getting error %+v", err)
						return false
					}
					if updatedJob.Status == test.expectJobStatus{
						log.Print("BatchInferenceJob expected status is detected")
						return true
					}
					log.Printf("BatchInferenceJob actual status %+v, but expected %+v",
						updatedJob.Status, test.expectJobStatus)
					return false
				},
				time.Minute,
				2*time.Second)

		})

	}
}

func runBatchController(t *testing.T, connGetter ConnGetter, podGetter PodGetter) (manager.Manager, func()) {
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
	req := require.New(t)
	req.NoError(err)

	stopChan := make(chan struct{})
	mgrStopped := &sync.WaitGroup{}
	mgrStopped.Add(1)

	batchConfig := config.NewDefaultConfig().Batch
	batchConfig.Namespace = testNamespace

	opts := BatchInferenceJobReconcilerOptions{
		Client: mgr.GetClient(),
		Schema: mgr.GetScheme(),
		ConnGetter:        connGetter,
		PodGetter:         podGetter,
		Cfg: batchConfig,
	}

	req.NoError(NewBatchInferenceJobReconciler(opts).SetupWithManager(mgr))

	go func() {
		defer mgrStopped.Done()
		req.NoError(mgr.Start(stopChan))
		log.Printf("Controller %p is started", mgr)
	}()

	return mgr, func() {
		close(stopChan)
		mgrStopped.Wait()
		log.Printf("Controller %p is stopped", mgr)
	}
}


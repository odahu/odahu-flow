/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

package batch

import (
	"context"
	"fmt"
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/types"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"time"
)
import kube_types "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"

var (
	log = logf.Log.WithName("batch-job-adapter")
)

func kubeStatusToAppStatus(runtime *kube_types.BatchInferenceJobStatus) (app api_types.InferenceJobStatus) {
	app.Reason = runtime.Reason
	app.Message = runtime.Message
	app.PodName = runtime.PodName
	switch runtimeState := runtime.State; runtimeState {
	case kube_types.BatchRunning:
		app.State = api_types.Running
	case kube_types.BatchScheduling:
		app.State = api_types.Scheduling
	case kube_types.BatchFailed:
		app.State = api_types.Failed
	case kube_types.BatchSucceeded:
		app.State = api_types.Succeeded
	case kube_types.BatchUnknown:
		app.State = api_types.Unknown
	default:
		log.Info("unable to match Kubernetes inference job state with application job state")
		app.State = api_types.Unknown
	}
	return app
}

func appSpecToKubeSpec(appJob *api_types.InferenceJobSpec,
	appService *api_types.InferenceServiceSpec) (runtime kube_types.BatchInferenceJobSpec) {

	// Fill parameters from service
	runtime.Image = appService.Image
	runtime.Command = appService.Command
	runtime.Args = appService.Args
	runtime.ModelConnection = appService.ModelSource.Connection
	runtime.ModelPath = appService.ModelSource.Path

	// Fill parameters from job
	runtime.InputConnection = appJob.DataSource.Connection
	runtime.InputPath = appJob.DataSource.Path
	runtime.OutputConnection = appJob.OutputDestination.Connection
	runtime.OutputPath = appJob.OutputDestination.Path
	runtime.Resources = appJob.Resources
	runtime.NodeSelector = appJob.NodeSelector
	runtime.BatchRequestID = appJob.BatchRequestID

	return runtime
}

type kubeClient interface {
	DeleteInferenceJob(name string) error
	GetInferenceJob(name string) (*kube_types.BatchInferenceJob, error)
	CreateInferenceJob(name string, spec *kube_types.BatchInferenceJobSpec) error
	ListInferenceJob() (kube_types.BatchInferenceJobList, error)
}

type apiServer interface {
	UpdateStatus(ctx context.Context, id string, status api_types.InferenceJobStatus) error
	Get(ctx context.Context, id string) (api_types.InferenceJob, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, options ...filter.ListOption) ([]api_types.InferenceJob, error)
}

type apiServerServiceAPI interface {
	Get(ctx context.Context, id string) (res api_types.InferenceService, err error)
}

type KubeEntity struct {
	obj        *kube_types.BatchInferenceJob
	kubeClient kubeClient
	apiServer  apiServer
}

func (k KubeEntity) GetID() string {
	return k.obj.Name
}

func (k KubeEntity) GetSpecHash() (uint64, error) {
	return hashutil.Hash(k.obj.Spec)
}

func (k KubeEntity) GetStatusHash() (uint64, error) {
	storageStatus := kubeStatusToAppStatus(&k.obj.Status)
	return hashutil.Hash(storageStatus)
}

// BatchInferenceJob in kubernetes operator is deleted immediately
func (k KubeEntity) IsDeleting() bool {
	return false
}

func (k KubeEntity) Delete() error {
	return k.kubeClient.DeleteInferenceJob(k.GetID())
}

func (k KubeEntity) ReportStatus() error {
	storageStatus := kubeStatusToAppStatus(&k.obj.Status)
	return k.apiServer.UpdateStatus(context.TODO(), k.GetID(), storageStatus)
}

// Storage Entity

type StorageEntity struct {
	obj        *api_types.InferenceJob
	service    *api_types.InferenceService
	kubeClient kubeClient
	apiServer  apiServer
}

func (s StorageEntity) GetID() string {
	return s.obj.ID
}

func (s StorageEntity) GetSpecHash() (uint64, error) {
	runtimeSpec := appSpecToKubeSpec(&s.obj.Spec, &s.service.Spec)
	return hashutil.Hash(runtimeSpec)
}

func (s StorageEntity) GetStatusHash() (uint64, error) {
	return hashutil.Hash(s.obj.Spec)
}

func (s StorageEntity) UpdateInRuntime() error {
	return fmt.Errorf("updating job is not permitted. InferenceJob is immutable")
}

func (s StorageEntity) CreateInRuntime() error {
	runtimeSpec := appSpecToKubeSpec(&s.obj.Spec, &s.service.Spec)
	return s.kubeClient.CreateInferenceJob(s.GetID(), &runtimeSpec)
}

func (s StorageEntity) DeleteInRuntime() error {
	return s.kubeClient.DeleteInferenceJob(s.GetID())
}

func (s StorageEntity) DeleteInDB() error {
	return s.apiServer.Delete(context.TODO(), s.GetID())
}

func (s StorageEntity) IsFinished() bool {
	return s.obj.Status.State == api_types.Succeeded || s.obj.Status.State == api_types.Failed
}

func (s StorageEntity) HasDeletionMark() bool {
	return s.obj.DeletionMark
}

type statusReconciler struct {
	syncHook   types.StatusPollingHookFunc
	kubeClient kubeClient
	apiServer  apiServer
	apiServerServiceAPI  apiServerServiceAPI
}

func (r *statusReconciler) Reconcile(request ctrl.Request) (ctrl.Result, error) {

	ID := request.Name
	eLog := log.WithValues("ID", ID)

	// For some cases we need to trigger status sync without updates in kubernetes
	// For example when some entity was deleted and created again immediately in Storage
	// So we need pull status from kubernetes because entity with the same specification is existed
	// (ODAHU Controller didn't get to delete that in Kubernetes before it was created again with the same spec)
	result := ctrl.Result{RequeueAfter: time.Second * 30}

	runObj, err := r.kubeClient.GetInferenceJob(ID)
	if err != nil && !odahu_errors.IsNotFoundError(err) {
		eLog.Error(err, "Unable to get entity from kubernetes")
		return result, err
	} else if err != nil && odahu_errors.IsNotFoundError(err) {
		return ctrl.Result{}, nil
	}

	seObj, err := r.apiServer.Get(context.TODO(), ID)
	if err != nil && !odahu_errors.IsNotFoundError(err) {
		eLog.Error(err, "Unable to get entity from storage")
		return result, err
	} else if err != nil && odahu_errors.IsNotFoundError(err) {
		return ctrl.Result{}, nil
	}
	servObj, err := r.apiServerServiceAPI.Get(context.TODO(), seObj.Spec.InferenceServiceID)
	if err != nil && !odahu_errors.IsNotFoundError(err) {
		eLog.Error(err, "Unable to get service from storage")
		return result, err
	} else if err != nil && odahu_errors.IsNotFoundError(err) {
		return ctrl.Result{}, nil
	}

	se := &StorageEntity{
		obj:        &seObj,
		service: &servObj,
		kubeClient: r.kubeClient,
		apiServer:  r.apiServer,
	}

	kubeEntity := &KubeEntity{
		obj:        runObj,
		apiServer:  r.apiServer,
		kubeClient: r.kubeClient,
	}

	return result, r.syncHook(kubeEntity, se)
}

type Adapter struct {
	mgr        ctrl.Manager
	kubeClient kubeClient
	apiServer apiServer
	apiServerServiceAPI apiServerServiceAPI
}

func NewAdapter(mgr ctrl.Manager,
	kubeClient kubeClient,
	apiServer apiServer,
	apiServerServiceAPI apiServerServiceAPI) *Adapter {
	return &Adapter{mgr: mgr, kubeClient: kubeClient, apiServer: apiServer, apiServerServiceAPI: apiServerServiceAPI}
}

func (a Adapter) ListStorage() ([]types.StorageEntity, error) {
	result := make([]types.StorageEntity, 0)
	enList, err := a.apiServer.List(context.TODO())
	if err != nil {
		return result, err
	}


	for _, e := range enList {

		s, err := a.apiServerServiceAPI.Get(context.TODO(), e.Spec.InferenceServiceID)
		if err != nil {
			log.Error(fmt.Errorf("unable to fetch service %s for job %s", e.Spec.InferenceServiceID, e.ID),
				"unable to fetch job's service", "originalErr", err)
			continue
		}

		result = append(result, &StorageEntity{
			obj:        &e,
			service: &s,
			kubeClient: a.kubeClient,
			apiServer:    a.apiServer,
		})
	}

	return result, nil
}

func (a Adapter) ListRuntime() ([]types.RuntimeEntity, error) {
	result := make([]types.RuntimeEntity, 0)
	enList, err := a.kubeClient.ListInferenceJob()
	if err != nil {
		return result, err
	}

	for i := range enList.Items {
		result = append(result, &KubeEntity{
			obj:        &enList.Items[i],
			apiServer:    a.apiServer,
			kubeClient: a.kubeClient,
		})
	}
	return result, nil
}

func (a Adapter) GetFromStorage(id string) (types.StorageEntity, error) {
	obj, err := a.apiServer.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}
	return &StorageEntity{
		obj:        &obj,
		kubeClient: a.kubeClient,
		apiServer:  a.apiServer,
	}, nil
}

func (a Adapter) GetFromRuntime(id string) (types.RuntimeEntity, error) {
	obj, err := a.kubeClient.GetInferenceJob(id)
	if err != nil {
		return nil, err
	}
	return &KubeEntity{
		obj:        obj,
		kubeClient: a.kubeClient,
		apiServer:  a.apiServer,
	}, nil
}

func (a Adapter) SubscribeRuntimeUpdates(hook types.StatusPollingHookFunc) error {
	sr := &statusReconciler{apiServer: a.apiServer, kubeClient: a.kubeClient, syncHook: hook,
		apiServerServiceAPI: a.apiServerServiceAPI}

	return ctrl.NewControllerManagedBy(a.mgr).
		For(&kube_types.BatchInferenceJob{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(event event.DeleteEvent) bool { return false },
		}).
		Complete(sr)
}

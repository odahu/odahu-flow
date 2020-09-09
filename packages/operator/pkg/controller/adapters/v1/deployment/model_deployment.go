package deployment

import (
	"context"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	odahuv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	deployment_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/types"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/deploymentclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/deployment"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"time"
)

var (
	log = logf.Log.WithName("model-deployment-adapter")
)

type KubeEntity struct {
	obj        *deployment_types.ModelDeployment
	service       deployment.Service
	kubeClient kube_client.Client
}

func (k KubeEntity) GetID() string {
	return k.obj.ID
}

func (k KubeEntity) GetSpecHash() (uint64, error) {
	return hashutil.Hash(k.obj.Spec)
}

func (k KubeEntity) GetStatusHash() (uint64, error) {
	return hashutil.Hash(k.obj.Status)
}

func (k KubeEntity) Delete() error {
	return k.kubeClient.DeleteModelDeployment(k.GetID())
}

func (k KubeEntity) ReportStatus() error {
	return k.service.UpdateModelDeploymentStatus(context.TODO(), k.obj.ID, k.obj.Status, k.obj.Spec)
}

func (k KubeEntity) IsDeleting() bool {
	return k.obj.Status.State == odahuv1alpha1.ModelDeploymentStateDeleting
}

type StorageEntity struct {
	obj        *deployment_types.ModelDeployment
	kubeClient kube_client.Client
	service       deployment.Service
}

func (s *StorageEntity) GetID() string {
	return s.obj.ID
}

func (s *StorageEntity) GetSpecHash() (uint64, error) {
	return hashutil.Hash(s.obj.Spec)
}

func (s *StorageEntity) GetStatusHash() (uint64, error) {
	return hashutil.Hash(s.obj.Status)
}

func (s *StorageEntity) IsFinished() bool {
	return false
}

func (s *StorageEntity) HasDeletionMark() bool {
	return s.obj.DeletionMark
}

func (s *StorageEntity) CreateInRuntime() error {
	return s.kubeClient.CreateModelDeployment(s.obj)
}

func (s *StorageEntity) UpdateInRuntime() error {
	return s.kubeClient.UpdateModelDeployment(s.obj)
}

func (s *StorageEntity) DeleteInRuntime() error {
	return s.kubeClient.DeleteModelDeployment(s.GetID())
}

func (s *StorageEntity) DeleteInDB() error {
	return s.service.DeleteModelDeployment(context.TODO(), s.GetID())
}


type statusReconciler struct {
	kubeClient kube_client.Client
	syncHook types.StatusPollingHookFunc
	service deployment.Service
}

func (r *statusReconciler) Reconcile(request ctrl.Request) (ctrl.Result, error) {

	ID := request.Name
	eLog := log.WithValues("ID", ID)

	// For some cases we need to trigger status sync without updates in kubernetes
	// For example when some entity was deleted and created again immediately in Storage
	// So we need pull status from kubernetes because entity with the same specification is existed
	// (ODAHU Controller didn't get to delete that in Kubernetes before it was created again with the same spec)
	result := ctrl.Result{RequeueAfter: time.Second * 30}

	reObj, err := r.kubeClient.GetModelDeployment(ID)
	if err != nil && !odahu_errors.IsNotFoundError(err) {
		eLog.Error(err, "Unable to get entity from kubernetes")
		return result, err
	} else if err != nil && odahu_errors.IsNotFoundError(err) {
		return ctrl.Result{}, nil
	}


	seObj, err := r.service.GetModelDeployment(context.TODO(), ID)
	if err != nil && !odahu_errors.IsNotFoundError(err) {
		eLog.Error(err, "Unable to get entity from storage")
		return result, err
	} else if err != nil && odahu_errors.IsNotFoundError(err) {
		return ctrl.Result{}, nil
	}

	se := &StorageEntity{
		obj:        seObj,
		kubeClient: r.kubeClient,
		service:       r.service,
	}

	kubeEntity := &KubeEntity{
		obj:        reObj,
		service:       r.service,
		kubeClient: r.kubeClient,
	}

	return result, r.syncHook(kubeEntity, se)

}


type Adapter struct {
	service       deployment.Service
	kubeClient kube_client.Client
	mgr        ctrl.Manager
}

func NewAdapter(service deployment.Service, kubeClient kube_client.Client, mgr ctrl.Manager) *Adapter {
	return &Adapter{
		service:       service,
		kubeClient: kubeClient,
		mgr: mgr,
	}
}

func (s *Adapter) ListStorage() ([]types.StorageEntity, error) {

	result := make([]types.StorageEntity, 0)
	enList, err := s.service.GetModelDeploymentList(context.TODO())
	if err != nil {
		return result, err
	}

	for i := range enList {

		result = append(result, &StorageEntity{
			obj:        &enList[i],
			kubeClient: s.kubeClient,
			service: 		s.service,
		})
	}

	return result, nil
}

func (s *Adapter) ListRuntime() ([]types.RuntimeEntity, error) {
	result := make([]types.RuntimeEntity, 0)
	enList, err := s.kubeClient.GetModelDeploymentList()
	if err != nil {
		return result, err
	}

	for i := range enList {
		result = append(result, &KubeEntity{
			obj:        &enList[i],
			service:       s.service,
			kubeClient: s.kubeClient,
		})
	}
	return result, nil
}

func (s *Adapter) GetFromRuntime(id string) (types.RuntimeEntity, error) {
	mt, err := s.kubeClient.GetModelDeployment(id)
	if err != nil {
		return nil, err
	}

	return &KubeEntity{
		obj:        mt,
		service:       s.service,
		kubeClient: s.kubeClient,
	}, nil
}

func (s *Adapter) GetFromStorage(id string) (types.StorageEntity, error) {
	mt, err := s.service.GetModelDeployment(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	return &StorageEntity{
		obj:        mt,
		kubeClient: s.kubeClient,
		service: 		s.service,
	}, nil
}

func (s *Adapter) SubscribeRuntimeUpdates(syncHook types.StatusPollingHookFunc) error {

	sr := &statusReconciler{service: s.service, kubeClient: s.kubeClient, syncHook: syncHook}
	return ctrl.NewControllerManagedBy(s.mgr).
		For(&odahuv1alpha1.ModelDeployment{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(event event.DeleteEvent) bool { return false },
		}).
		Complete(sr)
}

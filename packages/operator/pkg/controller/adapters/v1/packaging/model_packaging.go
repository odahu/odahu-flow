package packaging

import (
	"context"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	odahuv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	packaging_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/types"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/packagingclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/packaging"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"time"
)

var (
	log = logf.Log.WithName("model-packaging-adapter")
)

func isPackagingFinished(mp *packaging_types.ModelPackaging) bool {
	state := mp.Status.State
	return state == odahuv1alpha1.ModelPackagingSucceeded || state == odahuv1alpha1.ModelPackagingFailed
}

type KubeEntity struct {
	obj        *packaging_types.ModelPackaging
	service    packaging.Service
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
	return k.kubeClient.DeleteModelPackaging(k.GetID())
}

func (k KubeEntity) ReportStatus() error {
	return k.service.UpdateModelPackagingStatus(context.TODO(), k.obj.ID, k.obj.Status, k.obj.Spec)
}

func (k KubeEntity) IsDeleting() bool {
	// Packaging does not have dependents that we need to wait
	return false
}

type StorageEntity struct {
	obj        *packaging_types.ModelPackaging
	kubeClient kube_client.Client
	service    packaging.Service
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
	return isPackagingFinished(s.obj)
}

func (s *StorageEntity) HasDeletionMark() bool {
	return s.obj.DeletionMark
}

func (s *StorageEntity) CreateInRuntime() error {
	return s.kubeClient.CreateModelPackaging(s.obj)
}

func (s *StorageEntity) UpdateInRuntime() error {
	return s.kubeClient.UpdateModelPackaging(s.obj)
}

func (s *StorageEntity) DeleteInRuntime() error {
	return s.kubeClient.DeleteModelPackaging(s.GetID())
}

func (s *StorageEntity) DeleteInDB() error {
	return s.service.DeleteModelPackaging(context.TODO(), s.GetID())
}


type statusReconciler struct {
	kubeClient kube_client.Client
	syncHook   types.StatusPollingHookFunc
	service    packaging.Service
}

func (r *statusReconciler) Reconcile(request ctrl.Request) (ctrl.Result, error) {

	ID := request.Name
	eLog := log.WithValues("ID", ID)

	// For some cases we need to trigger status sync without updates in kubernetes
	// For example when some entity was deleted and created again immediately in Storage
	// So we need pull status from kubernetes because entity with the same specification is existed
	// (ODAHU Controller didn't get to delete that in Kubernetes before it was created again with the same spec)
	result := ctrl.Result{RequeueAfter: time.Second * 30}

	reObj, err := r.kubeClient.GetModelPackaging(ID)
	if err != nil && !odahu_errors.IsNotFoundError(err) {
		eLog.Error(err, "Unable to get entity from kubernetes")
		return result, err
	} else if err != nil && odahu_errors.IsNotFoundError(err) {
		return ctrl.Result{}, nil
	}


	seObj, err := r.service.GetModelPackaging(context.TODO(), ID)
	if err != nil && !odahu_errors.IsNotFoundError(err) {
		eLog.Error(err, "Unable to get entity from storage")
		return result, err
	} else if err != nil && odahu_errors.IsNotFoundError(err) {
		return ctrl.Result{}, nil
	}

	se := &StorageEntity{
		obj:        seObj,
		kubeClient: r.kubeClient,
		service:    r.service,
	}

	kubeEntity := &KubeEntity{
		obj:        reObj,
		service:    r.service,
		kubeClient: r.kubeClient,
	}

	return result, r.syncHook(kubeEntity, se)

}


type Adapter struct {
	service    packaging.Service
	kubeClient kube_client.Client
	mgr        ctrl.Manager
}

func NewAdapter(service packaging.Service, kubeClient kube_client.Client, mgr ctrl.Manager) *Adapter {
	return &Adapter{
		service:    service,
		kubeClient: kubeClient,
		mgr:        mgr,
	}
}

func (s *Adapter) ListStorage() ([]types.StorageEntity, error) {

	result := make([]types.StorageEntity, 0)
	enList, err := s.service.GetModelPackagingList(context.TODO())
	if err != nil {
		return result, err
	}

	for i := range enList {

		result = append(result, &StorageEntity{
			obj:        &enList[i],
			kubeClient: s.kubeClient,
			service:    s.service,
		})
	}

	return result, nil
}

func (s *Adapter) ListRuntime() ([]types.RuntimeEntity, error) {
	result := make([]types.RuntimeEntity, 0)
	enList, err := s.kubeClient.GetModelPackagingList()
	if err != nil {
		return result, err
	}

	for i := range enList {
		result = append(result, &KubeEntity{
			obj:        &enList[i],
			service:    s.service,
			kubeClient: s.kubeClient,
		})
	}
	return result, nil
}

func (s *Adapter) GetFromRuntime(id string) (types.RuntimeEntity, error) {
	mt, err := s.kubeClient.GetModelPackaging(id)
	if err != nil {
		return nil, err
	}

	return &KubeEntity{
		obj:        mt,
		service:    s.service,
		kubeClient: s.kubeClient,
	}, nil
}

func (s *Adapter) GetFromStorage(id string) (types.StorageEntity, error) {
	mt, err := s.service.GetModelPackaging(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	return &StorageEntity{
		obj:        mt,
		kubeClient: s.kubeClient,
		service:    s.service,
	}, nil
}

func (s *Adapter) SubscribeRuntimeUpdates(syncHook types.StatusPollingHookFunc) error {

	sr := &statusReconciler{service: s.service, kubeClient: s.kubeClient, syncHook: syncHook}
	return ctrl.NewControllerManagedBy(s.mgr).
		For(&odahuv1alpha1.ModelPackaging{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(event event.DeleteEvent) bool { return false },
		}).
		Complete(sr)
}
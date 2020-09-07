package training

import (
	odahuv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	training_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/controller/types"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	ctrl "sigs.k8s.io/controller-runtime"
	hashutil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"time"
)

var (
	log = logf.Log.WithName("model-training-adapter")
)

func isTrainingFinished(mt *training_types.ModelTraining) bool {
	state := mt.Status.State
	return state == odahuv1alpha1.ModelTrainingSucceeded || state == odahuv1alpha1.ModelTrainingFailed
}


type KubeEntity struct {
	obj        *training_types.ModelTraining
	repo       training.Repository
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
	return k.kubeClient.DeleteModelTraining(k.GetID())
}

func (k KubeEntity) ReportStatus() error {
	return k.repo.UpdateModelTrainingStatus(k.obj.ID, k.obj.Status)
}

func (k KubeEntity) IsDeleting() bool {
	// Training does not have dependents that we need to wait
	return false
}

type StorageEntity struct {
	obj        *training_types.ModelTraining
	kubeClient kube_client.Client
	repo       training.Repository
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
	return isTrainingFinished(s.obj)
}

func (s *StorageEntity) HasDeletionMark() bool {
	return s.obj.DeletionMark
}

func (s *StorageEntity) CreateInRuntime() error {
	return s.kubeClient.CreateModelTraining(s.obj)
}

func (s *StorageEntity) UpdateInRuntime() error {
	return s.kubeClient.UpdateModelTraining(s.obj)
}

func (s *StorageEntity) DeleteInRuntime() error {
	return s.kubeClient.DeleteModelTraining(s.GetID())
}

func (s *StorageEntity) DeleteInDB() error {
	return s.repo.DeleteModelTraining(s.GetID())
}


type statusReconciler struct {
	syncHook types.StatusPollingHookFunc
	kubeClient kube_client.Client
	repo       training.Repository
}

func (r *statusReconciler) Reconcile(request ctrl.Request) (ctrl.Result, error) {

	ID := request.Name
	eLog := log.WithValues("ID", ID)

	// For some cases we need to trigger status sync without updates in kubernetes
	// For example when some entity was deleted and created again immediately in Storage
	// So we need pull status from kubernetes because entity with the same specification is existed
	// (ODAHU Controller didn't get to delete that in Kubernetes before it was created again with the same spec)
	result := ctrl.Result{RequeueAfter: time.Second * 30}

	runObj, err := r.kubeClient.GetModelTraining(ID)
	if err != nil && !odahu_errors.IsNotFoundError(err) {
		eLog.Error(err, "Unable to get entity from kubernetes")
		return result, err
	} else if err != nil && odahu_errors.IsNotFoundError(err) {
		return ctrl.Result{}, nil
	}

	seObj, err := r.repo.GetModelTraining(ID)
	if err != nil && !odahu_errors.IsNotFoundError(err) {
		eLog.Error(err, "Unable to get entity from storage")
		return result, err
	} else if err != nil && odahu_errors.IsNotFoundError(err) {
		return ctrl.Result{}, nil
	}

	se := &StorageEntity{
		obj:        seObj,
		kubeClient: r.kubeClient,
		repo:       r.repo,
	}

	kubeEntity := &KubeEntity{
		obj:        runObj,
		repo:       r.repo,
		kubeClient: r.kubeClient,
	}

	return result, r.syncHook(kubeEntity, se)
}


type Adapter struct {
	repo       training.Repository
	kubeClient kube_client.Client
	mgr        ctrl.Manager
}

func NewAdapter(repo training.Repository, kubeClient kube_client.Client, mgr ctrl.Manager) *Adapter {
	return &Adapter{
		repo:       repo,
		kubeClient: kubeClient,
		mgr: mgr,
	}
}

func (s *Adapter) ListStorage() ([]types.StorageEntity, error) {

	result := make([]types.StorageEntity, 0)
	enList, err := s.repo.GetModelTrainingList()
	if err != nil {
		return result, err
	}

	for i := range enList {

		result = append(result, &StorageEntity{
			obj:        &enList[i],
			kubeClient: s.kubeClient,
			repo: 		s.repo,
		})
	}

	return result, nil
}

func (s *Adapter) ListRuntime() ([]types.RuntimeEntity, error) {
	result := make([]types.RuntimeEntity, 0)
	enList, err := s.kubeClient.GetModelTrainingList()
	if err != nil {
		return result, err
	}

	for i := range enList {
		result = append(result, &KubeEntity{
			obj:        &enList[i],
			repo:       s.repo,
			kubeClient: s.kubeClient,
		})
	}
	return result, nil
}

func (s *Adapter) GetFromRuntime(id string) (types.RuntimeEntity, error) {
	mt, err := s.kubeClient.GetModelTraining(id)
	if err != nil {
		return nil, err
	}

	return &KubeEntity{
		obj:        mt,
		repo:       s.repo,
		kubeClient: s.kubeClient,
	}, nil
}

func (s *Adapter) GetFromStorage(id string) (types.StorageEntity, error) {
	mt, err := s.repo.GetModelTraining(id)
	if err != nil {
		return nil, err
	}

	return &StorageEntity{
		obj:        mt,
		kubeClient: s.kubeClient,
		repo: 		s.repo,
	}, nil
}

func (s *Adapter) SubscribeRuntimeUpdates(syncHook types.StatusPollingHookFunc) error {

	sr := &statusReconciler{repo: s.repo, kubeClient: s.kubeClient, syncHook: syncHook}

	return ctrl.NewControllerManagedBy(s.mgr).
		For(&odahuv1alpha1.ModelTraining{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(event event.DeleteEvent) bool { return false },
		}).
		Complete(sr)
}
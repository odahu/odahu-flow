package training

import (
	"github.com/mitchellh/hashstructure"
	odahuv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	training_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/workers/v1/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	return hashstructure.Hash(k.obj.Spec, nil)
}

func (k KubeEntity) GetStatusHash() (uint64, error) {
	return hashstructure.Hash(k.obj.Status, nil)
}

func (k KubeEntity) Delete() error {
	return k.kubeClient.DeleteModelTraining(k.GetID())
}

func (k KubeEntity) ReportStatus() error {
	return k.repo.UpdateModelTraining(k.obj)
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
	return hashstructure.Hash(s.obj.Spec, nil)
}

func (s *StorageEntity) GetStatusHash() (uint64, error) {
	return hashstructure.Hash(s.obj.Status, nil)
}

func (s *StorageEntity) IsFinished() bool {
	return isTrainingFinished(s.obj)
}

func (s *StorageEntity) HasDeletionMark() bool {
	return s.obj.DeletionMark
}

func (s *StorageEntity) CreateInService() error {
	return s.kubeClient.CreateModelTraining(s.obj)
}

func (s *StorageEntity) UpdateInService() error {
	return s.kubeClient.UpdateModelTraining(s.obj)
}

func (s *StorageEntity) DeleteInService() error {
	return s.kubeClient.DeleteModelTraining(s.GetID())
}

func (s *StorageEntity) DeleteInDB() error {
	return s.repo.DeleteModelTraining(s.GetID())
}


type Syncer struct {
	repo       training.Repository
	kubeClient kube_client.Client
}

func NewSyncer(repo training.Repository, kubeClient kube_client.Client) *Syncer {
	return &Syncer{
		repo:       repo,
		kubeClient: kubeClient,
	}
}

func (s *Syncer) AttachController(mgr ctrl.Manager, rec reconcile.Reconciler) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&odahuv1alpha1.ModelTraining{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(event event.DeleteEvent) bool { return false },
		}).
		Complete(rec)
}

func (s *Syncer) ListStorage() ([]types.StorageEntity, error) {

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

func (s *Syncer) ListService() ([]types.ExternalServiceEntity, error) {
	result := make([]types.ExternalServiceEntity, 0)
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

func (s *Syncer) GetFromService(id string) (types.ExternalServiceEntity, error) {
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

func (s *Syncer) GetFromStorage(id string) (types.StorageEntity, error) {
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

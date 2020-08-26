package training

import (
	"github.com/mitchellh/hashstructure"
	odahuv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	training_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/appserver/workers/v1/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type KubeEntity struct {
	obj     *training_types.ModelTraining
	storage training.Repository
	service training.Service
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

func (k KubeEntity) SaveResultInStorage() error {
	return k.storage.UpdateModelTraining(k.obj)
}

func (k KubeEntity) DeleteInService() error {
	return k.service.DeleteModelTraining(k.obj.ID)
}

type StorageEntity struct {
	obj     *training_types.ModelTraining
	service training.Service
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

func (s *StorageEntity) CreateInService() error {
	// Err then service.create
	// If not err then calc hashes
	// If hashes are equal do nothing
	// If hashes are not equal then service.update
	return s.service.CreateModelTraining(s.obj)
}

func (s *StorageEntity) UpdateInService() error {
	return s.service.UpdateModelTraining(s.obj)
}

func (s *StorageEntity) DeleteInKube() error {
	return s.service.DeleteModelTraining(s.obj.ID)
}

type Syncer struct {
	storage training.Repository
	service training.Service
}

func NewSyncer(storage training.Repository, service training.Service) *Syncer {
	return &Syncer{
		storage: storage,
		service: service,
	}
}

func (s *Syncer) AttachController(mgr ctrl.Manager, rec reconcile.Reconciler) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&odahuv1alpha1.ModelTraining{}).
		Complete(rec)
}

func (s *Syncer) ListStorage() ([]types.StorageEntity, error) {

	result := make([]types.StorageEntity, 0)
	enList, err := s.storage.GetModelTrainingList()
	if err != nil {
		return result, err
	}

	for i := range enList {

		result = append(result, &StorageEntity{
			obj:     &enList[i],
			service: s.service,
		})
	}

	return result, nil
}

func (s *Syncer) ListService() ([]types.ExternalServiceEntity, error) {
	result := make([]types.ExternalServiceEntity, 0)
	enList, err := s.service.GetModelTrainingList()
	if err != nil {
		return result, err
	}

	for i := range enList {
		result = append(result, &KubeEntity{
			obj:     &enList[i],
			service: s.service,
			storage: s.storage,
		})
	}
	return result, nil
}

func (s *Syncer) GetFromService(id string) (types.ExternalServiceEntity, error) {
	mt, err := s.service.GetModelTraining(id)
	if err != nil {
		return nil, err
	}

	return &KubeEntity{
		obj:     mt,
		storage: s.storage,
		service: s.service,
	}, nil
}

func (s *Syncer) GetFromStorage(id string) (types.StorageEntity, error) {
	mt, err := s.storage.GetModelTraining(id)
	if err != nil {
		return nil, err
	}

	return &StorageEntity{
		obj:     mt,
		service: s.service,
	}, nil
}

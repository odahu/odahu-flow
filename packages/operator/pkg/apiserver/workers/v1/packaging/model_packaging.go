package packaging

import (
	"github.com/mitchellh/hashstructure"
	odahuv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	packaging_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/packagingclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/workers/v1/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type KubeEntity struct {
	obj        *packaging_types.ModelPackaging
	repo       packaging.Repository
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

func (k KubeEntity) SaveResultInStorage() error {
	return k.repo.UpdateModelPackaging(k.obj)
}

func (k KubeEntity) DeleteInService() error {
	return k.kubeClient.DeleteModelPackaging(k.obj.ID)
}

type StorageEntity struct {
	obj        *packaging_types.ModelPackaging
	kubeClient kube_client.Client
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
	return s.kubeClient.CreateModelPackaging(s.obj)
}

func (s *StorageEntity) UpdateInService() error {
	return s.kubeClient.UpdateModelPackaging(s.obj)
}

func (s *StorageEntity) DeleteInKube() error {
	return s.kubeClient.DeleteModelPackaging(s.obj.ID)
}

type Syncer struct {
	repo       packaging.Repository
	kubeClient kube_client.Client
}

func NewSyncer(repo packaging.Repository, kubeClient kube_client.Client) *Syncer {
	return &Syncer{
		repo:       repo,
		kubeClient: kubeClient,
	}
}

func (s *Syncer) AttachController(mgr ctrl.Manager, rec reconcile.Reconciler) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&odahuv1alpha1.ModelPackaging{}).
		Complete(rec)
}

func (s *Syncer) ListStorage() ([]types.StorageEntity, error) {

	result := make([]types.StorageEntity, 0)
	enList, err := s.repo.GetModelPackagingList()
	if err != nil {
		return result, err
	}

	for i := range enList {

		result = append(result, &StorageEntity{
			obj:        &enList[i],
			kubeClient: s.kubeClient,
		})
	}

	return result, nil
}

func (s *Syncer) ListService() ([]types.ExternalServiceEntity, error) {
	result := make([]types.ExternalServiceEntity, 0)
	enList, err := s.kubeClient.GetModelPackagingList()
	if err != nil {
		return result, err
	}

	for i := range enList {
		result = append(result, &KubeEntity{
			obj:        &enList[i],
			kubeClient: s.kubeClient,
			repo:       s.repo,
		})
	}
	return result, nil
}

func (s *Syncer) GetFromService(id string) (types.ExternalServiceEntity, error) {
	mt, err := s.kubeClient.GetModelPackaging(id)
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
	mt, err := s.repo.GetModelPackaging(id)
	if err != nil {
		return nil, err
	}

	return &StorageEntity{
		obj:        mt,
		kubeClient: s.kubeClient,
	}, nil
}

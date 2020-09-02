package deployment

import (
	"github.com/mitchellh/hashstructure"
	odahuv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	deployment_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/workers/v1/types"
	kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/deploymentclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type KubeEntity struct {
	obj        *deployment_types.ModelDeployment
	repo       deployment.Repository
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
	return k.kubeClient.DeleteModelDeployment(k.GetID())
}

func (k KubeEntity) ReportStatus() error {
	return k.repo.UpdateModelDeployment(k.obj)
}

func (k KubeEntity) IsDeleting() bool {
	return k.obj.Status.State == odahuv1alpha1.ModelDeploymentStateDeleting
}

type StorageEntity struct {
	obj        *deployment_types.ModelDeployment
	kubeClient kube_client.Client
	repo       deployment.Repository
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
	return false
}

func (s *StorageEntity) HasDeletionMark() bool {
	return s.obj.DeletionMark
}

func (s *StorageEntity) CreateInService() error {
	return s.kubeClient.CreateModelDeployment(s.obj)
}

func (s *StorageEntity) UpdateInService() error {
	return s.kubeClient.UpdateModelDeployment(s.obj)
}

func (s *StorageEntity) DeleteInService() error {
	return s.kubeClient.DeleteModelDeployment(s.GetID())
}

func (s *StorageEntity) DeleteInDB() error {
	return s.repo.DeleteModelDeployment(s.GetID())
}

type Syncer struct {
	repo       deployment.Repository
	kubeClient kube_client.Client
}

func NewSyncer(repo deployment.Repository, kubeClient kube_client.Client) *Syncer {
	return &Syncer{
		repo:       repo,
		kubeClient: kubeClient,
	}
}

func (s *Syncer) AttachController(mgr ctrl.Manager, rec reconcile.Reconciler) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&odahuv1alpha1.ModelDeployment{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(event event.DeleteEvent) bool { return false },
		}).
		Complete(rec)
}

func (s *Syncer) ListStorage() ([]types.StorageEntity, error) {

	result := make([]types.StorageEntity, 0)
	enList, err := s.repo.GetModelDeploymentList()
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
	enList, err := s.kubeClient.GetModelDeploymentList()
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
	mt, err := s.kubeClient.GetModelDeployment(id)
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
	mt, err := s.repo.GetModelDeployment(id)
	if err != nil {
		return nil, err
	}

	return &StorageEntity{
		obj:        mt,
		kubeClient: s.kubeClient,
		repo: 		s.repo,
	}, nil
}

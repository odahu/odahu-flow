package types

import (
	"sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Entity that represent state of remote process in persistent storage
type StorageEntity interface {
	GetID() string
	GetSpecHash() (uint64, error)
	GetStatusHash() (uint64, error)

	UpdateInService() error
	CreateInService() error
}

// Entity that represent process on some external service
type ExternalServiceEntity interface {
	GetID() string
	GetSpecHash() (uint64, error)
	GetStatusHash() (uint64, error)

	SaveResultInStorage() error
	DeleteInService() error
}

type StorageServiceSyncer interface {
	AttachController(mgr controllerruntime.Manager, rec reconcile.Reconciler) error
	ListStorage() ([]StorageEntity, error)
	ListService() ([]ExternalServiceEntity, error)
	GetFromStorage(id string) (StorageEntity, error)
	GetFromService(id string) (ExternalServiceEntity, error)
}

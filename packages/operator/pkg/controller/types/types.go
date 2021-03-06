package types

// Entity that represent state of runtime process in persistent storage
type StorageEntity interface {
	GetID() string
	// If RuntimeEntity And StorageEntity have different structures that describe specification
	// then GetSpecHash MUST return hash of RuntimeEntity specification that correspond this StorageEntity
	// specification
	// to have an ability to compare and detect changes
	GetSpecHash() (uint64, error)
	GetStatusHash() (uint64, error)

	UpdateInRuntime() error
	CreateInRuntime() error
	DeleteInRuntime() error
	DeleteInDB() error

	IsFinished() bool
	HasDeletionMark() bool
}

// Entity that represent process on some runtime
type RuntimeEntity interface {
	GetID() string
	GetSpecHash() (uint64, error)
	// If RuntimeEntity And StorageEntity have different structures that describe status
	// then GetSpecHash MUST return hash of StorageEntity status that correspond this RuntimeEntity status
	// to have an ability to compare and detect changes
	GetStatusHash() (uint64, error)
	IsDeleting() bool
	Delete() error
	ReportStatus() error
}

// Hook that must be called on each update from runtime
type StatusPollingHookFunc func(RuntimeEntity, StorageEntity) error

// RuntimeAdapter to interact with runtime
type RuntimeAdapter interface {
	// List of items in storage
	ListStorage() ([]StorageEntity, error)
	// List of items in runtime
	ListRuntime() ([]RuntimeEntity, error)
	// Get entity from storage by ID
	GetFromStorage(id string) (StorageEntity, error)
	// Get entity from runtime by ID
	GetFromRuntime(id string) (RuntimeEntity, error)
	// Subscribe on runtime updates
	SubscribeRuntimeUpdates(hook StatusPollingHookFunc) error
}



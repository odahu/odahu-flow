package event

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"time"
)

type Type string
type Group string

const ModelRouteCreatedEventType Type = "ModelRouteCreated"

const ModelRouteDeletedEventType Type = "ModelRouteDeleted"

const ModelRouteUpdatedEventType Type = "ModelRouteUpdate"

const ModelRouteStatusUpdatedEventType Type = "ModelRouteStatusUpdated"

const ModelRouteDeletionMarkIsSetEventType Type = "ModelRouteDeletionMarkIsSet"

const ModelRouteEventGroup Group = "ModelRoute"

const ModelDeploymentCreatedEventType Type = "ModelDeploymentCreated"

const ModelDeploymentDeletedEventType Type = "ModelDeploymentDeleted"

const ModelDeploymentUpdatedEventType Type = "ModelDeploymentUpdated"

const ModelDeploymentStatusUpdatedEventType Type = "ModelDeploymentStatusUpdated"

const ModelDeploymentDeletionMarkIsSetEventType Type = "ModelDeploymentDeletionMarkIsSet"

const ModelDeploymentEventGroup Group = "ModelDeployment"

// This event is used for publishing
type Event struct {
	EntityID   string
	EventType  Type
	EventGroup Group
	Datetime   time.Time
	Payload    interface{}
}

type RouteEvent struct {
	// EntityID contains ID of ModelRoute for ModelRouteDeleted and ModelRouteDeletionMarkIsSet
	// event types
	// Does not make sense in case of ModelRouteUpdate, ModelRouteCreate, ModelRouteStatusUpdated events
	EntityID string `json:"entityID"`
	// Payload contains ModelRoute for ModelRouteUpdate, ModelRouteCreate,
	// ModelRouteStatusUpdated  events.
	// Does not make sense in case of ModelRouteDelete, ModelRouteDeletionMarkIsSet events
	Payload deployment.ModelRoute `json:"payload"`
	// Possible values: ModelRouteCreate, ModelRouteUpdate, ModelRouteDeleted,
	// ModelRouteDeletionMarkIsSet, ModelRouteStatusUpdated
	EventType Type `json:"type"`
	// When event is raised
	Datetime time.Time `json:"datetime"`
}

type DeploymentEvent struct {
	// EntityID contains ID of ModelDeployment for ModelDeploymentDeleted and ModelDeploymentDeletionMarkIsSet
	// event types
	// Does not make sense in case of ModelDeploymentUpdate, ModelDeploymentCreate, ModelDeploymentStatusUpdated events
	EntityID string `json:"entityID"`
	// Payload contains ModelDeployment for ModelDeploymentUpdate, ModelDeploymentCreate,
	// ModelDeploymentStatusUpdated  events.
	// Does not make sense in case of ModelDeploymentDelete, ModelDeploymentDeletionMarkIsSet events
	Payload deployment.ModelDeployment `json:"payload"`
	// Possible values: ModelDeploymentCreate, ModelDeploymentUpdate, ModelRouteDeleted,
	// ModelDeploymentDeletionMarkIsSet, ModelDeploymentStatusUpdated
	EventType Type `json:"type"`
	// When event is raised
	Datetime time.Time `json:"datetime"`
}

type LatestRouteEvents struct {
	Events []RouteEvent `json:"events"`
	Cursor int          `json:"cursor"`
}

type LatestDeploymentEvents struct {
	Events []DeploymentEvent `json:"events"`
	Cursor int               `json:"cursor"`
}

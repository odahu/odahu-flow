package model_route

import (
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"time"
)

type RouteEventReader struct {
	DB *sql.DB
}

type RouteEvent struct {
	// EntityID contains ID of ModelRoute for ModelRouteDeleted event type
	// Does not make sense in case of ModelRouteUpdate and ModelRouteCreate events
	EntityID string `json:"entityID"`
	// EntityID contains ModelRoute for ModelRouteUpdate and ModelRouteCreate events
	// Does not make sense in case of ModelRouteDelete event
	Payload deployment.ModelRoute  `json:"payload"`
	// Possible values: ModelRouteCreate, ModelRouteUpdate, ModelRouteDeleted
	EventType string `json:"type"`
	// When event is raised
	Datetime time.Time  `json:"datetime"`
}

func (rer RouteEventReader) Get(cursor int) ([]RouteEvent, int, error)  {
	return nil, 0, nil
}
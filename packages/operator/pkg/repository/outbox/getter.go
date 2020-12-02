package outbox

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"time"
	sq "github.com/Masterminds/squirrel"
)

type RouteEventGetter struct {
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
	EventType EventType `json:"type"`
	// When event is raised
	Datetime time.Time  `json:"datetime"`
}

func (rer RouteEventGetter) Get(ctx context.Context, cursor int) (routes []RouteEvent, newCursor int, err error)  {
	stmt, args, err := sq.
		Select(IDCol, EntityIDCol, EventTypeCol, PayloadCol, DatetimeCol).
		From(Table).
		Where(sq.Eq{EventGroupCol: ModelRouteEventGroup}).
		Where(sq.Gt{IDCol: cursor}).
		PlaceholderFormat(sq.Dollar).
		OrderBy(IDCol).
		ToSql()

	if err != nil {
		return
	}

	var rows *sql.Rows

	rows, err = rer.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return
	}

	for rows.Next() {
		var e RouteEvent
		if err = rows.Scan(&newCursor, &e.EntityID, &e.EventType, &e.Payload, &e.Datetime); err != nil {
			log.Error(err, "Unable to scan route")
			return routes, newCursor, err
		}
		if e.EventType != ModelRouteCreatedEventType && e.EventType != ModelRouteUpdatedEventType &&
			e.EventType != ModelRouteDeletedEventType {
			log.Error(fmt.Errorf("unknown event for ModelRouteEventGroup: %v", e.EventType), "")
			continue
		}

		routes = append(routes, e)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("error during rows iteration: %v", err)
	}

	return routes, newCursor, err
}

type DeploymentEventGetter struct {
	DB *sql.DB
}

type DeploymentEvent struct {
	// EntityID contains ID of ModelDeployment for ModelDeploymentDeleted event type
	// Does not make sense in case of ModelDeploymentUpdate and ModelDeploymentCreate events
	EntityID string `json:"entityID"`
	// EntityID contains ModelDeployment for ModelDeploymentUpdate and ModelDeploymentCreate events
	// Does not make sense in case of ModelDeploymentDelete event
	Payload deployment.ModelDeployment  `json:"payload"`
	// Possible values: ModelDeploymentCreate, ModelDeploymentUpdate, ModelRouteDeleted
	EventType EventType `json:"type"`
	// When event is raised
	Datetime time.Time  `json:"datetime"`
}

func (rer DeploymentEventGetter) Get(ctx context.Context, cursor int) (routes []DeploymentEvent,
	newCursor int, err error)  {

	stmt, args, err := sq.
		Select(IDCol, EntityIDCol, EventTypeCol, PayloadCol, DatetimeCol).
		From(Table).
		Where(sq.Eq{EventGroupCol: ModelDeploymentEventGroup}).
		Where(sq.Gt{IDCol: cursor}).
		PlaceholderFormat(sq.Dollar).
		OrderBy(IDCol).
		ToSql()

	if err != nil {
		return
	}

	var rows *sql.Rows

	rows, err = rer.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return
	}

	for rows.Next() {
		var e DeploymentEvent
		if err = rows.Scan(&newCursor, &e.EntityID, &e.EventType, &e.Payload, &e.Datetime); err != nil {
			log.Error(err, "Unable to scan route")
			return
		}
		if e.EventType != ModelDeploymentCreatedEventType && e.EventType != ModelDeploymentUpdatedEventType &&
			e.EventType != ModelDeploymentDeletedEventType {
			log.Error(fmt.Errorf("unknown event for ModelDeploymentEventGroup: %v", e.EventType), "")
		}

		routes = append(routes, e)
	}
	if err = rows.Err(); err != nil {
		err = fmt.Errorf("error during rows iteration: %v", err)
	}
	return routes, newCursor, err
}
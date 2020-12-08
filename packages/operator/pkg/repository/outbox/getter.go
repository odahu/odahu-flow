package outbox

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
)

type RouteEventGetter struct {
	DB *sql.DB
}

func (rer RouteEventGetter) Get(
	ctx context.Context, cursor int) (routes []event.RouteEvent, newCursor int, err error)  {
	stmt, args, err := sq.
		Select(IDCol, EntityIDCol, EventTypeCol, PayloadCol, DatetimeCol).
		From(Table).
		Where(sq.Eq{EventGroupCol: event.ModelRouteEventGroup}).
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
		var e event.RouteEvent
		if err = rows.Scan(&newCursor, &e.EntityID, &e.EventType, &e.Payload, &e.Datetime); err != nil {
			log.Error(err, "Unable to scan route")
			return routes, newCursor, err
		}
		availableTypes := []event.Type{
			event.ModelRouteCreatedEventType, event.ModelRouteUpdatedEventType, event.ModelRouteDeletedEventType,
			event.ModelRouteDeletionMarkIsSetEventType, event.ModelRouteStatusUpdatedEventType,
		}
		if !EventTypeOK(availableTypes, e.EventType) {
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

func (rer DeploymentEventGetter) Get(ctx context.Context, cursor int) (routes []event.DeploymentEvent,
	newCursor int, err error)  {

	stmt, args, err := sq.
		Select(IDCol, EntityIDCol, EventTypeCol, PayloadCol, DatetimeCol).
		From(Table).
		Where(sq.Eq{EventGroupCol: event.ModelDeploymentEventGroup}).
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
		var e event.DeploymentEvent
		if err = rows.Scan(&newCursor, &e.EntityID, &e.EventType, &e.Payload, &e.Datetime); err != nil {
			log.Error(err, "Unable to scan route")
			return
		}

		availableTypes := []event.Type{
			event.ModelDeploymentCreatedEventType, event.ModelDeploymentUpdatedEventType, event.ModelDeploymentDeletedEventType,
			event.ModelDeploymentStatusUpdatedEventType, event.ModelDeploymentDeletionMarkIsSetEventType,
		}
		if !EventTypeOK(availableTypes, e.EventType) {
			log.Error(fmt.Errorf("unknown event for ModelDeploymentEventGroup: %v", e.EventType), "")
			continue
		}
		routes = append(routes, e)
	}
	if err = rows.Err(); err != nil {
		err = fmt.Errorf("error during rows iteration: %v", err)
	}
	return routes, newCursor, err
}
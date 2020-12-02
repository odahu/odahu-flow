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
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err, "Unable to close rows (release connection)")
		}
	}()

	for rows.Next() {
		var e RouteEvent
		if err = rows.Scan(&newCursor, &e.EntityID, &e.EventType, &e.Payload, &e.Datetime); err != nil {
			log.Error(err, "Unable to scan route")
			return
		}
		if e.EventType != ModelRouteCreatedEventType && e.EventType != ModelRouteUpdatedEventType &&
			e.EventType != ModelRouteDeletedEventType {
			log.Error(fmt.Errorf("unknown event for ModelRouteEventGroup: %v", e.EventType), "")
		}

		routes = append(routes, e)
	}
	return
}
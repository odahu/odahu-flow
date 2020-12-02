package outbox

import (
	"context"
	"database/sql"
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
		Select(EntityIDCol, EventTypeCol, PayloadCol, DatetimeCol).
		From(Table).
		Where(sq.Eq{EventGroupCol: ModelRouteEventGroup}, sq.Gt{IDCol: cursor}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return routes, newCursor, err
	}

	rer.DB.QueryContext(ctx, stmt, args...)

	return nil, 0, nil
}
package outbox

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	db_utils "github.com/odahu/odahu-flow/packages/operator/pkg/utils/db"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)


type EventType string
type EventGroup string

const (
	Table = "odahu_outbox"

	ModelRouteCreatedEventType           EventType  = "ModelRouteCreated"
	ModelRouteDeletedEventType           EventType  = "ModelRouteDeleted"
	ModelRouteUpdatedEventType           EventType  = "ModelRouteUpdate"
	ModelRouteStatusUpdatedEventType     EventType  = "ModelRouteStatusUpdated"
	ModelRouteDeletionMarkIsSetEventType EventType  = "ModelRouteDeletionMarkIsSet"
	ModelRouteEventGroup                 EventGroup = "ModelRoute"

	ModelDeploymentCreatedEventType           EventType  = "ModelDeploymentCreated"
	ModelDeploymentDeletedEventType           EventType  = "ModelDeploymentDeleted"
	ModelDeploymentUpdatedEventType           EventType  = "ModelDeploymentUpdated"
	ModelDeploymentStatusUpdatedEventType     EventType  = "ModelDeploymentStatusUpdated"
	ModelDeploymentDeletionMarkIsSetEventType EventType  = "ModelDeploymentDeletionMarkIsSet"
	ModelDeploymentEventGroup                 EventGroup = "ModelDeployment"
)


func EventTypeOK(available []EventType, actual EventType) bool {
	for _, e := range available {
		if e == actual {
			return true
		}
	}
	return false
}


var txOptions = &sql.TxOptions{
	Isolation: sql.LevelRepeatableRead,
	ReadOnly:  false,
}

var log = logf.Log.WithName("event-repository")

const (
	EntityIDCol = "entity_id"
	EventTypeCol = "event_type"
	EventGroupCol = "event_group"
	DatetimeCol = "datetime"
	PayloadCol = "payload"
)

type Event struct {
	EntityID string
	EventType EventType
	EventGroup EventGroup
	Datetime time.Time
	Payload  interface{}
}

type EventPublisher struct {
	DB *sql.DB
}

const (
	IDCol = "id"
)

func (repo EventPublisher) PublishEvent(ctx context.Context, tx *sql.Tx, event Event) (err error) {

	if tx == nil {
		tx, err = repo.DB.BeginTx(ctx, txOptions)
		if err != nil {
			return
		}
		defer func(){db_utils.FinishTx(tx, err, log)}()
	}

	// First, delete previous event with the same ID, because outbox table currently stores
	// Only last event with EntityID
	stmt, args, err := sq.
		Delete(Table).
		Where(sq.Eq{EntityIDCol: event.EntityID, EventGroupCol: event.EventGroup}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return
	}

	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return
	}

	if event.Datetime.IsZero() {
		event.Datetime = time.Now()
	}

	stmt, args, err = sq.
		Insert(Table).
		Columns(EntityIDCol, EventTypeCol, EventGroupCol, DatetimeCol, PayloadCol).
		Values(event.EntityID, event.EventType, event.EventGroup, event.Datetime, event.Payload).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return
	}

	_, err = tx.ExecContext(
		ctx,
		stmt,
		args...
	)
	if err != nil {
		return
	}

	return err
}


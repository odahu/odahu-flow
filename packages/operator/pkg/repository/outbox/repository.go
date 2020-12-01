package outbox

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"time"
)

const (
	Table = "odahu_outbox"

	ModelRouteCreatedEventType = "ModelRouteCreated"
	ModelRouteDeletedEventType = "ModelRouteDeleted"
	ModelRouteUpdatedEventType = "ModelRouteUpdate"
	ModelRouteEventGroup       = "ModelRoute"

	ModelDeploymentCreatedEventType = "ModelDeploymentCreated"
	ModelDeploymentDeletedEventType = "ModelDeploymentDeleted"
	ModelDeploymentUpdatedEventType = "ModelDeploymentUpdated"
	ModelDeploymentEventGroup       = "ModelDeployment"
)

var txOptions = &sql.TxOptions{
	Isolation: sql.LevelRepeatableRead,
	ReadOnly:  false,
}


type Event struct {
	EntityID string
	EventType string
	EventGroup string
	Datetime time.Time
	Payload  interface{}
}

type EventRepository struct {
	DB *sql.DB
}

type EventRecord struct {
	id int
	event Event
}

func (repo EventRepository) RaiseEvent(ctx context.Context, tx *sql.Tx, event Event) (err error) {

	if tx != nil {
		tx, err = repo.DB.BeginTx(ctx, txOptions)
		if err != nil {
			return
		}
	}

	// First, delete previous event with the same ID, because outbox table currently stores
	// Only last event with EntityID
	stmt, args, err := sq.
		Delete(Table).
		Where(sq.Eq{"id": event.EntityID, "event_group": event.EventGroup}).
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
		Columns("entity_id", "event_type", "event_group", "datetime", "payload").
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

	return
}


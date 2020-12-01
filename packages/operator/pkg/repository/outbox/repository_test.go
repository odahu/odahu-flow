package outbox_test

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/outbox"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	"log"
	"os"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

var (
	db *sql.DB
)

func Wrapper(m *testing.M) int {
	// Setup Test DB

	var closeDB func() error
	var err error
	db, _, closeDB, err = testenvs.SetupTestDB()
	defer func() {
		if err := closeDB(); err != nil {
			log.Print("Error during release test DB resources")
		}
	}()
	if err != nil {
		return -1
	}

	return m.Run()
}

func TestMain(m *testing.M) {

	os.Exit(Wrapper(m))

}

type CustomPayload struct {
	name string
}

func TestRaiseEvent(t *testing.T) {
	eventRepo := &outbox.EventRepository{DB: db}
	event1 := outbox.Event{
		EntityID:   "CustomID",
		EventType:  "Create",
		EventGroup: "CustomEventGroup",
		Datetime:   time.Time{},
		Payload:    CustomPayload{name: "test"},
	}
	err := eventRepo.RaiseEvent(context.Background(), nil, event1)
	assert.NoError(t, err)
	event2 := outbox.Event{
		EntityID:   "CustomID-2",
		EventType:  "Create",
		EventGroup: "CustomEventGroup",
		Datetime:   time.Time{},
		Payload:    CustomPayload{name: "test"},
	}
	err = eventRepo.RaiseEvent(context.Background(), nil, event2)
	assert.NoError(t, err)
	event1Delete := outbox.Event{
		EntityID:   "CustomID",
		EventType:  "Delete",
		EventGroup: "CustomEventGroup",
		Datetime:   time.Time{},
		Payload:    CustomPayload{name: "test"},
	}
	err = eventRepo.RaiseEvent(context.Background(), nil, event1Delete)
	assert.NoError(t, err)

	stmt, _, err := sq.
		Select(outbox.EntityIDCol, outbox.EventTypeCol, outbox.EventGroupCol, outbox.DatetimeCol, outbox.PayloadCol).
		From(outbox.Table).
		OrderBy(outbox.IDCol).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return
	}

	rows, err := db.Query(stmt)
	assert.NoError(t, err)
	defer func() {_ = rows.Close()}()

	assert.True(t, rows.Next())
	e := outbox.Event{}
	assert.NoError(t, rows.Scan(&e.EntityID, &e.EventType, &e.EventGroup, &e.Datetime, &e.Payload))
	assert.Equal(t, event2, e)

	assert.False(t, rows.Next())
	e = outbox.Event{}
	assert.NoError(t, rows.Scan(&e.EntityID, &e.EventType, &e.EventGroup, &e.Datetime, &e.Payload))
	assert.Equal(t, event1Delete, e)

}


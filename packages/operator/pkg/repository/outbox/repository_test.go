package outbox_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/outbox"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
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
	Name string `json:"name"`
}

func (c *CustomPayload) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	res := json.Unmarshal(b, &c)
	return res
}

func (c CustomPayload) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func TestRaiseEvent(t *testing.T) {
	eventRepo := &outbox.EventRepository{DB: db}
	event1 := outbox.Event{
		EntityID:   "CustomID",
		EventType:  "Create",
		EventGroup: "CustomEventGroup",
		Datetime:   time.Now().UTC(),
		Payload:    CustomPayload{Name: "test"},
	}
	err := eventRepo.RaiseEvent(context.Background(), nil, event1)
	assert.NoError(t, err)
	event2 := outbox.Event{
		EntityID:   "CustomID-2",
		EventType:  "Create",
		EventGroup: "CustomEventGroup",
		Datetime:   time.Now().UTC(),
		Payload:    CustomPayload{Name: "test"},
	}
	err = eventRepo.RaiseEvent(context.Background(), nil, event2)
	assert.NoError(t, err)
	event1Delete := outbox.Event{
		EntityID:   "CustomID",
		EventType:  "Delete",
		EventGroup: "CustomEventGroup",
		Datetime:   time.Now().UTC(),
		Payload:    CustomPayload{Name: "test"},
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
	p := CustomPayload{}
	assert.NoError(t, rows.Scan(&e.EntityID, &e.EventType, &e.EventGroup, &e.Datetime, &p))
	e.Datetime = e.Datetime.UTC()
	e.Payload = p
	assert.Equal(t, event2, e)

	assert.True(t, rows.Next())
	assert.NoError(t, rows.Scan(&e.EntityID, &e.EventType, &e.EventGroup, &e.Datetime, &p))
	e.Datetime = e.Datetime.UTC()
	e.Payload = p
	assert.Equal(t, event1Delete, e)

	assert.False(t, rows.Next())

}


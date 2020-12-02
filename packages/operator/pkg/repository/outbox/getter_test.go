package outbox_test

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/outbox"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	sq "github.com/Masterminds/squirrel"
)

func TestModelRouteGet(t *testing.T) {
	eventPublisher := &outbox.EventPublisher{DB: db}
	routeEventGetter := &outbox.RouteEventGetter{DB: db}


	payload1 := deployment.ModelRoute{Spec: v1alpha1.ModelRouteSpec{URLPrefix: "Test"}}
	event1 := outbox.Event{
		EntityID:   "event1", EventType: outbox.ModelRouteCreatedEventType,
		EventGroup: outbox.ModelRouteEventGroup, Datetime:   time.Now().UTC(),
		Payload:    payload1,
	}
	err := eventPublisher.PublishEvent(context.Background(), nil, event1)
	assert.NoError(t, err)

	routes, newC, err := routeEventGetter.Get(context.TODO(), 0)
	assert.NoError(t, err)
	assert.Len(t, routes, 1)
	r := routes[0]
	assert.Equal(t, event1.EntityID, r.EntityID)
	assert.Equal(t, event1.EventType, r.EventType)
	assert.Equal(t, payload1, r.Payload)
	assert.Equal(t, event1.Datetime, r.Datetime.UTC())

	payload2 := deployment.ModelRoute{Spec: v1alpha1.ModelRouteSpec{URLPrefix: "Test"}}
	event2 := outbox.Event{
		EntityID:   "event2", EventType: outbox.ModelRouteCreatedEventType,
		EventGroup: outbox.ModelRouteEventGroup, Datetime:   time.Now().UTC(),
		Payload:    payload2,
	}
	err = eventPublisher.PublishEvent(context.Background(), nil, event2)
	assert.NoError(t, err)

	payload3 := deployment.ModelRoute{Spec: v1alpha1.ModelRouteSpec{URLPrefix: "Test"}}
	event3 := outbox.Event{
		EntityID:   "event3", EventType: outbox.ModelRouteCreatedEventType,
		EventGroup: outbox.ModelRouteEventGroup, Datetime:   time.Now().UTC(),
		Payload:    payload3,
	}
	err = eventPublisher.PublishEvent(context.Background(), nil, event3)
	assert.NoError(t, err)

	// Fetch next route updates
	routes, newC, err = routeEventGetter.Get(context.TODO(), newC)
	assert.NoError(t, err)
	assert.Len(t, routes, 2)
	r = routes[0]
	assert.Equal(t, event2.EntityID, r.EntityID)
	assert.Equal(t, event2.EventType, r.EventType)
	assert.Equal(t, payload2, r.Payload)
	assert.Equal(t, event2.Datetime, r.Datetime.UTC())
	r = routes[1]
	assert.Equal(t, event3.EntityID, r.EntityID)
	assert.Equal(t, event3.EventType, r.EventType)
	assert.Equal(t, payload3, r.Payload)
	assert.Equal(t, event3.Datetime, r.Datetime.UTC())

	stmt, _, _ := sq.Delete(outbox.Table).ToSql()
	_, _ = db.Exec(stmt)

}

func TestModelDeploymentGet(t *testing.T) {
	eventPublisher := &outbox.EventPublisher{DB: db}
	deploymentEventGetter := &outbox.DeploymentEventGetter{DB: db}


	payload1 := deployment.ModelDeployment{Spec: v1alpha1.ModelDeploymentSpec{Image: "Image"}}
	event1 := outbox.Event{
		EntityID:   "event1", EventType: outbox.ModelDeploymentCreatedEventType,
		EventGroup: outbox.ModelDeploymentEventGroup, Datetime:   time.Now().UTC(),
		Payload:    payload1,
	}
	err := eventPublisher.PublishEvent(context.Background(), nil, event1)
	assert.NoError(t, err)

	routes, newC, err := deploymentEventGetter.Get(context.TODO(), 0)
	assert.NoError(t, err)
	assert.Len(t, routes, 1)
	r := routes[0]
	assert.Equal(t, event1.EntityID, r.EntityID)
	assert.Equal(t, event1.EventType, r.EventType)
	assert.Equal(t, payload1, r.Payload)
	assert.Equal(t, event1.Datetime, r.Datetime.UTC())

	payload2 := deployment.ModelDeployment{Spec: v1alpha1.ModelDeploymentSpec{Image: "Image"}}
	event2 := outbox.Event{
		EntityID:   "event2", EventType: outbox.ModelDeploymentCreatedEventType,
		EventGroup: outbox.ModelDeploymentEventGroup, Datetime:   time.Now().UTC(),
		Payload:    payload2,
	}
	err = eventPublisher.PublishEvent(context.Background(), nil, event2)
	assert.NoError(t, err)

	payload3 := deployment.ModelDeployment{Spec: v1alpha1.ModelDeploymentSpec{Image: "Image"}}
	event3 := outbox.Event{
		EntityID:   "event3", EventType: outbox.ModelDeploymentCreatedEventType,
		EventGroup: outbox.ModelDeploymentEventGroup, Datetime:   time.Now().UTC(),
		Payload:    payload3,
	}
	err = eventPublisher.PublishEvent(context.Background(), nil, event3)
	assert.NoError(t, err)

	// Fetch next route updates
	routes, newC, err = deploymentEventGetter.Get(context.TODO(), newC)
	assert.NoError(t, err)
	assert.Len(t, routes, 2)
	r = routes[0]
	assert.Equal(t, event2.EntityID, r.EntityID)
	assert.Equal(t, event2.EventType, r.EventType)
	assert.Equal(t, payload2, r.Payload)
	assert.Equal(t, event2.Datetime, r.Datetime.UTC())
	r = routes[1]
	assert.Equal(t, event3.EntityID, r.EntityID)
	assert.Equal(t, event3.EventType, r.EventType)
	assert.Equal(t, payload3, r.Payload)
	assert.Equal(t, event3.Datetime, r.Datetime.UTC())

	stmt, _, _ := sq.Delete(outbox.Table).ToSql()
	_, _ = db.Exec(stmt)

}
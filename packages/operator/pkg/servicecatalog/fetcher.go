package servicecatalog

import (
	event_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
)



type RouteEventsAPIClient interface {
	GetLastEvents(cursor int) (event_types.LatestRouteEvents, error)
}

type RouteEventFetcher struct {
	APIClient RouteEventsAPIClient
}

func (d RouteEventFetcher) GetLastEvents(cursor int) (LatestGenericEvents, error) {
	var generic LatestGenericEvents

	routeEvents, err := d.APIClient.GetLastEvents(cursor)
	if err != nil {
		return generic, err
	}
	generic.Cursor = routeEvents.Cursor
	for _, de := range routeEvents.Events {
		generic.Events = append(generic.Events, GenericEvent{
			EntityID: de.EntityID,
			Event: de,
		})
	}

	return generic, nil
}


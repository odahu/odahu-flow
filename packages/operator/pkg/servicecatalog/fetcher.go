package servicecatalog

import (
	event_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
)

type deploymentEventsAPIClient interface {
	GetLastEvents(cursor int) (event_types.LatestDeploymentEvents, error)
}

type DeploymentEventFetcher struct {
	APIClient deploymentEventsAPIClient
}

func (d DeploymentEventFetcher) GetLastEvents(cursor int) (latestGenericEvents, error) {
	var generic latestGenericEvents

	deploymentEvents, err := d.APIClient.GetLastEvents(cursor)
	if err != nil {
		return generic, err
	}
	generic.Cursor = deploymentEvents.Cursor
	for _, de := range deploymentEvents.Events {
		generic.Events = append(generic.Events, de)
	}

	return generic, nil
}



type routeEventsAPIClient interface {
	GetLastEvents(cursor int) (event_types.LatestRouteEvents, error)
}

type RouteEventFetcher struct {
	APIClient routeEventsAPIClient
}

func (d RouteEventFetcher) GetLastEvents(cursor int) (latestGenericEvents, error) {
	var generic latestGenericEvents

	routeEvents, err := d.APIClient.GetLastEvents(cursor)
	if err != nil {
		return generic, err
	}
	generic.Cursor = routeEvents.Cursor
	for _, de := range routeEvents.Events {
		generic.Events = append(generic.Events, de)
	}

	return generic, nil
}


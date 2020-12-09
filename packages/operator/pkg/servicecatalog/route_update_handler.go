package servicecatalog

import (
	"errors"
	"fmt"
	deployment_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	event_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	"go.uber.org/zap"
)

// ModelServerDiscoverer gets URL prefix for a deployed model
// It return Metadata and Swagger2 if passed prefix is served
// by MLServer that is known for ModelServerDiscoverer otherwise
// returns error
type ModelServerDiscoverer interface {
	// Discover ML Server endpoints under the prefix
	// and return model Metadata and Swagger if it is possible
	// return error if Web API behind prefix is not compatible with
	// MLServer
	Discover(prefix string) (model_types.DeployedModel, error)
	// Returns ML Server name for this discoverer
	GetMLServerName() string
}

type Catalog interface {
	Delete(RouteID string)
	CreateOrUpdate(route route) error
}

type UpdateHandler struct {
	Log             *zap.SugaredLogger
	Discoverers     []ModelServerDiscoverer
	Storage         Catalog
	GetDefaultRoute func(ModelDeploymentID string) (deployment_types.ModelRoute, error)
}

type route struct {
	id string
	prefix string

	isDefault bool
	// Only for default route
	model model_types.DeployedModel
}

func (r UpdateHandler) discoverModel(
	prefix string, log *zap.SugaredLogger) (model model_types.DeployedModel, err error) {

	for _, d := range r.Discoverers {

		mlServer := d.GetMLServerName()
		log.Infof("Check whether the model is served by %s", mlServer)
		model, err = d.Discover(prefix)
		if err != nil {
			log.Infow("Discovery is failed",
				zap.Error(err), "MLServer", d.GetMLServerName())
			continue
		}

		log.Info("MLServer behind the prefix is identified. Model metadata is fetched",
			"MLServer", d.GetMLServerName())
		break
	}

	if (model == model_types.DeployedModel{}) {
		err = errors.New("unable to identify MLServer that serves the model")
		log.Error(zap.Error(err))
	}

	return model, err
}

func (r UpdateHandler) Handle(object interface{}) (err error) {

	event, ok := object.(event_types.RouteEvent)
	if !ok {
		return fmt.Errorf("unable to assert received object as event_types.DeploymentEvent: %+v", object)
	}

	log := r.Log.With("EntityID", event.EntityID,
		"EventType", event.EventType,
		"EventDatetime", event.Datetime.String())

	if event.EventType == event_types.ModelRouteDeletedEventType {
		r.Storage.Delete(event.EntityID)
		log.Info("Route was deleted")
		return nil
	}

	route := route{
		id:        event.EntityID,
		prefix:    event.Payload.Spec.URLPrefix,
		isDefault: event.Payload.Default,
	}

	route.model, err = r.discoverModel(route.prefix, log)

	if err != nil {
		return err
	}

	if err != r.Storage.CreateOrUpdate(route) {
		return err
	}

	return nil

}


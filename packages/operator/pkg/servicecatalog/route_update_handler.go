package servicecatalog

import (
	"errors"
	"fmt"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	deployment_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	event_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	"go.uber.org/multierr"
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
	Discover(prefix string, log *zap.SugaredLogger) (model_types.ServedModel, error)
	// Returns ML Server name for this discoverer
	GetMLServerName() model_types.MLServerName
}

type Catalog interface {
	Delete(RouteID string)
	CreateOrUpdate(route route) error
}

type UpdateHandler struct {
	Log             *zap.SugaredLogger
	Discoverers     []ModelServerDiscoverer
	Catalog         Catalog
}


func extractModelDeploymentID(modelRoute deployment_types.ModelRoute) string {
	if modelRoute.Default && len(modelRoute.Spec.ModelDeploymentTargets) > 0 {
		target := modelRoute.Spec.ModelDeploymentTargets[0]
		return target.Name
	}
	return ""
}

func (r UpdateHandler) discoverModel(
	prefix string, log *zap.SugaredLogger) (model model_types.ServedModel, err error) {

	found := false
	var tempErrors error
	for _, d := range r.Discoverers {

		mlServer := d.GetMLServerName()
		log.Infof("Check whether the model is served by %s", mlServer)
		model, err = d.Discover(prefix, log)
		if err != nil {
			log.Infow("Discovery is failed",
				zap.Error(err), "MLServer", d.GetMLServerName())
			if !IsTemporary(err) {
				return model, err  // stop immediately
			}
			tempErrors = multierr.Append(tempErrors, err)
			continue
		}

		log.Info("MLServer behind the prefix is identified. Model metadata is fetched",
			"MLServer", d.GetMLServerName())
		found = true
		break
	}

	if !found {
		err = temporaryErr{
			error: multierr.Combine(errors.New("unable to identify ML server"), tempErrors),
		}
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
		r.Catalog.Delete(event.EntityID)
		log.Info("Route was deleted")
		return nil
	}

	if event.Payload.Status.State != odahuflowv1alpha1.ModelRouteStateReady {
		return temporaryErr{
			error:    fmt.Errorf("ModelRoute %s has not a Ready state", event.EntityID),
		}
	}

	route := route{
		id:        event.EntityID,
		prefix:    event.Payload.Spec.URLPrefix,
		isDefault: event.Payload.Default,
	}

	var servedModel model_types.ServedModel
	servedModel, err = r.discoverModel(route.prefix, log)
	if err != nil {
		return err
	}
	deployedModel := model_types.DeployedModel{
		DeploymentID: extractModelDeploymentID(event.Payload),
		ServedModel:  servedModel,
	}

	route.model = deployedModel

	if err != r.Catalog.CreateOrUpdate(route) {
		return err
	}

	return nil

}


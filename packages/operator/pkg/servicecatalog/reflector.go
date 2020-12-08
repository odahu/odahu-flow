package servicecatalog

import (
	"errors"
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	deployment_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	event_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
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
	Discover(prefix string) (model_types.Info, error)
	// Returns ML Server name for this discoverer
	GetMLServerName() string
}

type Storage interface {
	DeleteByOwnerDeployment(ModelDeploymentID string)
	CreateOrUpdate(route model_types.Info) error
}

type Reflector struct {
	Log         *zap.SugaredLogger
	Discoverers []ModelServerDiscoverer
	Storage     Storage
	GetDefaultRoute func(ModelDeploymentID string) (deployment_types.ModelRoute, error)
}

func (r Reflector) Handle(event event_types.DeploymentEvent) (err error) {

	log := r.Log.With("EntityID", event.EntityID,
		"EventType", event.EventType,
		"EventDatetime", event.Datetime.String())

	if event.EventType == event_types.ModelDeploymentDeletedEventType {
		r.Storage.DeleteByOwnerDeployment(event.EntityID)
		log.Info("Info was deleted from storage")
		return nil
	}

	DefRoute, err := r.GetDefaultRoute(event.EntityID)
	if err != nil {
		log.Errorw("Unable to get default info for ModelDeployment", zap.Error(err))
		return err
	}

	log = log.With("Default ModelRoute prefix", DefRoute.Spec.URLPrefix)

	var info model_types.Info
	for _, d := range r.Discoverers {

		mlServer := d.GetMLServerName()
		log.Infof("Check whether the model is served by %s", mlServer)
		info, err = d.Discover(event.EntityID)
		if err != nil {
			log.Infow("Discovery is failed",
				zap.Error(err), "MLServer", d.GetMLServerName())
			continue
		}

		log.Info("MLServer behind the prefix is identified. Model metadata is fetched",
			"MLServer", d.GetMLServerName())
		if err = r.Storage.CreateOrUpdate(info); err != nil {
			log.Errorw("Unable to create or update info in storage", zap.Error(err))
		}
		break

	}

	if (info == model_types.Info{}) {
		err = errors.New("unable to identify MLServer that serves the model")
		log.Error(zap.Error(err))
		return err
	}

	return nil

}
package servicecatalog

import (
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	"go.uber.org/zap"
)

type temporaryError interface {
	Temporary() bool
}

func IsTemporary(err error) bool {
	tErr, ok := err.(temporaryError)
	return ok && tErr.Temporary()
}

type Route struct {
	ID     string
	Prefix string

	IsDefault bool
	// If route serves more than one model than model is chosen randomly
	Model model_types.DeployedModel
}

type GenericEvent struct {
	EntityID string
	Event    interface{}
}

type LatestGenericEvents struct {
	Cursor int
	Events []GenericEvent
}

type Catalog interface {
	Delete(RouteID string, log *zap.SugaredLogger)
	CreateOrUpdate(route Route) error
}

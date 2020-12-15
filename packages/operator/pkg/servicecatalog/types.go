package servicecatalog

import model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"

type temporaryError interface {
	Temporary() bool
}

func IsTemporary(err error) bool {
	tErr, ok := err.(temporaryError)
	return ok && tErr.Temporary()
}


type temporaryErr struct {
	error
}

func (t temporaryErr) Temporary() bool  {
	return true
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
	Embedded interface{}
}

type LatestGenericEvents struct {
	Cursor int
	Events []GenericEvent
}

type Catalog interface {
	Delete(RouteID string)
	CreateOrUpdate(route Route) error
}
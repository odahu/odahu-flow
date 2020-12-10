package servicecatalog

import model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"

type temporaryErr struct {
	error
}

func (t temporaryErr) Temporary() bool  {
	return true
}

type route struct {
	id string
	prefix string

	isDefault bool
	model model_types.DeployedModel
}

type latestGenericEvents struct {
	Cursor int
	Events []interface{}
}
package model

type Metadata struct {
}

type Swagger2 struct {
}

type MLServerType string

const (
	MLServerODAHU MLServerType = "ODAHU"
	MLServerTriton MLServerType = "Triton"
)

type ServedModel struct {
	// Metadata of a model
	Metadata Metadata
	// MLServer that serves a model
	MLServer MLServerType
}

type Info struct {
	// Info belongs to ModelDeployment with this ID
	DeploymentID string
	// OpenAPI 2.0 Specification for Endpoint API. Nil if model does not support it
	// Swagger URLs already prefixed by URL prefix of default route
	Swagger  *Swagger2
	// WebService contains info about ML Server and serving ML Model
	ServedModel ServedModel
}
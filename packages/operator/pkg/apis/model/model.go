package model

type Metadata struct {
}

type Swagger2 struct {
	Content map[string]interface{}
}

type MLServerName string

const (
	MLServerODAHU  MLServerName = "ODAHU"
	MLServerTriton MLServerName = "Triton"
)

type ServedModel struct {
	// Metadata of a model
	Metadata Metadata
	// MLServer that serves a model
	MLServer MLServerName
	// OpenAPI 2.0 Specification for Endpoint API. Nil if model does not support it
	// Swagger URLs already prefixed by URL prefix of default route
	Swagger Swagger2
}

type DeployedModel struct {
	// DeployedModel belongs to ModelDeployment with this ID
	DeploymentID string
	// WebService contains info about ML Server and serving ML Model
	ServedModel ServedModel
}

package model

// Metadata of a model
type Metadata struct {
	// Optional metadata key, value
	Others map[string]string `json:"others"`
}

// Swagger is base64 encoded OpenAPI 2.0 definition of MLServer API
type Swagger2 struct {
	// Base64 encoded OpenAPI 2.0 definition of MLServer API
	Raw []byte `json:"raw" swaggertype:"string" format:"base64"`
}

type MLServerName string

const (
	MLServerODAHU  MLServerName = "ODAHU"
	MLServerTriton MLServerName = "Triton"
)

// ServedModel contains information about served model
type ServedModel struct {
	Metadata Metadata `json:"metadata"`
	// Possible values: ODAHU, Triton
	Swagger Swagger2 `json:"swagger2"`
}

// DeployedModel contains information about deployed model
type DeployedModel struct {
	// deploymentID is ModelDeployment that deploys this model
	DeploymentID string      `json:"deploymentID"`
	ServedModel  ServedModel `json:"servedModel"`
}

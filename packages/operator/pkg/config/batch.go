package config

type BatchConfig struct {
	// Kubernetes namespace, where BatchInferenceService and BatchInferenceJob will be created
	Namespace string `json:"namespace"`
	// Enable batch API/operator
	Enabled  bool                          `json:"enabled"`
}


func NewDefaultBatchConfig() BatchConfig {
	return BatchConfig{
		Namespace: "odahu-flow-batch",
		Enabled:   true,
	}
}
package config

import (
	corev1 "k8s.io/api/core/v1"
	"os"
	"time"
)

type BatchConfig struct {
	// Kubernetes namespace, where BatchInferenceService and BatchInferenceJob will be created
	Namespace string `json:"namespace"`
	// Enable batch API/operator
	Enabled  bool                          `json:"enabled"`
	// Node pools to run deployments
	NodePools []NodePool `json:"nodePools"`
	// Kubernetes tolerations for batch jobs
	Tolerations []corev1.Toleration        `json:"tolerations,omitempty"`
	// Timeout for full training process
	Timeout time.Duration `json:"timeout"`
	// RClone image that will be used to sync data with object storage
	RCloneImage string  `json:"rcloneImage"`
	// ODAHU tools image
	ToolsImage string  `json:"toolsImage"`
	// ODAHU tools secret name with config
	ToolsSecret string  `json:"toolsSecret"`
}


func NewDefaultBatchConfig() BatchConfig {
	return BatchConfig{
		Namespace: "odahu-flow-batch",
		Enabled:   true,
		Timeout: 4 * time.Hour,
		RCloneImage: "rclone/rclone",
		// workaround https://github.com/spf13/viper/issues/761
		ToolsImage: os.Getenv("ODAHU_TOOLS_IMAGE"),
		ToolsSecret: "tools-config",
	}
}
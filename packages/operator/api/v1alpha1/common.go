package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ResourceList struct {
	// Read more about GPU resource here https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/#using-device-plugins
	GPU *string `json:"gpu,omitempty"`
	// Read more about CPU resource here https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#meaning-of-cpu
	CPU *string `json:"cpu,omitempty"`
	// Read more about memory resource here https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#meaning-of-memory
	Memory *string `json:"memory,omitempty"`
}

type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	Limits *ResourceList `json:"limits,omitempty"`
	// Requests describes the minimum amount of compute resources required.
	Requests *ResourceList `json:"requests,omitempty"`
}

type EnvironmentVariable struct {
	// Name of an environment variable
	Name string `json:"name"`
	// Value of an environment variable
	Value string `json:"value"`
}

type Modifiable struct {
	CreatedAt *metav1.Time `json:"createdAt,omitempty"`
	UpdatedAt *metav1.Time `json:"updatedAt,omitempty"`
}

type VCS struct {
	ConnName  string `json:"connName,omitempty"`
	Reference string `json:"reference,omitempty"`
}

type ObjectStorage struct {
	ConnName string `json:"connName,omitempty"`
	Path     string `json:"path,omitempty"`
}

type AlgorithmSource struct {
	VCS           VCS           `json:"vcs,omitempty"`
	ObjectStorage ObjectStorage `json:"objectStorage,omitempty"`
}

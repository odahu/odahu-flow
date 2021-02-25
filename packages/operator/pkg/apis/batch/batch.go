package batch

import "time"

type PredictorWebhookTrigger struct {
	Enabled bool
}

type InferenceServiceTriggers struct {
	Webhook PredictorWebhookTrigger
}

type InferenceServiceSpec struct {
	// Docker image
	Image string `json:"image"`
	// Entrypoint array. Not executed within a shell.
	// The docker image's ENTRYPOINT is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax
	// can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	Command []string `json:"command"`
	// Arguments to the entrypoint.
	// The docker image's CMD is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax
	// can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	Args []string `json:"args"`
	InputConnection string
	OutputConnection string
	ModelConnection string
	ModelPath string
	Triggers InferenceServiceTriggers
}


type InferenceJobID string
type InferenceJobRun struct {
	LastState string
	Date time.Time
}

type InferenceServiceStatus struct {
	Triggers InferenceServiceTriggers
	Runs map[InferenceJobID]InferenceJobRun
}

type InferenceService struct {
	ID string `json:"id"`
	Spec InferenceServiceSpec
	Status InferenceServiceStatus
}

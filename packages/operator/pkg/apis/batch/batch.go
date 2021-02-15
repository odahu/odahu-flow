package batch

type PredictorWebhookTrigger struct {
	Enabled bool
}

type PredictorWebhookTriggerStatus struct {
	Endpoint string
}

type PredictorTriggers struct {
	Webhook PredictorWebhookTrigger
}

type PredictorTriggersStatus struct {
	Webhook PredictorWebhookTriggerStatus
}

type PredictorSpec struct {
	Image string `json:"image"`
	Entrypoint string
	Cmd []string
	InputConnection string
	OutputConnection string
	ModelConnection string
	Triggers PredictorTriggers
}

type PredictorStatus struct {
	Triggers PredictorTriggersStatus
}

type Predictor struct {
	ID string `json:"id"`
	Spec PredictorSpec
	Status PredictorStatus
}

type PredictionJobSpec struct {
	Entrypoint string
	Cmd []string
	InputPath string
	OutputPath string
	ModelPath string
}

type PredictionJobStatus struct {
	State string
	Trigger string
	TriggerDate string
}

type PredictionJob struct {
	ID string `json:"id"`
	Spec PredictionJobSpec
	Status PredictionJobStatus
}
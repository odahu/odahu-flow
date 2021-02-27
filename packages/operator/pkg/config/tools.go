package config

type FluentdConfig struct {
	BaseURL string `json:"baseURL"`
}

type FeedbackConfig struct {
	Fluentd FluentdConfig `json:"fluentd"`
}

type ToolsConfig struct {
	Auth AuthConfig `json:"auth"`
	Feedback FeedbackConfig `json:"feedback"`
}
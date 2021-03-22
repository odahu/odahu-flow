package feedback

type RequestResponse struct {
	RequestID           string            `msg:"request_id"`
	RequestHttpHeaders  map[string]string `msg:"request_http_headers"`
	RequestContent      string            `msg:"request_content"`
	RequestUri          string            `msg:"request_uri"`
	ResponseStatus      string            `msg:"response_status"`
	ResponseHttpHeaders map[string]string `msg:"response_http_headers"`
	RequestHost         string            `msg:"request_host"`
	ModelVersion        string            `msg:"model_version"`
	ModelName           string            `msg:"model_name"`
	RequestHttpMethod   string            `msg:"request_http_method"`
}

type ResponseBody struct {
	RequestID       string `msg:"request_id"`
	ModelVersion    string `msg:"model_version"`
	ModelName       string `msg:"model_name"`
	ResponseContent string `msg:"response_content"`
}

type ModelFeedback struct {
	RequestID    string                 `msg:"request_id"`
	ModelVersion string                 `msg:"model_version"`
	ModelName    string                 `msg:"model_name"`
	Payload      map[string]interface{} `msg:"payload"`
}


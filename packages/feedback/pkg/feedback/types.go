package feedback

const (
	ModelNameHeaderKey          = "model-name"
	ModelVersionHeaderKey       = "model-version"
	ModelEndpointHeaderKey      = "model-endpoint"
	RequestIdHeaderKey          = "x-request-id"
	OdahuFlowRequestIdHeaderKey = "request-id"
	HttpMethodHeaderKey         = ":method"
	OriginalUriHeaderKey        = "x-original-uri"
	StatusHeaderKey             = ":status"
	ForwardedHostHeaderKey      = "x-forwarded-host"
)

type DataLogging interface {
	Post(tag string, message interface{}) error
	Close() error
}

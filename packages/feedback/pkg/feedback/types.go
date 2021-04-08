package feedback

const (
	ModelNameHeaderKey               = "model-name"
	ModelVersionHeaderKey            = "model-version"
	OdahuFlowRequestIdHeaderKey      = "request-id"
	HttpMethodHeaderKey              = ":method"
	OriginalUriHeaderKey             = "x-original-uri"
	StatusHeaderKey                  = ":status"
	ForwardedHostHeaderKey           = "x-forwarded-host"
	EnvoyInternalRoutingHeader       = "x-envoy-decorator-operation"
	EnvoyInternalRoutingHeaderPrefix = "istio-ingressgateway"
)

type DataLogging interface {
	Post(tag string, message interface{}) error
	Close() error
}

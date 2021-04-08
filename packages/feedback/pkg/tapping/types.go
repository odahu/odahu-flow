package tapping

const (
	CfgEnvoyHost     = "envoy.host"
	CfgEnvoyPort     = "envoy.port"
	CfgEnvoyConfigId = "envoy.config_id"
)

// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/tap/v3/common.proto#config-tap-v3-matchpredicate
type MatchPredicate struct {
	OrMatch                  MatchSet         `yaml:"or_match,omitempty"`
	AndMatch                 MatchSet         `yaml:"and_match,omitempty"`
	NotMatch                 *MatchPredicate  `yaml:"not_match,omitempty"`
	AnyMatch                 bool             `yaml:"any_match,omitempty"`
	HttpRequestHeadersMatch  HttpHeadersMatch `yaml:"http_request_headers_match,omitempty"`
	HttpResponseHeadersMatch HttpHeadersMatch `yaml:"http_response_headers_match,omitempty"`
}

// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/tap/v3/common.proto#envoy-v3-api-msg-config-tap-v3-matchpredicate-matchset
type MatchSet struct {
	Rules []MatchPredicate `yaml:"rules"`
}

type HttpHeadersMatch struct {
	Headers []HeaderMatcher `yaml:"headers"`
}

// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#config-route-v3-headermatcher
type HeaderMatcher struct {
	Name          string `yaml:"name"`
	ExactMatch    string `yaml:"exact_match,omitempty"`
	RegexMatch    string `yaml:"regex_match,omitempty"`
	RangeMatch    string `yaml:"range_match,omitempty"`
	PresentMatch  bool   `yaml:"present_match,omitempty"`
	PrefixMatch   string `yaml:"prefix_match,omitempty"`
	SuffixMatch   string `yaml:"suffix_match,omitempty"`
	ContainsMatch string `yaml:"contains_match,omitempty"`
	InvertMatch   bool   `yaml:"invert_match,omitempty"`
}

type TapSink struct {
	StreamingAdmin map[string]string `yaml:"streaming_admin"`
}

type TapRequest struct {
	ConfigID  string `yaml:"config_id"`
	TapConfig struct {
		MatchConfig  MatchPredicate `yaml:"match_config"`
		OutputConfig struct {
			Sinks              []TapSink `yaml:"sinks"`
			MaxBufferedTxBytes int32     `yaml:"max_buffered_tx_bytes"`
			MaxBufferedRxBytes int32     `yaml:"max_buffered_rx_bytes"`
		} `yaml:"output_config"`
	} `yaml:"tap_config"`
}

type Trace struct {
	Headers []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"headers"`
	Body struct {
		Truncated bool   `json:"truncated"`
		AsBytes   string `json:"as_bytes"`
	} `json:"body"`
}

type Message struct {
	HttpBufferedTrace struct {
		Request  Trace `json:"request,omitempty"`
		Response Trace `json:"response,omitempty"`
	} `json:"http_buffered_trace"`
}

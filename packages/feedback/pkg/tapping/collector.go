package tapping

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/odahu/odahu-flow/packages/feedback/pkg/feedback"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	commons_feedback "odahu-commons/feedback"
	"odahu-commons/predictors"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

const (
	TapUrl            = "/tap"
	filterHeaderKey   = ":path"
	defaultBufferSize = 1024 * 1024
)

var log = logf.Log.WithName("collector")

type RequestCollector struct {
	envoyUrl            string
	feedbackRequestYaml []byte
	logger              feedback.DataLogging
	prohibitedHeaders   map[string]string
}

func NewRequestCollector(
	envoyHost string,
	envoyPort int,
	configId string,
	logger feedback.DataLogging,
	prohibitedHeaders []string,
) (*RequestCollector, error) {
	feedbackRequest := TapRequest{
		ConfigID: configId,
	}

	predictorRules := make([]MatchPredicate, 0, len(predictors.Predictors))
	for _, predictor := range predictors.Predictors {
		predictorRules = append(predictorRules, MatchPredicate{
			HttpRequestHeadersMatch: HttpHeadersMatch{
				Headers: []HeaderMatcher{{
					Name:       filterHeaderKey,
					RegexMatch: predictor.InferenceEndpointRegex,
				}},
			},
		})
	}

	feedbackRequest.TapConfig.MatchConfig = MatchPredicate{
		AndMatch: MatchSet{Rules: []MatchPredicate{
			{
				OrMatch: MatchSet{Rules: predictorRules},
			},
			{
				// This rule is necessary because every inference request goes through 2 VirtualServices:
				// ODAHU's one and Knative's one. So the request goes through the Istio Ingress twice.
				// The second iteration is internal and should be filtered out.
				// The transparent way to do that is to filter it basing on Knative-specific headers, which
				// appear in only the second iteration, but for some reason Envoy is unable to work with them.
				// Seems like a bug in Envoy. Example of more self-explaining rule:
				//HttpRequestHeadersMatch: HttpHeadersMatch{Headers: []HeaderMatcher{
				//	{
				//		Name: "knative-serving-revision",
				//		PresentMatch: true,
				//		InvertMatch: true,
				//	},
				//}},
				// During experiments with different rules this one was discovered as a working one
				HttpRequestHeadersMatch: HttpHeadersMatch{Headers: []HeaderMatcher{
					{
						Name:        feedback.EnvoyInternalRoutingHeader,
						PrefixMatch: feedback.EnvoyInternalRoutingHeaderPrefix,
					},
				}},
			},
		}},
	}

	feedbackRequest.TapConfig.OutputConfig.Sinks = append(
		feedbackRequest.TapConfig.OutputConfig.Sinks,
		TapSink{StreamingAdmin: map[string]string{}},
	)
	feedbackRequest.TapConfig.OutputConfig.MaxBufferedRxBytes = defaultBufferSize
	feedbackRequest.TapConfig.OutputConfig.MaxBufferedTxBytes = defaultBufferSize
	feedbackRequestYaml, err := yaml.Marshal(feedbackRequest)
	if err != nil {
		log.Error(err, "Tapping request")

		return nil, err
	}
	log.Info("generated tapping request", "request_yaml", string(feedbackRequestYaml))

	prohibitedHeadersMap := make(map[string]string, len(prohibitedHeaders))
	for _, header := range prohibitedHeaders {
		prohibitedHeadersMap[header] = ""
	}

	return &RequestCollector{
		envoyUrl:            fmt.Sprintf("%s:%d", envoyHost, envoyPort),
		feedbackRequestYaml: feedbackRequestYaml,
		logger:              logger,
		prohibitedHeaders:   prohibitedHeadersMap,
	}, nil
}

func (rc *RequestCollector) convertToFeedback(message *Message) (*commons_feedback.RequestResponse, *commons_feedback.ResponseBody, error) {
	responseBody := &commons_feedback.ResponseBody{}
	requestResponse := &commons_feedback.RequestResponse{}

	requestHeaders := make(map[string]string, len(message.HttpBufferedTrace.Request.Headers))
	for _, header := range message.HttpBufferedTrace.Request.Headers {
		if _, ok := rc.prohibitedHeaders[header.Key]; ok {
			continue
		}

		switch header.Key {
		case feedback.HttpMethodHeaderKey:
			requestResponse.RequestHttpMethod = header.Value

		case feedback.OriginalUriHeaderKey:
			requestResponse.RequestUri = header.Value

		case feedback.ForwardedHostHeaderKey:
			requestResponse.RequestHost = header.Value

		case feedback.OdahuFlowRequestIdHeaderKey:
			responseBody.RequestID = header.Value
			requestResponse.RequestID = header.Value
		}

		requestHeaders[header.Key] = header.Value
	}

	responseHeaders := make(map[string]string, len(message.HttpBufferedTrace.Response.Headers))
	for _, header := range message.HttpBufferedTrace.Response.Headers {
		if _, ok := rc.prohibitedHeaders[header.Key]; ok {
			continue
		}

		switch header.Key {
		case feedback.ModelNameHeaderKey:
			responseBody.ModelName = header.Value
			requestResponse.ModelName = header.Value

		case feedback.ModelVersionHeaderKey:
			responseBody.ModelVersion = header.Value
			requestResponse.ModelVersion = header.Value

		case feedback.StatusHeaderKey:
			requestResponse.ResponseStatus = header.Value
		}

		responseHeaders[header.Key] = header.Value
	}

	requestResponse.RequestHttpHeaders = requestHeaders
	requestResponse.ResponseHttpHeaders = responseHeaders

	responseBytes, err := base64.StdEncoding.DecodeString(message.HttpBufferedTrace.Response.Body.AsBytes)
	if err != nil {
		log.Error(err, "Encoding response body")

		return nil, nil, err
	}
	responseBody.ResponseContent = string(responseBytes)

	requestBytes, err := base64.StdEncoding.DecodeString(message.HttpBufferedTrace.Request.Body.AsBytes)
	if err != nil {
		log.Error(err, "Encode request body")

		return nil, nil, err
	}
	requestResponse.RequestContent = string(requestBytes)

	return requestResponse, responseBody, nil
}

func (rc *RequestCollector) TraceRequests() error {
	for {
		if err := rc.tapTraffic(); err != nil {
			errorTapping.Add(1)

			log.Error(err, "Traffic tapping")
			time.Sleep(1 * time.Second)
		}
	}
}

func (rc *RequestCollector) tapTraffic() error {
	req := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "http",
			Host:   rc.envoyUrl,
			Path:   TapUrl,
		},
		Body: ioutil.NopCloser(bytes.NewBuffer(rc.feedbackRequestYaml)),
	}
	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   0,
	}

	log.Info("tap request dump", "yaml", string(rc.feedbackRequestYaml))

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)
	for dec.More() {
		collectedRequests.Add(1)
		var message Message

		err := dec.Decode(&message)
		if err != nil {
			return err
		}

		requestResponse, responseBody, err := rc.convertToFeedback(&message)
		if err != nil {
			return err
		}

		log.Info("logged request", "modelName", requestResponse.ModelName,
			"modelVersion", requestResponse.ModelVersion,
			"request_id", responseBody.RequestID)

		err = rc.logger.Post(viper.GetString(feedback.CfgRequestResponseTag), *requestResponse)
		if err != nil {
			return err
		}

		err = rc.logger.Post(viper.GetString(feedback.CfgResponseBodyTag), *responseBody)
		if err != nil {
			return err
		}
	}

	return err
}

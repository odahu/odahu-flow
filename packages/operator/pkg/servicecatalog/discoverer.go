package servicecatalog

import (
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type OdahuMLServerDiscoverer struct {
	EdgeHost   string
	EdgeURL    url.URL
	HTTPClient httpClient
}

// ODAHU ML Server has not hardcoded swagger spec. But claims that "GET /" request
// return Swagger 2.0 spec for current deployed Model
func (o OdahuMLServerDiscoverer) discoverSwagger(
	prefix string, log *zap.SugaredLogger) (swagger model_types.Swagger2, err error) {
	modelRequest := o.generateModelRequest(prefix)

	var response *http.Response
	response, err = o.HTTPClient.Do(modelRequest)
	if err != nil {
		log.Error(
			err, "Can not get swagger response for prefix",
		)
		return swagger, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Errorw("Unable to close response body", zap.Error(err))
		}
	}()

	rawBody, err := ioutil.ReadAll(response.Body)
	log.Info("Get response from model", "content", string(rawBody))


	swagger = model_types.Swagger2{Raw: rawBody}
	return swagger, nil

}


// ODAHU ML Server currently does not support metadata endpoints
func (o OdahuMLServerDiscoverer) discoverMetadata(
	prefix string, log *zap.SugaredLogger) (metadata model_types.Metadata, err error) {
	return metadata, nil
}

func (o OdahuMLServerDiscoverer) Discover(
	prefix string, log *zap.SugaredLogger) (model model_types.ServedModel, err error) {

	model.Swagger, err = o.discoverSwagger(prefix, log)
	if err != nil {
		return
	}
	model.MLServer = o.GetMLServerName()

	model.Metadata, err = o.discoverMetadata(prefix, log)

	return

}

func (o OdahuMLServerDiscoverer) GetMLServerName() model_types.MLServerName {
	 return model_types.MLServerODAHU
}


func (o OdahuMLServerDiscoverer) generateModelRequest(prefix string) *http.Request {

	MlServerURL := url.URL{
		Scheme: o.EdgeURL.Scheme,
		Host: o.EdgeURL.Host,
		Path: path.Join(o.EdgeURL.Path, prefix),

	}

	return &http.Request{
		Method: http.MethodGet,
		URL:    &MlServerURL,
		Host:   o.EdgeHost,
	}
}
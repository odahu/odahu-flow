package event

import (
	"encoding/json"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ModelRouteEventClient struct {
	HTTPClient httpClient
	Log        *zap.SugaredLogger
}


func (m ModelRouteEventClient) GetLastEvents(cursor int) (events event.LatestRouteEvents, err error) {
	

	var buf []byte
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{
			Path:     "/model/route-events",
		},
	}
	q := req.URL.Query()
	q.Add("cursor", strconv.Itoa(cursor))
	req.URL.RawQuery = q.Encode()
	response, err := m.HTTPClient.Do(req)
	if err != nil {
		return events, err
	}
	defer response.Body.Close()
	buf, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return events, err
	}

	err = json.Unmarshal(buf, &events)
	if err != nil {
		m.Log.Error("Unable to unmarshall ModelRoute events", zap.Error(err))
	}

	return events, err

}





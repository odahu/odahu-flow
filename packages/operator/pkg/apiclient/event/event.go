package event

import (
	"encoding/json"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/event"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type httpClient interface {
	DoRequestGetBody(httpMethod, path string, body interface{}) ([]byte, error)
}

type ModelRouteEventClient struct {
	HTTPClient httpClient
	Log        *zap.SugaredLogger
}


func (m ModelRouteEventClient) GetLastEvents(cursor int) (events event.LatestRouteEvents, err error) {
	

	var buf []byte
	buf, err = m.HTTPClient.DoRequestGetBody(
		http.MethodGet,
		strings.Replace("/model/route-events?cursor=:cursor", ":cursor", strconv.Itoa(cursor), 1),
		nil,
	)
	if err != nil {
		m.Log.Errorw("Unable to fetch last ModelRoute events", zap.Error(err))
		return events, err
	}

	err = json.Unmarshal(buf, &events)
	if err != nil {
		m.Log.Error("Unable to unmarshall ModelRoute events", zap.Error(err))
	}

	return events, err

}





package event

import (
	"encoding/json"
	"fmt"
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

type temporaryErr struct {
	error
}

func (t temporaryErr) Temporary() bool {
	return true
}


func isTemporaryStatusCode(code int) bool {
	for tempCode := range []int{
		http.StatusRequestTimeout,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,

		// Client can be used in reactive workloads (background workers) that suppose backoff retries
		// So workload can retry attempt to get data after error on server will be
		// fixed
		http.StatusInternalServerError,
	} {
		if code == tempCode {
			return true
		}
	}

	return false
}

func (m ModelRouteEventClient) GetLastEvents(cursor int) (events event.LatestRouteEvents, err error) {
	
	var response *http.Response

	u := url.URL{
		Path:       "/model/route-events",
	}
	q := u.Query()
	q.Set("cursor", strconv.Itoa(cursor))
	u.RawQuery = q.Encode()

	response, err = m.HTTPClient.Do(&http.Request{
		Method:           http.MethodGet,
		URL:             &u,
	})
	if err != nil {
		m.Log.Errorw("Retrieving of the ModelRoute events is failed", zap.Error(err))
		return events, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			m.Log.Error("Unable to close connection", zap.Error(err))
		}
	}()

	var bytes []byte
	bytes, err =  ioutil.ReadAll(response.Body)
	if err != nil {
		m.Log.Errorw("Unable to read response body", zap.Error(err))
		return events, err
	}

	if response.StatusCode >= 400 {
		if isTemporaryStatusCode(response.StatusCode) {
			return events, temporaryErr{
				error: fmt.Errorf(
					"not correct status code: %d. Maybe temporary. Response body: %s",
					response.StatusCode, bytes,
				),
			}
		}
		return events, fmt.Errorf(
			"not correct status code: %d. Response body: %s",
			response.StatusCode, bytes,
		)
	}



	err = json.Unmarshal(bytes, &events)
	if err != nil {
		m.Log.Error("Unable to unmarshall ModelRoute events", zap.Error(err))
	}

	return events, err

}





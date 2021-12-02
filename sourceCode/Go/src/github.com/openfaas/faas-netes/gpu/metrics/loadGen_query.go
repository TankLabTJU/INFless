package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// LoadGenQuery represents parameters for querying LoadGen
type LoadGenQuery struct {
	Port   int
	Host   string
	Client *http.Client
}

type LoadGenQueryFetcher interface {
	Fetch(query string) (*[]*LoadGenQueryResponse, error)
}

// NewPrometheusQuery create a NewPrometheusQuery
func NewLoadGenQuery(host string, port int, client *http.Client) LoadGenQuery {
	return LoadGenQuery{
		Client: client,
		Host:   host,
		Port:   port,
	}
}

// Fetch queries aggregated stats
func (q LoadGenQuery) Fetch() (*[]LoadGenQueryResponse, error) {

	req, reqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/loadGen/getLoaderGenQuery.do", q.Host, q.Port), nil)
	if reqErr != nil {
		return nil, reqErr
	}

	res, getErr := q.Client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	bytesOut, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code from LoadGen want: %d, got: %d, body: %s\n", http.StatusOK, res.StatusCode, string(bytesOut))
	}
	var values []LoadGenQueryResponse

	unmarshalErr := json.Unmarshal(bytesOut, &values)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("Error unmarshaling result: %s, '%s'\n", unmarshalErr, string(bytesOut))
	}

	return &values, nil
}
type LoadGenQueryResponse struct {
	ServiceRate float64 `json:"windowAvgServiceRate"`
	RealRps int32 `json:"realRps"`
	QueryTime99th float64 `json:"queryTime99th"`
	LoaderName string `json:"loaderName"`
}

// Fetch queries aggregated stats
/*func (q PrometheusQuery) Delete(query string) (*VectorQueryResponse, error) {

	req, reqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/api/v1/admin/tsdb/delete_series?match[]=%s", q.Host, q.Port, query), nil)
	if reqErr != nil {
		return nil, reqErr
	}

	res, getErr := q.Client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	bytesOut, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code from Prometheus want: %d, got: %d, body: %s", http.StatusOK, res.StatusCode, string(bytesOut))
	}

	var values VectorQueryResponse
	unmarshalErr := json.Unmarshal(bytesOut, &values)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("Error unmarshaling result: %s, '%s'", unmarshalErr, string(bytesOut))
	}

	return &values, nil
}
type VectorDeleteResponse struct {
	Data struct {
		Result []struct {
			Metric struct {
				Code         string `json:"code"`
				FunctionName string `json:"function_name"`
				PodName      string `json:"pod"`
			}
			Value []interface{} `json:"value"`
		}
	}
}/*
// Fetch queries aggregated stats
/*func (q PrometheusQuery) FetchPod(query string) (*VectorQueryPodResponse, error) {

	req, reqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/api/v1/query?query=%s", q.Host, q.Port, query), nil)
	if reqErr != nil {
		return nil, reqErr
	}

	res, getErr := q.Client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	bytesOut, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code from Prometheus want: %d, got: %d, body: %s", http.StatusOK, res.StatusCode, string(bytesOut))
	}

	var values VectorQueryPodResponse

	unmarshalErr := json.Unmarshal(bytesOut, &values)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("Error unmarshaling result: %s, '%s'", unmarshalErr, string(bytesOut))
	}

	return &values, nil
}*/

/*
type VectorQueryPodResponse struct {
	Data struct {
		Result []struct {
			Metric struct {
				PodName string `json:"pod"`
			}
			Value []interface{} `json:"value"`
		}
	}
}
*/
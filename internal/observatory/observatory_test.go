package observatory

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateObservatory(t *testing.T) {
	observatory := New(newStubHTTPClient(""))
	assert.NotEqual(t, nil, observatory)
}

func TestGetSPWRealTimeTide(t *testing.T) {
	observatory := New(newStubHTTPClient(`
QUB,2019-04-20 14:40,0.97,
SPW,2019-04-20 14:40,0.77,
TBT,2019-04-20 14:40,1.59,
TMW,2019-04-20 14:40,0.73,
TPK,2019-04-20 14:40,1.02,
WAG,2019-04-20 14:40,----,M
CCH,2019-04-20 14:40,0.87,
CLK,2019-04-20 14:40,1.29,
KCT,2019-04-20 14:40,0.90,
KLW,2019-04-20 14:40,----,
MWC,2019-04-20 14:40,1.04,
TAO,2019-04-20 14:40,0.82,
`))

	data, _ := observatory.TidalDataAsOf(time.Unix(0, 0))
	assert.Equal(t, TidalData{
		Seconds: 1555742400,
		Nanos:   0,
		Meters:  0.77,
	}, data)
}

func TestSPWDataDoesNotExist(t *testing.T) {
	observatory := New(newStubHTTPClient("QUB,2019-04-20 14:40,0.97,"))
	_, err := observatory.TidalDataAsOf(time.Unix(0, 0))
	assert.EqualError(t, err, "Tidal data unavailable")
}

func TestSPWDataExistsButUnavailable(t *testing.T) {
	observatory := New(newStubHTTPClient("SPW,2019-04-20 14:40,----,"))
	_, err := observatory.TidalDataAsOf(time.Unix(0, 0))
	assert.EqualError(t, err, "Tidal data unavailable")
}

func TestGetPredictions(t *testing.T) {
	observatory := New(newStubHTTPClient(`
2019-04-22 22:00,----,1.35
2019-04-22 23:00,----,1.52
`))
	april22, _ := time.Parse("2006-01-02 15:04", "2019-04-22 00:00")
	predictions, _ := observatory.TidalPredictionsAsOf(april22)
	expected := []TidalData{
		{
			Seconds: 1555941600,
			Meters:  1.35,
		},
		{
			Seconds: 1555945200,
			Meters:  1.52,
		},
	}
	assert.Equal(t, expected, predictions)
}

func TestGetPredictionsShouldSkipDataNotOnGivenDate(t *testing.T) {
	observatory := New(newStubHTTPClient(`
2019-04-23 22:00,----,2.0
2019-04-23 23:00,----,2.1
	`))
	april22, _ := time.Parse("2006-01-02 15:04", "2019-04-22 00:00")
	predictions, _ := observatory.TidalPredictionsAsOf(april22)
	assert.Empty(t, predictions)
}

func newStubHTTPClient(responseBody string) *http.Client {
	return &http.Client{
		Transport: &mockTransport{
			responseBody: responseBody,
		},
	}
}

type mockTransport struct {
	responseBody string
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}
	response.Body = ioutil.NopCloser(strings.NewReader(t.responseBody))
	return response, nil
}

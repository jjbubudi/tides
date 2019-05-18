package tides

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jjbubudi/protos-go/tides"
	"github.com/jjbubudi/tides/internal/observatory"
	"github.com/stretchr/testify/assert"
)

func TestCreateService(t *testing.T) {
	service := New(nilData, nilPredictions, noOpPublisher)
	assert.NotNil(t, service)
}

func TestPublishRealTimeTidalData(t *testing.T) {
	var publishedSubject string
	var publishedData tides.TideRecorded

	var tidalData = func(when time.Time) (observatory.TidalData, error) {
		return observatory.TidalData{
			Seconds: 10000,
			Nanos:   0,
			Meters:  2.0,
		}, nil
	}

	var publisher = func(subject string, data []byte) error {
		publishedSubject = subject
		proto.Unmarshal(data, &publishedData)
		return nil
	}

	service := New(tidalData, nilPredictions, publisher)
	service.PublishRealTimeTidalData(time.Unix(0, 0))

	expectedData := tides.TideRecorded{
		Time: &timestamp.Timestamp{
			Seconds: 10000,
			Nanos:   0,
		},
		Meters: 2.0,
	}

	assert.Equal(t, "tides", publishedSubject)
	assert.Equal(t, expectedData, publishedData)
}

func TestPublishTidalPredictions(t *testing.T) {
	var publishedSubject string
	var publishedData tides.TidePredicted

	var tidalPredictions = func(when time.Time) ([]observatory.TidalData, error) {
		return []observatory.TidalData{
			observatory.TidalData{
				Seconds: 15000,
				Nanos:   0,
				Meters:  2.0,
			},
		}, nil
	}

	var publisher = func(subject string, data []byte) error {
		publishedSubject = subject
		proto.Unmarshal(data, &publishedData)
		return nil
	}

	service := New(nilData, tidalPredictions, publisher)
	service.PublishTidalPredictions(time.Unix(10000, 0))

	expectedData := tides.TidePredicted{
		Time: &timestamp.Timestamp{
			Seconds: 10000,
			Nanos:   0,
		},
		Predictions: []*tides.TidePredicted_Prediction{
			&tides.TidePredicted_Prediction{
				Time: &timestamp.Timestamp{
					Seconds: 15000,
					Nanos:   0,
				},
				Meters: 2.0,
			},
		},
	}

	assert.Equal(t, "tides_predictions", publishedSubject)
	assert.Equal(t, expectedData, publishedData)
}

func nilData(when time.Time) (observatory.TidalData, error) {
	return observatory.TidalData{}, nil
}

func nilPredictions(when time.Time) ([]observatory.TidalData, error) {
	return nil, nil
}

func noOpPublisher(subject string, data []byte) error {
	return nil
}

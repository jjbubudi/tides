package tides

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jjbubudi/protos-go/tides"
	"github.com/jjbubudi/tides/pkg/observatory"
)

const tidesChannel = "tides"
const tidesPredictionsChannel = "tides_predictions"

type tidalData func(time.Time) (observatory.TidalData, error)
type tidalPredictions func(time.Time) ([]observatory.TidalData, error)
type publisher func(string, []byte) error

// Service loads and publishes tidal data
type Service struct {
	tidalData        tidalData
	tidalPredictions tidalPredictions
	publish          publisher
}

// New returns a new Service instance
func New(tidalData tidalData, tidalPredictions tidalPredictions, publisher publisher) *Service {
	return &Service{
		tidalData:        tidalData,
		tidalPredictions: tidalPredictions,
		publish:          publisher,
	}
}

// PublishRealTimeTidalData loads and publishes tidal data
func (s *Service) PublishRealTimeTidalData(when time.Time) error {
	tidalData, err := s.tidalData(when)
	if err != nil {
		return err
	}

	event := &tides.TideRecorded{
		Time: &timestamp.Timestamp{
			Seconds: tidalData.Seconds,
			Nanos:   int32(tidalData.Nanos),
		},
		Meters: tidalData.Meters,
	}
	bytes, _ := proto.Marshal(event)

	if publishErr := s.publish(tidesChannel, bytes); publishErr != nil {
		return publishErr
	}

	return nil
}

// PublishTidalPredictions loads and publishes tidal predictions as of the given date
func (s *Service) PublishTidalPredictions(when time.Time) error {
	predictions, err := s.tidalPredictions(when)
	if err != nil {
		return err
	}

	var predictionsProto []*tides.TidePredicted_Prediction
	for _, prediction := range predictions {
		predictionsProto = append(predictionsProto, &tides.TidePredicted_Prediction{
			Time: &timestamp.Timestamp{
				Seconds: prediction.Seconds,
				Nanos:   int32(prediction.Nanos),
			},
			Meters: prediction.Meters,
		})
	}

	event := &tides.TidePredicted{
		Time: &timestamp.Timestamp{
			Seconds: when.Unix(),
			Nanos:   int32(when.Nanosecond()),
		},
		Predictions: predictionsProto,
	}
	bytes, _ := proto.Marshal(event)

	if err := s.publish(tidesPredictionsChannel, bytes); err != nil {
		return err
	}

	return nil
}

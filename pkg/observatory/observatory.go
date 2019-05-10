package observatory

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"
	"time"
)

const defaultStation = "SPW"
const realTimeURL = "https://www.hko.gov.hk/tide/marine/data/ALL.txt"
const predictionsURL = "https://www.hko.gov.hk/tide/marine/data/SPW.tab.txt"

var hongKong, _ = time.LoadLocation("Asia/Hong_Kong")

// TidalData represents the tidal data we get from HKO
type TidalData struct {
	Seconds int64
	Nanos   uint32
	Meters  float64
}

// Observatory provides tidal data from the HKO
type Observatory interface {
	TidalDataAsOf(when time.Time) (TidalData, error)
	TidalPredictionsAsOf(when time.Time) ([]TidalData, error)
}

// HTTPObservatory is a type that can query tidal data from the HKO
type HTTPObservatory struct {
	client *http.Client
}

// New returns a new Observatory given an http client
func New(client *http.Client) Observatory {
	return &HTTPObservatory{
		client: client,
	}
}

// TidalDataAsOf returns the real time tide data as of the given time
func (o *HTTPObservatory) TidalDataAsOf(when time.Time) (TidalData, error) {
	currentTime := strconv.FormatInt(when.UnixNano()/int64(time.Millisecond), 10)
	response, responseError := o.client.Get(realTimeURL + "?_=" + currentTime)
	if responseError != nil {
		return TidalData{}, responseError
	}
	defer response.Body.Close()

	reader := csv.NewReader(response.Body)
	records, csvError := reader.ReadAll()
	if csvError != nil {
		return TidalData{}, csvError
	}

	for _, record := range records {
		station := record[0]
		if station != defaultStation {
			continue
		}

		asOf, _ := time.ParseInLocation("2006-01-02 15:04", record[1], hongKong)
		meters, parseError := strconv.ParseFloat(record[2], 64)
		if parseError != nil {
			continue
		}

		return TidalData{
			Seconds: asOf.Unix(),
			Nanos:   uint32(asOf.Nanosecond()),
			Meters:  meters,
		}, nil
	}

	return TidalData{}, errors.New("Tidal data unavailable")
}

// TidalPredictionsAsOf returns the tidal predictions as of the given date
func (o *HTTPObservatory) TidalPredictionsAsOf(when time.Time) ([]TidalData, error) {
	response, responseError := o.client.Get(predictionsURL)
	if responseError != nil {
		return nil, responseError
	}
	defer response.Body.Close()

	reader := csv.NewReader(response.Body)
	records, csvError := reader.ReadAll()
	if csvError != nil {
		return nil, csvError
	}

	var data []TidalData
	for _, record := range records {
		asOf, _ := time.ParseInLocation("2006-01-02 15:04", record[0], hongKong)
		meters, _ := strconv.ParseFloat(record[2], 64)

		if asOf.Year() != when.Year() || asOf.YearDay() != when.YearDay() {
			continue
		}

		data = append(data, TidalData{
			Seconds: asOf.Unix(),
			Nanos:   uint32(asOf.Nanosecond()),
			Meters:  meters,
		})
	}

	return data, nil
}

package request

import (
	"time"

	"github.com/conradhodge/next-tram/src/lib/traveline"
	"github.com/google/uuid"
)

// Requests represents the requests that can be made using the Traveline API
type Requests interface {
	GetNextTramTime(naptanCode string, when time.Time) (string, error)
}

// requests is used to make API requests using the Traveline API
type requests struct {
	API traveline.API
}

// NewRequest returns a new instance of the Traveline API requests
func NewRequest(api traveline.API) Requests {
	return &requests{API: api}
}

// GetNextTramTime returns the time of the next tram at the stop that the NaPTAN code represents
func (req *requests) GetNextTramTime(naptanCode string, when time.Time) (string, error) {
	request, err := req.API.BuildRequest(uuid.New().String(), naptanCode, when)
	if err != nil {
		return "", err
	}

	response, err := req.API.Send(request)
	if err != nil {
		return "", err
	}

	responseInfo, err := req.API.ParseResponse(response)
	if err != nil {
		if _, ok := err.(*traveline.NoTimesFoundError); ok {
			return err.Error(), nil
		}
		return "", err
	}

	// Format the Alexa response
	aimedDepartureTime := responseInfo.AimedDepartureTime.Format(time.Kitchen)
	message := "Your next tram to " + responseInfo.DirectionName + " is due at " + aimedDepartureTime

	expectedDepartureTime := responseInfo.ExpectedDepartureTime.Format(time.Kitchen)
	if aimedDepartureTime != expectedDepartureTime {
		message = message + ", but is expected at " + expectedDepartureTime
	}

	// Convert to a simple readable time
	return message, nil
}

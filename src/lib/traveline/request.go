package traveline

import (
	"time"

	"github.com/google/uuid"
)

// Request is used to make API requests using the Traveline API
type Request struct {
	API API
}

// GetNextTram returns the time of the next tram at the stop that the NaPTAN code represents
func (req *Request) GetNextTram(naptanCode string, when time.Time) (string, error) {
	request, err := req.API.BuildRequest(uuid.New().String(), naptanCode, when)
	if err != nil {
		return "", err
	}

	response, err := req.API.Send(request)
	if err != nil {
		return "", err
	}

	aimedDepartureTime, err := req.API.ParseResponse(response)
	if err != nil {
		return "", err
	}

	// Convert to a simple readable time
	return aimedDepartureTime.Format(time.Kitchen), nil
}

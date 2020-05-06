package transport

import (
	"time"

	"github.com/conradhodge/next-tram/src/lib/traveline"
	"github.com/google/uuid"
)

// Traveline is used to make transport requests using the Traveline API
type Traveline struct {
	API traveline.API
}

// NewTraveline returns the implementation of the transport API for the Traveline API
func NewTraveline(api traveline.API) *Traveline {
	return &Traveline{API: api}
}

// GetNextTramTime returns the time of the next tram at the stop that the NaPTAN code represents
func (c *Traveline) GetNextTramTime(naptanCode string, when time.Time) (*NextTramInfo, error) {
	request, err := c.API.BuildServiceRequest(uuid.New().String(), naptanCode, when)
	if err != nil {
		return nil, err
	}

	response, err := c.API.Send(request)
	if err != nil {
		return nil, err
	}

	monitoredVehicleJourney, err := c.API.ParseServiceDelivery(response)
	if err != nil {
		return nil, err
	}

	nextTramInfo := NextTramInfo{
		LineName:      monitoredVehicleJourney.PublishedLineName,
		VehicleMode:   monitoredVehicleJourney.VehicleMode,
		DirectionName: monitoredVehicleJourney.DirectionName,
	}

	// Convert aimed departure time to time.Time
	aimedDepartureTime, err := convertDepartureTime(monitoredVehicleJourney.MonitoredCall.AimedDepartureTime)
	if err != nil {
		return nil, err
	}
	nextTramInfo.AimedDepartureTime = &aimedDepartureTime

	// Convert expected departure time to time.Time
	if len(monitoredVehicleJourney.MonitoredCall.ExpectedDepartureTime) > 0 {
		expectedDepartureTime, err := convertDepartureTime(monitoredVehicleJourney.MonitoredCall.ExpectedDepartureTime)
		if err != nil {
			return nil, err
		}
		nextTramInfo.ExpectedDepartureTime = &expectedDepartureTime
	}

	return &nextTramInfo, nil
}

func convertDepartureTime(departureTime string) (time.Time, error) {
	convertedDepartureTime, err := time.Parse(time.RFC3339, departureTime)
	if err != nil {
		return time.Time{}, &InvalidTimeFoundError{
			Time:   departureTime,
			Reason: err.Error(),
		}
	}

	return convertedDepartureTime, nil
}

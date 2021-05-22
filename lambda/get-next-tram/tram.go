package main

import (
	"fmt"
	"time"

	"github.com/conradhodge/travel-api-client/transport"
)

func GetNextTram(req transport.API, naptanCode string) (string, error) {
	nextTramInfo, err := req.GetNextTramTime(naptanCode, TimeNow())
	if err != nil {
		return "", err
	}

	// Format the Alexa response
	aimedDepartureTime := nextTramInfo.AimedDepartureTime.Format(time.Kitchen)
	message := fmt.Sprintf("Your next %s %s to %s is due at %s",
		nextTramInfo.LineName,
		nextTramInfo.VehicleMode,
		nextTramInfo.DirectionName,
		aimedDepartureTime,
	)

	if nextTramInfo.ExpectedDepartureTime != nil {
		expectedDepartureTime := nextTramInfo.ExpectedDepartureTime.Format(time.Kitchen)
		if aimedDepartureTime != expectedDepartureTime {
			message = fmt.Sprintf("%s, but is expected at %s", message, expectedDepartureTime)
		}
	}

	return message, nil
}

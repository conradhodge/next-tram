package main

import (
	"fmt"
	"time"

	"github.com/conradhodge/travel-api-client/transport"
)

func GetNextTram(req transport.API, naptanCode string) (string, error) {
	currentTime := TimeNow()

	nextTramInfo, err := req.GetNextDepartureTime(naptanCode, currentTime)
	if err != nil {
		return "", err
	}

	departureTime := getDepartureTime(nextTramInfo)
	diff := departureTime.Sub(currentTime)

	// Format the Alexa response
	message := fmt.Sprintf("Your next %s %s to %s is due in %d minutes at %s",
		nextTramInfo.LineName,
		nextTramInfo.VehicleMode,
		nextTramInfo.DirectionName,
		int64(diff.Minutes()),
		departureTime.Format(time.Kitchen),
	)

	return message, nil
}

func getDepartureTime(nextTramInfo *transport.DepartureInfo) *time.Time {
	if nextTramInfo.ExpectedDepartureTime != nil &&
		nextTramInfo.ExpectedDepartureTime != nextTramInfo.AimedDepartureTime {

		return nextTramInfo.ExpectedDepartureTime
	}

	return nextTramInfo.AimedDepartureTime
}

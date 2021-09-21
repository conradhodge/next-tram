package main

import (
	"fmt"
	"time"

	"github.com/conradhodge/travel-api-client/transport"
)

func GetNextTram(req transport.API, naptanCode string) (string, error) {
	nextTramInfo, err := req.GetNextDepartureTime(naptanCode, TimeNow())
	if err != nil {
		return "", err
	}

	departureTime := getDepartureTime(nextTramInfo)

	// Format the Alexa response
	message := fmt.Sprintf("Your next %s %s to %s is due %s at %s",
		nextTramInfo.LineName,
		nextTramInfo.VehicleMode,
		nextTramInfo.DirectionName,
		getDue(departureTime),
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

func getDue(departureTime *time.Time) string {
	currentTime := TimeNow()
	diff := departureTime.Sub(currentTime)

	due := int64(diff.Minutes())

	if due == 0 {
		return "now"
	} else if due == 1 {
		return "in one minute"
	}

	return fmt.Sprintf("in %d minutes", due)
}

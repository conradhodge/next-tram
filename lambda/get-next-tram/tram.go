package main

import (
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/conradhodge/travel-api-client/transport"
)

var LocalTimezone = "Local"

func GetNextTram(req transport.API, naptanCode string) (string, error) {
	nextTramInfo, err := req.GetNextDepartureTime(naptanCode, TimeNow())
	if err != nil {
		return "", err
	}

	departureTime := getDepartureTime(nextTramInfo)

	// Convert to local time
	location, err := time.LoadLocation(LocalTimezone)
	if err != nil {
		return "", err
	}

	localDepartureTime := departureTime.In(location)

	// Format the Alexa response
	message := fmt.Sprintf("Your next %s %s to %s is due %s at %s",
		nextTramInfo.LineName,
		nextTramInfo.VehicleMode,
		nextTramInfo.DirectionName,
		getDue(&localDepartureTime),
		localDepartureTime.Format(time.Kitchen),
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

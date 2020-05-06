package transport

import (
	"time"
)

// NextTramInfo represents the details for the next tram time
type NextTramInfo struct {
	VehicleMode           string
	LineName              string
	DirectionName         string
	AimedDepartureTime    *time.Time
	ExpectedDepartureTime *time.Time
}

// API represents an API to get travel times for public transport
type API interface {
	GetNextTramTime(naptanCode string, when time.Time) (*NextTramInfo, error)
}

package traveline

import (
	"time"
)

// API represents the interface to the Traveline API
type API interface {
	BuildServiceRequest(requestRef string, naptanCode string, when time.Time) (string, error)
	ParseServiceDelivery(response string) (*MonitoredVehicleJourney, error)
	Send(request string) (string, error)
}

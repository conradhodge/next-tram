package traveline

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// API represents the interface to the Traveline API
type API interface {
	BuildRequest(requestRef string, naptanCode string, when time.Time) (string, error)
	ParseResponse(response string) (time.Time, error)
	Send(request string) (string, error)
}

// API stores the details required to use the Traveline API
type api struct {
	Username string
	Password string
	Client   *http.Client
}

// NewAPI returns a new instance of the Traveline API
func NewAPI(username string, password string, httpClient *http.Client) API {
	return &api{
		Username: username,
		Password: password,
		Client:   httpClient,
	}
}

// Build will return the XML Siri request for the stop that the NaPTAN code represents
func (a *api) BuildRequest(requestRef string, naptanCode string, when time.Time) (string, error) {
	request := &SiriRequest{
		Version:                                siriVersion,
		XMLNS:                                  siriXMLNS,
		ServiceRequestRequestTimestamp:         when.Format(time.RFC3339),
		ServiceRequestRequestorRef:             a.Username,
		StopMonitoringRequestRequestTimestamp:  when.Format(time.RFC3339),
		StopMonitoringRequestMessageIdentifier: requestRef,
		StopMonitoringRequestMonitoringRef:     naptanCode,
	}

	requestBody, err := xml.Marshal(request)
	if err != nil {
		return "", err
	}

	return string(requestBody), nil
}

// Parse the response from the Traveline API and return the time of the next tram
func (a *api) ParseResponse(response string) (time.Time, error) {
	siriResponse := SiriResponse{}
	err := xml.Unmarshal([]byte(response), &siriResponse)
	if err != nil {
		return time.Time{}, err
	}

	log.Printf("RequestMessageRef: %s", siriResponse.ServiceDelivery.StopMonitoringDelivery.RequestMessageRef)

	monitorStopVisits := siriResponse.ServiceDelivery.StopMonitoringDelivery.MonitoredStopVisit
	if len(monitorStopVisits) == 0 {
		return time.Time{}, &NoTimesFoundError{}
	}

	for i, monitorStopVisit := range monitorStopVisits {
		log.Printf(
			"Index: %d, Vehicle: %s, Line: %s, Direction: %s, Time: %s",
			i,
			monitorStopVisit.MonitoredVehicleJourney.VehicleMode,
			monitorStopVisit.MonitoredVehicleJourney.PublishedLineName,
			monitorStopVisit.MonitoredVehicleJourney.DirectionName,
			monitorStopVisit.MonitoredVehicleJourney.MonitoredCall.AimedDepartureTime,
		)
	}

	// Convert next departure time to time.Time
	aimedDepartureTime, err := time.Parse(
		time.RFC3339,
		monitorStopVisits[0].MonitoredVehicleJourney.MonitoredCall.AimedDepartureTime,
	)
	if err != nil {
		return time.Time{}, err
	}

	return aimedDepartureTime, nil
}

// Send will send the request to Traveline API
func (a *api) Send(request string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(request))
	req.Header.Set("Content-type", contentType)
	req.SetBasicAuth(a.Username, a.Password)

	resp, err := a.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response from API: %v", resp)
		return string(body), errors.Errorf("Error status from API: %d", resp.StatusCode)
	}

	return string(body), nil
}

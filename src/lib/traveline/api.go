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
	ParseResponse(response string) (*ResponseInfo, error)
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
	siriRequest := &SiriRequest{
		Version:                                siriVersion,
		XMLNS:                                  siriXMLNS,
		ServiceRequestRequestTimestamp:         when.Format(time.RFC3339),
		ServiceRequestRequestorRef:             a.Username,
		StopMonitoringRequestRequestTimestamp:  when.Format(time.RFC3339),
		StopMonitoringRequestMessageIdentifier: requestRef,
		StopMonitoringRequestMonitoringRef:     naptanCode,
	}

	log.Printf("StopMonitoringRequestRequestTimestamp: %s", siriRequest.StopMonitoringRequestRequestTimestamp)
	log.Printf("StopMonitoringRequestMessageIdentifier: %s", siriRequest.StopMonitoringRequestMessageIdentifier)
	log.Printf("StopMonitoringRequestMonitoringRef: %s", siriRequest.StopMonitoringRequestMonitoringRef)

	requestBody, err := xml.Marshal(siriRequest)
	if err != nil {
		return "", err
	}

	return string(requestBody), nil
}

// ResponseInfo represents the details for the next stop response
type ResponseInfo struct {
	DirectionName         string
	AimedDepartureTime    time.Time
	ExpectedDepartureTime time.Time
}

// Parse the response from the Traveline API and return the time of the next tram
func (a *api) ParseResponse(response string) (*ResponseInfo, error) {
	siriResponse := SiriResponse{}
	err := xml.Unmarshal([]byte(response), &siriResponse)
	if err != nil {
		return nil, err
	}

	log.Printf("RequestMessageRef: %s", siriResponse.ServiceDelivery.StopMonitoringDelivery.RequestMessageRef)

	monitorStopVisits := siriResponse.ServiceDelivery.StopMonitoringDelivery.MonitoredStopVisit
	if len(monitorStopVisits) == 0 {
		return nil, &NoTimesFoundError{}
	}

	for i, monitorStopVisit := range monitorStopVisits {
		log.Printf("MonitoringRef: %s", monitorStopVisit.MonitoringRef)
		log.Printf(
			"Index: %d, Vehicle: %s, Line: %s, Direction: %s, Aimed Departure Time: %s, Expected Departure Time: %s",
			i,
			monitorStopVisit.MonitoredVehicleJourney.VehicleMode,
			monitorStopVisit.MonitoredVehicleJourney.PublishedLineName,
			monitorStopVisit.MonitoredVehicleJourney.DirectionName,
			monitorStopVisit.MonitoredVehicleJourney.MonitoredCall.AimedDepartureTime,
			monitorStopVisit.MonitoredVehicleJourney.MonitoredCall.ExpectedDepartureTime,
		)
	}

	monitoredVehicleJourney := monitorStopVisits[0].MonitoredVehicleJourney

	// Convert aimed departure time to time.Time
	aimedDepartureTime, err := a.convertDepartureTime(monitoredVehicleJourney.MonitoredCall.AimedDepartureTime)
	if err != nil {
		return nil, err
	}

	// Convert expected departure time to time.Time
	expectedDepartureTime, err := a.convertDepartureTime(monitoredVehicleJourney.MonitoredCall.ExpectedDepartureTime)
	if err != nil {
		return nil, err
	}

	return &ResponseInfo{
		DirectionName:         monitoredVehicleJourney.DirectionName,
		AimedDepartureTime:    aimedDepartureTime,
		ExpectedDepartureTime: expectedDepartureTime,
	}, nil
}

func (a *api) convertDepartureTime(departureTime string) (time.Time, error) {
	convertedDepartureTime, err := time.Parse(time.RFC3339, departureTime)
	if err != nil {
		return time.Time{}, &InvalidTimeFoundError{
			Time:   departureTime,
			Reason: err.Error(),
		}
	}

	return convertedDepartureTime, nil
}

// Send will send the request to Traveline API
func (a *api) Send(request string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(request))
	if err != nil {
		return "", err
	}

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

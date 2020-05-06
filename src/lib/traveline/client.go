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

// Client stores the details required to access the Traveline API
type Client struct {
	Username string
	Password string
	Client   *http.Client
}

// NewClient returns the client to access the Traveline API
func NewClient(username string, password string, httpClient *http.Client) API {
	return &Client{
		Username: username,
		Password: password,
		Client:   httpClient,
	}
}

// BuildServiceRequest will return the XML for the request for the stop that the NaPTAN code represents
func (c *Client) BuildServiceRequest(requestRef string, naptanCode string, when time.Time) (string, error) {
	serviceRequest := &ServiceRequest{
		Version:                                siriVersion,
		XMLNS:                                  siriXMLNS,
		ServiceRequestRequestTimestamp:         when.Format(time.RFC3339),
		ServiceRequestRequestorRef:             c.Username,
		StopMonitoringRequestRequestTimestamp:  when.Format(time.RFC3339),
		StopMonitoringRequestMessageIdentifier: requestRef,
		StopMonitoringRequestMonitoringRef:     naptanCode,
	}

	log.Printf("StopMonitoringRequestRequestTimestamp: %s", serviceRequest.StopMonitoringRequestRequestTimestamp)
	log.Printf("StopMonitoringRequestMessageIdentifier: %s", serviceRequest.StopMonitoringRequestMessageIdentifier)
	log.Printf("StopMonitoringRequestMonitoringRef: %s", serviceRequest.StopMonitoringRequestMonitoringRef)

	requestBody, err := xml.Marshal(serviceRequest)
	if err != nil {
		return "", err
	}

	return string(requestBody), nil
}

// ParseServiceDelivery the response from the Traveline API and return the time of the next tram
func (c *Client) ParseServiceDelivery(response string) (*MonitoredVehicleJourney, error) {
	serviceDelivery := ServiceDelivery{}
	err := xml.Unmarshal([]byte(response), &serviceDelivery)
	if err != nil {
		return nil, err
	}

	log.Printf("RequestMessageRef: %s", serviceDelivery.ServiceDelivery.StopMonitoringDelivery.RequestMessageRef)

	monitorStopVisits := serviceDelivery.ServiceDelivery.StopMonitoringDelivery.MonitoredStopVisit
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

	return &monitorStopVisits[0].MonitoredVehicleJourney, nil
}

// Send will send the request to Traveline API
func (c *Client) Send(request string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(request))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-type", contentType)
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if berr := resp.Body.Close(); berr != nil {
			err = berr
		}
	}()

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

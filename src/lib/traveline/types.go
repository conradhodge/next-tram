package traveline

import "encoding/xml"

// ServiceRequest represents the Siri Service Request XML
type ServiceRequest struct {
	XMLName                                xml.Name `xml:"Siri"`
	Version                                string   `xml:"version,attr"`
	XMLNS                                  string   `xml:"xmlns,attr"`
	ServiceRequestRequestTimestamp         string   `xml:"ServiceRequest>RequestTimestamp"`
	ServiceRequestRequestorRef             string   `xml:"ServiceRequest>RequestorRef"`
	StopMonitoringRequestRequestTimestamp  string   `xml:"ServiceRequest>StopMonitoringRequest>RequestTimestamp"`
	StopMonitoringRequestMessageIdentifier string   `xml:"ServiceRequest>StopMonitoringRequest>MessageIdentifier"`
	StopMonitoringRequestMonitoringRef     string   `xml:"ServiceRequest>StopMonitoringRequest>MonitoringRef"`
}

// ServiceDelivery represents the Siri Service Delivery XML response
type ServiceDelivery struct {
	XMLName         xml.Name `xml:"Siri"`
	Version         string   `xml:"version,attr"`
	XMLNS           string   `xml:"xmlns,attr"`
	ServiceDelivery struct {
		ResponseTimestamp      string `xml:"ResponseTimestamp"`
		StopMonitoringDelivery struct {
			ResponseTimestamp  string `xml:"ResponseTimestamp"`
			RequestMessageRef  string `xml:"RequestMessageRef"`
			MonitoredStopVisit []struct {
				RecordedAtTime          string                  `xml:"RecordedAtTime"`
				MonitoringRef           string                  `xml:"MonitoringRef"`
				MonitoredVehicleJourney MonitoredVehicleJourney `xml:"MonitoredVehicleJourney"`
			} `xml:"MonitoredStopVisit"`
		} `xml:"StopMonitoringDelivery"`
	} `xml:"ServiceDelivery"`
}

// MonitoredVehicleJourney represents the Siri Monitored Vehicle Journey XML
type MonitoredVehicleJourney struct {
	FramedVehicleJourneyRef struct {
		DataFrameRef           string `xml:"DataFrameRef"`
		DatedVehicleJourneyRef string `xml:"DatedVehicleJourneyRef"`
	} `xml:"FramedVehicleJourneyRef"`
	VehicleMode       string `xml:"VehicleMode"`
	PublishedLineName string `xml:"PublishedLineName"`
	DirectionName     string `xml:"DirectionName"`
	OperatorRef       string `xml:"OperatorRef"`
	MonitoredCall     struct {
		AimedDepartureTime    string `xml:"AimedDepartureTime"`
		ExpectedDepartureTime string `xml:"ExpectedDepartureTime"`
	} `xml:"MonitoredCall"`
}

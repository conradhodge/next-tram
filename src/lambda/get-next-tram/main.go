package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arienmalec/alexa-go"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/conradhodge/next-tram/src/lib/transport"
	"github.com/conradhodge/next-tram/src/lib/traveline"
)

// Handler is the lambda hander
func Handler() (alexa.Response, error) {
	req := transport.NewTraveline(
		traveline.NewClient(
			os.Getenv("TRAVELINE_API_USERNAME"),
			os.Getenv("TRAVELINE_API_PASSWORD"),
			&http.Client{},
		),
	)

	nextTramInfo, err := req.GetNextTramTime(os.Getenv("NAPTAN_CODE"), time.Now())
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return alexa.NewSimpleResponse("Error", "Something went wrong"), err
	}

	// Format the Alexa response
	aimedDepartureTime := nextTramInfo.AimedDepartureTime.Format(time.Kitchen)
	message := fmt.Sprintf("Your next %s tram to %s is due at %s",
		nextTramInfo.LineName,
		nextTramInfo.DirectionName,
		aimedDepartureTime,
	)

	if nextTramInfo.ExpectedDepartureTime != nil {
		expectedDepartureTime := nextTramInfo.ExpectedDepartureTime.Format(time.Kitchen)
		if aimedDepartureTime != expectedDepartureTime {
			message = fmt.Sprintf("%s, but is expected at %s", message, expectedDepartureTime)
		}
	}

	return alexa.NewSimpleResponse("Time of next tram", message), nil
}

func main() {
	lambda.Start(Handler)
}

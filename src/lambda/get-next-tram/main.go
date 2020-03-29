package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arienmalec/alexa-go"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/conradhodge/next-tram/src/lib/traveline"
)

const naptanMalinBridge = "37090179"

// Handler is the lambda hander
func Handler() (alexa.Response, error) {
	api := traveline.NewAPI(
		os.Getenv("TRAVELINE_API_USERNAME"),
		os.Getenv("TRAVELINE_API_PASSWORD"),
		&http.Client{},
	)
	request := traveline.Request{API: api}

	nextTramTime, err := request.GetNextTram(naptanMalinBridge, time.Now())
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return alexa.NewSimpleResponse("Error", "Something went wrong"), err
	}

	message := "Conrad, you're next tram is at " + nextTramTime
	return alexa.NewSimpleResponse("Time of next tram", message), nil
}

func main() {
	lambda.Start(Handler)
}

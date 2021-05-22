package main

import (
	"net/http"
	"os"

	"github.com/arienmalec/alexa-go"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/conradhodge/travel-api-client/transport"
	"github.com/conradhodge/travel-api-client/traveline"
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

	message, err := GetNextTram(req, os.Getenv("NAPTAN_CODE"))
	if err != nil {
		return alexa.NewSimpleResponse("Error", "Something went wrong"), err
	}

	return alexa.NewSimpleResponse("Time of next tram", message), nil
}

func main() {
	lambda.Start(Handler)
}

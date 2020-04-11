package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arienmalec/alexa-go"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/conradhodge/next-tram/src/lib/request"
	"github.com/conradhodge/next-tram/src/lib/traveline"
)

// Handler is the lambda hander
func Handler() (alexa.Response, error) {
	req := request.NewRequest(
		traveline.NewAPI(
			os.Getenv("TRAVELINE_API_USERNAME"),
			os.Getenv("TRAVELINE_API_PASSWORD"),
			&http.Client{},
		),
	)

	message, err := req.GetNextTramTime(os.Getenv("NAPTAN_CODE"), time.Now())
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return alexa.NewSimpleResponse("Error", "Something went wrong"), err
	}

	return alexa.NewSimpleResponse("Time of next tram", message), nil
}

func main() {
	lambda.Start(Handler)
}

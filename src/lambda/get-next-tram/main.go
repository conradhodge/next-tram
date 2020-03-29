package main

import (
	"fmt"
	"os"

	"github.com/arienmalec/alexa-go"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is the lambda hander
func Handler() (alexa.Response, error) {
	response := "Hello, Conrad"
	fmt.Fprintln(os.Stdout, response)
	return alexa.NewSimpleResponse("Saying Hello", response), nil
}

func main() {
	lambda.Start(Handler)
}

package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle the lambda request
func Handler() (string, error) {
	return fmt.Sprintf("Hello World"), nil
}

func main() {
	lambda.Start(Handler)
}

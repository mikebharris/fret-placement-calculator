package main

import (
	"fret-placement-calculator-lambda/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.Handler{}.HandleRequest)
}

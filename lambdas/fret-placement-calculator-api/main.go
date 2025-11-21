package main

import (
	"main/lambdas/fret-placement-calculator-api/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.Handler{}.HandleRequest)
}

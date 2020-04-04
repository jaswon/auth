package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("method=%s path=%s resource=%s body=%s", ev.HTTPMethod, ev.Path, ev.Resource, ev.Body),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

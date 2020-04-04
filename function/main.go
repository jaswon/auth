package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Method   string
	Resource string
	Path     string
	Body     string
}

func HandleRequest(ev events.APIGatewayProxyRequest) (Response, error) {
	return Response{
		ev.HTTPMethod,
		ev.Resource,
		ev.Path,
		ev.Body,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

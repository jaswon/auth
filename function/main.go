package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
)

var signKey *rsa.PrivateKey

func init() {
	pf, _ := ioutil.ReadFile("sign.key")
	signKey, _ = jwt.ParseRSAPrivateKeyFromPEM(pf)
}

func HandleRequest(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("method=%s path=%s body=%s", ev.HTTPMethod, ev.Path, ev.Body),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

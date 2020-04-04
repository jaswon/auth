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
var hashedSecret []byte

func init() {
	var err error
	hashedSecret, err = ioutil.ReadFile("secret")
	if err != nil || len(hashedSecret) == 0 {
		panic("unable to read secret file")
	}
	pf, err := ioutil.ReadFile("sign.key")
	if err != nil {
		panic("unable to read key file")
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(pf)
	if err != nil {
		panic("unable to parse key file")
	}
}

func HandleRequest(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ev.Body != "" {

	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("method=%s path=%s body=%s", ev.HTTPMethod, ev.Path, ev.Body),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

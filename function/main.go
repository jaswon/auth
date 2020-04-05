package main

import (
	"crypto/rsa"
	"io/ioutil"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var signKey *rsa.PrivateKey
var hashedSecret []byte

func init() {
	var err error
	hashedSecret, err = ioutil.ReadFile("./secret")
	if err != nil {
		log.Fatal("unable to read secret file", err)
	}
	pf, err := ioutil.ReadFile("./signkey")
	if err != nil {
		log.Fatal("unable to read key file", err)
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(pf)
	if err != nil {
		log.Fatal("unable to parse key file", err)
	}
}

func HandleRequest(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if bcrypt.CompareHashAndPassword(hashedSecret, []byte(ev.Body)) == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "auth success",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 403,
		Body:       "auth fail",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

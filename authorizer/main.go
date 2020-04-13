package main

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
)

var verifyKey *rsa.PublicKey

func init() {
	vf, err := ioutil.ReadFile("verifykey")
	if err != nil {
		log.Fatal("unable to read key file", err)
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(vf)
	if err != nil {
		log.Fatal("unable to parse key file", err)
	}
}

func HandleRequest(ev events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	policy := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				events.IAMPolicyStatement{
					Action:   []string{"execute-api:Invoke"},
					Resource: []string{strings.Split(ev.MethodArn, "/")[0] + "/*/*/*"},
				},
			},
		},
		Context: map[string]interface{}{},
	}

	token, err := jwt.Parse(ev.AuthorizationToken, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if token.Valid {
		policy.PolicyDocument.Statement[0].Effect = "ALLOW"
		return policy, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		policy.PolicyDocument.Statement[0].Effect = "DENY"
		policy.Context["error"] = ve.Error()
		return policy, nil
	} else {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}
}

func main() {
	lambda.Start(HandleRequest)
}

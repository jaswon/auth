package main

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var signKey *rsa.PrivateKey
var hashedSecret []byte

// refresh tokens last for one week
var refresh_ttl time.Duration = time.Hour * 24 * 7

// access tokens last for 5 minutes
var access_ttl time.Duration = time.Minute * 5

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
	if bcrypt.CompareHashAndPassword(hashedSecret, []byte(ev.Body)) != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       "Unauthorized",
		}, nil
	}

	now := time.Now()
	refresh_expire := now.Add(refresh_ttl)
	access_expire := now.Add(access_ttl)

	refresh_token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		ExpiresAt: refresh_expire.Unix(),
	})

	signed_refresh, err := refresh_token.SignedString(signKey)
	if err != nil {
		log.Println("unable to sign token", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "unable to sign token",
		}, nil
	}

	refresh_cookie := http.Cookie{
		Name:     "jwon_refresh",
		Value:    signed_refresh,
		Secure:   true,
		HttpOnly: true,
		Expires:  refresh_expire,
		Path:     "/",
	}

	access_token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		ExpiresAt: access_expire.Unix(),
	})

	signed_access, err := access_token.SignedString(signKey)
	if err != nil {
		log.Println("unable to sign token", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "unable to sign token",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       signed_access,
		Headers: map[string]string{
			"Set-Cookie": refresh_cookie.String(),
		},
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}

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

const cookie_name = "jwon_refresh"

var signKey *rsa.PrivateKey
var verifyKey *rsa.PublicKey
var hashedSecret []byte

// refresh tokens last for one week
var refresh_ttl time.Duration = time.Hour * 24 * 7

// access tokens last for 5 minutes
var access_ttl time.Duration = time.Minute * 5

func init() {
	var err error
	hashedSecret, err = ioutil.ReadFile("secret")
	if err != nil {
		log.Fatal("unable to read secret file", err)
	}

	sf, err := ioutil.ReadFile("signkey")
	if err != nil {
		log.Fatal("unable to read key file", err)
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(sf)
	if err != nil {
		log.Fatal("unable to parse key file", err)
	}

	vf, err := ioutil.ReadFile("verifykey")
	if err != nil {
		log.Fatal("unable to read key file", err)
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(vf)
	if err != nil {
		log.Fatal("unable to parse key file", err)
	}
}

func authorize(ev events.APIGatewayProxyRequest) error {
	if ev.Body != "" {
		err := bcrypt.CompareHashAndPassword(hashedSecret, []byte(ev.Body))
		if err == nil {
			return nil
		} else {
			log.Println(err)
		}
	}

	refresh_token, err := (&http.Request{Header: http.Header(ev.MultiValueHeaders)}).Cookie(cookie_name)
	if err != nil {
		return err
	}

	token, err := jwt.Parse(refresh_token.String(), func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if token.Valid {
		return nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		return ve
	} else {
		return err
	}
}

func HandleRequest(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if err := authorize(ev); err != nil {
		log.Println(err)
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
		Name:     cookie_name,
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

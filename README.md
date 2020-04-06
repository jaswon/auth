# personal auth service

Deploys an AWS stack which provides JWT access + refresh tokens authenticated by a single master password

The deployed service has two endpoints:
- `GET /pubkey` : serves the public key part used to verify signed tokens
- `POST /token` : upon successful authentication, sets a new refresh token cookie, and returns a new access token in the response body
    - authentication succeeds if the master password is supplied in the request body or a valid refresh token is present

## setup

### prerequisites
- have a domain
- create a certificate for this domain using aws certificate-manager

### steps

1. clone this repo
2. `cd auth`
3. `npm install`
4. modify `config.mk`
5. `make deploy`

### regenerate master password
```
make secret
```

### regenerate jwt keys
```
make key
```

## files

- `bin/auth.ts` - AWS Cloudformation stack defined with AWS CDK
- `function/main.go` - AWS Lambda function definition (deployed to `/token` endpoint)
- `function/gensecret/main.go` - utility for hashing the master password

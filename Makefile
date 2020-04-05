include config.mk

ASSETS = $(addprefix function/bin/,main secret signkey verifykey)
CDK_ENV = CERT_ARN=$(CERT_ARN) DOMAIN=$(DOMAIN)

synth: bin/auth.js $(ASSETS)
	$(CDK_ENV) cdk synth

deploy: bin/auth.js $(ASSETS)
	$(CDK_ENV) cdk deploy

clean:
	rm -r function/bin

bin/auth.js: bin/auth.ts
	npm run build

function/bin:
	mkdir function/bin

function/bin/main: function/main.go | function/bin
	cd function && GOOS=linux GOARCH=amd64 go build -o bin/main -ldflags="-X 'main.CookieName=$(COOKIE_NAME)'" main.go

secret function/bin/secret: | function/bin
	cd function && go run gensecret/main.go

key function/bin/signkey: | function/bin
	openssl genrsa -out function/bin/signkey 4096
	chmod 644 function/bin/signkey

function/bin/verifykey: function/bin/signkey
	openssl rsa -in function/bin/signkey -pubout -out function/bin/verifykey
	chmod 644 function/bin/verifykey

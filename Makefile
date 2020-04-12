include config.mk

HANDLER_ASSETS = main secret verifykey signkey
CDK_ENV = CERT_ARN=$(CERT_ARN) DOMAIN=$(DOMAIN)

synth: auth.js bin/handler.zip
	$(CDK_ENV) cdk synth

deploy: auth.js bin/handler
	$(CDK_ENV) cdk deploy

clean:
	rm -r bin/

auth.js: auth.ts
	npm run build

bin/handler.zip: $(addprefix bin/handler/,$(HANDLER_ASSETS))
	zip -j bin/handler bin/handler/*

bin/handler:
	mkdir -p bin/handler

bin/handler/main: handler/main.go | bin/handler
	GOOS=linux GOARCH=amd64 go build -o bin/handler/main -ldflags="-X 'main.CookieName=$(COOKIE_NAME)'" handler/main.go

secret bin/handler/secret: | bin/handler
	go run gensecret/main.go

key bin/handler/signkey: | bin/handler
	openssl genrsa -out bin/handler/signkey 4096
	chmod 644 bin/handler/signkey

bin/handler/verifykey: bin/handler/signkey
	openssl rsa -in bin/handler/signkey -pubout -out bin/handler/verifykey
	chmod 644 bin/handler/verifykey

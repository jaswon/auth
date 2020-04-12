include config.mk

HANDLER_ASSETS = main secret verifykey signkey
AUTHORIZER_ASSETS = main verifykey
CDK_ENV = CERT_ARN=$(CERT_ARN) DOMAIN=$(DOMAIN)

synth: auth.js bin/handler.zip bin/authorizer.zip
	$(CDK_ENV) cdk synth

deploy: auth.js bin/handler.zip bin/authorizer.zip
	$(CDK_ENV) cdk deploy

clean:
	rm -r bin/

auth.js: auth.ts
	npm run build

bin/authorizer.zip: $(addprefix bin/authorizer/,$(AUTHORIZER_ASSETS))
	zip -j bin/authorizer bin/authorizer/*

bin/authorizer:
	mkdir -p bin/authorizer

bin/authorizer/main: authorizer/main.go | bin/authorizer
	GOOS=linux GOARCH=amd64 go build -o bin/authorizer/main authorizer/main.go

bin/authorizer/verifykey: bin/handler/verifykey
	cp bin/handler/verifykey bin/authorizer/verifykey

bin/handler.zip: $(addprefix bin/handler/,$(HANDLER_ASSETS))
	zip -j bin/handler bin/handler/*

bin/handler:
	mkdir -p bin/handler

bin/handler/main: handler/main.go | bin/handler
	GOOS=linux GOARCH=amd64 go build -o bin/handler/main -ldflags="-X 'main.CookieName=$(COOKIE_NAME)'" handler/main.go

secret bin/handler/secret: | bin/handler
	read -sp 'enter new secret: ' && echo $$REPLY | go run gensecret/main.go > bin/handler/secret

key bin/handler/signkey: | bin/handler
	openssl genrsa -out bin/handler/signkey 4096
	chmod 644 bin/handler/signkey

bin/handler/verifykey: bin/handler/signkey
	openssl rsa -in bin/handler/signkey -pubout -out bin/handler/verifykey
	chmod 644 bin/handler/verifykey

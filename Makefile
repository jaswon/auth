ASSETS=$(addprefix assets/,signkey.pub function.zip)

synth: bin/auth.js $(ASSETS)
	cdk synth

deploy: bin/auth.js $(ASSETS)
	cdk deploy

bin/auth.js: bin/auth.ts
	npm run build

function/bin:
	mkdir function/bin

function/bin/main: function/main.go | function/bin
	cd function && GOOS=linux GOARCH=amd64 go build -o bin/main main.go

secret function/bin/secret: | function/bin
	cd function && go run gensecret/main.go

key function/bin/signkey: | function/bin
	openssl genrsa -out function/bin/signkey 4096

assets:
	mkdir assets

assets/signkey.pub: function/bin/signkey | assets
	openssl rsa -in function/bin/signkey -pubout -out assets/signkey.pub

assets/function.zip: $(addprefix function/bin/,main signkey secret) | assets
	zip -j assets/function function/bin/*

synth: lib/auth-stack.js function/bin/main
	cdk synth

deploy: lib/auth-stack.js function/bin/main
	cdk deploy

lib/auth-stack.js: lib/auth-stack.ts
	npm run build

function/bin/main: function/main.go
	cd function && GOOS=linux GOARCH=amd64 go build -o bin/main main.go

key:
	ssh-keygen -t rsa -b 2048 -m PEM -N "" -f function/bin/sign.key
	openssl rsa -in function/bin/sign.key -pubout -outform PEM -out function/sign.key.pub

secret:
	cd function && go run gensecret/main.go

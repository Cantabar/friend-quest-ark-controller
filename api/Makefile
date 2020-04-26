.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/list list/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/start start/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/stop stop/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

deploy-client:
	sls client deploy

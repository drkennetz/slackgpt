test:
	go test -v -cover -coverprofile=coverage.out --json ./...

coverage:
	go tool cover -func=coverage.out

build:
	go build -o ./bin/slackgpt

.PHONY: test coverage build

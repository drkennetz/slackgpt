test:
	go test -v -cover -coverprofile=coverage.out --json ./...

coverage:
	go tool cover -func=coverage.out

.PHONY: test coverage

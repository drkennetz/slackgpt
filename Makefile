test:
	go test -v -cover ./...

coverage:
	go tool cover -func=coverage.out

.PHONY: test coverage

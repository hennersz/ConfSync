.PHONY: build
build: 
	go build -o bin/conf-sync cmd/confSync/main.go

.PHONY: test
test:
	go test -v -race ./internal/... -coverprofile=coverage.out

.PHONY: coverage
coverage: test
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

.PHONY: install 
install:
	go get -v -t -d ./...

.PHONY: build
build: 
	go build -o bin/conf-sync cmd/confSync/main.go

.PHONY: test
test:
	go test -v ./internal/... -coverprofile=coverage.out

.PHONY: coverage
coverage: test
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: fix
fix:
	golangci-lint run --fix

.PHONY: release
release:
	npm install @codedependant/semantic-release-docker 
	npx semantic-release
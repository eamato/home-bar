.PHONY: install
install:
	@echo "installing dependencies"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3
	@go mod tidy

.PHONY: lint
lint:
	@echo "starting lint"
	@golangci-lint run

.PHONY: update
update:
	@echo "updating dependencies"
	@go get -u -all
	@go mod tidy
	@go mod vendor

.PHONY: run
run:
	@echo "running app"
	@go run ./cmd/main.go
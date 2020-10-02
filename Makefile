GOPATH:=$(shell go env GOPATH)
APP?=organization

.PHONY: format
## format: format files
format:
	@go get golang.org/x/tools/cmd/goimports
	goimports -local github.com/entrydsm -w .
	gofmt -s -w .
	go mod tidy

.PHONY: test
## test: run tests
test:
	@go get github.com/rakyll/gotest
	gotest -p 1 -race -cover -v ./...
	go mod tidy

.PHONY: lint
## lint: check everything's okay
lint:
	@go get github.com/kyoh86/scopelint
	golangci-lint run ./...
	scopelint --set-exit-status ./...
	go mod verify
	go mod tidy

.PHONY: generate
## generate: generate source code for mocking
generate:
	@go get golang.org/x/tools/cmd/stringer
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	go generate ./...
	go mod tidy

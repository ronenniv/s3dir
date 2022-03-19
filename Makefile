.PHONY: install test clean linter goimports govet golint docker-build docker-push cover

install:
	go build -v ./...

test:
	go test ./...

clean:
	go clean

linter:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.42.1 golangci-lint run -v

goimports:
	go get golang.org/x/tools/cmd/goimports
	goimports -w .

govet: goimports
	go vet ./...

golint: govet
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39.0
	golangci-lint run
	go mod tidy

# API_DEFINITIONS_SHA=$(shell git log --oneline | grep Regenerated | head -n1 | cut -d ' ' -f 5)
# docker-build:
# 	docker build -t ronenniv/s3dir .
# 	docker tag ronenniv/s3dir ronenniv/s3dir:${GITHUB_TAG}
# 	docker tag ronenniv/s3dir ronenniv/s3dir:latest

# docker-push:
# 	docker push ronenniv/s3dir:${GITHUB_TAG}
# 	docker push ronenniv/s3dir:latest

GO_DIRS = $(shell go list ./... | grep -v /rest/ | grep -v /form )
cover:
	go test ${GO_DIRS} -coverprofile coverage.out
	go test ${GO_DIRS} -json > test-report.out
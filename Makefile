# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
SERVICE_NAME = cla-api
BUILD_TIME=`date +%FT%T%z`
VERSION := $(shell sh -c 'git describe --always --tags')
BRANCH := $(shell sh -c 'git rev-parse --abbrev-ref HEAD')
COMMIT := $(shell sh -c 'git rev-parse --short HEAD')
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.branch=$(BRANCH) -X main.buildDate=$(BUILD_TIME)"
BUILD_TAGS=-tags aws_lambda

LINT_TOOL=$(shell go env GOPATH)/bin/golangci-lint
GO_PKGS=$(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/)
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')
TEST_ENV=AWS_REGION=us-east-1 CLA_SERVICE_AWS_ACCESS_KEY_ID=test-env-aws-access-key-id CLA_SERVICE_AWS_SECRET_ACCESS_KEY=test-env-aws-secret-access-key

.PHONY: generate setup setup_dev setup_deploy clean swagger up fmt test run deps build build-mac build_aws_lambda qc lint

generate: swagger

setup: $(LINT_TOOL) setup_dev
	mkdir -p bin

setup_dev:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/dep/cmd/dep	
	go get -u github.com/stripe/safesql
	sudo curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/download/v1.7.0/dbmate-linux-amd64
	sudo chmod +x /usr/local/bin/dbmate

clean:
	rm -rf ./gen ./bin

swagger: clean
	mkdir gen
	swagger -q generate server -t gen -f swagger/cla.yaml --exclude-main -A cla -P user.CLAUser

swagger-validate:
	swagger validate swagger/cla.yaml

up:
	dbmate -d ".build/db/migrations" -s ".build/db/schema.sql" up

fmt:
	@gofmt -w -l -s $(GO_FILES)
	@goimports -w -l $(GO_FILES)

test:
	@ $(TEST_ENV) go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /node_modules/) -coverprofile=cover.out

run:
	go run main.go

deps:
	dep ensure -v

build: deps
	mkdir -p bin
	env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(SERVICE_NAME) main.go
	chmod +x bin/$(SERVICE_NAME)

build-mac: deps
	mkdir -p bin
	env GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(SERVICE_NAME) main.go
	chmod +x bin/$(SERVICE_NAME)

build-aws-lambda: deps
	mkdir -p bin
	env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) $(BUILD_TAGS) -o bin/$(SERVICE_NAME) main.go
	chmod +x bin/$(SERVICE_NAME)

$(LINT_TOOL):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.18.0

lint: $(LINT_TOOL)
	$(LINT_TOOL) run --config=.golangci.yaml ./...

safesql:
	safesql -v $(GO_FILES)
	#$(GOPATH)/src/github.com/communitybridge/safesql/bin/safesql -v $(GO_PKGS)

all: clean swagger swagger-validate deps fmt build-mac test lint

POSTGRES_DATA_DIR=postgresdata
POSTGRES_DOCKER_IMAGE=circleci/postgres:9.6-alpine
POSTGRES_PORT=5432
POSTGRES_DB_NAME=civil_crawler
POSTGRES_USER=docker
POSTGRES_PSWD=docker

PUBSUB_SIM_DOCKER_IMAGE=kinok/google-pubsub-emulator:latest

GOVERSION=go1.12.7

GOCMD=go
GOGEN=$(GOCMD) generate
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCOVER=$(GOCMD) tool cover

GO:=$(shell command -v go 2> /dev/null)
DOCKER:=$(shell command -v docker 2> /dev/null)
APT:=$(shell command -v apt-get 2> /dev/null)
GOVERCURRENT=$(shell go version |awk {'print $$3'})

# curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin vX.Y.Z
GOLANGCILINT_URL=https://install.goreleaser.com/github.com/golangci/golangci-lint.sh
GOLANGCILINT_VERSION_TAG=v1.16.0

## Reliant on go and $GOPATH being set
.PHONY: check-go-env
check-go-env:
ifndef GO
	$(error go command is not installed or in PATH)
endif
ifndef GOPATH
	$(error GOPATH is not set)
endif
ifneq ($(GOVERCURRENT), $(GOVERSION))
	$(error Incorrect go version, needs $(GOVERSION))
endif

## NOTE: If installing on a Mac, use Docker for Mac, not Docker toolkit
## https://www.docker.com/docker-mac
.PHONY: check-docker-env
check-docker-env:
ifndef DOCKER
	$(error docker command is not installed or in PATH)
endif

.PHONY: install-gobin
install-gobin: check-go-env ## Installs gobin tool
	@GO111MODULE=off go get -u github.com/myitcv/gobin

# Use commit until they cut a release with our fix
.PHONY: install-conform
install-conform: install-gobin ## Installs conform
	@gobin github.com/autonomy/conform@7bed9129bc73b5a6fd2b6d8c12f7c024ea4b7107

.PHONY: install-linter
install-linter: check-go-env ## Installs linter
	@curl -sfL $(GOLANGCILINT_URL) | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCILINT_VERSION_TAG)
ifdef APT
	@sudo apt-get install golang-race-detector-runtime || true
endif

.PHONY: install-cover
install-cover: check-go-env ## Installs code coverage tool
	@gobin -u golang.org/x/tools/cmd/cover

.PHONY: install-gorunpkg
install-gorunpkg: ## Installs the gorunpkg command
	@gobin -u github.com/vektah/gorunpkg

.PHONY: install-gorm
install-gorm: ## Installs gorm package
	@$(GOGET) -u github.com/jinzhu/gorm

.PHONY: install-gormmigrate
install-gormmigrate:
	@$(GOGET) -u gopkg.in/gormigrate.v1

.PHONY: install-pq
install-pq: ## Installs gorm package
	@$(GOGET) -u github.com/lib/pq

.PHONY: setup-githooks
setup-githooks: ## Setups any git hooks in githooks
	@curl -sfL https://raw.githubusercontent.com/joincivil/go-common/master/githooks/commit-msg -o .git/hooks/commit-msg
	@chmod 755 .git/hooks/commit-msg

.PHONY: setup
setup: check-go-env install-conform install-linter install-cover install-gorunpkg install-pq install-gorm install-gormmigrate setup-githooks ## Sets up the tooling.

.PHONY: postgres-setup-launch
postgres-setup-launch:
ifeq ("$(wildcard $(POSTGRES_DATA_DIR))", "")
	mkdir -p $(POSTGRES_DATA_DIR)
	docker run \
		-v $$PWD/$(POSTGRES_DATA_DIR):/tmp/$(POSTGRES_DATA_DIR) -i -t $(POSTGRES_DOCKER_IMAGE) \
		/bin/bash -c "cp -rp /var/lib/postgresql /tmp/$(POSTGRES_DATA_DIR)"
endif
	docker run -e "POSTGRES_USER="$(POSTGRES_USER) -e "POSTGRES_PASSWORD"=$(POSTGRES_PSWD) -e "POSTGRES_DB"=$(POSTGRES_DB_NAME) \
	    -v $$PWD/$(POSTGRES_DATA_DIR)/postgresql:/var/lib/postgresql -d -p $(POSTGRES_PORT):$(POSTGRES_PORT) \
		$(POSTGRES_DOCKER_IMAGE);

.PHONY: postgres-check-available
postgres-check-available:
	@for i in `seq 1 10`; \
	do \
		nc -z localhost 5432 2> /dev/null && exit 0; \
		sleep 3; \
	done; \
	exit 1;

.PHONY: postgres-start
postgres-start: check-docker-env postgres-setup-launch postgres-check-available ## Starts up a development PostgreSQL server
	@echo "Postgresql launched and available"

.PHONY: postgres-stop
postgres-stop: check-docker-env ## Stops the development PostgreSQL server
	@docker stop `docker ps -q`
	@echo 'Postgres stopped'

## golangci-lint config in .golangci.yml
.PHONY: lint
lint: check-go-env  ## Runs linting.
	@golangci-lint run ./...

.PHONY: conform
conform: check-go-env ## Runs conform (commit message linting)
	@conform enforce

.PHONY: build
build: check-go-env ## Builds the repo, mainly to ensure all the files will build properly
	$(GOBUILD) ./...

.PHONY: test
test: check-go-env ## Runs unit tests and tests code coverage
	@echo 'mode: atomic' > coverage.txt && $(GOTEST) -covermode=atomic -coverprofile=coverage.txt -v -race -timeout=60s ./...

.PHONY: test-integration
test-integration: check-go-env ## Runs tagged integration tests
	@echo 'mode: atomic' > coverage.txt && PUBSUB_EMULATOR_HOST=localhost:8042 $(GOTEST) -covermode=atomic -coverprofile=coverage.txt -v -race -timeout=60s -tags=integration ./...

.PHONY: cover
cover: test ## Runs unit tests, code coverage, and runs HTML coverage tool.
	@$(GOCOVER) -html=coverage.txt

.PHONY: clean
clean: ## go clean and clean up of artifacts
	@$(GOCLEAN) ./... || true
	@rm coverage.txt || true
	@rm build || true

## Some magic from http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


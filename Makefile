go_source_files := $(wildcard **/*.go)
ts_source_files := $(wildcard ui/src/*/*.ts)
tsx_source_files := $(wildcard ui/src/*/*.tsx)
GOCMD?=go

.PHONY: all

all: build-backend build-ui

start-local: all
	LOCAL_START=1 bin/kpack-ui

build-ui: $(ts_source_files) $(tsx_source_files) ui/package.json ui-deps
	cd ui && npm run build

build-backend: $(go_source_files)
	mkdir -p bin
	go build -o bin/kpack-ui ./cmd/kpack-ui/main.go

docker:
	docker build -t kpack-ui .

ui-deps:
	cd ui && npm install


unit-test: unit-test-go unit-test-ui

unit-test-go:
	go test ./...

unit-test-ui: ui-deps
	cd ui && npm test

install-goimports:
	@echo "> Installing goimports..."
	cd tools; $(GOCMD) install golang.org/x/tools/cmd/goimports

format: install-goimports
	@echo "> Formating code..."
	@goimports -l -w -local ${PACKAGE_BASE} ${SRC}

install-golangci-lint:
	@echo "> Installing golangci-lint..."
	cd tools; $(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-golangci-lint
	@echo "> Linting code..."
	@golangci-lint run -c golangci.yaml

lint: lint-go lint-ui

lint-ui: ui-deps
	cd ui && npm run lint

lint-go: install-golangci-lint
	@echo "> Linting code..."
	@golangci-lint run -c golangci.yaml

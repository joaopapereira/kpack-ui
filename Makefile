go_source_files := $(wildcard **/*.go)
ts_source_files := $(wildcard web/src/*/*.ts)
tsx_source_files := $(wildcard web/src/*/*.tsx)
GOCMD?=go

.PHONY: all

all: build-backend build_web

start-local: all
	LOCAL_START=1 bin/kpack-ui

build_web: $(ts_source_files) $(tsx_source_files) web/package.json web_deps
	cd web && npm run build

build-backend: $(go_source_files)
	mkdir -p bin
	go build -o bin/kpack-ui ./cmd/kpack-ui/main.go

docker:
	docker build -t kpack-ui .

web_deps:
	cd web && npm install


unit-test: unit-test-go unit_test_web

unit-test-go:
	go test ./...

unit_test_web: web_deps
	cd web && npm test

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

lint: lint_go lint_web

lint_web: web_deps
	cd web && npm run lint

lint_go: install-golangci-lint
	@echo "> Linting code..."
	@golangci-lint run -c golangci.yaml

generate_icons: fyne_app_installed
	@echo "> Regenerating icon resources..."
	@fyne bundle -package static -prefix icon static/icons/ > static/bundled-icons.go

fyne_app_installed:
	@echo "> Installing fyne app..."
	@go get fyne.io/fyne/cmd/fyne
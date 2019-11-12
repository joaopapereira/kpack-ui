go_source_files := $(wildcard **/*.go)
ts_source_files := $(wildcard ui/src/*/*.ts)
tsx_source_files := $(wildcard ui/src/*/*.tsx)

.PHONY: all

all: build-backend build-ui

start-local: all
	LOCAL_START=1 bin/kpack-ui

build-ui: $(ts_source_files) $(tsx_source_files) ui/package.json
	cd ui && npm install
	cd ui && npm run build

build-backend: $(go_source_files)
	mkdir -p bin
	go build -o bin/kpack-ui ./cmd/kpack-ui/main.go

docker:
	docker build -t kpack-ui .
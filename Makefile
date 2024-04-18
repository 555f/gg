GO_BIN = $(shell go env GOPATH)/bin
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

build-server-linux:
	GOOS=linux GOARCH=amd64 go build ./cmd/selfupdate-server

sync:
	rsync -a --delete selfupdate-server ./public vitaly@51.250.88.10:~/

build:
	@VERSION=$(shell gitver show) ./build.sh

install: build
	@cp ./build/${GOOS}-${GOARCH} ${GO_BIN}/gg

chglog:
	@git-chglog -o CHANGELOG.md

check:
	@go vet ./...
	@go test -v ./...

css:
	@npx tailwindcss -i ./internal/plugin/http/files/tailwind.css -o ./internal/plugin/http/files/style.css  -c ./internal/plugin/http/files/tailwind.config.js --watch

gen-examples: gen-examples-rest-service-echo gen-examples-rest-service-chi

gen-examples-rest-service-echo:
	 go run cmd/gg/main.go run --config examples/rest-service-echo/gg.yaml

gen-examples-rest-service-chi:
	 go run cmd/gg/main.go run --config examples/rest-service-chi/gg.yaml

gen-examples-grpc:
	 go run cmd/gg/main.go run --config examples/grpc-service/gg.yaml

gen-examples-pwa:
	 go run cmd/gg/main.go run --config examples/pwa/gg.yaml

gen-examples-pwa-build:
	GOARCH=wasm GOOS=js go build -o ./examples/pwa/web/app.wasm ./examples/pwa/cmd/pwa/main.go
	go build -o ./examples/pwa ./examples/pwa/cmd/pwa

gen-examples-pwa-dev: gen-examples-pwa-build
	$(shell cd ./examples/pwa && ./pwa)

.PHONY: build

default: build
GO_BIN = $(shell go env GOPATH)/bin
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)

build-server-linux:
	GOOS=linux GOARCH=amd64 go build  -o ./selfupdate-server ./cmd/selfupdate

sync:
	rsync -a selfupdate-server ./public vitaly@51.250.88.10:~/

build:
	@./build.sh

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
	 go run cmd/gg/main.go run  --config examples/rest-service-echo/gg.yaml

gen-examples-rest-service-chi:
	 go run cmd/gg/main.go run  --config examples/rest-service-chi/gg.yaml

gen-examples-grpc:
	 go run cmd/gg/main.go run  --config examples/grpc-service/gg.yaml

.PHONY: build

default: build
//go:generate gg -p ./examples/metrics-middleware/internal/client/... --plugins middleware-output=./examples/metrics-middleware/internal/client/middleware.go,metrics-output=./examples/metrics-middleware/internal/client/metrics.go run

//go:generate gg -w ./examples/rest-service-chi -p ./... --plugins http-error-wrapper=./pkg/errors/ErrorWrapper,http-error-default=./pkg/errors/DefaultError,http-openapi-tpl=./docs,http-server-output=./internal/server/server.go,http-client-output=./pkg/client/client.go,http-openapi-output=./docs,http-apidoc-output=./docs/api.html,config-doc-output=./docs/CONFIG.md,middleware-output=./internal/middleware/middleware.go,klog-output=./internal/logging/logging.go,config-output=./internal/config/config_loader.go run
package gg

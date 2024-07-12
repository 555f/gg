//go:generate gg -p ./examples/metrics-middleware/internal/client/... --plugins middleware-output=./examples/metrics-middleware/internal/middleware/middleware.go,metrics-output=./examples/metrics-middleware/internal/metrics/metrics.go run

package gg

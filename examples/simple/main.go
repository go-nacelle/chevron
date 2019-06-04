package main

import (
	"compress/gzip"
	"context"
	"net/http"

	"github.com/efritz/response"
	"github.com/go-nacelle/chevron/middleware"
	"github.com/go-nacelle/nacelle"

	"github.com/go-nacelle/chevron"
)

type TestResource struct {
	*chevron.EmptySpec

	Logger nacelle.Logger `service:"logger"`
}

func (tr *TestResource) Get(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return response.Respond([]byte("abcdefghijklmnopqrstuvwxyz\n\n"))
}

func setupRoutes(config nacelle.Config, router chevron.Router) error {
	router.AddMiddleware(middleware.NewLogging())
	router.AddMiddleware(middleware.NewRequestID())
	router.AddMiddleware(middleware.NewGzip(middleware.WithGzipLevel(gzip.BestCompression)))

	router.MustRegister("/", &TestResource{})

	return nil
}

func main() {
	chevron.BootAndExit("app", chevron.RouteInitializerFunc(setupRoutes))
}

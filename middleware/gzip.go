package middleware

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/efritz/nacelle"
	"github.com/efritz/response"

	"github.com/efritz/chevron"
)

type GzipMiddleware struct {
	level int
}

func NewGzip(configs ...GzipMiddlewareConfigFunc) chevron.Middleware {
	m := &GzipMiddleware{
		level: gzip.DefaultCompression,
	}

	for _, f := range configs {
		f(m)
	}

	return m
}

func (m *GzipMiddleware) Convert(f chevron.Handler) (chevron.Handler, error) {
	if m.level < gzip.HuffmanOnly || m.level > gzip.BestCompression {
		return nil, fmt.Errorf("gzip: invalid compression level: %d", m.level)
	}

	handler := func(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
		resp := f(ctx, req, logger)

		if !shouldGzip(req, resp) {
			return resp
		}

		if resp.StatusCode() == http.StatusNoContent {
			resp.SetHeader("Content-Encoding", "")
		} else {
			resp.SetHeader("Content-Encoding", "gzip")
		}

		resp.SetHeader("Content-Length", "")
		resp.SetHeader("Vary", "Accept-Encoding")
		resp.SetHeader("Content-Type", "application/octet-stream")

		return resp.DecorateWriter(func(w io.Writer) io.Writer {
			// TODO - need to close the inner one too?
			gzipWriter, _ := gzip.NewWriterLevel(w, m.level)
			return gzipWriter
		})
	}

	return handler, nil
}

func shouldGzip(req *http.Request, resp response.Response) bool {
	// Skip if the request doesn't accept gzipped responses
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}

	// Skip if handler already set an encoding
	return resp.Header("Content-Encoding") == ""
}

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDGetRequestID(t *testing.T) {
	var (
		val = "x"
		ctx = context.WithValue(context.Background(), TokenRequestID, val)
	)

	// Present
	assert.Equal(t, val, GetRequestID(ctx))

	// Missing
	assert.Empty(t, GetRequestID(context.Background()))
}

func TestRequestIDNewRequestIDGenerated(t *testing.T) {
	var ctxVal string
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		ctxVal = GetRequestID(ctx)
		return response.Empty(http.StatusNoContent)
	}

	wrapped, err := NewRequestID().Convert(bare)
	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())

	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
	assert.Equal(t, ctxVal, resp.Header("X-Request-ID"))
	assert.Len(t, resp.Header("X-Request-ID"), 36)
}

func TestRequestIDNewRequestIDSuppliedByClient(t *testing.T) {
	var ctxVal string
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		ctxVal = GetRequestID(ctx)
		return response.Empty(http.StatusNoContent)
	}

	wrapped, err := NewRequestID().Convert(bare)
	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("X-Request-ID", "1234")
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())

	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
	assert.Equal(t, ctxVal, resp.Header("X-Request-ID"))
	assert.Equal(t, "1234", resp.Header("X-Request-ID"))
}

func TestRequestIDRequestIDFailure(t *testing.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.Empty(http.StatusNoContent)
	}

	expectedResp := response.JSON(map[string]string{
		"message": "entropy whoopsie",
	})

	generator := func() (string, error) {
		return "", fmt.Errorf("utoh")
	}

	errorFactory := func(err error) response.Response {
		return expectedResp
	}

	wrapped, err := NewRequestID(
		WithRequestIDGenerator(generator),
		WithRequestIDErrorFactory(errorFactory),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.Equal(t, expectedResp, resp)
}

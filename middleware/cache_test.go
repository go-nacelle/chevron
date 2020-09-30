package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"testing"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	"github.com/stretchr/testify/assert"

	"github.com/go-nacelle/chevron/middleware/mocks"
)

func TestCacheMiddlewareCached(t *testing.T) {
	var (
		cache  = mocks.NewMockCache()
		called = false
	)

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return nil
	}

	cache.GetValueFunc.SetDefaultHook(func(key string) (string, error) {
		assert.Equal(t, "foo.bar", key)

		resp := response.Respond([]byte("foobar"))
		resp.SetStatusCode(http.StatusCreated)
		return serialize(resp)
	})

	wrapped, err := NewResponseCache(cache).Convert(bare)
	assert.Nil(t, err)
	assert.False(t, called)

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())

	// Check returned value
	_, body, err := response.Serialize(resp)
	assert.Nil(t, err)
	assert.Equal(t, "foobar", string(body))
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
}

func TestCacheMiddlewareCacheReadError(t *testing.T) {
	var (
		cache        = mocks.NewMockCache()
		called       = false
		expectedResp = response.JSON(map[string]string{
			"message": "cache whoopsie",
		})
	)

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return nil
	}

	// TODO - can simplify these
	cache.GetValueFunc.SetDefaultHook(func(key string) (string, error) {
		return "", fmt.Errorf("utoh")
	})

	errorFactory := func(err error) response.Response {
		return expectedResp
	}

	wrapped, err := NewResponseCache(
		cache,
		WithCacheErrorFactory(errorFactory),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.False(t, called)
	assert.Equal(t, expectedResp, resp)
}

func TestCacheMiddlewareCacheWriteError(t *testing.T) {
	var (
		cache        = mocks.NewMockCache()
		called       = false
		expectedResp = response.JSON(map[string]string{
			"message": "cache whoopsie",
		})
	)

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return response.Empty(http.StatusNoContent)
	}

	cache.SetValueFunc.SetDefaultHook(func(key, value string, keys ...string) error {
		return fmt.Errorf("utoh")
	})

	errorFactory := func(err error) response.Response {
		return expectedResp
	}

	wrapped, err := NewResponseCache(
		cache,
		WithCacheErrorFactory(errorFactory),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.True(t, called)
	assert.Equal(t, expectedResp, resp)
}

func TestCacheMiddlewareCacheReadCacheJunkData(t *testing.T) {
	var (
		cache        = mocks.NewMockCache()
		called       = false
		expectedResp = response.JSON(map[string]string{
			"message": "cache whoopsie",
		})
	)

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return response.Empty(http.StatusNoContent)
	}

	cache.GetValueFunc.SetDefaultHook(func(key string) (string, error) {
		return "foobar", nil
	})

	errorFactory := func(err error) response.Response {
		return expectedResp
	}

	wrapped, err := NewResponseCache(
		cache,
		WithCacheErrorFactory(errorFactory),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.False(t, called)
	assert.Equal(t, expectedResp, resp)
}

func TestCacheMiddlewareWritesToCache(t *testing.T) {
	var (
		cache = mocks.NewMockCache()
	)

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		resp := response.Respond([]byte("foobar"))
		resp.SetStatusCode(http.StatusCreated)
		return resp
	}

	wrapped, err := NewResponseCache(
		cache,
		WithCacheTags("foo", "bar", "baz"),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())

	// Check returned value
	_, body, err := response.Serialize(resp)
	assert.Nil(t, err)
	assert.Equal(t, "foobar", string(body))
	assert.Equal(t, http.StatusCreated, resp.StatusCode())

	// Check cached values
	assert.Equal(t, 1, len(cache.SetValueFunc.History()))
	params := cache.SetValueFunc.History()[0]
	assert.Equal(t, "foo.bar", params.Arg0)
	vs := params.Arg2 // TODO - rename
	sort.Strings(vs)
	assert.Equal(t, []string{"bar", "baz", "foo"}, vs)

	deserialized, err := deserialize(params.Arg1)
	assert.Nil(t, err)

	_, body2, err := response.Serialize(deserialized)
	assert.Nil(t, err)
	assert.Equal(t, body2, body)
}

func TestCacheMiddlewareShouldNotCache(t *testing.T) {
	var (
		cache  = mocks.NewMockCache()
		called = false
	)

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return nil
	}

	wrapped, err := NewResponseCache(cache).Convert(bare)
	assert.Nil(t, err)
	assert.False(t, called)

	r, _ := http.NewRequest("POST", "/foo/bar", nil)
	wrapped(context.Background(), r, nacelle.NewNilLogger())

	assert.Equal(t, 0, len(cache.GetValueFunc.History()))
	assert.Equal(t, 0, len(cache.SetValueFunc.History()))
}

func TestCacheMiddlewareNoCacheInstance(t *testing.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		resp := response.Respond([]byte("foobar"))
		resp.SetStatusCode(http.StatusCreated)
		return resp
	}

	wrapped, err := NewResponseCache(nil).Convert(bare)
	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	wrapped(context.Background(), r, nacelle.NewNilLogger())
}

func TestCacheMiddlewareReserializeJSON(t *testing.T) {
	testReserialize(t, response.JSON(map[string]string{"foo": "bar", "baz": "bonk"}))
}

func TestCacheMiddlewareReserializeReader(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"foo": "bar", "baz": "bonk"}`))
	testReserialize(t, response.Stream(ioutil.NopCloser(reader)))
}

func testReserialize(t *testing.T, resp response.Response) {
	resp.SetStatusCode(http.StatusCreated)
	resp.AddHeader("X-Order", "a")
	resp.AddHeader("X-Order", "b")
	resp.AddHeader("X-Order", "c")

	serialized, err := serialize(resp)
	assert.Nil(t, err)

	deserialized, err := deserialize(serialized)
	assert.Nil(t, err)

	header, body, err := response.Serialize(deserialized)
	assert.Contains(t, header, "X-Order")
	assert.Equal(t, []string{"a", "b", "c"}, header["X-Order"])
	assert.JSONEq(t, `{"foo": "bar", "baz": "bonk"}`, string(body))
	assert.Equal(t, http.StatusCreated, deserialized.StatusCode())
}

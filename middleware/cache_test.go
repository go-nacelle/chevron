package middleware

//go:generate go-mockgen github.com/efritz/gache -i Cache -d mocks

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aphistic/sweet"
	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	. "github.com/onsi/gomega"

	"github.com/go-nacelle/chevron/middleware/mocks"
)

type CacheSuite struct{}

func (s *CacheSuite) TestCached(t sweet.T) {
	var (
		cache  = mocks.NewMockCache()
		called = false
	)

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return nil
	}

	cache.GetValueFunc = func(key string) (string, error) {
		Expect(key).To(Equal("foo.bar"))

		resp := response.Respond([]byte("foobar"))
		resp.SetStatusCode(http.StatusCreated)
		return serialize(resp)
	}

	wrapped, err := NewResponseCache(cache).Convert(bare)
	Expect(err).To(BeNil())
	Expect(called).To(BeFalse())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())

	// Check returned value
	_, body, err := response.Serialize(resp)
	Expect(err).To(BeNil())
	Expect(string(body)).To(Equal("foobar"))
	Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
}

func (s *CacheSuite) TestCacheReadError(t sweet.T) {
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

	cache.GetValueFunc = func(key string) (string, error) {
		return "", fmt.Errorf("utoh")
	}

	errorFactory := func(err error) response.Response {
		return expectedResp
	}

	wrapped, err := NewResponseCache(
		cache,
		WithCacheErrorFactory(errorFactory),
	).Convert(bare)

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(called).To(BeFalse())
	Expect(resp).To(Equal(expectedResp))
}

func (s *CacheSuite) TestCacheWriteError(t sweet.T) {
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

	cache.SetValueFunc = func(key, value string, keys ...string) error {
		return fmt.Errorf("utoh")
	}

	errorFactory := func(err error) response.Response {
		return expectedResp
	}

	wrapped, err := NewResponseCache(
		cache,
		WithCacheErrorFactory(errorFactory),
	).Convert(bare)

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(called).To(BeTrue())
	Expect(resp).To(Equal(expectedResp))
}

func (s *CacheSuite) TestCacheReadCacheJunkData(t sweet.T) {
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

	cache.GetValueFunc = func(key string) (string, error) {
		return "foobar", nil
	}

	errorFactory := func(err error) response.Response {
		return expectedResp
	}

	wrapped, err := NewResponseCache(
		cache,
		WithCacheErrorFactory(errorFactory),
	).Convert(bare)

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(called).To(BeFalse())
	Expect(resp).To(Equal(expectedResp))
}

func (s *CacheSuite) TestWritesToCache(t sweet.T) {
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

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())

	// Check returned value
	_, body, err := response.Serialize(resp)
	Expect(err).To(BeNil())
	Expect(string(body)).To(Equal("foobar"))
	Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

	// Check cached values
	Expect(cache.SetValueFuncCallCount()).To(Equal(1))
	params := cache.SetValueFuncCallParams()[0]
	Expect(params.Arg0).To(Equal("foo.bar"))
	Expect(params.Arg2).To(ConsistOf([]string{"foo", "bar", "baz"}))

	deserialized, err := deserialize(params.Arg1)
	Expect(err).To(BeNil())

	_, body2, err := response.Serialize(deserialized)
	Expect(err).To(BeNil())
	Expect(body).To(Equal(body2))
}

func (s *CacheSuite) TestShouldNotCache(t sweet.T) {
	var (
		cache  = mocks.NewMockCache()
		called = false
	)

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return nil
	}

	wrapped, err := NewResponseCache(cache).Convert(bare)
	Expect(err).To(BeNil())
	Expect(called).To(BeFalse())

	r, _ := http.NewRequest("POST", "/foo/bar", nil)
	wrapped(context.Background(), r, nacelle.NewNilLogger())

	Expect(cache.GetValueFuncCallCount()).To(Equal(0))
	Expect(cache.SetValueFuncCallCount()).To(Equal(0))
}

func (s *CacheSuite) TestNoCacheInstance(t sweet.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		resp := response.Respond([]byte("foobar"))
		resp.SetStatusCode(http.StatusCreated)
		return resp
	}

	wrapped, err := NewResponseCache(nil).Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	wrapped(context.Background(), r, nacelle.NewNilLogger())
}

func (s *CacheSuite) TestReserializeJSON(t sweet.T) {
	testReserialize(response.JSON(map[string]string{"foo": "bar", "baz": "bonk"}))
}

func (s *CacheSuite) TestReserializeReader(t sweet.T) {
	reader := bytes.NewReader([]byte(`{"foo": "bar", "baz": "bonk"}`))
	testReserialize(response.Stream(ioutil.NopCloser(reader)))
}

func testReserialize(resp response.Response) {
	resp.SetStatusCode(http.StatusCreated)
	resp.AddHeader("X-Order", "a")
	resp.AddHeader("X-Order", "b")
	resp.AddHeader("X-Order", "c")

	serialized, err := serialize(resp)
	Expect(err).To(BeNil())

	deserialized, err := deserialize(serialized)
	Expect(err).To(BeNil())

	header, body, err := response.Serialize(deserialized)
	Expect(header).To(HaveKey("X-Order"))
	Expect(header["X-Order"]).To(Equal([]string{"a", "b", "c"}))
	Expect(body).To(MatchJSON(`{"foo": "bar", "baz": "bonk"}`))
	Expect(deserialized.StatusCode()).To(Equal(http.StatusCreated))
}

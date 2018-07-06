package middleware

import "github.com/efritz/gache"

// CacheMiddlewareConfigFunc is a function used to initialize a new cache
// middleware instance.
type CacheMiddlewareConfigFunc func(*CacheMiddleware)

// WithCacheInstance sets the cache instance. The cache instance can also
// be set by injection via a nacelle service container.
func WithCacheInstance(cache gache.Cache) CacheMiddlewareConfigFunc {
	return func(c *CacheMiddleware) { c.Cache = cache }
}

// WithCacheTags sets the tags applied to keys written to the cache by the
// associated middleware.
func WithCacheTags(tags ...string) CacheMiddlewareConfigFunc {
	return func(c *CacheMiddleware) { c.tags = tags }
}

func WithCacheErrorFactory(factory ErrorFactory) CacheMiddlewareConfigFunc {
	return func(m *CacheMiddleware) { m.errorFactory = factory }
}

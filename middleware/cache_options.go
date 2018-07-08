package middleware

// CacheMiddlewareConfigFunc is a function used to initialize a new cache
// middleware instance.
type CacheMiddlewareConfigFunc func(*CacheMiddleware)

// WithCacheTags sets the tags applied to keys written to the cache by the
// associated middleware.
func WithCacheTags(tags ...string) CacheMiddlewareConfigFunc {
	return func(c *CacheMiddleware) { c.tags = tags }
}

func WithCacheErrorFactory(factory ErrorFactory) CacheMiddlewareConfigFunc {
	return func(m *CacheMiddleware) { m.errorFactory = factory }
}

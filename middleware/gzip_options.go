package middleware

type GzipMiddlewareConfigFunc func(*GzipMiddleware)

func WithGzipLevel(level int) GzipMiddlewareConfigFunc {
	return func(c *GzipMiddleware) { c.level = level }
}

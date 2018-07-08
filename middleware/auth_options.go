package middleware

type AuthMiddlewareConfigFunc func(*AuthMiddleware)

func WithAuthErrorFactory(factory ErrorFactory) AuthMiddlewareConfigFunc {
	return func(m *AuthMiddleware) { m.errorFactory = factory }
}

func WithAuthForbiddenResponseFactory(factory ResponseFactory) AuthMiddlewareConfigFunc {
	return func(m *AuthMiddleware) { m.forbiddenResponseFactory = factory }
}

func WithAuthUnauthorizedResponseFactory(factory ErrorFactory) AuthMiddlewareConfigFunc {
	return func(m *AuthMiddleware) { m.unauthorizedResponseFactory = factory }
}

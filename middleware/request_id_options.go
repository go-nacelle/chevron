package middleware

type RequestIDConfigFunc func(m *RequestIDMiddleware)

func WithRequestIDErrorFactory(factory ErrorFactory) RequestIDConfigFunc {
	return func(m *RequestIDMiddleware) { m.errorFactory = factory }
}

func WithRequestIDGenerator(generator RequestIDGenerator) RequestIDConfigFunc {
	return func(m *RequestIDMiddleware) { m.requestIDGenerator = generator }
}

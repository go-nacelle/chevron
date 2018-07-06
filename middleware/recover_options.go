package middleware

type (
	RecoverConfigFunc func(m *RecoverMiddleware)
)

func WithRecoverErrorFactory(factory PanicErrorFactory) RecoverConfigFunc {
	return func(m *RecoverMiddleware) { m.errorFactory = factory }
}

func WithRecoverStackBufferSize(stackBufferSize int) RecoverConfigFunc {
	return func(m *RecoverMiddleware) { m.stackBufferSize = stackBufferSize }
}

func WithRecoverLogAllGoroutines(logAllGoroutines bool) RecoverConfigFunc {
	return func(m *RecoverMiddleware) { m.logAllGoroutines = logAllGoroutines }
}

package middleware

import "github.com/derision-test/glock"

type LoggingConfigFunc func(m *LoggingMiddleware)

func WithLoggingClock(clock glock.Clock) LoggingConfigFunc {
	return func(m *LoggingMiddleware) { m.clock = clock }
}

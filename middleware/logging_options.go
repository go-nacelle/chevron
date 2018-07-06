package middleware

import "github.com/efritz/glock"

type LoggingConfigFunc func(m *LoggingMiddleware)

func WithLoggingClock(clock glock.Clock) LoggingConfigFunc {
	return func(m *LoggingMiddleware) { m.clock = clock }
}

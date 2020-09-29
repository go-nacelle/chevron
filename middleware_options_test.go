package chevron

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithMiddleware(t *testing.T) {
	var (
		hm       = makeTestHandlerMap()
		numCalls = 0
	)

	middleware := MiddlewareFunc(func(h Handler) (Handler, error) {
		numCalls++
		return makeEmptyHandler(106), nil
	})

	// Apply the middleware config
	assert.Nil(t, WithMiddleware(middleware)(hm))

	assert.Equal(t, 6, numCalls)
	assert.Equal(t, 106, hm[MethodGet](nil, nil, nil).StatusCode())
	assert.Equal(t, 106, hm[MethodOptions](nil, nil, nil).StatusCode())
	assert.Equal(t, 106, hm[MethodPost](nil, nil, nil).StatusCode())
	assert.Equal(t, 106, hm[MethodPut](nil, nil, nil).StatusCode())
	assert.Equal(t, 106, hm[MethodPatch](nil, nil, nil).StatusCode())
	assert.Equal(t, 106, hm[MethodDelete](nil, nil, nil).StatusCode())
}

func TestWithMiddlewareError(t *testing.T) {
	middleware := MiddlewareFunc(func(h Handler) (Handler, error) {
		return nil, fmt.Errorf("utoh")
	})

	// Apply the middleware config
	assert.EqualError(t, WithMiddleware(middleware)(makeTestHandlerMap()), "utoh")
}

func TestWithMiddlewareFor(t *testing.T) {
	var (
		hm       = makeTestHandlerMap()
		numCalls = 0
	)

	middleware := MiddlewareFunc(func(h Handler) (Handler, error) {
		numCalls++
		return makeEmptyHandler(106), nil
	})

	// Apply the middleware config
	assert.Nil(t, WithMiddlewareFor(middleware, MethodGet, MethodPatch)(hm))

	assert.Equal(t, 2, numCalls)
	assert.Equal(t, 106, hm[MethodGet](nil, nil, nil).StatusCode())
	assert.Equal(t, 101, hm[MethodOptions](nil, nil, nil).StatusCode())
	assert.Equal(t, 102, hm[MethodPost](nil, nil, nil).StatusCode())
	assert.Equal(t, 103, hm[MethodPut](nil, nil, nil).StatusCode())
	assert.Equal(t, 106, hm[MethodPatch](nil, nil, nil).StatusCode())
	assert.Equal(t, 105, hm[MethodDelete](nil, nil, nil).StatusCode())
}

//
//

func makeTestHandlerMap() handlerMap {
	return handlerMap{
		MethodGet:     makeEmptyHandler(100),
		MethodOptions: makeEmptyHandler(101),
		MethodPost:    makeEmptyHandler(102),
		MethodPut:     makeEmptyHandler(103),
		MethodPatch:   makeEmptyHandler(104),
		MethodDelete:  makeEmptyHandler(105),
	}
}

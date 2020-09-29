package chevron

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceHandle(t *testing.T) {
	m := map[Method]int{
		MethodGet:  http.StatusFound,
		MethodPost: http.StatusMovedPermanently,
	}

	handlerMap := map[Method]Handler{}
	for method, status := range m {
		handlerMap[method] = makeEmptyHandler(status)
	}

	r := &resource{
		hm: handlerMap,
		router: &router{
			notImplementedHandler: defaultNotImplementedHandler,
		},
	}

	for _, method := range allMethods {
		expected, ok := m[method]
		if !ok {
			expected = http.StatusMethodNotAllowed
		}

		assert.Equal(t, expected, r.Handle(nil, makeEmptyRequest(method), nil).StatusCode())
	}
}

//
//

func makeEmptyRequest(method Method) *http.Request {
	req, _ := http.NewRequest(method.String(), "", nil)
	return req
}

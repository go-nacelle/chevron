package chevron

import (
	"net/http"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ResourceSuite struct{}

func (s *ResourceSuite) TestHandle(t sweet.T) {
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

		Expect(r.Handle(nil, makeEmptyRequest(method), nil).StatusCode()).To(Equal(expected))
	}
}

//
//

func makeEmptyRequest(method Method) *http.Request {
	req, _ := http.NewRequest(method.String(), "", nil)
	return req
}

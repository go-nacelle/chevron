package chevron

import (
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type MiddlewareOptionsSuite struct{}

func (s *MiddlewareOptionsSuite) TestWithMiddleware(t sweet.T) {
	var (
		hm       = makeTestHandlerMap()
		numCalls = 0
	)

	middleware := MiddlewareFunc(func(h Handler) (Handler, error) {
		numCalls++
		return makeEmptyHandler(106), nil
	})

	// Apply the middleware config
	Expect(WithMiddleware(middleware)(hm)).To(BeNil())

	Expect(numCalls).To(Equal(6))
	Expect(hm[MethodGet](nil, nil, nil).StatusCode()).To(Equal(106))
	Expect(hm[MethodOptions](nil, nil, nil).StatusCode()).To(Equal(106))
	Expect(hm[MethodPost](nil, nil, nil).StatusCode()).To(Equal(106))
	Expect(hm[MethodPut](nil, nil, nil).StatusCode()).To(Equal(106))
	Expect(hm[MethodPatch](nil, nil, nil).StatusCode()).To(Equal(106))
	Expect(hm[MethodDelete](nil, nil, nil).StatusCode()).To(Equal(106))
}

func (s *MiddlewareOptionsSuite) TestWithMiddlewareError(t sweet.T) {
	middleware := MiddlewareFunc(func(h Handler) (Handler, error) {
		return nil, fmt.Errorf("utoh")
	})

	// Apply the middleware config
	Expect(WithMiddleware(middleware)(makeTestHandlerMap())).To(MatchError("utoh"))
}

func (s *MiddlewareOptionsSuite) TestWithMiddlewareFor(t sweet.T) {
	var (
		hm       = makeTestHandlerMap()
		numCalls = 0
	)

	middleware := MiddlewareFunc(func(h Handler) (Handler, error) {
		numCalls++
		return makeEmptyHandler(106), nil
	})

	// Apply the middleware config
	Expect(WithMiddlewareFor(middleware, MethodGet, MethodPatch)(hm)).To(BeNil())

	Expect(numCalls).To(Equal(2))
	Expect(hm[MethodGet](nil, nil, nil).StatusCode()).To(Equal(106))
	Expect(hm[MethodOptions](nil, nil, nil).StatusCode()).To(Equal(101))
	Expect(hm[MethodPost](nil, nil, nil).StatusCode()).To(Equal(102))
	Expect(hm[MethodPut](nil, nil, nil).StatusCode()).To(Equal(103))
	Expect(hm[MethodPatch](nil, nil, nil).StatusCode()).To(Equal(106))
	Expect(hm[MethodDelete](nil, nil, nil).StatusCode()).To(Equal(105))
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

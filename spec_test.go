package chevron

import (
	"context"
	"net/http"

	"github.com/aphistic/sweet"
	"github.com/efritz/nacelle"
	"github.com/efritz/response"
	. "github.com/onsi/gomega"
)

type SpecSuite struct{}

func (s *SpecSuite) TestEmptySpecComposition(t sweet.T) {
	// Static test for conformance to interface
	var ts ResourceSpec = &TestSpec{}

	// Show that we can "override" a method
	resp := ts.Get(testBackground(), nil, nil)
	Expect(resp.StatusCode()).To(Equal(http.StatusOK))

	// Check response
	data, err := response.Serialize(resp)
	Expect(err).To(BeNil())
	Expect(data).To(MatchJSON(`["foo", "bar", "baz"]`))

	// The remaining methods should fall through to the default

	for _, handler := range []Handler{ts.Options, ts.Post, ts.Put, ts.Patch, ts.Delete} {
		Expect(handler(testBackground(), nil, nil).StatusCode()).To(Equal(http.StatusMethodNotAllowed))
	}
}

func testBackground() context.Context {
	return setNotImplementedHandler(context.Background(), defaultNotImplementedHandler)
}

//
//

type TestSpec struct {
	*EmptySpec
}

func (ts *TestSpec) Get(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
	return response.JSON([]string{"foo", "bar", "baz"})
}

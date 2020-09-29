package chevron

import (
	"context"
	"net/http"
	"testing"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	"github.com/stretchr/testify/assert"
)

func TestEmptySpecComposition(t *testing.T) {
	// Static test for conformance to interface
	var ts ResourceSpec = &TestSpec{}

	// Show that we can "override" a method
	resp := ts.Get(testBackground(), nil, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Check response
	_, data, err := response.Serialize(resp)
	assert.Nil(t, err)
	assert.JSONEq(t, `["foo", "bar", "baz"]`, string(data))

	// The remaining methods should fall through to the default

	for _, handler := range []Handler{ts.Options, ts.Post, ts.Put, ts.Patch, ts.Delete} {
		assert.Equal(t, http.StatusMethodNotAllowed, handler(testBackground(), nil, nil).StatusCode())
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

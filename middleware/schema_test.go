package middleware

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"

	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
)

// Assigned at setup
var TestSchemaPath string

func TestMain(t *testing.M) {
	wd, _ := os.Getwd()
	TestSchemaPath = filepath.Join(wd, "test-schemas")
	os.Exit(t.Run())
}

func TestSchemaGetJSONData(t *testing.T) {
	var (
		val = []byte("[1, 2, 3]")
		ctx = context.WithValue(context.Background(), TokenJSONData, val)
	)

	// Present
	assert.Equal(t, val, GetJSONData(ctx))

	// Missing
	assert.Empty(t, GetJSONData(context.Background()))
}

func TestSchemaValidateInput(t *testing.T) {
	var ctxVal []byte
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		ctxVal = GetJSONData(ctx)
		return response.Empty(http.StatusNoContent)
	}

	wrapped, err := NewSchemaMiddleware(TestSchemaPath + "/point.json").Convert(bare)
	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/", strings.NewReader(`{
		"x": 1,
		"y": 2,
		"z": 3
	}`))

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.JSONEq(t, `{"x": 1, "y": 2, "z": 3}`, string(ctxVal))
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
}

func TestSchemaValidateInputYAMLSchema(t *testing.T) {
	var ctxVal []byte
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		ctxVal = GetJSONData(ctx)
		return response.Empty(http.StatusNoContent)
	}

	wrapped, err := NewSchemaMiddleware(TestSchemaPath + "/point.yaml").Convert(bare)
	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/", strings.NewReader(`{
		"x": 1,
		"y": 2,
		"z": 3
	}`))

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.JSONEq(t, `{"x": 1, "y": 2, "z": 3}`, string(ctxVal))
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
}

func TestSchemaBadRequest(t *testing.T) {
	var (
		called       = false
		expectedResp = response.JSON(map[string]string{
			"message": "json whoopsie",
		})
	)

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return response.Empty(http.StatusNoContent)
	}

	badRequestFactory := func() response.Response {
		return expectedResp
	}

	wrapped, err := NewSchemaMiddleware(
		TestSchemaPath+"/point.json",
		WithSchemaBadRequestFactory(badRequestFactory),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/", strings.NewReader(`not even json`))
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.False(t, called)
	assert.Equal(t, expectedResp, resp)
}

func TestSchemaUnprocessableEntity(t *testing.T) {
	var (
		called = false
	)

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		called = true
		return response.Empty(http.StatusNoContent)
	}

	unprocessableEntityFactory := func(errors []gojsonschema.ResultError) response.Response {
		errs := []string{}
		for _, err := range errors {
			errs = append(errs, err.Description())
		}

		resp := response.JSON(errs)
		resp.SetStatusCode(http.StatusUnprocessableEntity)
		return resp
	}

	wrapped, err := NewSchemaMiddleware(
		TestSchemaPath+"/point.json",
		WithSchemaUnprocessableEntityFactory(unprocessableEntityFactory),
	).Convert(bare)

	assert.Nil(t, err)

	r, _ := http.NewRequest("GET", "/", strings.NewReader(`{
		"x": 1,
		"y": 3.5
	}`))

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	assert.False(t, called)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode())

	_, body, err := response.Serialize(resp)
	assert.Nil(t, err)

	expected := `[
		"z is required",
		"Invalid type. Expected: integer, given: number"
	]`
	assert.JSONEq(t, expected, string(body))
}

func TestSchemaMissingSchema(t *testing.T) {
	_, err := NewSchemaMiddleware(TestSchemaPath + "/missing.json").Convert(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to load schema")
}

func TestSchemaBadSchema(t *testing.T) {
	_, err := NewSchemaMiddleware(TestSchemaPath + "/malformed.json").Convert(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid schema")
}

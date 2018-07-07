package middleware

import (
	"context"
	"go/build"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/aphistic/sweet"
	"github.com/efritz/nacelle"
	"github.com/efritz/response"
	. "github.com/onsi/gomega"
)

type SchemaSuite struct{}

// Assigned at setup
var TestSchemaPath string

func (s *SchemaSuite) TestGetJSONData(t sweet.T) {
	var (
		val = []byte("[1, 2, 3]")
		ctx = context.WithValue(context.Background(), TokenJSONData, val)
	)

	// Present
	Expect(GetJSONData(ctx)).To(Equal(val))

	// Missing
	Expect(GetJSONData(context.Background())).To(BeEmpty())
}

func (s *SchemaSuite) TestValidateInput(t sweet.T) {
	var ctxVal []byte
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		ctxVal = GetJSONData(ctx)
		return response.Empty(http.StatusNoContent)
	}

	wrapped, err := NewSchemaMiddleware(TestSchemaPath + "/point.json").Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/", strings.NewReader(`{
		"x": 1,
		"y": 2,
		"z": 3
	}`))

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(ctxVal).To(MatchJSON(`{"x": 1, "y": 2, "z": 3}`))
	Expect(resp.StatusCode()).To(Equal(http.StatusNoContent))
}

func (s *SchemaSuite) TestValidateInputYAMLSchema(t sweet.T) {
	var ctxVal []byte
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		ctxVal = GetJSONData(ctx)
		return response.Empty(http.StatusNoContent)
	}

	wrapped, err := NewSchemaMiddleware(TestSchemaPath + "/point.yaml").Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/", strings.NewReader(`{
		"x": 1,
		"y": 2,
		"z": 3
	}`))

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(ctxVal).To(MatchJSON(`{"x": 1, "y": 2, "z": 3}`))
	Expect(resp.StatusCode()).To(Equal(http.StatusNoContent))
}

func (s *SchemaSuite) TestBadRequest(t sweet.T) {
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

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/", strings.NewReader(`not even json`))
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(called).To(BeFalse())
	Expect(resp).To(Equal(expectedResp))
}

func (s *SchemaSuite) TestUnprocessableEntity(t sweet.T) {
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

	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/", strings.NewReader(`{
		"x": 1,
		"y": 3.5
	}`))

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	Expect(called).To(BeFalse())
	Expect(resp.StatusCode()).To(Equal(http.StatusUnprocessableEntity))

	_, body, err := response.Serialize(resp)
	Expect(err).To(BeNil())
	Expect(body).To(MatchJSON(`[
		"z is required",
		"Invalid type. Expected: integer, given: number"
	]`))
}

func (s *SchemaSuite) TestMissingSchema(t sweet.T) {
	_, err := NewSchemaMiddleware(TestSchemaPath + "/missing.json").Convert(nil)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("failed to load schema"))
}

func (s *SchemaSuite) TestBadSchema(t sweet.T) {
	_, err := NewSchemaMiddleware(TestSchemaPath + "/malformed.json").Convert(nil)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("invalid schema"))
}

//
// Point to correct schema path

func (s *SchemaSuite) SetUpSuite() {
	TestSchemaPath = filepath.Join(gopath(), "src", "github.com/efritz/chevron/middleware", "test-schemas")
}

func gopath() string {
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return gopath
	}

	return build.Default.GOPATH
}

//
//

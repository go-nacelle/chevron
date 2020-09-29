package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/efritz/response"
	"github.com/ghodss/yaml"
	"github.com/go-nacelle/nacelle"
	"github.com/xeipuuv/gojsonschema"

	"github.com/go-nacelle/chevron"
)

type SchemaMiddleware struct {
	path                       string
	errorFactory               ErrorFactory
	badRequestFactory          SchemaBadRequestFactory
	unprocessableEntityFactory SchemaUnprocessableEntityFactory
}

type SchemaBadRequestFactory func() response.Response
type SchemaUnprocessableEntityFactory func([]gojsonschema.ResultError) response.Response

type tokenJSONData string

var TokenJSONData = tokenJSONData("chevron.middleware.json_data")

func GetJSONData(ctx context.Context) []byte {
	if val, ok := ctx.Value(TokenJSONData).([]byte); ok {
		return val
	}

	return nil
}

func NewSchemaMiddleware(path string, configs ...SchemaConfigFunc) chevron.Middleware {
	m := &SchemaMiddleware{
		path:                       path,
		errorFactory:               defaultErrorFactory,
		badRequestFactory:          defaultBadRequestFactory,
		unprocessableEntityFactory: defaultUnprocessableEntityFactory,
	}

	for _, f := range configs {
		f(m)
	}

	return m
}

func (m *SchemaMiddleware) Convert(f chevron.Handler) (chevron.Handler, error) {
	schema, err := loadSchema(m.path)
	if err != nil {
		return nil, err
	}

	handler := func(ctx context.Context, req *http.Request, logger nacelle.Logger) response.Response {
		defer req.Body.Close()

		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Error("Failed to read request body (%s)", err.Error())
			return m.errorFactory(err)
		}

		if !isJSON(data) {
			return m.badRequestFactory()
		}

		result, err := schema.Validate(gojsonschema.NewStringLoader(string(data)))
		if err != nil {
			logger.Error("Failed to load json schema", err.Error())
			return m.errorFactory(err)
		}

		if !result.Valid() {
			return m.unprocessableEntityFactory(result.Errors())
		}

		return f(context.WithValue(ctx, TokenJSONData, data), req, logger)
	}

	return handler, nil
}

//
// Helpers

func isJSON(data []byte) bool {
	return json.Unmarshal(data, &json.RawMessage{}) == nil
}

func loadSchema(path string) (*gojsonschema.Schema, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("failed to load schema %s (%s)", path, err.Error())
	}

	loader, err := getJSONLoader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load schema %s (%s)", path, err.Error())
	}

	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return nil, fmt.Errorf("invalid schema %s (%s)", path, err.Error())
	}

	return schema, nil
}

func getJSONLoader(path string) (loader gojsonschema.JSONLoader, err error) {
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		loader, err = loadYAMLSchema(path)
		return
	}

	loader = loadJSONSchema(path)
	return
}

func loadJSONSchema(path string) gojsonschema.JSONLoader {
	return gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", path))
}

func loadYAMLSchema(path string) (gojsonschema.JSONLoader, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	json, err := yaml.YAMLToJSON(content)
	if err != nil {
		return nil, err
	}

	return gojsonschema.NewBytesLoader(json), nil
}

func defaultBadRequestFactory() response.Response {
	return response.Empty(http.StatusBadRequest)
}

func defaultUnprocessableEntityFactory([]gojsonschema.ResultError) response.Response {
	return response.Empty(http.StatusUnprocessableEntity)
}

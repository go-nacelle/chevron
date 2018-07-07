package middleware

type (
	SchemaConfigFunc func(m *SchemaMiddleware)
)

func WithSchemaErrorFactory(factory ErrorFactory) SchemaConfigFunc {
	return func(m *SchemaMiddleware) { m.errorFactory = factory }
}

func WithSchemaBadRequestFactory(factory SchemaBadRequestFactory) SchemaConfigFunc {
	return func(m *SchemaMiddleware) { m.badRequestFactory = factory }
}

func WithSchemaUnprocessableEntityFactory(factory SchemaUnprocessableEntityFactory) SchemaConfigFunc {
	return func(m *SchemaMiddleware) { m.unprocessableEntityFactory = factory }
}

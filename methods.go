package chevron

// Method is an enumeration of HTTP methods.
type Method int

const (
	// MethodGet represents the GET HTTP method.
	MethodGet Method = iota

	// MethodOptions represents the OPTIONS HTTP method.
	MethodOptions

	// MethodPost represents the POST HTTP method.
	MethodPost

	// MethodPut represents the PUT HTTP method.
	MethodPut

	// MethodPatch represents the PATCH HTTP method.
	MethodPatch

	// MethodDelete represents the DELETE HTTP method.
	MethodDelete
)

var allMethods = []Method{
	MethodGet,
	MethodOptions,
	MethodPost,
	MethodPut,
	MethodPatch,
	MethodDelete,
}

var methodMap = map[string]Method{
	"GET":     MethodGet,
	"OPTIONS": MethodOptions,
	"POST":    MethodPost,
	"PUT":     MethodPut,
	"PATCH":   MethodPatch,
	"DELETE":  MethodDelete,
}

var methodStrings = map[Method]string{
	MethodGet:     "GET",
	MethodOptions: "OPTIONS",
	MethodPost:    "POST",
	MethodPut:     "PUT",
	MethodPatch:   "PATCH",
	MethodDelete:  "DELETE",
}

// String returns the uppercased HTTP method name.
func (m Method) String() string {
	return methodStrings[m]
}

package chevron

type Method int

const (
	MethodGet Method = iota
	MethodOptions
	MethodPost
	MethodPut
	MethodPatch
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

func (m Method) String() string {
	return methodStrings[m]
}

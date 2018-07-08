package middleware

import (
	"net/http"

	"github.com/efritz/response"
)

type (
	ResponseFactory   func() response.Response
	ErrorFactory      func(error) response.Response
	PanicErrorFactory func(interface{}) response.Response
)

func defaultErrorFactory(val error) response.Response {
	return response.Empty(http.StatusInternalServerError)
}

func defaultPanicErrorFactory(val interface{}) response.Response {
	return response.Empty(http.StatusInternalServerError)
}

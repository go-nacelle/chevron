package middleware

import (
	"net/http"

	"github.com/efritz/response"
)

type ResponseFactory func() response.Response

type ErrorFactory func(error) response.Response

func defaultErrorFactory(val error) response.Response {
	return response.Empty(http.StatusInternalServerError)
}

type PanicErrorFactory func(interface{}) response.Response

func defaultPanicErrorFactory(val interface{}) response.Response {
	return response.Empty(http.StatusInternalServerError)
}

package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go/request"
)

type jwtQueryExtractor struct {
	name string
}

type jwtHeaderExtractor struct {
	name   string
	scheme string
}

func NewJWTQueryExtractor(name string) request.Extractor {
	return &jwtQueryExtractor{name: name}
}

func NewJWTHeaderExtractor(name, scheme string) request.Extractor {
	return &jwtHeaderExtractor{name: name, scheme: scheme}
}

func (e *jwtQueryExtractor) ExtractToken(req *http.Request) (string, error) {
	if token := req.URL.Query().Get(e.name); token != "" {
		return token, nil
	}

	return "", request.ErrNoTokenInRequest
}

func (e *jwtHeaderExtractor) ExtractToken(req *http.Request) (string, error) {
	value := req.Header.Get(e.name)
	if strings.HasPrefix(value, e.scheme) {
		return value[len(e.scheme)+1:], nil
	}

	return "", request.ErrNoTokenInRequest
}

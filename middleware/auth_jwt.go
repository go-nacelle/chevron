package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type (
	jwtAuthorizer struct {
		keyfunc       jwt.Keyfunc
		extractor     request.Extractor
		claimsFactory JWTClaimsFactory
	}

	JWTClaimsFactory func() jwt.Claims
)

func GetJWTClaims(ctx context.Context) jwt.Claims {
	if claims, ok := ctx.Value(TokenAuthPayload).(jwt.Claims); ok {
		return claims
	}

	return nil
}

func NewJWTAuthorizer(keyfunc jwt.Keyfunc, configs ...JWTAuthorizerConfigFunc) Authorizer {
	a := &jwtAuthorizer{
		keyfunc:       keyfunc,
		extractor:     NewJWTQueryExtractor("jwt"),
		claimsFactory: defaultJWTClaimsFactory,
	}

	for _, f := range configs {
		f(a)
	}

	return a
}

func (a *jwtAuthorizer) Authorize(ctx context.Context, req *http.Request) (AuthResult, interface{}, error) {
	token, err := request.ParseFromRequest(
		req,
		a.extractor,
		a.wrappedKeyFunc,
		request.WithClaims(a.claimsFactory()),
	)

	if err != nil {
		if err == request.ErrNoTokenInRequest {
			return AuthResultUnauthorized, nil, err
		}

		if vErr, ok := err.(*jwt.ValidationError); ok {
			return AuthResultUnauthorized, nil, vErr
		}

		return AuthResultInvalid, nil, err
	}

	return AuthResultOK, token.Claims, nil
}

func (a *jwtAuthorizer) wrappedKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return a.keyfunc(token)
}

func defaultJWTClaimsFactory() jwt.Claims {
	return jwt.MapClaims{}
}

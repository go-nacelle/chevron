package middleware

import "github.com/dgrijalva/jwt-go/request"

type JWTAuthorizerConfigFunc func(*jwtAuthorizer)

func WithJWTAuthExtractor(extractor request.Extractor) JWTAuthorizerConfigFunc {
	return func(m *jwtAuthorizer) { m.extractor = extractor }
}

func WithJWTAuthClaimsFactory(factory JWTClaimsFactory) JWTAuthorizerConfigFunc {
	return func(m *jwtAuthorizer) { m.claimsFactory = factory }
}

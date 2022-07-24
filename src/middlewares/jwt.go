package middlewares

import (
	"github.com/ppcamp/go-auth/jwt"
	"github.com/tchap/go-patricia/v2/patricia"
)

type JwtMiddleware struct {
	jwt.Jwt

	Trie     *patricia.Trie
	Patterns map[string]any
}

func (s *JwtMiddleware) Data(signedToken string) (any, error) {
	return s.Session(signedToken)
}

func (s *JwtMiddleware) Allow(path string) bool {
	return true
}

package auth

import (
	"github.com/ppcamp/go-auth-microservice/src/repositories/cache"
	"github.com/ppcamp/go-auth-microservice/src/services"
	"github.com/ppcamp/go-auth/jwt"
)

type InvalidateTokensIn struct {
	User  string
	Token string
}

type InvalidateTokensOut struct{}

type invalidateTokensService[In, Out any] struct {
	services.BaseBusiness

	cache  cache.Auth
	signer jwt.Jwt
}

// NewInvalidateTokensService creates a service that get user password, check it, and
// return a valid JWT token
func NewInvalidateTokensService(
	repo cache.Auth,
	signer jwt.Jwt,
) services.IBaseBusiness[InvalidateTokensIn, InvalidateTokensOut] {
	return &invalidateTokensService[InvalidateTokensIn, InvalidateTokensOut]{
		cache:  repo,
		signer: signer,
	}
}

func (s *invalidateTokensService[In, Out]) Execute(in InvalidateTokensIn) (*InvalidateTokensOut, error) {
	if _, err := s.signer.Session(in.Token); err != nil {
		return nil, err
	}
	err := s.cache.InvalidateAll(s.Context, in.User)
	return new(InvalidateTokensOut), err
}

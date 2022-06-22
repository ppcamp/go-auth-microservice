package auth

import (
	"authentication/helpers/services"
	"authentication/repositories/cache"
	"authentication/utils/jwt"
)

type InvalidateTokensIn struct {
	User  string
	Token string
}

type InvalidateTokensOut struct{}

type invalidateTokensService[In, Out any] struct {
	services.BaseBusiness

	cache cache.Auth
}

// NewInvalidateTokensService creates a service that get user password, check it, and
// return a valid JWT token
func NewInvalidateTokensService(
	repo cache.Auth,
) services.IBaseBusiness[InvalidateTokensIn, InvalidateTokensOut] {
	return &invalidateTokensService[InvalidateTokensIn, InvalidateTokensOut]{cache: repo}
}

func (s *invalidateTokensService[In, Out]) Execute(in InvalidateTokensIn) (*InvalidateTokensOut, error) {
	if _, err := jwt.Signer.Session(in.Token); err != nil {
		return nil, err
	}
	err := s.cache.InvalidateAll(s.Context, in.User)
	return new(InvalidateTokensOut), err
}

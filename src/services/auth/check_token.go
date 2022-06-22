package auth

import (
	"authentication/helpers/services"
	"authentication/repositories/cache"
)

type CheckTokenIn struct {
	Token string
}

type CheckTokenOut struct{}

type checkTokenService[In, Out any] struct {
	services.BaseBusiness

	cache cache.Auth
}

// NewCheckTokenService creates a service that get user password, check it, and
// return a valid JWT token
func NewCheckTokenService(repo cache.Auth) services.IBaseBusiness[CheckTokenIn, CheckTokenOut] {
	return &checkTokenService[CheckTokenIn, CheckTokenOut]{cache: repo}
}

func (s *checkTokenService[In, Out]) Execute(in CheckTokenIn) (*CheckTokenOut, error) {
	err := s.cache.Valid(s.Context, in.Token, in.Token)
	if err != nil {
		return nil, err
	}
	return new(CheckTokenOut), nil
}

package auth

import (
	"time"

	"authentication/helpers/services"
	"authentication/repositories/cache"
	"authentication/utils/jwt"
)

type RefreshTokenIn struct {
	Token string
}

type RefreshTokenOut struct {
	Token string
	Exp   time.Time
}

type refreshTokenService[In, Out any] struct {
	services.BaseBusiness

	cache     cache.Auth
	signerExp time.Duration
}

// NewRefreshTokenService creates a service that get user password, check it, and
// return a valid JWT token
func NewRefreshTokenService(
	repo cache.Auth,
	exp time.Duration,
) services.IBaseBusiness[RefreshTokenIn, RefreshTokenOut] {
	return &refreshTokenService[RefreshTokenIn, RefreshTokenOut]{cache: repo, signerExp: exp}
}

func (s *refreshTokenService[In, Out]) Execute(in RefreshTokenIn) (*RefreshTokenOut, error) {
	exp := time.Now().Add(s.signerExp)

	session, err := jwt.Signer.Session(in.Token)
	if err != nil {
		return nil, err
	}

	err = s.cache.Invalidate(s.Context, *session.UserId, in.Token)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Signer.Generate(session, s.signerExp)
	if err != nil {
		return nil, err
	}

	err = s.cache.Authorize(s.Context, *session.UserId, token, exp)
	return new(RefreshTokenOut), err
}

package auth

import (
	"context"
	"errors"
	"time"

	handlers "github.com/ppcamp/go-auth-microservice/src/http"
	auth "github.com/ppcamp/go-auth-microservice/src/services/authentication"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errNotImplemented error = errors.New("not implemented")

const tokenExp time.Duration = time.Second * 60 * 60 * 60

type AuthService struct {
	UnsafeAuthServiceServer

	*handlers.Handler
	tokenExp time.Duration
}

func NewAuthService(handler *handlers.Handler) AuthServiceServer {
	return &AuthService{Handler: handler, tokenExp: tokenExp}
}

func (s *AuthService) Login(ctx context.Context, pl *LoginInput) (*AuthOutput, error) {
	input := auth.LoginIn{User: pl.User, Password: pl.Password}

	service := auth.NewLoginService(s.Cache, s.tokenExp, s.Signer)

	response, err := handlers.Handle[auth.LoginIn, auth.LoginOut](ctx, s.Handler, input, service)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		Token: response.Token,
		Exp:   timestamppb.New(response.Exp),
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, pl *TokenInput) (*AuthOutput, error) {
	input := auth.RefreshTokenIn{Token: pl.Token}

	service := auth.NewRefreshTokenService(s.Cache, s.tokenExp, s.Signer)

	response, err := handlers.Handle(ctx, s.Handler, input, service)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		Token: response.Token,
		Exp:   timestamppb.New(response.Exp),
	}, nil
}

func (s *AuthService) InvalidateAll(ctx context.Context, pl *SessionsInput) (*empty.Empty, error) {
	input := auth.InvalidateTokensIn{User: pl.User, Token: pl.Token}

	service := auth.NewInvalidateTokensService(s.Cache, s.Signer)

	_, err := handlers.Handle(ctx, s.Handler, input, service)
	if err != nil {
		return nil, err
	}

	return new(empty.Empty), nil
}

func (s *AuthService) IsValid(ctx context.Context, pl *TokenInput) (*empty.Empty, error) {
	input := auth.CheckTokenIn{Token: pl.Token}

	service := auth.NewCheckTokenService(s.Cache)

	_, err := handlers.Handle(ctx, s.Handler, input, service)
	if err != nil {
		return nil, err
	}

	return new(empty.Empty), nil
}

func (s *AuthService) ActiveSessions(context.Context, *SessionsInput) (*SessionsOutput, error) {
	return nil, errNotImplemented
}

func (s *AuthService) Invalidate(ctx context.Context, in *TokenInput) (*empty.Empty, error) {
	return new(empty.Empty), errNotImplemented
}

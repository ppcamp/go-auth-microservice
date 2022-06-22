package auth

import (
	"context"
	"errors"
	"time"

	"authentication/helpers/handlers"
	"authentication/services/auth"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrNotImplemented error = errors.New("not implemented")

type AuthService struct {
	UnsafeAuthServiceServer

	handler  *handlers.Handler
	tokenExp time.Duration
}

func NewAuthService(handler *handlers.Handler) AuthServiceServer {
	return &AuthService{handler: handler, tokenExp: time.Second * 60 * 60 * 60}
}

func (s *AuthService) Login(
	ctx context.Context, in *LoginInput) (*AuthOutput, error) {
	pl := auth.LoginIn{User: in.User, Password: in.Password}
	service := auth.NewLoginService(s.handler.Cache, s.tokenExp)

	response, err := handlers.Handle[auth.LoginIn, auth.LoginOut](ctx, pl, service)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		Token: response.Token,
		Exp:   timestamppb.New(response.Exp),
	}, nil
}

func (s *AuthService) Refresh(
	ctx context.Context, in *TokenInput) (*AuthOutput, error) {
	pl := auth.RefreshTokenIn{Token: in.Token}
	service := auth.NewRefreshTokenService(s.handler.Cache, s.tokenExp)

	response, err := handlers.Handle(ctx, pl, service)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		Token: response.Token,
		Exp:   timestamppb.New(response.Exp),
	}, nil
}

func (s *AuthService) InvalidateAll(
	ctx context.Context, in *SessionsInput) (*empty.Empty, error) {
	pl := auth.InvalidateTokensIn{User: in.User, Token: in.Token}
	service := auth.NewInvalidateTokensService(s.handler.Cache)

	_, err := handlers.Handle(ctx, pl, service)
	if err != nil {
		return nil, err
	}

	return new(empty.Empty), nil
}

func (s *AuthService) IsValid(
	ctx context.Context, in *TokenInput) (*empty.Empty, error) {
	pl := auth.CheckTokenIn{Token: in.Token}
	service := auth.NewCheckTokenService(s.handler.Cache)

	_, err := handlers.Handle(ctx, pl, service)
	if err != nil {
		return nil, err
	}

	return new(empty.Empty), nil
}

func (s *AuthService) ActiveSessions(
	context.Context, *SessionsInput) (*SessionsOutput, error) {
	return nil, ErrNotImplemented
}

func (s *AuthService) Invalidate(
	ctx context.Context, in *TokenInput) (*empty.Empty, error) {
	return new(empty.Empty), ErrNotImplemented
}

package auth

import (
	"time"

	"github.com/ppcamp/go-auth/src/repositories/cache"
	"github.com/ppcamp/go-auth/src/repositories/database"
	"github.com/ppcamp/go-auth/src/services"
	"github.com/ppcamp/go-auth/src/utils/jwt"

	"golang.org/x/crypto/bcrypt"
)

// validatePassword takes a password and a hashedPassword and check if its valid or not
func validatePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type LoginIn struct {
	User     string
	Password string
}

type LoginOut struct {
	Token string
	Exp   time.Time
}

type loginService[In, Out any] struct {
	services.TransactionBusiness[database.AuthStorage]

	Storage   database.UserStorage
	cache     cache.Auth
	signerExp time.Duration
	signer    jwt.Jwt
}

// NewLoginService creates a service that get user password, check it, and
// return a valid JWT token
func NewLoginService(
	repo cache.Auth,
	exp time.Duration,
	signer jwt.Jwt,
) services.ITransactionBusiness[LoginIn, LoginOut] {
	return &loginService[LoginIn, LoginOut]{
		cache:     repo,
		signerExp: exp,
		signer:    signer,
	}
}

func (s *loginService[In, Out]) Execute(in LoginIn) (*LoginOut, error) {
	hash, err := s.Storage.GetUserPassword(in.User)
	if err != nil {
		return nil, err
	}

	if err = validatePassword(in.Password, hash); err != nil {
		return nil, err
	}

	exp := time.Now().Add(s.signerExp)
	token, err := s.signer.Generate(&jwt.Session{}, s.signerExp)
	if err != nil {
		return nil, err
	}

	response := &LoginOut{token, exp}

	if err = s.cache.Authorize(s.Context, in.User, token, exp); err != nil {
		return nil, err
	}

	return response, nil
}

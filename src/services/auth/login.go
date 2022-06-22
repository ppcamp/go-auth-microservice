package auth

import (
	"time"

	"authentication/helpers/services"
	"authentication/repositories/cache"
	"authentication/repositories/database"
	"authentication/utils/jwt"

	"golang.org/x/crypto/bcrypt"
)

func validatePassword(pswd, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pswd))
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
}

// NewLoginService creates a service that get user password, check it, and
// return a valid JWT token
func NewLoginService(
	repo cache.Auth,
	exp time.Duration,
) services.ITransactionBusiness[LoginIn, LoginOut] {
	return &loginService[LoginIn, LoginOut]{
		cache:     repo,
		signerExp: exp,
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
	token, err := jwt.Signer.Generate(&jwt.Session{}, s.signerExp)
	if err != nil {
		return nil, err
	}

	response := &LoginOut{token, exp}

	if err = s.cache.Authorize(s.Context, in.User, token, exp); err != nil {
		return nil, err
	}

	return response, nil
}

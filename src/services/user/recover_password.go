package user

import (
	"time"

	"authentication/helpers/services"
	"authentication/repositories/cache"
	"authentication/repositories/database"

	"github.com/ppcamp/go-lib/random"
)

type RecoverPasswordIn struct {
	Email string
}

type RecoverPasswordOut struct {
	Secret string
}

type recoverPasswordService struct {
	services.TransactionBusiness[database.UserStorage]

	cache cache.UserData
	exp   time.Duration
}

// NewRecoverPasswordService creates a service that generates a token to recover some user
// password
func NewRecoverPasswordService(
	cache cache.UserData,
) services.ITransactionBusiness[RecoverPasswordIn, RecoverPasswordOut] {
	return &recoverPasswordService{cache: cache}
}

func (s *recoverPasswordService) Execute(in RecoverPasswordIn) (*RecoverPasswordOut, error) {
	secret := random.String(30)

	err := s.cache.StoreSecret(s.Context, in.Email, secret, s.exp)
	if err != nil {
		return nil, err
	}

	return &RecoverPasswordOut{secret}, err
}

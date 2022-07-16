package user

import (
	"github.com/ppcamp/go-auth-microservice/src/repositories/cache"
	"github.com/ppcamp/go-auth-microservice/src/repositories/database"
	"github.com/ppcamp/go-auth-microservice/src/services"
)

type ActivateAccountIn struct {
	Secret string
}

type ActivateAccountOut struct{}

type activateAccountService struct {
	services.TransactionBusiness[database.UserStorage]

	cache cache.UserData
}

// NewActivateAccountService creates an UserActivateAccout service, which is responsible to
// validate some secret and turn the user on
func NewActivateAccountService(
	cache cache.UserData,
) services.ITransactionBusiness[ActivateAccountIn, ActivateAccountOut] {
	return &activateAccountService{cache: cache}
}

func (s *activateAccountService) Execute(in ActivateAccountIn) (*ActivateAccountOut, error) {
	user, err := s.cache.UserFromSecret(in.Secret)
	if err != nil {
		return nil, err
	}

	err = s.Storage.ActivateUser(user)
	return new(ActivateAccountOut), err
}

package user

import (
	"github.com/ppcamp/go-auth-microservice/src/repositories/cache"
	"github.com/ppcamp/go-auth-microservice/src/repositories/database"
	"github.com/ppcamp/go-auth-microservice/src/services"
)

type UpdatePasswordIn struct {
	RecoverToken string
	Password     string
}

type UpdatePasswordOut struct{}

type updatePasswordService struct {
	services.TransactionBusiness[database.UserStorage]

	cache cache.UserData
}

// NewUpdatePasswordService creates a service that get user password, check it, and
// return a valid JWT token
func NewUpdatePasswordService(
	cache cache.UserData,
) services.ITransactionBusiness[UpdatePasswordIn, UpdatePasswordOut] {
	return &updatePasswordService{cache: cache}
}

func (s *updatePasswordService) Execute(in UpdatePasswordIn) (*UpdatePasswordOut, error) {
	user, err := s.cache.UserFromSecret(in.RecoverToken)
	if err != nil {
		return nil, err
	}

	err = s.Storage.UpdatePassword(user, in.Password)
	return new(UpdatePasswordOut), err
}

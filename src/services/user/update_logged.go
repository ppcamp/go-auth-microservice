package user

import (
	"authentication/helpers/services"
	"authentication/repositories/cache"
	"authentication/repositories/database"
	"authentication/utils/jwt"
)

type UpdateLoggedIn struct {
	JwtToken string
	Password string
}

type UpdateLoggedOut struct{}

type updateLoggedService struct {
	services.TransactionBusiness[database.UserStorage]

	cache cache.UserData
}

// NewUpdateLoggedService creates a service that get the current logged user, check if is
// valid, and update its password
func NewUpdateLoggedService(
	cache cache.UserData,
) services.ITransactionBusiness[UpdateLoggedIn, UpdateLoggedOut] {
	return &updateLoggedService{cache: cache}
}

func (s *updateLoggedService) Execute(in UpdateLoggedIn) (*UpdateLoggedOut, error) {
	session, err := jwt.Signer.Session(in.JwtToken)
	if err != nil {
		return nil, err
	}

	if err = s.Storage.UpdatePassword(*session.UserId, in.Password); err != nil {
		return nil, err
	}

	return new(UpdateLoggedOut), nil
}

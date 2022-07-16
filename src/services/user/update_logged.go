package user

import (
	"github.com/ppcamp/go-auth-microservice/src/repositories/cache"
	"github.com/ppcamp/go-auth-microservice/src/repositories/database"
	"github.com/ppcamp/go-auth-microservice/src/services"
	"github.com/ppcamp/go-auth-microservice/src/utils/jwt"
)

type UpdateLoggedIn struct {
	JwtToken string
	Password string
}

type UpdateLoggedOut struct{}

type updateLoggedService struct {
	services.TransactionBusiness[database.UserStorage]

	cache  cache.UserData
	signer jwt.Jwt
}

// NewUpdateLoggedService creates a service that get the current logged user, check if is
// valid, and update its password
func NewUpdateLoggedService(
	cache cache.UserData,
	signer jwt.Jwt,
) services.ITransactionBusiness[UpdateLoggedIn, UpdateLoggedOut] {
	return &updateLoggedService{cache: cache, signer: signer}
}

func (s *updateLoggedService) Execute(in UpdateLoggedIn) (*UpdateLoggedOut, error) {
	session, err := s.signer.Session(in.JwtToken)
	if err != nil {
		return nil, err
	}

	if err = s.Storage.UpdatePassword(*session.UserId, in.Password); err != nil {
		return nil, err
	}

	return new(UpdateLoggedOut), nil
}

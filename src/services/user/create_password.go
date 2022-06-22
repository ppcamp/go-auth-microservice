package user

import (
	"database/sql"
	"time"

	"authentication/helpers/services"
	"authentication/repositories/cache"
	"authentication/repositories/database"
	"authentication/utils/password"

	"errors"

	"github.com/ppcamp/go-lib/random"
	"github.com/sirupsen/logrus"
)

var (
	ErrUserAlreadyExist error = errors.New("user already exist and it's active")
)

type CreatePasswordIn struct {
	User     string
	Password string
}

type CreatePasswordOut struct {
	ActivateToken string
}

type createPasswordService struct {
	services.TransactionBusiness[database.UserStorage]

	cache cache.UserData
	exp   time.Duration
}

// NewCreatePasswordService creates a service that creates a new login
func NewCreatePasswordService(
	cache cache.UserData,
) services.ITransactionBusiness[CreatePasswordIn, CreatePasswordOut] {
	return &createPasswordService{cache: cache, exp: time.Hour * 24}
}

func (s *createPasswordService) Execute(in CreatePasswordIn) (*CreatePasswordOut, error) {
	logrus.Info("getting password")
	_, err := s.Storage.GetUserPassword(in.User)
	if err == nil {
		return nil, ErrUserAlreadyExist
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	logrus.Info("hashing password")
	hashedPassword, err := password.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	logrus.Info("storing password at database")
	err = s.Storage.CreateUserPassword(in.User, hashedPassword)
	if err != nil {
		return nil, err
	}

	unlockSecret := random.String(30)

	logrus.Info("storing password at cache")
	err = s.cache.StoreSecret(s.Context, in.User, unlockSecret, s.exp)
	if err != nil {
		return nil, err
	}

	logrus.Info("returning")
	return &CreatePasswordOut{unlockSecret}, nil
}

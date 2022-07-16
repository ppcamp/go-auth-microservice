package user

import (
	"database/sql"
	"time"

	"github.com/ppcamp/go-auth-microservice/src/repositories/cache"
	"github.com/ppcamp/go-auth-microservice/src/repositories/database"
	"github.com/ppcamp/go-auth-microservice/src/services"

	"errors"

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
	hashedPassword, err := hashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	logrus.Info("storing password at database")
	err = s.Storage.CreateUserPassword(in.User, hashedPassword)
	if err != nil {
		return nil, err
	}

	unlockSecret := newSecret()

	logrus.Info("storing password at cache")
	err = s.cache.StoreSecret(s.Context, in.User, unlockSecret, s.exp)
	if err != nil {
		return nil, err
	}

	logrus.Info("returning")
	return &CreatePasswordOut{unlockSecret}, nil
}

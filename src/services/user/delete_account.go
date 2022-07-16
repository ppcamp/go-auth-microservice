package user

import (
	"github.com/ppcamp/go-auth-microservice/src/repositories/database"
	"github.com/ppcamp/go-auth-microservice/src/services"
)

type DeleteAccountIn struct {
	Token string
}

type DeleteAccountOut struct{}

type deleteAccountService struct {
	services.TransactionBusiness[database.UserStorage]
}

func NewDeleteAccountService() services.ITransactionBusiness[DeleteAccountIn, DeleteAccountOut] {
	return new(deleteAccountService)
}

func (s *deleteAccountService) Execute(in DeleteAccountIn) (*DeleteAccountOut, error) {
	return new(DeleteAccountOut), nil
}

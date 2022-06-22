package user

import (
	"authentication/helpers/services"
	"authentication/repositories/database"
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

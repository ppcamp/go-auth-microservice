package user_password

import (
	context "context"

	"authentication/helpers/handlers"
	"authentication/services/user"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
)

type UserPasswordService struct {
	UnsafeUserPasswordServiceServer

	handler *handlers.Handler
}

func NewUserPasswordService(
	handler *handlers.Handler) UserPasswordServiceServer {
	return &UserPasswordService{handler: handler}
}

func (u *UserPasswordService) Create(
	ctx context.Context, in *CreateInput) (*CreateOutput, error) {
	logrus.WithField("in", in).Info("calling create")

	pl := user.CreatePasswordIn{User: in.GetUser(), Password: in.GetPassword()}

	service := user.NewCreatePasswordService(u.handler.Cache)
	response, err := handlers.
		Handle[user.CreatePasswordIn, user.CreatePasswordOut](ctx, pl, service)
	if err != nil {
		return nil, err
	}

	logrus.Info("returning create")

	return &CreateOutput{Secret: response.ActivateToken}, nil
}

func (u *UserPasswordService) Activate(
	ctx context.Context, in *ActivateInput) (*empty.Empty, error) {
	pl := user.ActivateAccountIn{Secret: in.GetSecret()}
	service := user.NewActivateAccountService(u.handler.Cache)

	_, err := handlers.
		Handle[user.ActivateAccountIn, user.ActivateAccountOut](ctx, pl, service)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

func (u *UserPasswordService) Recover(
	ctx context.Context, in *RecoverInput) (*RecoverOutput, error) {
	pl := user.RecoverPasswordIn{Email: in.GetEmail()}
	service := user.NewRecoverPasswordService(u.handler.Cache)

	response, err := handlers.
		Handle[user.RecoverPasswordIn, user.RecoverPasswordOut](ctx, pl, service)
	if err != nil {
		return nil, err
	}
	return &RecoverOutput{Secret: response.Secret}, nil
}

func (u *UserPasswordService) Update(
	ctx context.Context, in *UpdateInput) (*empty.Empty, error) {
	pl := user.UpdatePasswordIn{
		RecoverToken: in.GetSecret(),
		Password:     in.GetPassword(),
	}
	service := user.NewUpdatePasswordService(u.handler.Cache)

	_, err := handlers.
		Handle[user.UpdatePasswordIn, user.UpdatePasswordOut](ctx, pl, service)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

func (u *UserPasswordService) UpdateByToken(
	ctx context.Context, in *UpdateInput) (*empty.Empty, error) {
	pl := user.UpdateLoggedIn{
		JwtToken: in.GetSecret(),
		Password: in.GetPassword(),
	}
	service := user.NewUpdateLoggedService(u.handler.Cache)

	_, err := handlers.
		Handle[user.UpdateLoggedIn, user.UpdateLoggedOut](ctx, pl, service)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

func (u *UserPasswordService) Delete(
	ctx context.Context, in *DeleteInput) (*empty.Empty, error) {
	pl := user.DeleteAccountIn{Token: in.GetToken()}
	service := user.NewDeleteAccountService()

	_, err := handlers.
		Handle[user.DeleteAccountIn, user.DeleteAccountOut](ctx, pl, service)
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

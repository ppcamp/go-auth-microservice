package user_password

import (
	context "context"

	handlers "github.com/ppcamp/go-auth/src/http"
	"github.com/ppcamp/go-auth/src/services/user"

	empty "github.com/golang/protobuf/ptypes/empty"
)

type server struct {
	UnsafeUserPasswordServiceServer
	*handlers.Handler
}

func NewUserPasswordService(handler *handlers.Handler) UserPasswordServiceServer {
	return &server{Handler: handler}
}

func (u *server) Create(ctx context.Context, pl *CreateInput) (*CreateOutput, error) {
	input := user.CreatePasswordIn{User: pl.GetUser(), Password: pl.GetPassword()}

	service := user.NewCreatePasswordService(u.Cache)

	response, err := handlers.Handle[
		user.CreatePasswordIn,
		user.CreatePasswordOut](ctx, u.Handler, input, service)

	if err != nil {
		return nil, err
	}

	return &CreateOutput{Secret: response.ActivateToken}, nil
}

func (u *server) Activate(ctx context.Context, pl *ActivateInput) (*empty.Empty, error) {
	input := user.ActivateAccountIn{Secret: pl.GetSecret()}

	service := user.NewActivateAccountService(u.Cache)

	_, err := handlers.Handle[
		user.ActivateAccountIn,
		user.ActivateAccountOut](ctx, u.Handler, input, service)

	if err != nil {
		return nil, err
	}

	return new(empty.Empty), nil
}

func (u *server) Recover(ctx context.Context, pl *RecoverInput) (*RecoverOutput, error) {
	input := user.RecoverPasswordIn{Email: pl.GetEmail()}

	service := user.NewRecoverPasswordService(u.Cache)

	response, err := handlers.Handle[
		user.RecoverPasswordIn,
		user.RecoverPasswordOut](ctx, u.Handler, input, service)

	if err != nil {
		return nil, err
	}

	return &RecoverOutput{Secret: response.Secret}, nil
}

func (u *server) Update(ctx context.Context, pl *UpdateInput) (*empty.Empty, error) {
	input := user.UpdatePasswordIn{RecoverToken: pl.GetSecret(), Password: pl.GetPassword()}

	service := user.NewUpdatePasswordService(u.Cache)

	_, err := handlers.Handle[
		user.UpdatePasswordIn,
		user.UpdatePasswordOut](ctx, u.Handler, input, service)

	if err != nil {
		return nil, err
	}

	return new(empty.Empty), nil
}

func (u *server) UpdateByToken(ctx context.Context, pl *UpdateInput) (*empty.Empty, error) {
	input := user.UpdateLoggedIn{JwtToken: pl.GetSecret(), Password: pl.GetPassword()}

	service := user.NewUpdateLoggedService(u.Cache, u.Signer)

	_, err := handlers.Handle[
		user.UpdateLoggedIn,
		user.UpdateLoggedOut](ctx, u.Handler, input, service)

	if err != nil {
		return nil, err
	}

	return new(empty.Empty), nil
}

func (u *server) Delete(ctx context.Context, pl *DeleteInput) (*empty.Empty, error) {
	input := user.DeleteAccountIn{Token: pl.GetToken()}

	service := user.NewDeleteAccountService()

	_, err := handlers.Handle[
		user.DeleteAccountIn,
		user.DeleteAccountOut](ctx, u.Handler, input, service)

	if err != nil {
		return nil, err
	}

	return new(empty.Empty), nil
}

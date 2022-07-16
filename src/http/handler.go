package http

import (
	"context"

	"github.com/ppcamp/go-auth-microservice/src/repositories/cache"
	"github.com/ppcamp/go-auth-microservice/src/repositories/database"
	"github.com/ppcamp/go-auth-microservice/src/services"
	"github.com/ppcamp/go-auth-microservice/src/utils/jwt"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	Cache    cache.Cache
	Database database.Connection
	Signer   jwt.Jwt
}

// Handle is responsible to identify the service type and assign the correct
// workflow to it. Furthermore, it's also responsible to commit/rollback
// transactions when needed and set another responses parameters if needed.
func Handle[In, Out any](
	ctx context.Context,
	handler *Handler,
	input In,
	service services.IBaseBusiness[In, Out],
) (*Out, error) {
	service.SetContext(ctx)

	switch service := service.(type) {
	case services.ITransactionBusiness[In, Out]:
		return handleTransactionService(handler, input, service)
	default:
		return handleBaseService(handler, input, service)
	}
}

func handleBaseService[In, Out any](
	h *Handler,
	input In,
	service services.IBaseBusiness[In, Out],
) (response *Out, err error) {
	return service.Execute(input)
}

func handleTransactionService[In, Out any](
	h *Handler,
	input In,
	service services.ITransactionBusiness[In, Out],
) (response *Out, err error) {

	tr, err := h.Database.StartTransaction()
	if err != nil {
		return nil, err
	}

	defer func() {
		panicked := recover()

		if err != nil || panicked != nil {
			logrus.WithFields(logrus.Fields{
				"panic": panicked,
				"err":   err,
			}).Warn("err or panicked")
			if err := tr.Rollback(); err != nil {
				logrus.WithFields(logrus.Fields{"input": input}).
					Error("fail to rollback transaction")
			}

			return
		}

		if err := tr.Commit(); err != nil {
			logrus.WithFields(logrus.Fields{"input": input}).
				Error("fail to commit transaction")
		}
	}()

	service.SetTransaction(tr)

	response, err = service.Execute(input)
	return response, err
}

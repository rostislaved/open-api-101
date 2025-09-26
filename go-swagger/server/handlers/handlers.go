package handlers

import (
	"context"
	"errors"

	"github.com/go-openapi/runtime/middleware"

	"server/generated/models"
	"server/generated/restapi/operations"
	"server/usecases"
)

type Handlers struct {
	useCases UseCases
}

type UseCases interface {
	GetUser(ctx context.Context, id int) (usecases.User, error)
	CreateUsers(ctx context.Context, userRequests usecases.CreateUserRequestDTO) (int, error)
}

func New(useCases UseCases) *Handlers {
	return &Handlers{
		useCases: useCases,
	}
}

func (h *Handlers) GetUsers(params operations.GetUserByIDParams) middleware.Responder {
	user, err := h.useCases.GetUser(params.HTTPRequest.Context(), int(params.ID))
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrNotFound):
			resp := operations.
				NewGetUserByIDNotFound().
				WithPayload(
					&models.ErrorResponse{
						Code:  ToPtr(int64(404)),
						Error: ToPtr("Not Found"),
					},
				)

			return resp
		case errors.Is(err, usecases.ErrNotPublic1):
			resp := operations.
				NewGetUserByIDInternalServerError().
				WithPayload(
					&models.ErrorResponse{
						Code:  ToPtr(int64(1)),
						Error: ToPtr("Internal Server Error 1"),
					},
				)

			return resp
		case errors.Is(err, usecases.ErrNotPublic2):
			resp := operations.
				NewGetUserByIDInternalServerError().
				WithPayload(
					&models.ErrorResponse{
						Code:  ToPtr(int64(2)),
						Error: ToPtr("Internal Server Error 2"),
					},
				)

			return resp
		default:
			resp := operations.
				NewGetUserByIDInternalServerError().
				WithPayload(
					&models.ErrorResponse{
						Code:  ToPtr(int64(-1)),
						Error: ToPtr("Internal Server Error"),
					},
				)

			return resp
		}
	}

	resp := operations.
		NewGetUserByIDOK().
		WithPayload(
			&models.GetUserByIDResponse{
				ID:   ToPtr(int64(user.ID)),
				Name: ToPtr(user.Name),
			},
		)

	return resp
}

func (h *Handlers) CreateUsers(params operations.CreateUserParams) middleware.Responder {
	createUserRequestDTO := usecases.CreateUserRequestDTO{
		Name: *params.Body.Name,
	}

	id, err := h.useCases.CreateUsers(params.HTTPRequest.Context(), createUserRequestDTO)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrValidation):
			resp := operations.
				NewCreateUserBadRequest().
				WithPayload(
					&models.ErrorResponse{
						Code:  ToPtr(int64(3)),
						Error: ToPtr(err.Error()),
					},
				)

			return resp
		default:
			resp := operations.
				NewCreateUserInternalServerError().
				WithPayload(
					&models.ErrorResponse{
						Code:  ToPtr(int64(-1)),
						Error: ToPtr("Internal Server Error"),
					},
				)

			return resp
		}
	}

	resp := operations.
		NewCreateUserCreated().
		WithPayload(
			&models.CreateUserResponse{
				ID: ToPtr(int64(id)),
			},
		)

	return resp
}

func ToPtr[T any](v T) *T {
	return &v
}

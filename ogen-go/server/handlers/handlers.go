package handlers

import (
	"context"
	"errors"

	api "server/generated"
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

func (h *Handlers) GetUserById(ctx context.Context, params api.GetUserByIdParams) (api.GetUserByIdRes, error) {
	id := params.ID

	user, err := h.useCases.GetUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrNotFound):
			response := api.GetUserByIdNotFound{
				Code:  404,
				Error: "Not Found",
			}

			return &response, nil
		case errors.Is(err, usecases.ErrNotPublic1):
			response := api.GetUserByIdInternalServerError{
				Code:  1,
				Error: "Internal Server Error 1",
			}

			return &response, nil
		case errors.Is(err, usecases.ErrNotPublic2):
			response := api.GetUserByIdInternalServerError{
				Code:  2,
				Error: "Internal Server Error 2",
			}

			return &response, nil
		default:
			response := api.GetUserByIdInternalServerError{
				Code:  -1,
				Error: "Internal Server Error",
			}

			return &response, nil
		}
	}

	response := api.GetUserByIdResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	return &response, nil

}

func (h *Handlers) CreateUser(ctx context.Context, req *api.CreateUserRequest) (api.CreateUserRes, error) {
	createUserRequestDTO := usecases.CreateUserRequestDTO{
		Name: req.Name,
	}

	id, err := h.useCases.CreateUsers(ctx, createUserRequestDTO)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrValidation):
			response := api.CreateUserBadRequest{
				Code:  3,
				Error: err.Error(),
			}

			return &response, nil
		default:
			response := api.CreateUserInternalServerError{
				Code:  -1,
				Error: "Internal Server Error",
			}

			return &response, nil
		}
	}

	response := api.CreateUserResponse{
		ID: id,
	}

	return &response, nil
}

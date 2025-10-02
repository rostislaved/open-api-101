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

func (h *Handlers) GetUserById(ctx context.Context, request api.GetUserByIdRequestObject) (api.GetUserByIdResponseObject, error) {
	id := request.Id

	user, err := h.useCases.GetUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrNotFound):
			response := api.ErrorResponse{
				Code:  404,
				Error: "Not Found",
			}

			return api.GetUserById404JSONResponse(response), nil
		case errors.Is(err, usecases.ErrNotPublic1):
			response := api.ErrorResponse{
				Code:  1,
				Error: "Internal Server Error 1",
			}

			return api.GetUserById500JSONResponse(response), nil
		case errors.Is(err, usecases.ErrNotPublic2):
			response := api.ErrorResponse{
				Code:  2,
				Error: "Internal Server Error 2",
			}

			return api.GetUserById500JSONResponse(response), nil
		default:
			response := api.ErrorResponse{
				Code:  -1,
				Error: "Internal Server Error",
			}

			return api.GetUserById500JSONResponse(response), nil
		}
	}

	response := api.GetUserByIdResponse{
		Id:   user.ID,
		Name: user.Name,
	}

	return api.GetUserById200JSONResponse(response), nil

}

func (h *Handlers) CreateUser(ctx context.Context, request api.CreateUserRequestObject) (api.CreateUserResponseObject, error) {
	createUserRequestDTO := usecases.CreateUserRequestDTO{
		Name: request.Body.Name,
	}

	id, err := h.useCases.CreateUsers(ctx, createUserRequestDTO)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrValidation):
			response := api.ErrorResponse{
				Code:  3,
				Error: err.Error(),
			}

			return api.CreateUser400JSONResponse(response), nil
		default:
			response := api.ErrorResponse{
				Code:  -1,
				Error: "Internal Server Error",
			}

			return api.CreateUser500JSONResponse(response), nil
		}
	}

	response := api.CreateUserResponse{
		Id: id,
	}

	return api.CreateUser201JSONResponse(response), nil
}

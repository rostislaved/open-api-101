package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

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

func (h *Handlers) GetUserById(w http.ResponseWriter, r *http.Request, id int) {
	user, err := h.useCases.GetUser(r.Context(), id)
	if err != nil {
		var response api.ErrorResponse

		switch {
		case errors.Is(err, usecases.ErrNotFound):
			response = api.ErrorResponse{
				Code:  404,
				Error: "Not Found",
			}
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, usecases.ErrNotPublic1):
			response = api.ErrorResponse{
				Code:  1,
				Error: "Internal Server Error 1",
			}
			w.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, usecases.ErrNotPublic2):
			response = api.ErrorResponse{
				Code:  2,
				Error: "Internal Server Error 2",
			}
			w.WriteHeader(http.StatusInternalServerError)
		default:
			response = api.ErrorResponse{
				Code:  -1,
				Error: "Internal Server Error",
			}
			w.WriteHeader(http.StatusInternalServerError)
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(responseBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	response := api.GetUserByIdResponse{
		Id:   user.ID,
		Name: user.Name,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request api.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	createUserRequestDTO := usecases.CreateUserRequestDTO{
		Name: request.Name,
	}

	id, err := h.useCases.CreateUsers(r.Context(), createUserRequestDTO)
	if err != nil {
		var response api.ErrorResponse

		switch {
		case errors.Is(err, usecases.ErrValidation):
			response = api.ErrorResponse{
				Code:  3,
				Error: err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)
		default:
			response = api.ErrorResponse{
				Code:  -1,
				Error: "Internal Server Error",
			}
			w.WriteHeader(http.StatusInternalServerError)
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(responseBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	response := api.CreateUserResponse{
		Id: id,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

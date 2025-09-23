package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	api "server/generated"
	"server/usecases"
)

type Handlers struct {
	useCases UseCases
}

type UseCases interface {
	GetUser(ctx context.Context, id int) usecases.User
	CreateUsers(ctx context.Context, userRequests usecases.CreateUserRequestDTO) int
}

func New(useCases UseCases) *Handlers {
	return &Handlers{
		useCases: useCases,
	}
}

func (h *Handlers) GetUserById(w http.ResponseWriter, r *http.Request, id int) {
	fmt.Println(id)

	user := h.useCases.GetUser(r.Context(), id)

	response := api.GetUserByIdResponse{
		Id:   &user.ID,
		Name: &user.Name,
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

	fmt.Println(request)

	createUserRequestDTO := usecases.CreateUserRequestDTO{
		Name: request.Name,
	}

	id := h.useCases.CreateUsers(r.Context(), createUserRequestDTO)

	response := api.CreateUserResponse{
		Id: &id,
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

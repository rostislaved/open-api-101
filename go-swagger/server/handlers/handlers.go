package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"server/generated/models"
	"server/generated/restapi/operations"
)

type Handlers struct {
}

func New() Handlers {
	return Handlers{}
}

func (h *Handlers) GetUsers(params operations.GetUsersParams) middleware.Responder {
	list := []*models.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	getUsersResponse := operations.NewGetUsersOK()
	getUsersResponse.Payload = list

	return getUsersResponse
}

func (h *Handlers) CreateUsers(params operations.CreateUsersParams) middleware.Responder {
	user := &models.User{
		// ID:   3,
		Name: *params.Body.Name,
	}

	resp := operations.NewCreateUsersCreated()
	resp.Payload = user

	return resp
}

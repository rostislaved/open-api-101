package usecases

import (
	"context"
)

type UseCases struct {
}

func New() *UseCases {
	return &UseCases{}
}

func (u *UseCases) GetUser(ctx context.Context, id int) User {
	user := User{
		ID:   id,
		Name: "Alice",
	}

	return user
}

type User struct {
	ID   int
	Name string
}

type CreateUserRequestDTO struct {
	ID   int
	Name string
}

func (u *UseCases) CreateUsers(ctx context.Context, userRequests CreateUserRequestDTO) int {
	return 1
}

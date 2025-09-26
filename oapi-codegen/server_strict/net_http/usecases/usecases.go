package usecases

import (
	"context"
	"errors"
)

type UseCases struct {
}

func New() *UseCases {
	return &UseCases{}
}

var (
	ErrNotFound   = errors.New("not found")
	ErrNotPublic1 = errors.New("we can't expose this text 1")
	ErrNotPublic2 = errors.New("we can't expose this text 2")
	ErrUnknown    = errors.New("unknown error")
	ErrValidation = errors.New("validation error")
)

func (u *UseCases) GetUser(ctx context.Context, id int) (User, error) {
	switch id {
	case 1:
		user := User{
			ID:   id,
			Name: "Alice",
		}

		return user, nil
	case 2:
		return User{}, ErrNotFound
	case 3:
		return User{}, ErrNotPublic1
	case 4:
		return User{}, ErrNotPublic2
	default:
		return User{}, ErrUnknown
	}
}

type User struct {
	ID   int
	Name string
}

type CreateUserRequestDTO struct {
	Name string
}

func (u *UseCases) CreateUsers(ctx context.Context, createUserRequestDTO CreateUserRequestDTO) (int, error) {
	switch createUserRequestDTO.Name {
	case "Alice":
		return 10, nil
	case "":
		return 0, ErrValidation
	default:
		return 0, ErrUnknown
	}
}

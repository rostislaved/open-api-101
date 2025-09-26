package handlers

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	api "server/generated"
	"server/usecases"
)

func TestHandlers_GetUserById(t *testing.T) {
	type fields struct {
		setup func(t *testing.T) UseCases
	}
	type args struct {
		ctx     context.Context
		request api.GetUserByIdParams
	}
	
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.GetUserByIdRes
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						GetUser(mock.Anything, 1).
						Return(
							usecases.User{
								ID:   1,
								Name: "Alice",
							},
							nil,
						).
						Once()

					return m
				},
			},
			args: args{
				ctx:     context.Background(),
				request: api.GetUserByIdParams{ID: 1},
			},
			want: &api.GetUserByIdResponse{
				ID:   1,
				Name: "Alice",
			},
			wantErr: false,
		},
		{
			name: "not found",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						GetUser(mock.Anything, 2).
						Return(
							usecases.User{},
							usecases.ErrNotFound,
						).
						Once()

					return m
				},
			},
			args: args{
				ctx:     context.Background(),
				request: api.GetUserByIdParams{ID: 2},
			},
			want: &api.GetUserByIdNotFound{
				Code:  404,
				Error: "Not Found",
			},
			wantErr: false,
		},
		{
			name: "internal error 1",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)
					m.EXPECT().
						GetUser(mock.Anything, 3).
						Return(
							usecases.User{},
							usecases.ErrNotPublic1,
						).
						Once()

					return m
				},
			},
			args: args{
				ctx:     context.Background(),
				request: api.GetUserByIdParams{ID: 3},
			},
			want: &api.GetUserByIdInternalServerError{
				Code:  1,
				Error: "Internal Server Error 1",
			},
			wantErr: false,
		},
		{
			name: "internal error 2",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						GetUser(mock.Anything, 4).
						Return(
							usecases.User{},
							usecases.ErrNotPublic2,
						).
						Once()

					return m
				},
			},
			args: args{
				ctx:     context.Background(),
				request: api.GetUserByIdParams{ID: 4},
			},
			want: &api.GetUserByIdInternalServerError{
				Code:  2,
				Error: "Internal Server Error 2",
			},
			wantErr: false,
		},
		{
			name: "internal server error",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						GetUser(mock.Anything, 5).
						Return(
							usecases.User{},
							usecases.ErrValidation,
						).
						Once()

					return m
				},
			},
			args: args{
				ctx:     context.Background(),
				request: api.GetUserByIdParams{ID: 5},
			},
			want: &api.GetUserByIdInternalServerError{
				Code:  -1,
				Error: "Internal Server Error",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				useCases: tt.fields.setup(t),
			}

			got, err := h.GetUserById(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlers_CreateUser(t *testing.T) {
	type fields struct {
		setup func(t *testing.T) UseCases
	}
	type args struct {
		ctx     context.Context
		request *api.CreateUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.CreateUserRes
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						CreateUsers(
							mock.Anything,
							usecases.CreateUserRequestDTO{Name: "Alice"},
						).
						Return(10, nil).
						Once()

					return m
				},
			},
			args: args{
				ctx:     context.Background(),
				request: &api.CreateUserRequest{Name: "Alice"},
			},
			want: &api.CreateUserResponse{
				ID: 10,
			},
			wantErr: false,
		},
		{
			name: "validation error",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						CreateUsers(
							mock.Anything,
							usecases.CreateUserRequestDTO{Name: ""},
						).
						Return(0, usecases.ErrValidation).
						Once()

					return m
				},
			},
			args: args{
				ctx:     context.Background(),
				request: &api.CreateUserRequest{Name: ""},
			},
			want: &api.CreateUserBadRequest{
				Code:  3,
				Error: usecases.ErrValidation.Error(),
			},
			wantErr: false,
		},
		{
			name: "internal server error",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						CreateUsers(
							mock.Anything,
							usecases.CreateUserRequestDTO{Name: "Bob"},
						).
						Return(0, usecases.ErrNotPublic1).
						Once()

					return m
				},
			},
			args: args{
				ctx:     context.Background(),
				request: &api.CreateUserRequest{Name: "Bob"},
			},
			want: &api.CreateUserInternalServerError{
				Code:  -1,
				Error: "Internal Server Error",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				useCases: tt.fields.setup(t),
			}

			got, err := h.CreateUser(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

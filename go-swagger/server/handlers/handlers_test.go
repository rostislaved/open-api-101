package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/go-openapi/runtime"

	"server/generated/models"
	"server/generated/restapi/operations"
	"server/usecases"
)

func readJSONBody[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()

	var zero T

	bodyBytes := rr.Body.Bytes()
	if len(bodyBytes) == 0 {
		return zero
	}

	var got T

	err := json.Unmarshal(bodyBytes, &got)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v; raw: %q", err, string(bodyBytes))
	}

	return got
}

func TestHandlers_GetUsers(t *testing.T) {
	type fields struct {
		setup func(t *testing.T) UseCases
	}
	type args struct {
		id int64
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
		wantCT         string
		wantBody       any
	}{
		{
			name: "ok",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						GetUser(mock.Anything, 1).
						Return(usecases.User{ID: 1, Name: "Alice"}, nil).
						Once()

					return m
				},
			},
			args:           args{id: 1},
			wantStatusCode: http.StatusOK,
			wantCT:         runtime.JSONMime,
			wantBody: &models.GetUserByIDResponse{
				ID:   ToPtr(int64(1)),
				Name: ToPtr("Alice"),
			},
		},
		{
			name: "not found",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						GetUser(mock.Anything, 2).
						Return(usecases.User{}, usecases.ErrNotFound).
						Once()

					return m
				},
			},
			args:           args{id: 2},
			wantStatusCode: http.StatusNotFound,
			wantCT:         runtime.JSONMime,
			wantBody: &models.ErrorResponse{
				Code:  ToPtr(int64(404)),
				Error: ToPtr("Not Found"),
			},
		},
		{
			name: "internal error 1",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						GetUser(mock.Anything, 3).
						Return(usecases.User{}, usecases.ErrNotPublic1).
						Once()

					return m
				},
			},
			args:           args{id: 3},
			wantStatusCode: http.StatusInternalServerError,
			wantCT:         runtime.JSONMime,
			wantBody: &models.ErrorResponse{
				Code:  ToPtr(int64(1)),
				Error: ToPtr("Internal Server Error 1"),
			},
		},
		{
			name: "internal error 2",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						GetUser(mock.Anything, 4).
						Return(usecases.User{}, usecases.ErrNotPublic2).
						Once()

					return m
				},
			},
			args:           args{id: 4},
			wantStatusCode: http.StatusInternalServerError,
			wantCT:         runtime.JSONMime,
			wantBody: &models.ErrorResponse{
				Code:  ToPtr(int64(2)),
				Error: ToPtr("Internal Server Error 2"),
			},
		},
		{
			name: "internal server error (default)",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						GetUser(mock.Anything, 5).
						Return(usecases.User{}, errors.New("unexpected")).
						Once()

					return m
				},
			},
			args:           args{id: 5},
			wantStatusCode: http.StatusInternalServerError,
			wantCT:         runtime.JSONMime,
			wantBody: &models.ErrorResponse{
				Code:  ToPtr(int64(-1)),
				Error: ToPtr("Internal Server Error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				useCases: tt.fields.setup(t),
			}

			req := httptest.NewRequest(http.MethodGet, "/users/", nil)
			req = req.WithContext(context.Background())

			params := operations.GetUserByIDParams{
				HTTPRequest: req,
				ID:          tt.args.id,
			}

			responder := h.GetUsers(params)

			rr := httptest.NewRecorder()

			responder.WriteResponse(rr, runtime.JSONProducer())

			if rr.Code != tt.wantStatusCode {
				t.Fatalf("status code = %d, want %d; body: %s", rr.Code, tt.wantStatusCode, rr.Body.String())
			}

			switch want := tt.wantBody.(type) {
			case *models.GetUserByIDResponse:
				got := readJSONBody[models.GetUserByIDResponse](t, rr)
				if !reflect.DeepEqual(&got, want) {
					t.Fatalf("body = %#v, want %#v", got, *want)
				}
			case *models.ErrorResponse:
				got := readJSONBody[models.ErrorResponse](t, rr)
				if !reflect.DeepEqual(&got, want) {
					t.Fatalf("body = %#v, want %#v", got, *want)
				}
			default:
				t.Fatalf("unsupported wantBody type: %T", tt.wantBody)
			}
		})
	}
}

// ---------- CreateUsers ----------

func TestHandlers_CreateUsers(t *testing.T) {
	type fields struct {
		setup func(t *testing.T) UseCases
	}
	type args struct {
		name string
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
		wantCT         string
		wantBody       any
	}{
		{
			name: "created 201",
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
			args:           args{name: "Alice"},
			wantStatusCode: http.StatusCreated,
			wantCT:         runtime.JSONMime,
			wantBody: &models.CreateUserResponse{
				ID: ToPtr(int64(10)),
			},
		},
		{
			name: "validation error -> 400 (code=3, message from error)",
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
			args:           args{name: ""},
			wantStatusCode: http.StatusBadRequest,
			wantCT:         runtime.JSONMime,
			wantBody: &models.ErrorResponse{
				Code:  ToPtr(int64(3)),
				Error: ToPtr(usecases.ErrValidation.Error()),
			},
		},
		{
			name: "internal server error (default) -> 500",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						CreateUsers(
							mock.Anything,
							usecases.CreateUserRequestDTO{Name: "Bob"},
						).
						Return(0, errors.New("unexpected")).
						Once()

					return m
				},
			},
			args:           args{name: "Bob"},
			wantStatusCode: http.StatusInternalServerError,
			wantCT:         runtime.JSONMime,
			wantBody: &models.ErrorResponse{
				Code:  ToPtr(int64(-1)),
				Error: ToPtr("Internal Server Error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				useCases: tt.fields.setup(t),
			}

			body := &models.CreateUserRequest{
				Name: ToPtr(tt.args.name),
			}

			req := httptest.NewRequest(http.MethodPost, "/users", nil)
			req = req.WithContext(context.Background())

			params := operations.CreateUserParams{
				HTTPRequest: req,
				Body:        body,
			}

			responder := h.CreateUsers(params)

			rr := httptest.NewRecorder()

			responder.WriteResponse(rr, runtime.JSONProducer())

			if rr.Code != tt.wantStatusCode {
				t.Fatalf("status code = %d, want %d; body: %s", rr.Code, tt.wantStatusCode, rr.Body.String())
			}

			switch want := tt.wantBody.(type) {
			case *models.CreateUserResponse:
				got := readJSONBody[models.CreateUserResponse](t, rr)
				if !reflect.DeepEqual(&got, want) {
					t.Fatalf("body = %#v, want %#v", got, *want)
				}
			case *models.ErrorResponse:
				got := readJSONBody[models.ErrorResponse](t, rr)
				if !reflect.DeepEqual(&got, want) {
					t.Fatalf("body = %#v, want %#v", got, *want)
				}
			default:
				t.Fatalf("unsupported wantBody type: %T", tt.wantBody)
			}
		})
	}
}

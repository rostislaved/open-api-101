package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	api "server/generated"
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

func TestHTTPHandlers_GetUserById(t *testing.T) {
	type fields struct {
		setup func(t *testing.T) UseCases
	}
	type args struct {
		id int
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
			name: "happy path",
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
			wantCT:         "application/json",
			wantBody: api.GetUserByIdResponse{
				Id:   1,
				Name: "Alice",
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
			wantCT:         "application/json",
			wantBody: api.ErrorResponse{
				Code:  404,
				Error: "Not Found",
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
			wantCT:         "application/json",
			wantBody: api.ErrorResponse{
				Code:  1,
				Error: "Internal Server Error 1",
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
			wantCT:         "application/json",
			wantBody: api.ErrorResponse{
				Code:  2,
				Error: "Internal Server Error 2",
			},
		},
		{
			name: "internal server error (default)",
			fields: fields{
				setup: func(t *testing.T) UseCases {
					m := NewMockUseCases(t)

					m.EXPECT().
						GetUser(mock.Anything, 5).
						Return(usecases.User{}, usecases.ErrValidation).
						Once()

					return m
				},
			},
			args:           args{id: 5},
			wantStatusCode: http.StatusInternalServerError,
			wantCT:         "application/json",
			wantBody: api.ErrorResponse{
				Code:  -1,
				Error: "Internal Server Error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				useCases: tt.fields.setup(t),
			}

			req := httptest.NewRequest(http.MethodGet, "/users/", nil)
			rr := httptest.NewRecorder()

			h.GetUserById(rr, req, tt.args.id)

			if rr.Code != tt.wantStatusCode {
				t.Fatalf("status code = %d, want %d; body: %s", rr.Code, tt.wantStatusCode, rr.Body.String())
			}

			ct := rr.Header().Get("Content-Type")
			if ct != tt.wantCT {
				t.Fatalf("content-type = %q, want %q", ct, tt.wantCT)
			}

			switch want := tt.wantBody.(type) {
			case api.GetUserByIdResponse:
				got := readJSONBody[api.GetUserByIdResponse](t, rr)
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("body = %+v, want %+v", got, want)
				}
			case api.ErrorResponse:
				got := readJSONBody[api.ErrorResponse](t, rr)
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("body = %+v, want %+v", got, want)
				}
			default:
				t.Fatalf("unsupported wantBody type: %T", tt.wantBody)
			}
		})
	}
}

func TestHTTPHandlers_CreateUser(t *testing.T) {
	type fields struct {
		setup func(t *testing.T) UseCases
	}
	type args struct {
		body any
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
			args:           args{body: api.CreateUserRequest{Name: "Alice"}},
			wantStatusCode: http.StatusOK, // обычный сервер возвращает 200
			wantCT:         "application/json",
			wantBody: api.CreateUserResponse{
				Id: 10,
			},
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
			args:           args{body: api.CreateUserRequest{Name: ""}},
			wantStatusCode: http.StatusBadRequest,
			wantCT:         "application/json",
			wantBody: api.ErrorResponse{
				Code:  3,
				Error: usecases.ErrValidation.Error(),
			},
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
			args:           args{body: api.CreateUserRequest{Name: "Bob"}},
			wantStatusCode: http.StatusInternalServerError,
			wantCT:         "application/json",
			wantBody: api.ErrorResponse{
				Code:  -1,
				Error: "Internal Server Error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				useCases: tt.fields.setup(t),
			}

			bodyBytes, err := json.Marshal(tt.args.body)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
			req = req.WithContext(context.Background())
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			h.CreateUser(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Fatalf("status code = %d, want %d; body: %s", rr.Code, tt.wantStatusCode, rr.Body.String())
			}

			ct := rr.Header().Get("Content-Type")
			if ct != tt.wantCT {
				t.Fatalf("content-type = %q, want %q", ct, tt.wantCT)
			}

			switch want := tt.wantBody.(type) {
			case api.CreateUserResponse:
				got := readJSONBody[api.CreateUserResponse](t, rr)
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("body = %+v, want %+v", got, want)
				}
			case api.ErrorResponse:
				got := readJSONBody[api.ErrorResponse](t, rr)
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("body = %+v, want %+v", got, want)
				}
			default:
				t.Fatalf("unsupported wantBody type: %T", tt.wantBody)
			}
		})
	}
}

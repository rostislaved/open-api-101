package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	baseURL = "http://localhost:8080"
)

func main() {
	ctx := context.Background()

	created, err := createUser(ctx, baseURL, "Alice")
	if err != nil {
		panic(err)
	}

	fmt.Printf("created: %+v\n", created)

	got, err := getUserByID(ctx, baseURL, created.Id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("fetched: %+v\n", got)
}

func createUser(ctx context.Context, base string, name string) (CreateUserResponse, error) {
	var out CreateUserResponse

	u := strings.TrimRight(base, "/") + "/users"

	payload := CreateUserRequest{
		Name: name,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return out, err
	}
	bodyBytes = []byte("{}")

	fmt.Println(string(bodyBytes))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(bodyBytes))
	if err != nil {
		return out, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		err = json.Unmarshal(respBody, &out)
		if err != nil {
			return out, err
		}

		return out, nil
	case http.StatusBadRequest:
		err = errors.New(string(respBody))

		return out, err
	default:
		var e ErrorResponse

		err = json.Unmarshal(respBody, &e)
		if err != nil {
			return out, err
		}

		err = fmt.Errorf("create user failed: status=%d code=%d error=%s", resp.StatusCode, e.Code, e.Error)

		return out, err
	}
}

func getUserByID(ctx context.Context, base string, id int) (GetUserByIdResponse, error) {
	var out GetUserByIdResponse

	u := fmt.Sprintf("%s/users/%d", strings.TrimRight(base, "/"), id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return out, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(respBody, &out)
		if err != nil {
			return out, err
		}

		return out, nil
	}

	var e ErrorResponse

	err = json.Unmarshal(respBody, &e)
	if err != nil {
		return out, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return out, fmt.Errorf("user not found: code=%d error=%s", e.Code, e.Error)
	}

	return out, fmt.Errorf("get user failed: status=%d code=%d error=%s", resp.StatusCode, e.Code, e.Error)
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	Id int `json:"id"`
}

type GetUserByIdResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

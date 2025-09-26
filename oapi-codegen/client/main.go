package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	api "client/generated"
)

func main() {
	option1()
	option2()
}

const (
	host = "http://localhost:8080"
	id   = 1
)

// Обычный Client
func option1() {
	baseClient := http.Client{}

	client, err := api.NewClient(host, api.WithHTTPClient(&baseClient))
	if err != nil {
		panic(err)
	}

	response, err := client.GetUserById(context.Background(), id)
	if err != nil {
		panic(err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	switch response.StatusCode {
	case http.StatusOK:
		var r api.GetUserByIdResponse

		err = json.Unmarshal(bodyBytes, &r)
		if err != nil {
			panic(err)
		}

		fmt.Println(r)
	case http.StatusNotFound:
		var r api.ErrorResponse

		err = json.Unmarshal(bodyBytes, &r)
		if err != nil {
			panic(err)
		}

		fmt.Println(r)
	case http.StatusInternalServerError:
		var r api.ErrorResponse

		err = json.Unmarshal(bodyBytes, &r)
		if err != nil {
			panic(err)
		}

		fmt.Println(r)
	default:
		var r api.ErrorResponse

		err = json.Unmarshal(bodyBytes, &r)
		if err != nil {
			panic(err)
		}

		fmt.Println(r)
	}

}

// ClientWithResponses
func option2() {
	baseClient := http.Client{}

	client, err := api.NewClientWithResponses(host, api.WithHTTPClient(&baseClient))
	if err != nil {
		panic(err)
	}

	response, err := client.GetUserByIdWithResponse(context.Background(), id)
	if err != nil {
		panic(err)
	}

	switch response.StatusCode() {
	case http.StatusOK:
		fmt.Printf("%+v\n", *response.JSON200)
	case http.StatusNotFound:
		fmt.Printf("%+v\n", *response.JSON404)
	case http.StatusInternalServerError:
		fmt.Printf("%+v\n", *response.JSON500)
	default:
		fmt.Printf("%+v\n", string(response.Body))
	}
}

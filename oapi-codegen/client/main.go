package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	api "client/generated"
)

func main() {
	option2_ClientWithResponses()
}

func option1_Client() {
	baseClient := http.Client{}

	client, err := api.NewClient("http://localhost:8080", api.WithHTTPClient(&baseClient))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	response, err := client.GetUserById(ctx, 1)
	if err != nil {
		panic(err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bodyBytes))
}

func option2_ClientWithResponses() {
	baseClient := http.Client{}

	client, err := api.NewClientWithResponses("http://localhost:8080", api.WithHTTPClient(&baseClient))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	response, err := client.GetUserByIdWithResponse(ctx, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.StatusCode())
	fmt.Println(response.Status())
	fmt.Println(*response.JSON200.Id)
	fmt.Println(*response.JSON200.Name)
	fmt.Println(response.JSON404)
}

package main

import (
	"fmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"client/generated/client"
	"client/generated/client/operations"
	"client/generated/models"
)

func main() {
	option1()
	option2()
}

func option1() {
	transport := httptransport.New("localhost:8080", "", []string{"http"})
	apiClient := client.New(transport, strfmt.Default)

	params := operations.NewGetUserByIDParams()
	params.SetID(1)

	resp, err := apiClient.Operations.GetUserByID(params)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID: %v\n", *resp.Payload.ID)
	fmt.Printf("Name: %v\n", *resp.Payload.Name)
}

func option2() {
	transport := httptransport.New("localhost:8080", "", []string{"http"})
	apiClient := client.New(transport, strfmt.Default)

	name := "Alice"
	newUser := models.CreateUserRequest{Name: &name}

	params := operations.NewCreateUserParams()
	params.SetBody(&newUser)

	resp, err := apiClient.Operations.CreateUser(params)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID: %+v\n", *resp.Payload.ID)
}

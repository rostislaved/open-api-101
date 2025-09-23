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
	transport := httptransport.New("localhost:8080", "", []string{"http"})
	apiClient := client.New(transport, strfmt.Default)

	// request 1
	params := operations.NewGetUsersParams()

	resp, err := apiClient.Operations.GetUsers(params)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Users: %v\n", resp.Payload[0])

	// request 2
	createUsersRequest := operations.NewCreateUsersParams()

	name := "Bob"
	newUser := models.NewUser{Name: &name}

	createUsersRequest.SetBody(&newUser)

	resp2, err := apiClient.Operations.CreateUsers(createUsersRequest)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Users: %+v\n", resp2.Payload)

}

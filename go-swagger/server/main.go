package main

import (
	"github.com/go-openapi/loads"

	"server/generated/restapi"
	"server/generated/restapi/operations"
	"server/handlers"
	"server/usecases"
)

func main() {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		panic(err)
	}

	api := operations.NewUsersAPIAPI(swaggerSpec)

	useCases := usecases.New()
	handlers := handlers.New(useCases)

	api.GetUserByIDHandler = operations.GetUserByIDHandlerFunc(handlers.GetUsers)
	api.CreateUserHandler = operations.CreateUserHandlerFunc(handlers.CreateUsers)

	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.ConfigureFlags()
	server.Port = 8080
	server.ConfigureAPI()

	err = server.Serve()
	if err != nil {
		panic(err)
	}

}

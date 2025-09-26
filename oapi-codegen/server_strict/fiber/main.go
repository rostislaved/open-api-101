package main

import (
	"github.com/gofiber/fiber/v2"

	api "server/generated"
	"server/handlers"
	"server/usecases"
)

func main() {
	useCases := usecases.New()
	handlers := handlers.New(useCases)

	strictMux := api.NewStrictHandler(handlers, nil)

	mux := fiber.New()
	api.RegisterHandlers(mux, strictMux)

	err := mux.Listen(":8080")
	if err != nil {
		panic(err)
	}
}

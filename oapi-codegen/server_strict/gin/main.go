package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "server/generated"
	"server/handlers"
	"server/usecases"
)

func main() {
	useCases := usecases.New()
	handlers := handlers.New(useCases)

	strictMux := api.NewStrictHandler(handlers, nil)

	mux := gin.New()

	validationMW, err := addValidationMiddleware()
	if err != nil {
		panic(err)
	}

	mux.Use(validationMW)

	api.RegisterHandlers(mux, strictMux)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

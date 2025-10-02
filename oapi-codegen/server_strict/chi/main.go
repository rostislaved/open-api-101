package main

import (
	"net/http"

	api "server/generated"
	"server/handlers"
	"server/usecases"
)

func main() {
	useCases := usecases.New()
	handlers := handlers.New(useCases)

	strictMux := api.NewStrictHandler(handlers, nil)
	mux := api.Handler(strictMux)

	mux, err := addValidationMiddleware(mux)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

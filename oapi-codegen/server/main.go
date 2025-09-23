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

	mux := api.Handler(handlers)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

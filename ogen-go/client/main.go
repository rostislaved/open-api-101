package main

import (
	"context"
	"encoding/json"
	"fmt"

	api "client/generated"
)

func main() {
	client, err := api.NewClient("http://localhost:8080")
	if err != nil {
		panic(err)
	}

	params := api.GetUserByIdParams{ID: 1}

	response, err := client.GetUserById(context.Background(), params)
	if err != nil {
		panic(err)
	}

	switch p := response.(type) {
	case *api.GetUserByIdResponse:
		rr, err := json.MarshalIndent(p, " ", " ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(rr))
	case *api.GetUserByIdNotFound:
		rr, err := json.MarshalIndent(p, " ", " ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(rr))
	case *api.GetUserByIdInternalServerError:
		rr, err := json.MarshalIndent(p, " ", " ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(rr))
	default:
		rr, err := json.MarshalIndent(p, " ", " ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(rr))
	}
}

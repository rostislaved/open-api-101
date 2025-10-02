package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	middleware "github.com/oapi-codegen/nethttp-middleware"

	api "server/generated"
)

func addValidationMiddleware(mux http.Handler) (http.Handler, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}

	opts := middleware.Options{
		ErrorHandlerWithOpts: errorHandler,
		DoNotValidateServers: true,
	}

	validationMW := middleware.OapiRequestValidatorWithOptions(spec, &opts)

	mux = validationMW(mux)

	return mux, nil
}

func errorHandler(ctx context.Context, e error, w http.ResponseWriter, r *http.Request, eh middleware.ErrorHandlerOpts) {
	w.Header().Set("Content-Type", "application/json")

	statusCode := eh.StatusCode
	if statusCode == 0 {
		statusCode = http.StatusBadRequest
	}

	resp := ErrorResponse{
		Error: Error{
			Code:    statusCode,
			Message: getMsgFromErr(e),
		},
	}

	respJsonBytes, err := json.MarshalIndent(resp, " ", " ")
	if err != nil {
		http.Error(w, `{"code":500,"error":"internal error"}`, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(statusCode)

	_, err = w.Write(respJsonBytes)
	if err != nil {
		http.Error(w, `{"code":500,"error":"internal error"}`, http.StatusInternalServerError)

		return
	}
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func getMsgFromErr(err error) string {
	msg := err.Error()

	i := strings.IndexRune(msg, '\n')
	if i > 0 {
		msg = msg[:i]
	}

	return msg
}

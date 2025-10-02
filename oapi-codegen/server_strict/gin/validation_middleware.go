package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"

	api "server/generated"
)

func addValidationMiddleware() (gin.HandlerFunc, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}

	spec.Servers = nil

	opts := middleware.Options{
		ErrorHandler: errorHandler,
	}

	mw := middleware.OapiRequestValidatorWithOptions(spec, &opts)
	
	return mw, nil
}

func errorHandler(c *gin.Context, message string, statusCode int) {
	msg := getMsgFromErr(message)

	c.Header("Content-Type", "application/json")
	c.Status(statusCode)

	resp := ErrorResponse{
		Error: Error{
			Code:    statusCode,
			Message: msg,
		},
	}

	c.AbortWithStatusJSON(statusCode, resp)
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func getMsgFromErr(msg string) string {
	if i := strings.IndexRune(msg, '\n'); i > 0 {
		return msg[:i]
	}

	return msg
}

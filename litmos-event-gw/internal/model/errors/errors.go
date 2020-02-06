package errors

import (
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/logger"
	"net/http"
)

type ErrorType int

const (
	InternalError ErrorType = iota
	BadInput
	UnAuthorized
)

func HandleError(writer http.ResponseWriter, err error, errorType ErrorType) {
	if err != nil {
		logger.Logger.Errorw("Got error while handling event", "error", err)
	}
	switch errorType {
	case InternalError:
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	case BadInput:
		http.Error(writer, "Invalid input", http.StatusBadRequest)
	case UnAuthorized:
		http.Error(writer, "UnAuthorized", http.StatusUnauthorized)
	}
}

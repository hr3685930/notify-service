package errs

import (
	"net/http"
	"notify-service/pkg/http/gin"
	"runtime/debug"
)

func BadRequest(msg string) *gin.HttpError {
	return gin.NewError(http.StatusBadRequest, 4400, msg, debug.Stack())
}

func ResourceNotFound(msg string) *gin.HttpError {
	return gin.NewError(http.StatusNotFound, 4404, msg, debug.Stack())
}

func AuthenticationFailed(msg string) *gin.HttpError {
	return gin.NewError(http.StatusUnauthorized, 4401, msg, debug.Stack())
}

func AuthorizationFailed(msg string) *gin.HttpError {
	return gin.NewError(http.StatusForbidden, 4403, msg, debug.Stack())
}

func Conflict(msg string) *gin.HttpError {
	return gin.NewError(http.StatusMethodNotAllowed, 4405, msg, debug.Stack())
}

func ValidationFailed(msg string) *gin.HttpError {
	return gin.NewError(http.StatusUnprocessableEntity, 4422, msg, debug.Stack())
}

func InternalError(msg string) *gin.HttpError {
	return gin.NewError(http.StatusInternalServerError, 5500, msg, debug.Stack())
}

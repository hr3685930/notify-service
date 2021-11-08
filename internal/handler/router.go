package handler

import (
	"github.com/gin-gonic/gin"
	"notify-service/internal/server"
	gin2 "notify-service/pkg/http/gin"
)

func Route(e *gin.Engine) {
	e.Use(gin2.ErrHandler(HTTPErrorReport))
	e.GET("/test", server.Name)
}

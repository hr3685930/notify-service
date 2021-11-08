package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"notify-service/internal/errs"
	"notify-service/internal/service"
	"notify-service/internal/types"
)

func Name(c *gin.Context) {
	user, err := service.NewUserService().GetAll(c.Request.Context(),c.Query("after"), c.Query("before"))
	response := &types.TestResponse{}
	//bff层需要这个 微服务不需要
	err = copier.Copy(&response, &user)
	if err != nil {
		c.Error(errs.ResourceNotFound(err.Error()))
		return
	}
	c.JSON(200, response);return
}

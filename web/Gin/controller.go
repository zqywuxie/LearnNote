package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
}

func (c *UserController) GetUser(context *gin.Context) {
	//context.Writer.WriteString("controller")
	context.String(http.StatusOK, "hello")
	//c.Ctx.WriteString("你好")
}

package main

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestController(t *testing.T) {
	engine := gin.Default()
	c := &UserController{}
	engine.GET("/", func(context *gin.Context) {
		context.Writer.WriteString("你好哇")
	})
	engine.GET("/user", c.GetUser)
	engine.Run(":9090")
}

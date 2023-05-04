//go:build e2e

package main

import (
	"github.com/beego/beego/v2/server/web"
	"testing"
)

func TestController(t *testing.T) {
	c := &UserController{}
	web.Router("/user", c, "get:GetUser")
	web.Run(":9098")
}

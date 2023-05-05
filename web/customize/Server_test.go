// @Author: zqy
// @File: Server_test.go
// @Date: 2023/5/4 15:00

package customize

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewHttpServer()

	s.Get("/", func(ctx *Context) {
		fmt.Println("ok")
		fmt.Println("ok1")
	})
	s.Start(":9090")
}

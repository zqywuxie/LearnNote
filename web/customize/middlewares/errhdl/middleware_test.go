// @Author: zqy
// @File: middleware_test.go.go
// @Date: 2023/5/14 16:19
// @Description todo

package errhdl

import (
	"GoCode/web/customize"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.AddData(http.StatusNotFound, []byte(`
<html>
	<h1>抱歉你走丢了</h1>
</html>
`)).
		AddData(http.StatusBadGateway, []byte(`
<html>
	<h1>你搞什么,不允许</h1>
<html>`))
	server := customize.NewHttpServer(customize.ServerWithMiddleWare(builder.Build()))
	server.Start(":9090")
}

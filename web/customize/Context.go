// @Author: zqy
// @File: Context.go
// @Date: 2023/5/4 15:08

package customize

import "net/http"

type Context struct {
	req        *http.Request
	resp       http.ResponseWriter
	pathParams map[string]string
}

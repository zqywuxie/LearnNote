// @Author: zqy
// @File: Router_test.go
// @Date: 2023/5/4 20:40

package customize

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
	}

	mockHandler := func(ctx *Context) {}
	r := newRouter()
	for _, tr := range testRoutes {
		r.AddRoute(tr.method, tr.path, mockHandler)
	}

	// 断言路由树 进行测试
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path: "user",
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}

	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)
	// 不能直接断言，因为router里面有方法，是不可比的
}

func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的http method"), false
		}
		msg, ok := dst.equal(v)
		if !ok {
			return msg, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数不匹配"), false
	}
	if n.starChildren != nil {
		msg, ok := n.starChildren.equal(y.starChildren)
		if !ok {
			return msg, ok
		}
	}

	// 方法通过反射来比较
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintf("节点处理逻辑不匹配"), false
	}
	for k, v := range n.children {
		child, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("子节点不存在"), false
		}
		// 递归深入
		msg, ok := v.equal(child)
		if !ok {
			return msg, ok
		}
	}
	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		//{
		//	method: http.MethodGet,
		//	path:   "/",
		//},
		//{
		//	method: http.MethodPost,
		//	path:   "/order/get",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/order/*",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/*",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/*/*",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/*/abc",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/*/abc/*",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/user/*/home",
		//},
		//{
		//	method: http.MethodPost,
		//	path:   "/login/:username",
		//},
		{
			method: http.MethodPost,
			path:   "/logout/:id(^\\d{4}$)",
		},
		//{
		//	method: http.MethodPost,
		//	path:   "/login/*",
		//},
	}

	r := newRouter()

	mockHandler := func(ctx *Context) {}

	for _, t := range testRoutes {
		r.AddRoute(t.method, t.path, mockHandler)
	}
	testCases := []struct {
		name     string
		method   string
		path     string
		found    bool
		wantNode *matchInfo
	}{
		//{
		//	name:   "method not found",
		//	method: http.MethodOptions,
		//	//found: false,
		//},
		{
			name:   "root",
			method: http.MethodGet,
			path:   "/",
			found:  true,
			wantNode: &matchInfo{
				n: &node{
					path:    "/",
					handler: mockHandler,
				},
			},
		},
		{
			name:   "first level",
			method: http.MethodGet,
			path:   "/order",
			found:  true,
			wantNode: &matchInfo{
				n: &node{
					path: "order",
					children: map[string]*node{
						"get": &node{
							path:    "get",
							handler: mockHandler,
						},
					},
				},
			},
		},
		{
			name:   "leaf",
			method: http.MethodGet,
			path:   "/xx/abc",
			found:  true,
			wantNode: &matchInfo{
				n: &node{
					path:    "abc",
					handler: mockHandler,
				},
			},
		},
		{
			name:   "login username",
			method: http.MethodPost,
			path:   "/login/zqy",
			found:  true,
			wantNode: &matchInfo{
				n: &node{
					path:    "username",
					handler: mockHandler,
				},
				pathParam: map[string]string{
					"username": "zqy",
				},
			},
		},
		{
			name:   "reg logout",
			method: http.MethodPost,
			path:   "/logout/1234",
			found:  true,
			wantNode: &matchInfo{
				pathParam: map[string]string{
					"id": "1234",
				},
				n: &node{
					path:    "id(^\\d{4}$)",
					handler: mockHandler,
				},
			},
		},
		//{
		//	name:   "/order/*",
		//	method: http.MethodGet,
		//	path:   "/order/delete",
		//	found:  true,
		//	wantNode: &node{
		//		path:    "*",
		//		handler: mockHandler,
		//	},
		//},
		//{
		//	name:   "/user/*/help",
		//	method: http.MethodGet,
		//	path:   "/user/tome/help",
		//	found:  true,
		//	wantNode: &node{
		//		path:    "help",
		//		handler: mockHandler,
		//	},
		//},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.found, found)
			if !found {
				return
			}
			//assert.Equal(t, tc.wantNode.children, info.children)
			assert.Equal(t, info.pathParam, tc.wantNode.pathParam)
			msg, Hflag := info.n.equal(tc.wantNode.n)
			assert.True(t, Hflag, msg)
			fmt.Println(info.pathParam["username"])
		})
	}
}

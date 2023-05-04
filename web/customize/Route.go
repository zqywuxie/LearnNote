// @Author: zqy
// @File: Route.go
// @Date: 2023/5/4 20:34

package customize

import (
	"fmt"
	"strings"
)

type node struct {
	// children path => node
	children map[string]*node
	path     string

	// 到达叶子节点才执行
	handler HandleFunc
}

type router struct {
	// 使用http method 进行组织tree
	trees map[string]*node
}

// newRouter 创建路由树
func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

func (r *router) AddRoute(method string, path string, handle HandleFunc) {
	if path == "" {
		panic("web: 空路径")
	}
	if path[0] != '/' {
		panic("web: 路径必须以/开头")
	}
	if path[0] != '/' && path[len(path)-1] != '/' {
		panic("web : 路径不能以/结尾")
	}
	root, ok := r.trees[method]
	// 如果根不存在就创建
	if !ok {
		root = &node{path: "/"}
		r.trees[method] = root
	}

	if path == "/" {
		if root.handler != nil {
			// 如果handler存在，说明此路由被注册了，所以将其handlerFunc覆盖
			panic("web : 路由冲突")
		}
		root.handler = handle
		return
	}
	// /user/login => 空，user,login
	segs := strings.Split(path[1:], "/")
	for _, s := range segs {
		if s == "" {
			panic(fmt.Sprintf("非法路由不允许使用"))
		}
		root = root.childCreate(s)
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web : 路由冲突[%s]", path))
	}
	root.handler = handle

}

func (n *node) childCreate(path string) *node {

	// 不存在子树就新建
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]

	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}

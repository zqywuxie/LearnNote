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

	// 通配符匹配
	starChildren *node

	paramChildren *node

	path string

	// 到达叶子节点才执行
	handler HandleFunc
}

type matchInfo struct {
	n         *node
	pathParam map[string]string
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

	if path[0] == ':' {
		if n.starChildren != nil {
			panic(fmt.Sprintf("web：非法路由,不允许同时注册（通配符)"))
		}
		if n.paramChildren != nil {
			panic(fmt.Sprintf("web：路由冲突"))
		} else {
			n.paramChildren = &node{
				path: path[1:],
			}
		}
		return n.paramChildren
	}

	if path == "*" {
		// 避免重复注册
		if n.paramChildren != nil {
			panic(fmt.Sprintf("web：非法路由,不允许同时注册（参数路径)"))
		}

		if n.starChildren == nil {
			n.starChildren = &node{
				path: "*",
			}
		}
		return n.starChildren
	}

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

// 判断是否节点是否存在
// 第一个bool 判断参数是否命中
// 第二个bool 判断是否存在节点
func (n *node) childOf(path string) (*node, bool, bool) {

	// 如果子节点不存在，或者静态匹配不成功 都查看通配符是否存在
	if n.children == nil {
		if n.paramChildren != nil {
			return n.paramChildren, true, true
		}
		return n.starChildren, false, n.starChildren != nil
	}
	child, ok := n.children[path]
	// 优先静态匹配，没有就返回通配符匹配
	if !ok {
		if n.paramChildren != nil {
			return n.paramChildren, true, true
		}
		return n.starChildren, false, n.starChildren != nil
	}
	return child, false, ok
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	//如果是根就直接返回了
	if path == "/" {
		return &matchInfo{n: root}, true
	}
	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")

	// 创建matchInfo
	var pathParams map[string]string
	for _, s := range segs {
		var paramOk bool
		root, paramOk, ok = root.childOf(s)
		if !ok {
			return nil, false
		}
		if paramOk {
			pathParams = make(map[string]string)
			pathParams[root.path] = s
		}
	}
	mi := &matchInfo{
		n:         root,
		pathParam: pathParams,
	}
	return mi, true
}

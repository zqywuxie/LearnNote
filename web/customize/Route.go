// @Author: zqy
// @File: Route.go
// @Date: 2023/5/4 20:34

package customize

import (
	"fmt"
	"regexp"
	"strings"
)

type nodeType int

const (
	nodeTypeAny = iota
	nodeTypeParam
	nodeTypeReg
	nodeTypeStatic
)

// 静态 > 正则 > 路径参数 > 通配符
type node struct {
	// children path => node
	children map[string]*node

	// 通配符匹配
	starChildren *node

	paramChildren *node

	// 正则匹配
	regChildren *node
	regExpr     *regexp.Regexp

	//参数key
	paramName string

	route string

	path string

	// 到达叶子节点才执行
	handler HandleFunc
	typ     int
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
		root.route = "/"
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

	// route 获得完整匹配路径
	root.route = path

}

func (n *node) childCreate(path string) *node {

	if path[0] == ':' {
		paramName, reg, ok := n.parseParam(path)
		if !ok {
			return n.childOrCreateParam(path[1:], paramName)
		} else {
			return n.childOrCreateReg(path, paramName, reg)
		}
	}

	if path == "*" {
		// 避免重复注册
		if n.paramChildren != nil {
			panic(fmt.Sprintf("web：非法路由,不允许同时注册（参数路径)"))
		}

		if n.starChildren == nil {
			n.starChildren = &node{
				path: "*",
				typ:  nodeTypeAny,
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
		child = &node{
			path: path,
			typ:  nodeTypeStatic,
		}
		n.children[path] = child
	}
	return child
}

func (n *node) childOrCreateParam(path, paramName string) *node {
	if n.starChildren != nil || n.regChildren != nil {
		panic(fmt.Sprintf("web：非法路由,不允许同时注册（通配符,正则)"))
	}
	if n.paramChildren != nil {
		panic(fmt.Sprintf("web：路由冲突"))
	}

	n.paramChildren = &node{
		path:      path,
		paramName: paramName,
		typ:       nodeTypeParam,
	}
	return n.paramChildren
}

func (n *node) childOrCreateReg(path, paramName, reg string) *node {
	if n.starChildren != nil || n.paramChildren != nil {
		panic(fmt.Sprintf("web：非法路由,不允许同时注册（通配符,参数路径)"))
	}
	if n.regChildren != nil {
		// 判断是否存在
		if n.regChildren.regExpr.String() != reg || n.paramName != paramName {
			panic(fmt.Sprintf("web：路由冲突"))
		}
	} else {
		compile, err := regexp.Compile(reg)
		if err != nil {
			panic(fmt.Sprintf("正则表达式错误"))
		}
		n.regChildren = &node{
			path:      path[1:],
			paramName: paramName,
			regExpr:   compile,
			typ:       nodeTypeReg,
		}
	}
	return n.regChildren
}
func (n *node) parseParam(path string) (string, string, bool) {
	// 去除:
	path = path[1:]
	segs := strings.SplitN(path, "(", 2)
	// /:id(xxx)
	if len(segs) == 2 {
		expr := segs[1]
		if strings.HasSuffix(expr, ")") {
			return segs[0], expr[:len(expr)-1], true
		}
	}
	return path, "", false
}

// 判断是否节点是否存在
// 第一个bool 判断参数是否命中
// 第二个bool 判断是否存在节点
func (n *node) childOf(path string) (*node, bool) {

	// 如果子节点不存在，或者静态匹配不成功 都查看通配符是否存在
	if n.children == nil {
		return n.childOfNoStatic(path)
	}
	child, ok := n.children[path]
	// 优先静态匹配，没有就返回通配符匹配
	if !ok {
		return n.childOfNoStatic(path)
	}
	return child, ok
}
func (n *node) childOfNoStatic(path string) (*node, bool) {
	if n.regChildren != nil {
		// 如果匹配成功 就返回节点
		if n.regChildren.regExpr.MatchString(path) {
			return n.regChildren, true
		}
	}

	if n.paramChildren != nil {
		return n.paramChildren, true
	}
	// 最后返回通配符匹配
	return n.starChildren, n.starChildren != nil
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
	mi := &matchInfo{}
	//child := root
	for _, s := range segs {
		//var paramOk bool
		child, ok := root.childOf(s)
		if !ok {
			// 如果最后节点是通配符匹配
			// todo 注意这里的root还是child
			if root.typ == nodeTypeAny {
				mi.n = root
				// todo 直接返回？
				return mi, true
			}
			return nil, false
		}
		if child.paramName != "" {
			pathParams = make(map[string]string)
			pathParams[child.paramName] = s
		}
		root = child
	}
	mi.n = root
	mi.pathParam = pathParams
	return mi, true
}

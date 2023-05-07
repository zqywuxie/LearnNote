// @Author: zqy
// @File: Context.go
// @Date: 2023/5/4 15:08

package customize

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// Context
/*
在Go中，http.Request是一个结构体而http.ResponseWriter是一个接口。
这是因为http.Request包含了来自客户端的HTTP请求中的所有信息，
例如请求头、URL和请求体等，
因此它需要有一个结构体来存储这些数据。另一方面，
http.ResponseWriter只需要实现一些方法来写入响应正文到客户端，
具体的实现方式可以多种多样（例如输出到文件或内存缓冲区等），
因此使用一个接口来表示响应处理器更加灵活。
*/
type Context struct {
	req         *http.Request
	resp        http.ResponseWriter
	pathParams  map[string]string
	queryValues url.Values
}

var (
	JSONUseNumber = true
)

// BindJson 将HTTP请求的json数据绑定到一个对象模型上
func (c *Context) BindJson(val any) error {
	if val == nil {
		return errors.New("web:输入内容不能为nil")
	}
	if c.req.Body == nil {
		return errors.New("web:body 为 nil")
	}

	decoder := json.NewDecoder(c.req.Body)

	// 数组使用Number类型 -> string
	// 默认为float64
	//decoder.UseNumber()

	// 对未知字段检测报错
	// 如果解析为结构体user 只有name，如果多出了age就会报错
	//decoder.DisallowUnknownFields()

	// 如上，如果要对一些功能进行使用，是否要添加参数让用户来选择是否启动？
	// BindJsonOpt(val any,userNumber bool ..)
	// 关于这些方法的开发，要根据用户需求。如果一个小众需求，用户可以自己解决，那么就不要在框架
	// 核心进行设计
	return decoder.Decode(val)
}

// FormValue 获得表单数据
// 关于ParseForm 重复解析的问题，其方法会先查找是否有解析后的Form和PostForm,如果为nil才会解析
// 是幂等性
// 不打算提供其他数据类型的方法，因为数据类型较多，让用户自己进行转换
// func (c *Context) FormValueInt64(key string) (string, error)
func (c *Context) FormValue(key string) (string, error) {
	err := c.req.ParseForm()
	if err != nil {
		return "", err
	}
	return c.req.FormValue(key), nil
}

type StringValue struct {
	string
	error
}

// QueryValue 获得查询参数
func (c *Context) QueryValue(key string) (string, error) {
	//query := c.req.URL.Query()
	// 这里不能进行如上判断非空，查看源码得,每次都会进行ParseQuery，并且make，没有像form进行缓存
	/*
		func ParseQuery(query string) (Values, error) {
			m := make(Values)
			err := parseQuery(m, query)
			return m, err
		}
	*/
	if c.queryValues == nil {
		c.queryValues = c.req.URL.Query()
	}
	value, ok := c.queryValues[key]
	if !ok {
		return "", errors.New(" web:key不存在")
	}
	return value[0], nil
}

func (c *Context) PathValue(key string) (string, error) {
	value, ok := c.pathParams[key]
	if !ok {
		return "", errors.New("web : 没找到相关参数")
	}
	return value, nil
}

// PathValueAsString 将不同数据类型返回值进行封装
func (c *Context) PathValueAsString(key string) StringValue {
	value, ok := c.pathParams[key]
	if !ok {
		return StringValue{
			"", errors.New("web : 没找到相关参数"),
		}
	}
	return StringValue{value, nil}
}

// AsInt64 单独进行转换
func (s StringValue) AsInt64() (int64, error) {
	if s.error != nil {
		return -1, s.error
	}
	return strconv.ParseInt(s.string, 10, 64)
}
func (c *Context) test() {
	asInt64, err := c.PathValueAsString("123").AsInt64()
}

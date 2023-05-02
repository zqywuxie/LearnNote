# web

## 1.最简单的web服务器

~~~go
package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "主页")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "登录页")
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
~~~



## 2.request里面的相关方法属性

### Body/GetBody

#### Body

request里面的body，开发者只能获得一次。后续读取不会报错，但是什么都读不到。

~~~go
func GetBodyOnce(w http.ResponseWriter, r *http.Request) {
	all, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "第一次没读到body")
		return
	}
	fmt.Fprintf(w, "read the data :%s\n", string(all))
	readAll, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "第二次没读到body")
		return
	}
	fmt.Fprintf(w, "read the data :%s\n", string(readAll))
}
~~~

![image-20230501214920495](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230501214920495.png)

#### GetBody

原则上是可以多次读取，但是http库里面的GetBody没有进行赋值，是nil。

~~~go
//方法字段
GetBody func() (io.ReadCloser, error)
========


func GetBodyIsNil(w http.ResponseWriter, r *http.Request) {
	if r.GetBody == nil {
		fmt.Fprintf(w, "GetBody is nill")
		return
	} else {
		fmt.Fprintf(w, "not nil")
		return
	}
}
~~~

![image-20230501215414651](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230501215414651.png)



### header

注意go的规则，它会将header进行格式化，自动将header首字母大写了



### Form

> 1. 使用Form表单之前要使用request.ParseForm
> 2. 建议加上 Content-Type: application/x-www-form- 
>    urlencoded



## 3. type定义

抽象出方法，达到一个方法的内聚。

### 3.1 类型

> type A B 起了一个新的类型
>
> ​	- A不能用B的方法，除非强转
>
> type A = B 相当于别名
>
> ​	- A可以使用B的方法



总结：

- 鸭子类型，只要有个某个接口的全部方法，就实现了这个接口
- type作用，先有抽象再有实现，先写接口



### 3.2 结构自引用

结构体自引用，比如链式存储。只能使用指针类型，原因如下:

> 在Go语言中，结构体是一个值类型（value type），这意味着当你在函数中传递结构体时，会进行值复制。如果结构体包含自引用的成员变量，并且这些成员变量是非指针类型，则在赋值或者传递时会出现**无限递归的情况**。
>
> 使用指针可以避免这种情况，因为指针传递的是内存地址，而不是实际的值。这样，在修改指针所指向的值时，所有指向该值的指针都会受到影响。
>
> 此外，使用指针还有一个好处，就是**可以减少内存占用**。如果结构体中的某个成员变量是一个大型数据结构，直接将它包含在结构体中可能会导致大量的内存分配和复制操作。通过将该成员变量定义为指针类型，可以将内存占用降至最小。
>
> 因此，为了避免无限递归并减少内存占用，Go语言中的结构体自引用必须使用指针。



## 4.抽象使用

### Http Server实现

---

 **原始服务器的使用**

~~~go
http.HandleFunc("/", login)
http.ListenAndServe(":9090", nil)
~~~



#### 进行封装

`Server.go`

~~~go
package main

import "net/http"

type Server interface {
	RouteTable
	Start(address string) error
}

type sdkHttpServer struct {
	Name    string
	handler Handler
}

func (s *sdkHttpServer) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	//TODO implement me )
	s.handler.Route(method, pattern, handleFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	//http.Handle("/", s.handler)
    //路由处理器放下面？ todo
	return http.ListenAndServe(address, s.handler)
}

func NewHttpServer(name string) Server {
	return &sdkHttpServer{Name: name, handler: NewHandlerBaseOnMap()}
}

func SignUp(ctx *Context) {
	ctx.ReadJson(nil)
}

~~~

解释：

1.向下托付

~~~go
func (s *sdkHttpServer) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	//交给下一层进行实现，为的是不暴露具体实现细节
	s.handler.Route(method, pattern, handleFunc)
}
~~~





`Context.go`

将上下文一些操作进行了封装

~~~go
package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{W: w, R: r}
}

func (c *Context) ReadJson(obj interface{}) (err error) {
	all, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(all, obj)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) WriteJson(code int, resp interface{}) error {
	c.W.WriteHeader(code)
	marshal, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = c.W.Write(marshal)
	return err
}

func (c *Context) okJson(resp interface{}) error {
	return c.WriteJson(http.StatusOK, resp)
}

~~~

解释：

1.Context (上下文）是HTTP请求时传递数据和控制操作的。将其独立封装并且设计数据操作的方法



`map_based_handler.go`

~~~go
package main

import (
	"net/http"
)

type HandlerBaseOnMap struct {
	// handler method + url
	//  handlers 强耦合了
	handlers map[string]func(ctx *Context)
}
// 多处地方使用到了Route 进行提取
type RouteTable interface {
	Route(method string, pattern string, handleFunc func(ctx *Context))
}

// Handler 组合
type Handler interface {
	http.Handler
	RouteTable
}

// 是为了校验，HandlerBaseOnMap是否实现了Handler，如果没有实现完方法就会报错
var _ Handler = &HandlerBaseOnMap{}

// 注意一个问题，使用接口返回值，而不是具体实现类
func NewHandlerBaseOnMap() Handler {
	return &HandlerBaseOnMap{handlers: make(map[string]func(ctx *Context))}
}

func (h *HandlerBaseOnMap) Route(method string, pattern string, handleFunc func(ctx *Context)) {
	//TODO implement me
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		// 注册路由
		key := h.GetKey(method, pattern)
		h.handlers[key] = handleFunc
	})
}

// http.handler里的实现接口
func (h *HandlerBaseOnMap) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	key := h.GetKey(request.Method, request.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		handler(NewContext(writer, request))
	} else {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("Not Found"))
	}
}
func (h *HandlerBaseOnMap) GetKey(method string, patten string) string {
	return method + "#" + patten
}

~~~

解析:

1. 使用`map[string]func(ctx *Context)`，设计路由

2. > // 是为了校验，HandlerBaseOnMap是否实现了Handler，如果没有实现完方法就会报错
   > var _ Handler = &HandlerBaseOnMap{}

3. 因为多处使用到了`Route`接口，就提取为一个独立接口



## 5.错误处理

> 当你怀疑可以用error的时候，就不需要panic
>
> 一般情况下，只有快速失败的过程，才会考虑panic。比如服务器启动过程等

常见`errors`错误处理库

```go
type MyError struct {
}

func (m *MyError) Error() string {
	return "自定义error"
}
===
func ErrorsPkg() {
	myError := &MyError{}
   // 创建error
    errors.New()
    //fmt.Errorf 包装error
	err := fmt.Errorf("this is an wrapped error %w", myError)
	//解包
    err = errors.Unwrap(err)
	
    // err会一直解包，所以返回true
    if errors.Is(err, myError) {
		fmt.Println("自动解包")
	}


	copyErr := &MyError{}
	if errors.As(err, &copyErr) {
		
	}
}
```

## 6.闭包

~~~go
package main

import "fmt"

func ReturnClosure(name string) func() string {
	return func() string {
		return name + "hello"
	}
}

//闭包的延时绑定
func Delay() {
	fns := make([]func(), 0, 10)
	for i := 0; i < 10; i++ {
		fns = append(fns, func() {
			fmt.Println(i)
			fmt.Printf("hello,this is %d\n", i)
		})
	}
	for _, fn := range fns {
		fn()
	}
}
func main() {
	i := 123
	a := func() {
		fmt.Printf("i is %d \n", i)
	}
	a()
	fmt.Println(ReturnClosure("tom")())
	Delay()
}

~~~

解释：

1.//闭包的延时绑定，是使用了变量的引用。所以一直输出10

~~~go
for i := 0; i < 10; i++ {
    fns = append(fns, func() {
        // 只是得到了i的引用，所以for结束后，i的引用是10
        fmt.Println(i)
        fmt.Printf("hello,this is %d\n", i)
    })
}
~~~



### AOP：闭包实现责任链

AOP：横向关注点，覆盖多重业务的逻辑：日志，限流等

filter：真正请求前进行拦截掉不需要的东西



~~~go
// FilterBuilder 责任链
type FilterBuilder func(next Filter) Filter

type Filter func(ctx *Context)
~~~

因为拦截器在接收请求之前就设置了

~~~go
func NewHttpServer(name string, builders ...FilterBuilder) Server {
	handler := NewHandlerBaseOnMap()
	var root Filter = handler.ServeHTTP
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		// 形成了一个责任链 filterA(filterB(filterC))
		//func(func(func()))
		root = b(root)
	}
    // 添加一个root Filter字段，将filerchain根节点传递给它
	return &sdkHttpServer{Name: name, handler: NewHandlerBaseOnMap(), root: root}
}
~~~

然后启动时进行调用filerchain

~~~go
func (s *sdkHttpServer) Start(address string) error {
	//http.Handle("/", s.handler)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		c := NewContext(writer, request)

		// 开始执行
		s.root(c)
	})
	return http.ListenAndServe(address, nil)
}
~~~

## 7.sync

## 8.路由树

![image-20230502161228795](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230502161228795.png)
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

考虑什么时候使用`sync.map和map`，因为并发相关的东西性能都较低，因为考虑到安全问题。

~~~go
package main

import (
	"sync"
)

var mutex sync.Mutex
var rwmutex sync.RWMutex

func Mutex() {
	mutex.Lock()
	//推荐使用defer
	defer mutex.Unlock()

	// 如果在解锁前，panic了，那么就不会解锁了
	//mutex.Unlock()
}

// 不可重入
func Failed1() {
	mutex.Lock()
	defer mutex.Unlock()

	// 会死锁
	//如果只有一个goroutine会导致程序崩溃
	/*
		这段代码会导致死锁，因为在第一个mutex.Lock()被执行后，
		锁被占用，并且在第二个mutex.Lock()尝试获取这个锁时就会被阻塞
		，直到第一个锁被释放。
		而由于第一个锁的释放需要等到函数结束才能进行，
		所以第二个锁永远无法被获取到，就会一直阻塞在那里，导致死锁。
	*/
	mutex.Lock()
	defer mutex.Unlock()
}

// 不可升级,加了读锁再加写锁就报错
func Failed2() {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	rwmutex.Lock()
	defer rwmutex.Unlock()
}

func main() {
	Failed2()
	//s := sync.Map{}
	//s.Store("cat", "Tom")
	//value, ok := s.Load("cat")
	//if ok {
	//	// value.(string) 类型断言
	//	fmt.Println(len(value.(string)))
	//}
}

~~~



## 8.路由树

![image-20230502161228795](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230502161228795.png)

### v1 精确匹配

~~~go
package main

import "strings"

type HandlerBaseOnTree struct {
	root *Node
}

type Node struct {
	path     string
	children []*Node

	//如果是叶子节点，那么匹配上就可以调用该方法
	handler HandlerFunc
}

func newNode(path string) *Node {
	return &Node{
		path:     path,
		children: make([]*Node, 0, 2),
	}

}

/*
新增一条：/user/friends
步骤：
1.从根节点出发，作为当前节点
2.查找命中的子节点
3.将子节点作为当前节点，重复2
4.如果当前节点的子节点没有匹配下一段的，为下一
段路径创建子节点
5.如果路径还没结束，重复4
6.新增成功
*/
func (h *HandlerBaseOnTree) ServeHTTP(c *Context) {
	//TODO implement me
	panic("implement me")
}

func (h *HandlerBaseOnTree) Route(method string, pattern string, handleFunc HandlerFunc) {
	// /user/login/ 去掉斜杠 -> user/login
	path := strings.Trim(pattern, "/")
	paths := strings.Split(path, "/")
	root := h.root
	for index, path := range paths {
		matchChild, ok := root.findMatchChild(path)
		if ok {
			root = matchChild
		} else {
			// 没有找到就进行创建子节点(后面的一整段)
			root.CreateChild(root, paths[index:], handleFunc)
			return
		}
	}

}

// 该方法定义在node上较为合理，node自身去查找子节点
func (h *Node) findMatchChild(path string) (*Node, bool) {
	for _, child := range h.children {
		if child.path == path {
			return child, true
		}
	}
	return nil, false
}

func (h *Node) CreateChild(root *Node, paths []string, handleFunc HandlerFunc) {
	cur := root
	for _, path := range paths {
		node := newNode(path)
		cur.children = append(cur.children, node)
		cur = node
	}
	cur.handler = handleFunc
}

func main() {

}

~~~



### v2  /* 通配符匹配

### v3 路径参数

前面是根据路径进行不断匹配，那么如果有多条匹配规则，如`/*`,`/:username`,`/正则表达`。

```go
//if child.path == path && child.path != "*" {
//	return child, true
//}
//// 就不直接返回，要等后续路由匹配完全
//if child.path == "*" {
//	Node = child
//}
```



那么就可以将匹配规则进行抽象出来。并且将匹配路径的工作交给路由树的节点去做，所以抽象出节点

~~~go
package main

const (
	nodeTypeRoot = iota

	// /*
	nodeTypeAny

	//路径参数
	nodeTypeParam

	//正则表达
	nodeTypeReg

	//完全匹配
	nodeTypeStatic
)

type matchFunc func(path string, c *Context) bool
type Node struct {
	child    []*Node
	handler  HandlerFunc
	mathFunc matchFunc
	pattern  string
	nodeType int
}

// 将匹配方法抽象出来
func newStaticNode(path string) *Node {
	return &Node{
		child: make([]*Node, 0, 2),
		mathFunc: func(p string, c *Context) bool {
			return path != "*" && path == p
		},
		pattern:  path,
		nodeType: nodeTypeStatic,
	}
}
func newParamNode(path string) *Node {
	paramName := path[1:]
	return &Node{
		child: make([]*Node, 0, 2),
		mathFunc: func(p string, c *Context) bool {
			if c != nil {
				c.PathParams[paramName] = p
			}
			return p != "*"
		},
		pattern:  path,
		nodeType: nodeTypeParam,
	}
}
func newTypeAny() *Node {
	return &Node{
		child: make([]*Node, 0, 2),
		mathFunc: func(p string, c *Context) bool {
			return true
		},
		pattern:  "*",
		nodeType: nodeTypeAny,
	}
}

~~~

最后根据不同规则，一条路径可能匹配到的规则很多。所以就需要设计一个规则优先级，也就是type

~~~go
func (h *Node) findMatchChild(path string, c *Context) (*Node, bool) {
	// 还要对通配符进行校验
	// a1/*
	// a1/*/a2 | a1* 等不允许
	candidates := make([]*Node, 0, 2)
	for _, child := range h.child {
		if child.mathFunc(path, c) {
			candidates = append(candidates, child)
		}
	}
	if len(candidates) == 0 {
		return nil, false
	}
    // 进行优先级判定
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].nodeType < candidates[j].nodeType
	})
	return candidates[len(candidates)-1], true
}
~~~



### 自定义规则

开发者除了一些常见的匹配规则，还会创建出一些自定义的匹配规则。那么此时就不方便添加抽象，所以就提前进行判定。使用**工厂模式**

~~~go
type Factory func(path string) *Node

var factory Factory

func RegisterFactory(f Factory) {
    // 后续就使用 factory() 进行创建节点
	factory = f
}

func main() {

	// 先执行自定义规则
	RegisterFactory(func(path string) *Node {
		if strings.HasPrefix(path, ":dad") {
			return &Node{}

			// 如果不满足就委托给
		} else {
			return newNode(path)
		}
	})
}
~~~



## 9.优雅退出

![image-20230503110050506](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230503110050506.png)

5.（再次收到关闭请求，就强制关闭）

### 1.channel

有缓冲的channel

~~~go
func channelWithCache() {
	ch := make(chan string, 1)
	go func() {
		ch <- "1 message"
		time.Sleep(time.Second)
		ch <- "2 message"
	}()
    // 2s期间其实已经有了第一个缓冲，并且1s < 2s,所以第二条数据就在阻塞
	time.Sleep(time.Second * 2)
    // 所以当输出第一个后，第二个也会马上进行输出，输出间隔就小于1s
	msg := <-ch
	fmt.Println(time.Now().String() + msg)
	msg = <-ch
	fmt.Println(time.Now().String() + msg)
}
~~~

### 2.select

常与for组合使用

~~~go
func selectUse() {
	ch1 := make(chan string, 2)
	ch2 := make(chan string, 2)

	go func() {
		time.Sleep(time.Second)
		ch1 <- "from ch1"
	}()
	go func() {
		time.Sleep(time.Second)
		ch2 <- "from ch2"
	}()

	// 同时有数据的情况，顺序没有保证
    // 为什么是2，因为只进行了两条数据的输入
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println(msg)
		case msg := <-ch2:
			fmt.Println(msg)
		}
	}
}
~~~

### 3.hook

~~~go
type Hook func(c context.Context) error

func BuildCloseServerHook(servers ...http.Server) Hook {
	return func(c context.Context) error {
		wg := sync.WaitGroup{}
		doneCh := make(chan struct{})
		wg.Add(len(servers))
		for _, s := range servers {
			go func(srv http.Server) {
				err := srv.Shutdown(c)
				if err != nil {
					fmt.Printf("server shuwdown error:%v", err)
				}
				time.Sleep(time.Second)
				wg.Done()
			}(s)
		}

		//再开一个协程进行等待
		go func() {
			wg.Wait()
			// 如果关闭完就发送信息
			doneCh <- struct{}{}
		}()

		select {
		case <-c.Done():
			fmt.Println("closing servers timeout")
            return erros.New("timeout")
		case <-doneCh:
			fmt.Println("close all servers")
            return nil
		}
		return nil
	}
}
~~~

关于为什么doneCh设计为`chan struct{}`,参考了context的done

~~~go
//context
Done() <-chan struct{}
~~~



### 4.context

Go提供的线程安全工具，称为上下文

![image-20230503160131740](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230503160131740.png)

~~~go
func withTimeout() {
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	start := time.Now().Unix()
	<-timeout.Done()
	end := time.Now().Unix()
	fmt.Println(end - start)
}

func withValue() {
	parentKey := "parent"
	parent := context.WithValue(context.Background(), parentKey, "this is parent")
	sonKey := "son"
	son := context.WithValue(parent, sonKey, "this is son")

	fmt.Println(parent.Value(parentKey))
	fmt.Println(son.Value(parentKey))
	// 类似继承，子能查找到父，父查找不到子
	fmt.Println(parent.Value(sonKey))
}
~~~



[Java/Go] context与thread-local

> go团队推荐通过在函数的参数中显式的传递context来实现goroutine的local storage。

- Go官方没有支持thread-local
- 缺乏像Java的ThreadLocal，所以大多都是依赖context在方法直接传递。建议自己的方法签名，都把context.Context作为第一个参数

### 5.atomic包

原子性的基本数据类型操作

- Addxxxx
- Loadxxx 读取
- CompareAndSwapxxx CAS操作，比较并交换
- Storexx 写入一个值
- Swapxxx: 写入一个值，并返回旧的值。与CompareAndSwap的区别在于不关心旧的值是什么



## 10.静态资源

设计策略

1. 添加在Server里面，作为一个暴露的方法，对静态文件进行一些操作
2. 作为一个handler，设计静态路由的route



## 11.文件操作

- 读文件
  - ReadFile
  - Open/OpenFile

openFile可以设计一些权限

~~~go
func TestFile() {
	file, _ := os.OpenFile("E:/Code/GoCode/my.txt", os.O_APPEND, fs.ModeAppend)
	_, _ = file.WriteString("hello,")
}
~~~

> flag：打开文件的模式
>
> - os.O_WRONLY 只写
> - os.O_CREATE 创建文件
> - os.O_RDONLY 只读
> - os.O_RDWR 读写
> - os.O_TRUNC 清空
> - os.O_APPEND 追加

**注意相对路径定位是相对于当前工作目录，用`os.Getwd`查看**



## 12.Options模式

Go没有构造函数，也没有方法重载，所以使用Option模式设计一个Newxx方法。偶尔考虑Builder模式

~~~go
package main

type User struct {
	ID      string
	Name    string
	Address string
}

type Option func(u *User)

func newUser(ID, Name string, options ...Option) *User {
	u := &User{
		ID:   ID,
		Name: Name,
	}
	for _, o := range options {
		o(u)
	}

	return u
}

// withxx option命名
func WithAddress(address string) Option {
	return func(u *User) {
		u.Address = address
	}
}

~~~

## 13.复用Context

![image-20230503162638282](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230503162638282.png)

使用步骤

```go
package main

import "sync"

type User struct {
	Name string
	Age  int
}

// 加上Reset方法，用于重置对象
func (u *User) Reset(name string, age int) {
	u.Age = age
	u.Name = name
}

func main() {
    // 定义同步池
	pool := sync.Pool{
		New: func() interface{} {
			return &User{}
		},
	}
    // Get获取内容
	user := pool.Get().(*User)
    
    // 最后放回去
	defer pool.Put(user)
    // 重置user，然后进行业务处理
	user.Reset("zqy", 12)
    ...
}

```

## 14 接口实现的注册与查找

一般使用map进行索引与接口的映射

![image-20230503164414552](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230503164414552.png)

通过`import _ xxx文件`，进行init的调用然后对接口进行注册



## 15.泛型

~~~go
package main

import "fmt"

type Numeric interface {
	int | int64 | int32 | float32
}

type StringConstraint interface {
	string
}

type Constraint struct {
	Name string
}
// 可以to
type AllConstraint interface {
	StringConstraint | Numeric
}
type Constraint[T AllConstraint] struct {
	Name string
}
func Test[T AllConstraint]() {
	
}

type Sub[T any] struct {
	Name T
}
// 这种传递貌似没意义
type Parent[T any] struct {
	Sub[T]
}

// 约束必须是接口
//func Genertic[T Constraint](t T) {
//	fmt.Printf(t.Name)
//}

// 多个约束是 | ,而不是 &
func Sum[T Numeric | StringConstraint](values []T) T {
	var res T
	for _, val := range values {
		res += val
	}
	return res
}

// map key不能是any，而是comparable可比的
func PutIfAbsent[K comparable, T any](m map[K]T) {

}

func main() {
	strings := []string{"hello", "world"}
	sum := Sum(strings)
	fmt.Println(sum)
	//p := &Parent[string]{
	//	Sub[string]{
	//		Name: "hello",
	//	},
	//}
	//fmt.Println(p.Sub.Name)
}

~~~


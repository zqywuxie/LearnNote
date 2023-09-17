# YunGo(后续进行文档的优化)

> 适于初学者学习一站式Go框架
>
> 代码地址[zqywuxie/LearnNote at code (github.com)](https://github.com/zqywuxie/LearnNote/tree/code)
>
> author : zqy 
>
> updateTime : 2023年9月17日 15点28分
> - 自定义web框架（完结）
> - 自定义ORM框架（进行中)

# 自定义Web框架

- Route.go 
    - 路由树的查找和注册（静态匹配，通配符匹配，参数路径匹配）

- Server.go 
    - http.ListenAndServe自定义

- middleware
    - accesslog 日志打印
    - tracing（openTelemetry整合zipkin) 链路追踪
    - errhdl 错误信息的定义
    - Prometheus 整合Prometheus进行系统监控
    - recover panic的恢复
- Render 页面的渲染
- file.go 文件的上传下载
- session的操作
  - 支持redis的存储
  - 本地内存的存储
  - cookie的存储管理
- Context.go
  - 对context.Context进行二次封装




路由树的总结

注意事项

- 已经注册的路由，无法被覆盖，提示冲突
- path必须以/ 开始，结尾不能有/，中间不允许连续的/
- 同一个位置注册不同的参数路径,eg：/user/:username,/user/:id冲突
- 不能同一个位置同时注册参数路径和通配符路径
- 同名路径参数，在路由匹配时值被覆盖 /user/:id/abc/:id，/user/23/ab/34，最后id是34

为什么使用panic而不是error

- error需要用户去解决，麻烦
- 一般来讲，用户需要注册好路由后才能正常启动服务器，所以panic可以避免服务器启动后再发现错误。

路由树是线程安全的吗

不是线程安全的，因为map只是简单的类型，而不是sync.map，操作也没有进行锁的管理。

但我们启动服务器之前是依次注册路由，后续只会产生并发读是没有危险的。并发读写才会导致线程不安全

面试要点

- 性能受到什么影响:树的高度和宽度(如果子节点树用map，倒不用考虑，查找很快)

优化：

功能需求：

- 路由打印功能
- 路由测试功能，在不启动服务器时候，判断某个路径是否匹配某个路由

非功能需求：

- 性能：查找路由性能必须好，提供基准测试数据
- 拓展性：不支持开发者自定义路由匹配逻辑
- 安全性：检测开发者注册的路由，不符合条件或存在冲突就报错
- 文档：提供详细文档和错误排查手册



# 自定义ORM框架

## 需求背景

使用原生SQL处理方法问题：

- 手写sql：
  - 容易出错
  - 难以重构
- 手动处理结果集
  - 样板代码较为麻烦
  - 将时间花费在不必要的模板代码中

## 其他框架分析

### Beego的ORM

需要注册模型、驱动以及DB

#### 元数据

对模型的描述，在beego里面分成

modleInfo(对整个表的介绍) -> fields (表中数据的介绍)-> fieldInfo(单个字段的介绍)

**实际结构体和真实数据库之间的连接**

#### 查询接口

DQL(封装查询)

QueryBuilder(查询构建-自定义构建查询语句)

QuerySeter(一些选择，filter...)

#### 事务接口

- 普通的Begin,commit,rollback
- 闭包形式的Doxxx (进行封装,只需要关注于业务逻辑)
  - 通过txOrm进行处理



### GORM

#### 元数据

只有 Schema(数据库) -> field (子弹)两级

#### 查询接口

设计思路是将构建SQL的过程分装成单一组件,然后进行组装使用

#### 事务接口

- Begin,Commit,Rollback
- 闭包接口
- **savepoint (类似快照)**



### Ent

代码生成技术

## ORM框架

orm(object relation mapping):对象关系映射

帮助对象从对象到SQL,以及结果集到对象的映射工具

### 对象到SQL

![image-20230605152337836](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605152337836.png)

对象输入一个对象,能够产生对应的SQL;比如使用对象生成SQL于进行查询



### 结果集到对象

查询得到的结果集进行映射到结构体上



### 功能点

![image-20230605152144671](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605152144671.png)

SQL：必须要支持的就是增删改查，DDL 一般是作为 一个扩展功能，或者作为一个工具来提供。
映射：将结果集封装成对象，性能瓶颈。

事务：主要在于维护好事务状态。

元数据：SQL 和映射两个部分的基石。

AOP：处理横向关注点。
关联关系：部分 ORM 框架会提供，性价比低。

方言：兼容不同的数据库，至少要兼容 MySQL、 SQLite、 PostgreSQL。
其它：……



## select语句

### Beego的设计

1. 通过传入参数，指定的列，然后进行拼接sql语句

Beego 的这种构造形式是为了满足这些方法的需求。

- 优点：对于用户来说，这些 API 极其简单，使用无脑

- 缺点：代码耦合性强，扩展性差
  •   耦合性强： SQL 的构造和 SQL 的执行以及结果集的处理完全混在一起，职责不清，界限不明
  •   扩展性差： SELECT 是一个语法形式非常复杂的语句，这种处理方式难以支持完整的语法

2. 通过QueryBuilder

将SQL构建的各个部分进行分割成组件，最后进行组合使用
这种设计则具有很强的扩展性，而且是一个完全独立的组 件，和执行查询、处理结果集都解耦了。
但是 Beego 的这个实现有一个特性：用户必须完全按照 SQL 语句的顺序来调用这些方法。
这个特性优劣参半：
•   好处是用户完全明白他们在干什么
•   缺点就是约束太强，在一些灵活构建 SQL 的需求中使用 起来比较困难
你们可以注意到，后期出现的 ORM 框架大多都是类似设计。



### GORM设计

GORM  的完整机制看起来是稍微有点绕的。简单来说它有四 个主要抽象：
•   Builder：提供了最基本的构造方法
•   **Expression**：表达式，表达式和表达式可以组合成复合表 达式
•   Clause：按照特定需要组合而成的 SQL 的一个部分（from xxx，select xxx，where xxx 都称为一个clause）
•   Interface：构造它自身，以及和其它部分 Clause 组合
所以 GORM 的核心理念是：我不知道怎么构造 SQL，但是 你们知道。
这里的你们，就是指实现接口的人。
也可以认为 SQL 的不同部分分开构造，再组合在一起（有点像微服务的设计理念）。

![image-20230605154714481](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605154714481.png)

### Ent 的构造形式

Ent 是最为纯粹的 Builder 模式，可以 称为经典 Builder 模式。
和 Beego 比起来，**它不要求调用顺 序**；

和 GORM 比起来，它没有复杂的接 口机制，虽然也因此丧失了 GORM 的灵 活性。
你在设计业务系统、中间件的时候都要 平衡扩展性和系统复杂度的关系。
**往往高扩展性带来的就是复杂的接口机 制。任何非功能特性都是有代价的。**

### 定义核心接口

1. ORM 起步 —— 定义核心接口
   设计风格一：定义接口叫做 Orm，对应的是 Beego 的第一种 SQL 构造形式。
   •   Orm 的实例应该是**无状态的（安全的，没有数据存储）**、可以被复用的
   •   Orm 接口无法使用泛型
   •   缺点：Orm 是**大而全的接口**，但凡有新的需求，就往 Orm 里面加方法

![image-20230605155514216](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605155514216.png)

2. 设计风格二：大一统的 Query 接口，增删改查的方法都放在一起， Builder 模式。**可以使用泛型**

缺点：违背了单一职责接口的原则

![image-20230605155507538](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605155507538.png)

3. 设计风格三：直接定义 Selector 接口，需要构造复杂查询就往里面加方法。类似的还有 Deleter 接口、Updater 接口和 Inserter 接口。**单一职责的 Builder 模式。**

![image-20230605155544898](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605155544898.png)

4. 设计风格四：只定义 Builder 模式的终结方法。依旧是 Builder 模式。

![image-20230605160051907](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605160051907.png)



综上采用设计风格四

• 使用 Builder 模式：SQL 里面的不同语句有不同的实现。
•   使用泛型做类型约束：例如 SELECT 语句和 INSERT 语 句。
•   额外引入一个 QueryBuilder：作为**构建 SQL 这一个单独 步骤的顶级抽象**。

![image-20230605160154583](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605160154583.png)



```go
// @Author: zqy
// @File: types.go
// @Date: 2023/6/5 16:05
// @Description orm核心接口和模块

package customize

import (
	"context"
	"database/sql"
)

// Selector 单一查询接口
type Selector[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) (*T[], error)
}

// Executor 执行接口，用于update insert delete
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

// Query SQL语句参数
type Query struct {
	SQL  string
	args []any
}

// QueryBuilder SQL构建
type QueryBuilder interface {

	// Build ，Query 也可以，返回指针方便可以进行修改
	Build() (*Query, error)
}

```



```go
func (s *Selector[T]) Build() (*Query, error) {
	var sb strings.Builder
	// 处理空格问题
	sb.WriteString("SELECT * FROM ")
	// 通过反射获得表明，泛型的名称。默认使用结构体名作为表名
	// 对于带db的参数
	// 1. 让用户自己加入`
	// 2. 开发者自行切割
	sb.WriteByte('`')

	if s.table == "" {
		var t T
		tableName := reflect.TypeOf(t).Name()
		sb.WriteString(tableName)
	} else {
		segs := strings.Split(s.table, ".")
		sb.WriteString(segs[0])
		sb.WriteString("`.`")
		sb.WriteString(segs[1])

	}
	sb.WriteByte('`')

	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		args: nil,
	}, nil
}
```



### 语句规范

目前我们不考虑多方言兼容，只考虑 MySQL。

MySQL 语法规范中重要的部分：
•   FROM：普通表、子查询、JOIN 查询

```go
// From 兼容传入空字符串，如果加入校验会影响连调功能

func (s *Selector[T]) From(table string) *Selector[T] {

	s.table = table
	return s
}
```



•   WHERE：各种查询条件，以及由 AND、OR、NOT 混合在一起的复杂查询条件

设计一

![image-20230605205844507](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605205844507.png)

> 优点：
> •   API 简单明了，实现简单
> •   非常灵活
> 缺点：
> •   缺乏校验，用户容易写错，例如写错字段名，漏了括号等
> •   args 作为不定参数，容易误用切片

设计二

用户指定用于 WHERE 的列名，例如 Where(id, name) 则生成 WHERE id = ? and name =? 这种查询条件。
对 And、Or、Not 难以支持。也不方便进行设计`>/</>=/<=`
实际上，日常使用的 SQL 语句，最复杂的就是这个 WHERE，这种设计方案是绝对难以满足 的。

![image-20230605210112826](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605210112826.png)



设计三

使用gorm的设计理念，每个比较符号都有一个实现例如:Eq,IN

都被认为是一个Expression的集合。

> **Expression的抽象**，考虑and，or可以设计为二叉树
>
> 这种设计不同于 GORM 的设计，GORM 的设计构成一颗多叉树（切片代表了多叉树）。

![image-20230605213021052](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230605213021052.png)

•   HAVING：可以使用 WHERE 里面的条件，以及使用 聚合函数的条件
•   ORDER BY
•   GROUP BY
•   LIMIT 和 OFFSET



## 元数据

### 校验问题

当用户输入错误，该如何解决

- 方法一：**不做任何校验**，当SQL语句在数据库执行时返回错误
- 方法二：**尽早进行校验**
  - 方便测试，在不连接上数据库时进行校验
  - 尽早发现错误（在编写时就进行返回）

![image-20230606220200614](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230606220200614.png)



### 元数据概念

> ORM 框架需要解析模型以获得模型的元数据，这 些元数据将被用于构建 SQL、执行校验，以及用 于处理结果集。
>
> 模型：一般是指**对应到数据库表的 Go 结构体定义**，也被称为 Schema、Table 等。

![image-20230606220310071](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/05/image-20230606220310071.png)



> 设计总结：
> •   模型（表的抽象）：对应的表名、索引、主键、关联关系
> • 列：列名、Go 类型、数据库类型、是否主键、是否外键……

但是一开始结构体的字段都是慢慢演化，根据需求进行不断添加

![image-20230606220913209](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/06/image-20230606220913209.png)

目前进行设计简单的模型，使用**反射**进行获得

> 反射的相关 API 都在 reflect 包，最核心的两 个：
> •   reflect.Value：用于操作值，部分值是可以 被反射修改的
> •   reflect.Type：用于操作类信息，类信息是 **只能读取**,无法改变类型信息
>
> reflect.Type 可以通过 reflect.Value 得到，但 是反过来则不行。

![image-20230606230703189](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/06/06/image-20230606230703189.png)

## 反射

### 读写字段

```go
typeOf := reflect.TypeOf(entity) //从一个任何非接口类型的值创建一个 reflect.Type 值
valueOf := reflect.ValueOf(entity) //返回的是代表着此接口值的动态值的一个 reflect.Value 值
valueOf.IsZero() // 判断值是否为零值


// 层层解引用
// for死循环一直深入指针，直到指向内容
for typeOf.Kind() == reflect.Pointer {
    // Elem 获得切片等内部值，或指针指向的值
    typeOf = typeOf.Elem()
    valueOf = valueOf.Elem()
}

typeOf.NumField() //得到字段数

for i := 0; i < field; i++ {
    fieldType := typeOf.Field(i)
    fieldValue := valueOf.Field(i)
    // IsExported 判断这个字段是否为私有的
    if fieldType.IsExported() {
        // 通过Interface获得最终值
        res[fieldType.Name] = fieldValue.Interface()
    } else {
        res[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
    }
}


// 类型只读
// 字段值可读写，通过CanSet查看是否可写
if !fieldByName.CanSet() {
    return errors.New("该字段不可被修改")
}

```

### 方法

```go
typ := reflect.TypeOf(entity)
// 获得定义到结构体上的，指针上的得不到
numMethod := typ.NumMethod()


method := typ.Method(i) //获得实体类的某个方法
Func := method.Func // .Func
numIn := Func.Type().NumIn() // 方法接收参数数
numOut := Func.Type().NumOut() // 方法返回参数树

// 对于结构体实现接口，参数是结构体本身,所以第一个参数是结构体本身
// 否则下面只会遍历接入参数
InputValues = append(InputValues, reflect.ValueOf(entity))
InputArgs = append(InputArgs, reflect.TypeOf(entity))

// 调用该方法
resValues := Func.Call(InputValues)

```



### 遍历

```go
// 数组与切片遍历方式一样
func IterateArrayOrSlice(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		res = append(res, val.Index(i).Interface())
	}
	return res, nil

}



// 遍历Map，通过MapKeys，MapIndex获取值
func IterateMap(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	entityType := reflect.TypeOf(entity)
	if reflect.Map != entityType.Kind() {
		return nil, nil, errors.New("非Map")
	}
	MapKeys := make([]any, 0, val.Len())
	MapValues := make([]any, 0, val.Len())
	Keys := val.MapKeys()
	for _, key := range Keys {
		MapKeys = append(MapKeys, key.Interface())
		MapValues = append(MapValues, val.MapIndex(key).Interface())
	}
	return MapKeys, MapValues, nil
}
```



### 解析模型

**通过结构体来映射数据库**

```go
func ParseModel(val any) (*Model, error) {
    types := reflect.TypeOf(val)
    if types.Kind() != reflect.Pointer && types.Kind() != reflect.Struct {
        return nil, internal.ErrModelTypeSelect
    }
    // 一级指针获得指向的内容
    types = types.Elem()
    
    //结构体字段表字段
    numField := types.NumField()
    filedMap := make(map[string]*Filed, numField)
    for i := 0; i < numField; i++ {
        field := types.Field(i)
        filedMap[field.Name] = &Filed{ColName: underscoreName(field.Name)}
    }
    return &Model{
        
        // underscoreName 进行驼峰转下划线
        //FirstName -> first_name
        TableName: underscoreName(types.Name()),
        FiledMap:  filedMap,
    }, nil
}



// 驼峰转字符串
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
        // unicode.IsUpper(runne)
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)
}

```



### 利用元数据改造selector

使用解析模型方法，对之前使用结构体名和字段名进行更新

### delete

仿照select将内容进行修改

## 结果集处理

### Go对齐原则

![在这里插入图片描述](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/09/017/03dd4a0a58674fb8baf39be6100a78a4.png)

这么设计的目的，是减少 CPU 访问内存的次数，加大 CPU 访问内存的吞吐量。比如同样读取 8 个字节的数据，一次读取 4 个字节那么只需要读取 2 次。`CPU 始终以字长访问内存`。

32位一次读4个字节(1字长)，64位则是8个字节

```go
func PrintFiledOffset(entity any) {
	typeOf := reflect.TypeOf(entity)
	numField := typeOf.NumField()
	for i := 0; i < numField; i++ {
		filed := typeOf.Field(i)
		fmt.Println(filed.Offset)
	}
}
type User struct {
	Name string //16
	Age  int32 // 4
	//Alias int32 
	Hello string 
} 
// 0 16 24，因为一次性读8个字节，int32就剩余4个字节，但是下面string不够，所以就跳过int32读取的字节再读一个字长


type User struct {
	Name string //16
	Age  int32 // 4
	Alias int32 
	Hello string 
} 

// 16 20 24 Alias只需要4个字节，那么就直接使用Age剩余的字节了
```



### Unsafe使用

#### 读写值

```go
type UnsafeAccessor struct {
	fields  map[string]FieldMata
	address unsafe.Pointer
}

type FieldMata struct {
    // 离起始点的偏移量
	Offset uintptr
	Type   reflect.Type
}

```

**unsafe.pointer和uintptr**

> pointer是Go层面上的指针，uintptr只是一个指向地址的值
>
> 如果遇到GC处理后复活，那么地址就会进行改变，此时uintptr就会报错
>
> 但是pointer会自己维护好

为什么unsafe比反射高效，因为反射封装了unsafe，所以直接使用unsafe要省去封装的过程。unsafe更倾向于直击底层

```go
func NewUnsafeAccessor(entity any) UnsafeAccessor {
	typeof := reflect.TypeOf(entity)
	typeof = typeof.Elem()
	numField := typeof.NumField()
	filedMatas := make(map[string]FieldMata, numField)
	for i := 0; i < numField; i++ {
		field := typeof.Field(i)
       
		filedMatas[field.Name] = FieldMata{Offset: field.Offset, Type: field.Type}
	}
	valueOf := reflect.ValueOf(entity)
	return UnsafeAccessor{
		fields: filedMatas,
		//valueOf.UnsafeAddr() 不用这个是防止被GC等干扰导致出错
		// 使用指针就算地址出错，但还是会指向
		address: valueOf.UnsafePointer(),
	}
}
```



```go
func (u UnsafeAccessor) Field(filed string) (any, error) {
	val, ok := u.fields[filed]
	if !ok {
		return nil, errors.New("未知字段")
	}
	// 获得了当前地址
	cur := unsafe.Pointer(val.Offset + uintptr(u.address))

	//return *(*int)(cur), nil
	// 一般来讲不会直到值的确切类型 通过上面的反射拿到type
	//reflect.New/NewAt 创建指针 所以使用Elem
	return reflect.NewAt(val.Type, cur).Elem().Interface(), nil
}

func (u UnsafeAccessor) SetField(filed string, value any) error {
	val, ok := u.fields[filed]
	if !ok {
		return errors.New("未知字段")
	}
	// 获得了当前地址
	cur := unsafe.Pointer(val.Offset + uintptr(u.address))

	//*(*int)(cur) = value.(int)

	reflect.NewAt(val.Type, cur).Elem().Set(reflect.ValueOf(value))

	//return *(*int)(cur), nil
	// 一般来讲不会直到值的确切类型 通过上面的反射拿到type
	//reflect.New/NewAt 创建指针 所以使用Elem
	return nil
}
```


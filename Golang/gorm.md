# gorm

## 入门

### 概述

**全功能ORM库，方便操作数据库。**

#### 安装

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
```



#### 模型定义

~~~go
GORM 定义一个 gorm.Model 结构体，其包括字段 ID、CreatedAt、UpdatedAt、DeletedAt

// gorm.Model 的定义
type Model struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}
~~~

- ID默认主键值
- CreateAt 追踪创建时间
- UpdateAt 追踪更新时间
- DeleteAt 逻辑删除



#### 连接到数据库

##### Mysql

~~~go
import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)


dsn := "user:pass.@tcp(127.0.0.1:3306)/studb?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
~~~

> 如果正确处理`time.Time`，需要带上`parseTime`的参数，支持完整的UTF-8编码，将charset=utf8改为uft8mb4

##### 自定义驱动

~~~go
import (
_ "github.com/go-sql-driver/mysql"
"gorm.io/driver/mysql"
"gorm.io/gorm"
"time"
)

dsn := "root:wszqy123.@tcp(127.0.0.1:3306)/studb?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.New(mysql.Config{
    DriverName: "mysql",
    DSN:        dsn,
}), &gorm.Config{})
~~~



#### 连接池

~~~go
sqlDB, err := db.DB()
// 设置最大空闲连接
sqlDB.SetMaxIdleConns(10)
//设置连接存活时间
sqlDB.SetConnMaxLifetime(time.Hour)
//设置连接最大值
sqlDB.SetMaxOpenConns(100)
~~~



#### ID作为主键

GORM默认使用`ID`作为主键名

~~~go
type User struct {
  ID   string // 字段名 `ID` 将被作为默认的主键名
}

// 设置字段 `AnimalID` 为默认主键
type Animal struct {
  AnimalID int64 `gorm:"primary_key"`
  Name     string
  Age      int64
}
~~~



#### 表名

默认为结构体名称的复数形式

~~~go
type User struct {} // 默认的表名是 `users`

// 设置 `User` 的表名为 `profiles`
func (User) TableName() string {
  return "profiles"
}
~~~



## CRUD接口

### 创建

#### 创建记录

~~~go
user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

result := db.Create(&user) // 通过数据的指针来创建

user.ID             // 返回插入数据的主键
result.Error        // 返回 error
result.RowsAffected // 返回插入记录的条数
~~~

还可以创建多个数据

~~~go
users := []*User{
    User{Name: "Jinzhu", Age: 18, Birthday: time.Now()},
    User{Name: "Jackson", Age: 19, Birthday: time.Now()},
}

result := db.Create(users) // pass a slice to insert multiple row

result.Error        // returns error
result.RowsAffected // returns inserted records count
~~~



#### 指定字段创建记录

~~~go
db.Select("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`name`,`age`,`created_at`) VALUES ("jinzhu", 18, "2020-07-04 11:05:21.775"）
~~~

**省略某字段创建记录**

不添加Name，Age，CreatedAt字段

~~~go
db.Omit("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`birthday`,`updated_at`) VALUES ("2020-01-01 00:00:00.000", "2020-07-04 11:05:21.775"
~~~



#### 创建钩子

~~~go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    // 创建前判断字段内容
    if u.Role == "admin" {
        return errors.New("invalid role")
    }
    return
}
~~~

可以跳过钩子

~~~go
DB.Session(&gorm.Session{SkipHooks: true}).Create(&user)

DB.Session(&gorm.Session{SkipHooks: true}).Create(&users)

DB.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(users, 100)

~~~

gorm 的 CreateInBatches 方法可以将一批记录分批创建到数据库中，以减少内存使用和提高性能。



#### 根据Map创建

~~~go
db.Model(&User{}).Create(map[string]interface{}{
    "Name": "Zhangsan", "Age": 18, "Email": "xx@qq.com",
})
~~~



#### 默认值

~~~go
Age int64 `form:"Age" gorm:"default:12"`
~~~



### 查询

#### 检索单个对象

~~~go
// 获取第一条记录（主键升序）
db.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1;

// 获取一条记录，没有指定排序字段
db.Take(&user)
// SELECT * FROM users LIMIT 1;

// 获取最后一条记录（主键降序）
db.Last(&user)
// SELECT * FROM users ORDER BY id DESC LIMIT 1;

result := db.First(&user)
result.RowsAffected // 返回找到的记录数
result.Error        // returns error or nil

// 检查 ErrRecordNotFound 错误
errors.Is(result.Error, gorm.ErrRecordNotFound)

~~~

> 如果你想避免`ErrRecordNotFound`错误，你可以使用`Find`，比如`db.Limit(1).Find(&user)`，`Find`方法可以接受struct和slice的数据。

~~~go
var user []User
db.Limit(2).Find(&user)
fmt.Println(user) //查找两条数据
~~~



> 对单个对象使用`Find`而不带limit，`db.Find(&user)`**将会查询整个表**并且只返回第一个对象，这是性能不高并且不确定的。



根据主键查询第一条和最后一条记录，**仅当指向目标结构的指针作为参数传递给方法或使用指定模型时**，它们才有效。此外，如果没有为相关模型定义主键，则模型将按第一个字段排序。例如：`First``Last``db.Model()`

~~~go
var user User
var users []User


db.First(&user)
// SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1


result := map[string]interface{}{}
db.Model(&User{}).First(&result)
// SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1

// doesn't work
result := map[string]interface{}{}
db.Table("users").First(&result)

// works with Take
result := map[string]interface{}{}
db.Table("users").Take(&result)

// no primary key defined, results will be ordered by first field (i.e., `Code`)
//会按照结构体第一个参数进行排序
type Language struct {
  Code string
  Name string
}
db.First(&Language{})
// SELECT * FROM `languages` ORDER BY `languages`.`code` LIMIT 1
~~~



#### 根据字段检索

也是根据结构体的第一个参数进行查询（数字)

~~~go
db.First(&user, 10)
// SELECT * FROM users WHERE id = 10;

db.First(&user, "10")
// SELECT * FROM users WHERE id = 10;

db.Find(&users, []int{1,2,3})
// SELECT * FROM users WHERE id IN (1,2,3);

~~~

如果字段是字符串，查询方式如下编写

~~~go
db.Debug().First(&user, "Name=?", "hello")
//SELECT * FROM `user` WHERE Name='hello' AND `user`.`DeletedAt` IS NULL ORDER BY `user`.`id` LIMIT 1
~~~



当目标对象有一个主键值时，将使用主键构建查询条件，

~~~go
var user = User{ID: 10}
db.First(&user)
// SELECT * FROM users WHERE id = 10;

var result User
db.Model(User{ID: 10}).First(&result)
// SELECT * FROM users WHERE id = 10;
~~~



#### 检索全部对象

~~~go
var user []User
result := db.Find(&user)
fmt.Println(result.RowsAffected) //返回数据数
result.Error
~~~



#### 条件查询

~~~go
// 获得匹配信息的第一条
db.Where("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

// <> 不匹配
db.Where("name <> ?", "jinzhu").Find(&users)
// SELECT * FROM users WHERE name <> 'jinzhu';

// IN
db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
// SELECT * FROM users WHERE name IN ('jinzhu','jinzhu 2');

// LIKE 注意自己添加% %
db.Where("name LIKE ?", "%jin%").Find(&users)
// SELECT * FROM users WHERE name LIKE '%jin%';

// AND
db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;

// Time
db.Where("updated_at > ?", lastWeek).Find(&users)
// SELECT * FROM users WHERE updated_at > '2000-01-01 00:00:00';

// BETWEEN
db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
// SELECT * FROM users WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08
~~~

> 如果已设置对象的主键，则条件查询不会涵盖主键的值，而是将其用作“and”条件。例如：
>
> ```go
> var user = User{ID: 10}
> db.Where("id = ?", 20).First(&user)
> // SELECT * FROM users WHERE id = 10 and id = 20 ORDER BY id ASC LIMIT 1
> ```
>
> 此查询将给出错误。因此，在要使用变量（例如从数据库中获取新值）之前，将主键属性设置为 nil。`record not found``id``user`

#### Struct & Map 条件

~~~go
// Struct
db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;

// Map
db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;

// Slice of primary keys
db.Where([]int64{20, 21, 22}).Find(&users)
// SELECT * FROM users WHERE id IN (20, 21, 22);
~~~

> `db.Where([]int64{20, 21, 22}).Find(&users)`
>
> 也是根据结构体第一个字段进行查询；
>
> **注意**使用 struct 查询时，GORM 将仅使用非零字段进行查询，这意味着如果您的字段的值为 、 或其他[零值](https://tour.golang.org/basics/12)，则不会用于构建查询条件，例如：`0``''``false`

~~~go
db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu";
~~~

#### 指定结构体查询字段

~~~go
db.Debug().Where(&User{Name: "hello", Age: 12}, "Name", "Age").First(&user)
//SELECT * FROM `user` WHERE `user`.`name` = 'hello' AND `user`.`age` = 12 AND `user`.`DeletedAt` IS NULL ORDER BY `user`.`id` LIMIT 1
~~~

指定结构体的`Name`,`Age`查询，参数名**首字母忽略大小写**

#### 内联条件

查询条件可以内联到方法

~~~go
// Get by primary key if it were a non-integer type
db.First(&user, "id = ?", "string_primary_key")
// SELECT * FROM users WHERE id = 'string_primary_key';

// Plain SQL
db.Find(&user, "name = ?", "jinzhu")
// SELECT * FROM users WHERE name = "jinzhu";

db.Find(&users, "name <> ? AND age > ?", "jinzhu", 20)
// SELECT * FROM users WHERE name <> "jinzhu" AND age > 20;

// Struct
db.Find(&users, User{Age: 20})
// SELECT * FROM users WHERE age = 20;

// Map
db.Find(&users, map[string]interface{}{"age": 20})
// SELECT * FROM users WHERE age = 20;
~~~



#### Not / Or条件

~~~go
db.Not("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE NOT name = "jinzhu" ORDER BY id LIMIT 1;

// Not In
db.Not(map[string]interface{}{"name": []string{"jinzhu", "jinzhu 2"}}).Find(&users)
// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");

// Struct
db.Not(User{Name: "jinzhu", Age: 18}).First(&user)
// SELECT * FROM users WHERE name <> "jinzhu" AND age <> 18 ORDER BY id LIMIT 1;

// Not In slice of primary keys
db.Not([]int64{1,2,3}).First(&user)
// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;

============
db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';

// Struct
db.Where("name = 'jinzhu'").Or(User{Name: "jinzhu 2", Age: 18}).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);

// Map
db.Where("name = 'jinzhu'").Or(map[string]interface{}{"name": "jinzhu 2", "age": 18}).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);
~~~



#### 返回特定的字段

默认首字母小写查询

~~~go
db.Select("name", "age").Find(&users)
// SELECT name, age FROM users;

db.Select([]string{"name", "age"}).Find(&users)
// SELECT name, age FROM users;

db.Table("users").Select("COALESCE(age,?)", 42).Rows()
// SELECT COALESCE(age,'42') FROM users;
//coalesce 将空参数设置为42
~~~



#### 排序

`ASC` 顺序，`DESC` 逆序

~~~go
db.Order("age desc, name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;

// Multiple orders
db.Order("age desc").Order("name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;

db.Clauses(clause.OrderBy{
  Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{[]int{1, 2, 3}}, WithoutParentheses: true},
}).Find(&User{})
// SELECT * FROM users ORDER BY FIELD(id,1,2,3)
~~~



#### Limit & Offset

~~~go
db.Limit(3).Find(&users)
// SELECT * FROM users LIMIT 3;

// Cancel limit condition with -1
// -1 取消limit条件
db.Limit(10).Find(&users1).Limit(-1).Find(&users2)
// SELECT * FROM users LIMIT 10; (users1)
// SELECT * FROM users; (users2)

db.Offset(3).Find(&users)
// SELECT * FROM users OFFSET 3;

db.Limit(10).Offset(5).Find(&users)
// SELECT * FROM users OFFSET 5 LIMIT 10;

// Cancel offset condition with -1
//-1 取消offset条件
db.Offset(10).Find(&users1).Offset(-1).Find(&users2)
// SELECT * FROM users OFFSET 10; (users1)
// SELECT * FROM users; (users2)
~~~

#### Group By & Having

~~~go
type result struct {
Date  time.Time
Total int
}

db.Model(&User{}).Select("name, sum(age) as total").Where("name LIKE ?", "group%").Group("name").First(&result)
// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "group%" GROUP BY `name` LIMIT 1


db.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("name = ?", "group").Find(&result)
// SELECT name, sum(age) as total FROM `users` GROUP BY `name` HAVING name = "group"

rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()
defer rows.Close()
// 读取出结果集的内容
for rows.Next() {
    var total int64
    if err := rows.Scan(&total); err != nil {
        fmt.Println(err)
    }
    fmt.Println(total)
}

// 获得的结果直接赋值给变量
db.Debug().Model(&User{}).Select("sum(age) as total").Scan(&total)


rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
defer rows.Close()
for rows.Next() {
...
}

type Result struct {
Date  time.Time
Total int64
}
db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results）
~~~



#### Distinct

~~~go
db.Distinct("name", "age").Order("name, age desc").Find(&results)
db.Debug().Model(&User{}).Distinct("name").Select("Age,email").Find(&user)
~~~



#### joins

~~~go
type result struct {
Name  string
Email string
}

db.Model(&User{}).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id

rows, err := db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()
for rows.Next() {
...
}

var user User
db.Debug().Table("user").Select("user.name,teacher.Age").Joins("left join teacher on teacher.Age = user.ID").Find(&user)
fmt.Println(user)

db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)

// multiple joins with parameter
db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)
~~~



预加载

~~~go
type User struct {
  gorm.Model
  Name      string
  Email     string
  Addresses []Address // 一个用户有多个地址
}

type Address struct {
  gorm.Model
  UserID    uint
  Street    string
  City      string
  State     string
  ZipCode   string
}

// 使用 Joins 加载所有相关数据
db.Joins("Addresses").Find(&users)

// 使用 Preload 预加载单个关联
db.Preload("Addresses").Find(&users)

~~~

> 在GORM中，Joins和Preload都是进行关联查询的方法，但它们的查询方式有所不同。
>
> **Joins方法会将多个表连接起来，一次性查询所有相关数据，**并返回一个包含所有结果的结果集。这种方式在处理大量数据时可能会导致性能问题，因为它会生成一个包含所有结果的SQL语句。
>
> Preload方法则是通过**分别执行多条SQL语句来获取关联数据**，可以避免Join方法中的性能问题。具体来说，Preload方法会首先查询主表（即调用Preload方法的模型）的数据，然后再执行一条或多条SQL语句来获取与主表关联的其他表的数据。这种方式可以减少Join方法中不必要的数据传输，提高查询性能。另外，Preload还支持使用链式调用语法进行嵌套查询，可以方便地查询主表关联的多个表的数据。
>
> 因此，在处理较小量的数据时，两种方式的性能差异可能不明显，但在处理较大量的数据时，Preload方法往往比Join方法更快且更可靠。但需要注意的是，使用Preload方法时要注意避免N+1查询问题，即在查询关联数据时可能会产生大量额外的SQL查询。

还有联接衍生表

~~~go
type User struct {
    Id  int
    Age int
}

type Order struct {
    UserId     int
    FinishedAt *time.Time
}

query := db.Table("order").Select("MAX(order.finished_at) as latest").Joins("left join user user on order.user_id = user.id").Where("user.age > ?", 18).Group("order.user_id")

// 插入一个表query
db.Model(&Order{}).Joins("join (?) q on order.finished_at = q.latest", query).Scan(&results)
// SELECT `order`.`user_id`,`order`.`finished_at` FROM `order` join (SELECT MAX(order.finished_at) as latest FROM `order` left join user user on order.user_id = user.id WHERE user.age > 18 GROUP BY `order`.`user_id`) q on order.finished_at = q.latest
~~~





#### Scan

将结果扫描到结构中的工作方式类似于我们使用的方式`Find`

```go
type Result struct {
  Name string
  Age  int
}

var result Result
db.Table("users").Select("name", "age").Where("name = ?", "Antonio").Scan(&result)

// Raw SQL
db.Raw("SELECT name, age FROM users WHERE name = ?", "Antonio").Scan(&result)
```
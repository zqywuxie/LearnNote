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



### 高级查询

#### 智能选择字段

gorm可以通过`Select` 方法选择字段，如果Select选择的字段经常使用，那么可以单独提取出结构体进行使用。

~~~go
type User struct {
  ID     uint
  Name   string
  Age    int
  Gender string
  // 假设后面还有几百个字段...
}

type APIUser struct {
  ID   uint
  Name string
}

// 查询时会自动选择 `id`, `name` 字段
db.Model(&User{}).Limit(10).Find(&APIUser{})
// SELECT `id`, `name` FROM `users` LIMIT 10

~~~



#### 加锁Locking

~~~go
db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
// SELECT * FROM `users` FOR UPDATE
//当使用这个指令时，查询结果集中所涉及的所有行都会被锁定，防止其他事务对这些行进行修改，直到当前事务结束。


db.Clauses(clause.Locking{
  Strength: "SHARE",
  Table: clause.Table{Name: clause.CurrentTable},
}).Find(&users)
// SELECT * FROM `users` FOR SHARE OF `users`
// of  针对某张表

db.Clauses(clause.Locking{
  Strength: "UPDATE",
  Options: "NOWAIT",
}).Find(&users)
// SELECT * FROM `users` FOR UPDATE NOWAIT
//它表示在尝试获取行级锁时，如果当前行被其他事务占用，则不会等待，而是立即返回异常。
~~~



#### 子查询

~~~go
db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)
// SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");

subQuery := db.Select("AVG(age)").Where("name LIKE ?", "name%").Table("users")
db.Select("AVG(age) as avgage").Group("name").Having("AVG(age) > (?)", subQuery).Find(&results)
// SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%"
~~~

允许在`Table` 方法通过FROM子句使用子查询

~~~go
db.Table("(?) as u", db.Model(&User{}).Select("name", "age")).Where("age = ?", 18).Find(&User{})
// SELECT * FROM (SELECT `name`,`age` FROM `users`) as u WHERE `age` = 18

subQuery1 := db.Model(&User{}).Select("name")
subQuery2 := db.Model(&Pet{}).Select("name")
db.Table("(?) as u, (?) as p", subQuery1, subQuery2).Find(&User{})
// SELECT * FROM (SELECT `name` FROM `users`) as u, (SELECT `name` FROM `pets`) as p
~~~



####  带多个列的in

~~~go
db.Where("(name, age, role) IN ?", [][]interface{}{{"jinzhu", 18, "admin"}, {"jinzhu2", 19, "user"}}).Find(&users)
// SELECT * FROM users WHERE (name, age, role) IN (("jinzhu", 18, "admin"), ("jinzhu 2", 19, "user"));
~~~



#### Find至map

允许将查找到的结果赋值给map

~~~go
result := map[string]interface{}{}
db.Model(&User{}).First(&result, "id = ?", 1)

var results []map[string]interface{}
db.Table("users").Find(&results)
~~~



#### FirstOrInit

获取第一条匹配的记录，或者根据给定的条件初始化一个实例（仅支持 sturct 和 map 条件）

这里deleted_at ，gorm里面默认是记录删除时间。更改为`1/0`

~~~go
gorm.Model `gorm:"softDelete:flag"`
~~~

但修改后的问题，查询语句中对于删除的语句`deleted_at is NUll` ，会导致查询不出。所以使用下方的方式，添加一个全局作用域方法`Active`

~~~go
func Active() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NULL OR deleted_at = 0")
	}
}
~~~

`Unscoped` 会去除默认作用域，`Scopes`再添加一个作用域

~~~go
db.Debug().Unscoped().Model(&User{}).Scopes(Active()).
    Where("Name = ?", "zqy").FirstOrInit(&user)
~~~



如果没有找到记录，还可以使用`Attrs` 添加字段，但不会用于生成SQL语句

~~~go
// 未找到 user，则根据给定的条件以及 Attrs 初始化 user
db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrInit(&user)
// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
// user -> User{Name: "non_existing", Age: 20}
~~~

#### FirstOrCreate

获取匹配的第一条记录或者根据给定条件创建一条新纪录（仅 struct, map 条件有效），`RowsAffected` 返回创建、更新的记录数

~~~go
// 未找到 User，根据给定条件创建一条新纪录
result := db.FirstOrCreate(&user, User{Name: "non_existing"})
// INSERT INTO "users" (name) VALUES ("non_existing");
// user -> User{ID: 112, Name: "non_existing"}
// result.RowsAffected // => 1
~~~

如果没有找到记录，可以使用包含更多的属性的结构体创建记录，`Attrs` 不会被用于生成查询 SQL 。

~~~go
// 未找到 user，根据条件和 Assign 属性创建记录
db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
// user -> User{ID: 112, Name: "non_existing", Age: 20}
~~~

**不管是否找到记录，`Assign` 都会将属性赋值给 struct，并将结果写回数据库**



#### 迭代

GORM 支持通过行进行迭代

```go
rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Rows()
defer rows.Close()

for rows.Next() {
  var user User
  // ScanRows 方法用于将一行记录扫描至结构体
  db.ScanRows(rows, &user)

  // 业务逻辑...
}
```



#### FindInBatches

用于批量查询并处理记录

```go
// 每次批量处理 100 条
result := db.Where("processed = ?", false).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
  for _, result := range results {
    // 批量处理找到的记录
  }

  tx.Save(&results)

  tx.RowsAffected // 本次批量操作影响的记录数

  batch // Batch 1, 2, 3

  // 如果返回错误会终止后续批量操作
  return nil
})

result.Error // returned error
result.RowsAffected // 整个批量操作影响的记录数
```



#### Pluck

Pluck 用于从数据库查询单个列，并将结果扫描到切片。如果您想要查询多列，您应该使用 `Select` 和 `Scan`

~~~go
var ages []int64
db.Model(&users).Pluck("age", &ages)

var names []string
db.Model(&User{}).Pluck("name", &names)

db.Table("deleted_users").Pluck("name", &names)

// Distinct Pluck,去重
db.Model(&User{}).Distinct().Pluck("Name", &names)
// SELECT DISTINCT `name` FROM `users`

// 超过一列的查询，应该使用 `Scan` 或者 `Find`，例如：
db.Select("name", "age").Scan(&users)
db.Select("name", "age").Find(&users)
~~~



#### Scope

`Scopes` 允许你指定常用的查询，可以在调用方法时引用这些查询

~~~go
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
  return db.Where("amount > ?", 1000)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
  return db.Where("pay_mode_sign = ?", "C")
}

func PaidWithCod(db *gorm.DB) *gorm.DB {
  return db.Where("pay_mode_sign = ?", "C")
}

func OrderStatus(status []string) func (db *gorm.DB) *gorm.DB {
  return func (db *gorm.DB) *gorm.DB {
    return db.Where("status IN (?)", status)
  }
}

db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)
// 查找所有金额大于 1000 的信用卡订单

db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&orders)
// 查找所有金额大于 1000 的货到付款订单

db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
// 查找所有金额大于 1000 且已付款或已发货的订单
~~~



#### Count

**用于获得匹配的记录数**



### 更新

#### save 保存所有字段

`Save` 会保存所有的字段，即使字段是零值

~~~go
db.First(&user)

user.Name = "jinzhu 2"
user.Age = 100
db.Save(&user)
// UPDATE users SET name='jinzhu 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=111;
~~~

如果没有添加主键值，那么save便会执行create方法。添加主键值才会执行更新

#### update 更新单个列

~~~go
// Update with conditions
db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;
~~~

#### updates 更新多个列

```go
// Update attributes with `struct`, will only update non-zero fields
db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;


// Update attributes with `map`
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

```

如果没有主键为非空的记录，gorm就会执行批量更新`Model`

~~~go
// Update with struct
db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
// UPDATE users SET name='hello', age=18 WHERE role = 'admin';
// 就会把role = admin的数据都更新
~~~



#### select 更新选定字段

~~~go
// Select with Map
// User's ID is `111`:
db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello' WHERE id=111;
~~~



#### 阻止全局更新

如果没有条件执行批量更新，gorm不会允许，默认情况返回`ErrMissingWhereClause`。您必须使用某些条件或使用原始 SQL 或启用模式，例如：`AllowGlobalUpdate`

~~~go
db.Model(&User{}).Update("name", "jinzhu").Error // gorm.ErrMissingWhereClause

db.Model(&User{}).Where("1 = 1").Update("name", "jinzhu")
// UPDATE users SET `name` = "jinzhu" WHERE 1=1

db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&User{}).Update("name", "jinzhu")
// UPDATE users SET `name` = "jinzhu"
~~~



#### 更新的记录数

通过`RowsAffectd` 返回更新的记录数

~~~go
// Get updated records count with `RowsAffected`
result := db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
// UPDATE users SET name='hello', age=18 WHERE role = 'admin';

result.RowsAffected // returns updated records count
result.Error        // returns updating error
~~~



#### 高级选项

##### 使用SQL表达式更新

~~~go
// product's ID is `3`
db.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))
// UPDATE "products" SET "price" = price * 2 + 100, "updated_at" = '2013-11-17 21:34:10' WHERE "id" = 3;

db.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})
// UPDATE "products" SET "price" = price * 2 + 100, "updated_at" = '2013-11-17 21:34:10' WHERE "id" = 3;

db.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = 3;

db.Model(&product).Where("quantity > 1").UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = 3 AND quantity > 1;
~~~

还可以使用子查询进行更新

~~~go
db.Model(&user).Update("company_name", db.Model(&Company{}).Select("name").Where("companies.id = users.company_id"))
// UPDATE "users" SET "company_name" = (SELECT name FROM companies WHERE companies.id = users.company_id);

db.Table("users as u").Where("name = ?", "jinzhu").Update("company_name", db.Table("companies as c").Select("name").Where("c.id = u.company_id"))

db.Table("users as u").Where("name = ?", "jinzhu").Updates(map[string]interface{}{"company_name": db.Table("companies as c").Select("name").Where("c.id = u.company_id")})
~~~



##### 不使用Hook和时间跟踪

`UpdateColumn,UpdateColumns` 更新列，用法可与update(s)类似。但是不会使用Hook和时间跟踪(updatedat 不会被更新)

~~~go
//UPDATE `user` SET `age`=18,`name`='1212',`updated_at`='2023-04-28 09:19:07.802' WHERE age = 18 AND `user`.`deleted_at` IS NULL

//可见updated_at更新
db.Debug().Model(&User{}).Where("age = ?", 18).Update("name", "1212") 


// UPDATE `user` SET `name`='1212' WHERE age = 18 AND `user`.`deleted_at` IS NULL
// 不带updated_at
db.Debug().Model(&User{}).Where("age = ?", 18).UpdateColumn("name", "1212")
~~~



##### 检查字段是否有更新

在`BeforeUpdate`的Hook里面使用

~~~go
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
  // if Role changed
    if tx.Statement.Changed("Role") {
    return errors.New("role not allowed to change")
    }

  if tx.Statement.Changed("Name", "Admin") { // if Name or Role changed
    tx.Statement.SetColumn("Age", 18)
  }

  // if any fields changed
    if tx.Statement.Changed() {
        tx.Statement.SetColumn("RefreshedAt", time.Now())
    }
    return nil
}
~~~



##### 在update时 修改值

~~~go
func (user *User) BeforeSave(tx *gorm.DB) (err error) {
  if pw, err := bcrypt.GenerateFromPassword(user.Password, 0); err == nil {
    tx.Statement.SetColumn("EncryptedPassword", pw)
  }

  if tx.Statement.Changed("Code") {
    user.Age += 20
    tx.Statement.SetColumn("Age", user.Age)
  }
}

db.Model(&user).Update("Name", "jinzhu")
~~~



### 删除

#### 删除一条记录

删除对象必须指定主键，否则会触发批量删除

~~~go
// Email 的 ID 是 `10`
db.Delete(&email)
// DELETE from emails where id = 10;

// 带额外条件的删除
db.Where("name = ?", "jinzhu").Delete(&email) 
// DELETE from emails where id = 10 AND name = "jinzhu";
~~~



#### 根据主键删除

~~~go
db.Delete(&User{}, 10)
// DELETE FROM users WHERE id = 10;

db.Delete(&User{}, "10")
// id默认为int类型，转入sql为int
// DELETE FROM users WHERE id = 10;

db.Delete(&users, []int{1,2,3})
// DELETE FROM users WHERE id IN (1,2,3);

~~~



#### 批量删除

**指定值不存在主键，会执行批量删除**

可以将一个主键切片传递给`Delete` 方法，以便更高效的删除数据量大的记录

```go
var users = []User{{ID: 1}, {ID: 2}, {ID: 3}}
db.Delete(&users)
// DELETE FROM users WHERE id IN (1,2,3);

db.Delete(&users, "name LIKE ?", "%jinzhu%")
// DELETE FROM users WHERE name LIKE "%jinzhu%" AND id IN (1,2,3);
```



#### **阻止全局删除**

当你试图执行不带任何条件的批量删除时，GORM将不会运行并返回`ErrMissingWhereClause` 错误

如果一定要这么做，你必须添加一些条件，或者使用原生SQL，或者开启`AllowGlobalUpdate` 模式，如下例：

~~~go
db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})
// DELETE FROM users
~~~



#### 软删除

如果你的模型包含了 `gorm.DeletedAt`字段（该字段也被包含在`gorm.Model`中），那么该模型将会自动获得软删除的能力！

当调用`Delete`时，GORM并不会从数据库中删除该记录，而是将该记录的`DeleteAt`设置为当前时间，而后的一般查询方法将无法查找到此条记录。

##### 删除标志

`gorm.Model`使用`*time.Time`作为`DeletedAt` 的字段类型，不过软删除插件`gorm.io/plugin/soft_delete`同时也提供其他的数据格式支持`

###### 使用 `1` / `0` 作为 删除标志

~~~go
go get gorm.io/plugin/soft_delete
====
DeletedAt soft_delete.DeletedAt `gorm:"softDelete:flag"`
~~~

###### 混合模式

~~~go
type User struct {
  ID        uint
  Name      string
  DeletedAt time.Time
  IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
  // IsDel     soft_delete.DeletedAt `gorm:"softDelete:,DeletedAtField:DeletedAt"` // use `unix second`
  // IsDel     soft_delete.DeletedAt `gorm:"softDelete:nano,DeletedAtField:DeletedAt"` // use `unix nano second`
}

// 查询
SELECT * FROM users WHERE is_del = 0;

// 软删除
UPDATE users SET is_del = 1, deleted_at = /* current unix second */ WHERE ID = 1;
~~~



##### 被软删除记录的操作

`Unscoped`可以用来查询被软删除的记录，也可以用来永久删除匹配的记录

~~~go
db.Unscoped().Where("age = 20").Find(&users)
// SELECT * FROM users WHERE age = 20;

db.Unscoped().Delete(&order)
// DELETE FROM orders WHERE id=10;
~~~



### 原生sql与sql生成器

`Raw`和`Exec` 使用原生sql语句查询

~~~go
var result Result
db.Raw("SELECT id, name, age FROM users WHERE id = ?", 3).Scan(&result)

db.Exec("DROP TABLE users")
db.Exec("UPDATE orders SET shipped_at = ? WHERE id IN ?", time.Now(), []int64{1, 2, 3})
~~~



#### 命名参数

GORM 支持 [`sql.NamedArg`map[string]interface{}{}` 或 struct 形式的命名参数，例如：

~~~go
db.Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(&user)
// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"

db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu2"}).First(&result3)
// SELECT * FROM `users` WHERE name1 = "jinzhu2" OR name2 = "jinzhu2" ORDER BY `users`.`id` LIMIT 1

type NamedArgument struct {
    Name string
    Name2 string
}

db.Raw("SELECT * FROM users WHERE (name1 = @Name AND name3 = @Name) AND name2 = @Name2",
     NamedArgument{Name: "jinzhu", Name2: "jinzhu2"}).Find(&user)
// SELECT * FROM users WHERE (name1 = "jinzhu" AND name3 = "jinzhu") AND name2 = "jinzhu2"
~~~



#### DeyRun模式

不执行情况下生成SQL语句 ，为后续执行做准备

~~~go
stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
stmt.SQL.String() //=> SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
stmt.Vars         //=> []interface{}{1}
~~~



#### Row & Rows

获得`*sql.Row`

~~~go
// 使用 GORM API 构建 SQL
row := db.Table("users").Where("name = ?", "jinzhu").Select("name", "age").Row()
row.Scan(&name, &age)

// 使用原生 SQL
row := db.Raw("select name, age, email from users where name = ?", "jinzhu").Row()
row.Scan(&name, &age, &email)

~~~

~~~go
获取 *sql.Rows 结果

// 使用 GORM API 构建 SQL
rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows()
defer rows.Close()
for rows.Next() {
  rows.Scan(&name, &age, &email)

  // 业务逻辑...
}

// 原生 SQL
rows, err := db.Raw("select name, age, email from users where name = ?", "jinzhu").Rows()
defer rows.Close()
for rows.Next() {
  rows.Scan(&name, &age, &email)

  // 业务逻辑...
}
~~~



#### 将 `sql.Rows` 扫描至 model

使用 `ScanRows` 将一行记录扫描至 struct，例如：

```go
rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows() // (*sql.Rows, error)
defer rows.Close()

var user User
for rows.Next() {
  // ScanRows 将一行扫描至 user
  db.ScanRows(rows, &user)

  // 业务逻辑...
}
```


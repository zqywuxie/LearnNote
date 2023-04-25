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
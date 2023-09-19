// @Author: zqy
// @File: types.go
// @Date: 2023/6/5 16:05
// @Description orm核心接口和模块

package customize

import (
	"context"
	"database/sql"
)

// Queries 单一查询接口
type Queries[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
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

type TableName interface {
	TableName() string
}

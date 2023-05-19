// @Author: zqy
// @File: types.go
// @Date: 2023/5/17 11:21
// @Description todo

package session

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrKeyNotFound     = errors.New("session key not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionNotExist = errors.New("session不存在")
)

// Session /Provider 设置为接口,不自定义session，而是使用gorilla/sessions
// 内存存到context.Context里面
type Session interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value string) error
	//ID  获得session-id
	ID() string
}

type Store interface {
	// Generate 创建session
	Generate(ctx context.Context, id string) (Session, error)

	Remove(ctx context.Context, id string) error

	// Refresh 将内存中的session刷新到存储中
	Refresh(ctx context.Context, id string) error

	Get(ctx context.Context, id string) (Session, error)
}

// Propagator session-id 与 http之间联系
type Propagator interface {
	Inject(id string, w http.ResponseWriter) error
	Extract(w *http.Request) (string, error)
	Remove(w http.ResponseWriter) error
}

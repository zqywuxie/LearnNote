// @Author: zqy
// @File: session.go
// @Date: 2023/5/19 9:21
// @Description 基于内存的session管理

package memory

import (
	"GoCode/web/customize/session"
	"context"
	"fmt"
	cache "github.com/patrickmn/go-cache"
	"sync"
	"time"
)

// ErrKeyNotFound
// sentinel error 预定义错误
// 这种可以设置全局类中，也可以针对某个具体实现具体业务进行设计

type Store struct {
	sessions   *cache.Cache
	expiration time.Duration
	// mutex对于一些操作的安全性
	mutex sync.RWMutex
}

// NewStore 设置过期时间
func NewStore(expiration time.Duration) *Store {
	return &Store{
		sessions:   cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

type Session struct {
	// 这种操作性较强 因为可以控制锁
	// mutex sync.RWMutex
	// values map[string]any

	id     string
	values sync.Map
}

var _ session.Session = &Session{}
var _ session.Store = &Store{}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess := &Session{
		id: id,
	}
	s.sessions.Set(id, sess, s.expiration)
	return sess, nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// 无法判断session是否存在和是否删除
	s.sessions.Delete(id)
	return nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.sessions.Get(id)
	if !ok {
		return fmt.Errorf("session：该id %s 对应的session不存在", id)
	}
	s.sessions.Set(id, val, s.expiration)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	sess, ok := s.sessions.Get(id)
	if !ok {
		return nil, errSessionNotFound
	}

	return sess.(*Session), nil
}
func (s *Session) Get(ctx context.Context, key string) (any, error) {
	value, ok := s.values.Load(key)
	if !ok {
		return nil, fmt.Errorf("%s,key %s", errKeyNotFound, key)
	}
	return value, nil
}

func (s *Session) Set(ctx context.Context, key string, value string) error {
	s.values.Store(key, value)
	return nil
}

// ID
// 关于线程安全的问题，这里在正常业务当中是只读的，并不存在所谓的线程安全
func (s *Session) ID() string {

	return s.id
}

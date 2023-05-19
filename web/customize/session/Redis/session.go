// @Author: zqy
// @File: session.go
// @Date: 2023/5/19 10:26
// @Description 基于Redis的session管理

package Redis

import (
	"GoCode/web/customize/session"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Store struct {
	//map[string]map[string]string
	client     redis.Cmdable
	expiration time.Duration
	prefix     string
}

type StoreOption func(store *Store)

func NewStore(client redis.Cmdable, opts ...StoreOption) *Store {
	res := &Store{
		client:     client,
		expiration: time.Minute * 15,
		prefix:     "sessID",
	}

	for _, opt := range opts {
		opt(res)
	}
	return res
}

func StoreWithPrefixAndExpiration(prefix string, expiration ...time.Duration) StoreOption {
	return func(store *Store) {
		store.prefix = prefix
		if len(expiration) > 0 {
			store.expiration = expiration[0]
		}
	}
}

var _ session.Store = &Store{}
var _ session.Session = &Session{}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {

	//	const lua = `
	//redis.call("hset",KEYS[],ARGV[1],ARGV[2])
	//return redis.call("expire",KEYS[1],ARGV[1])
	//`

	key := PrefixKey(s.prefix, id)
	_, err := s.client.HSet(ctx, key, id, id).Result()
	if err != nil {
		return nil, err
	}
	_, err = s.client.Expire(ctx, key, s.expiration).Result()
	if err != nil {
		return nil, err
	}
	return &Session{
		id:     id,
		client: s.client,
		key:    PrefixKey(s.prefix, id),
	}, nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	cnt, err := s.client.Del(ctx, PrefixKey(s.prefix, id)).Result()
	if err != nil {
		return err
	}
	// 1 有数据
	if cnt != 1 {
		return errors.New("session: 删除session 失败")
	}
	return err
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	ok, err := s.client.Expire(ctx, PrefixKey(s.prefix, id), s.expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("session: redis 未找到该session")
	}
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	// 自由决策是否将session里面的数据拿出来
	// 1.都拿
	// 2.只拿热点数据
	// 3.都不拿（用户获得session后，根据自己需求自己拿）
	key := PrefixKey(s.prefix, id)
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if result != 1 {
		return nil, session.ErrSessionNotFound
	}

	// 关于这里直接创建session，传一个client即可，因为session是存储在redis里面的
	return &Session{
		client: s.client,
		id:     id,
		key:    key,
	}, nil
}

type Session struct {
	client redis.Cmdable
	id     string
	prefix string
	key    string
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	// 不用检测，因为拿不到也得不到数据
	return s.client.HGet(ctx, s.key, key).Result()
}

// Set 需要进行检测，判断是否key还存有数据
// 使用lua脚本，保证操作原子性
func (s *Session) Set(ctx context.Context, key string, value string) error {
	const lua = `
// 判断session是否存在
if redis.call("exists",KEYS[1])
then
	return redis.call("hset",KEYS[1],ARGV[1],ARGV[2])
else
	return -1
end
`
	// 将id进行转化，防止出现重复id
	res, err := s.client.Eval(ctx, lua, []string{s.key}, key, value).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return session.ErrSessionNotExist
	}
	return nil
}

func (s *Session) ID() string {
	return s.id
}

func PrefixKey(prefix, id string) string {
	return fmt.Sprintf("%s-%s", prefix, id)
}

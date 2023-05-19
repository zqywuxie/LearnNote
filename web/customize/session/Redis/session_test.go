// @Author: zqy
// @File: session_test.go.go
// @Date: 2023/5/19 15:04
// @Description todo

package Redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_Generate(t *testing.T) {
	s := newStore()

	ctx := context.Background()
	id := "sess_test_id"
	sess, err := s.Generate(ctx, id)
	require.NoError(t, err)
	//defer s.Remove(ctx, id)
	err = sess.Set(ctx, "key1", "123")
	require.NoError(t, err)
	value, err := sess.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, value, "123")
}

func newStore() *Store {
	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6667",
		Password: "123456",
		DB:       0,
	})
	// 测试连接是否成功
	pong, err := rc.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Redis connected successfully! Response:", pong)
	}
	PrefixOpt := StoreWithPrefixAndExpiration("zqy")
	return NewStore(rc, PrefixOpt)
}

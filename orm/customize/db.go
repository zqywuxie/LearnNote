package customize

import (
	"sync"
)

type DB struct {
	r *register
}

type DBOption func(db *DB)

func NewDB(ops ...DBOption) (*DB, error) {
	res := &DB{
		r: &register{
			models: sync.Map{},
		},
	}
	for _, op := range ops {
		op(res)
	}
	return res, nil
}

// MustNewDB 如果上面不想加error 就需要设置两个方法
func MustNewDB(ops ...DBOption) *DB {
	db, err := NewDB(ops...)
	if err != nil {
		panic(err)
	}
	return db
}

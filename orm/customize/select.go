// @Author: zqy
// @File: select.go
// @Date: 2023/6/5 16:17
// @Description todo

package customize

import (
	"context"
	"strings"
)

type Selector[T any] struct {
	table string
	Builder
	where  []Predicate
	Having []Predicate
	db     *DB
}

func newSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db:      db,
		Builder: Builder{sb: &strings.Builder{}},
	}
}

func (s *Selector[T]) Build() (*Query, error) {
	sb := s.sb
	var err error
	s.model, err = s.db.r.get(new(T))
	//s.model, err = s.db.r.ParseModel(new(T))
	if err != nil {
		return nil, err
	}
	// 处理空格问题
	sb.WriteString("SELECT * FROM ")
	// 通过反射获得表明，泛型的名称。默认使用结构体名作为表名
	// 对于带db的参数
	// 1. 让用户自己加入`
	// 2. 开发者自行切割
	sb.WriteByte('`')

	if s.table == "" {
		sb.WriteString(s.model.TableName)
	} else {
		if strings.Contains(s.table, ".") {
			segs := strings.Split(s.table, ".")
			sb.WriteString(segs[0])
			sb.WriteString("`.`")
			sb.WriteString(segs[1])
		} else {
			sb.WriteString(s.table)
		}

	}
	sb.WriteByte('`')

	if len(s.where) > 0 {
		sb.WriteString(" WHERE ")

		if s.buildPredicates(s.where) != nil {
			return nil, err
		}

		//sb.WriteString("" + args)
	}

	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		args: s.args,
	}, nil

}

// From 兼容传入空字符串，如果加入校验会影响连调功能

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Where(p ...Predicate) *Selector[T] {
	s.where = p
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

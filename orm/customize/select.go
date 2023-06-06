// @Author: zqy
// @File: select.go
// @Date: 2023/6/5 16:17
// @Description todo

package customize

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
	where []Predicate
	args  []any
	sb    *strings.Builder
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	sb := s.sb
	// 处理空格问题
	sb.WriteString("SELECT * FROM ")
	// 通过反射获得表明，泛型的名称。默认使用结构体名作为表名
	// 对于带db的参数
	// 1. 让用户自己加入`
	// 2. 开发者自行切割
	sb.WriteByte('`')

	if s.table == "" {
		var t T
		tableName := reflect.TypeOf(t).Name()
		sb.WriteString(tableName)
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
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}

		// 最后同一处理p，二叉树
		// 但是不知道left是什么类型
		err := s.buildExpression(p)
		if err != nil {
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

func (s *Selector[T]) buildExpression(e Expression) error {
	switch exp := e.(type) {
	case nil:
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')

		_, ok = exp.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(exp.name)
		s.sb.WriteByte('`')
	case Value:
		s.sb.WriteString("?")
		s.addArgs(exp.val)
	default:
		return fmt.Errorf("orm :不支持的表达式类型 %v", exp)
	}

	return nil
}

func (s *Selector[T]) addArgs(args any) {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, args)
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

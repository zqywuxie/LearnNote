package customize

import (
	"GoCode/orm/customize/internal"
	"fmt"
	"strings"
)

type Builder struct {
	args  []any
	sb    *strings.Builder
	model *Model
	table string
}

func (s *Builder) buildExpression(e Expression) error {
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
		fd, ok := s.model.FiledMap[exp.name]
		if !ok {
			return internal.NewErrorUnknownField(exp.name)
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(fd.ColName)
		s.sb.WriteByte('`')
	case Value:
		s.sb.WriteString("?")
		s.addArgs(exp.val)
	default:
		return fmt.Errorf("orm :不支持的表达式类型 %v", exp)
	}

	return nil
}

func (s *Builder) addArgs(args any) {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, args)
}

func (s *Builder) buildPredicates(ps []Predicate) error {
	p := ps[0]
	for i := 1; i < len(ps); i++ {
		p = p.And(ps[i])
	}

	// 最后同一处理p，二叉树
	// 但是不知道left是什么类型
	return s.buildExpression(p)
}

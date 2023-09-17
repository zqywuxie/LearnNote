package customize

import "strings"

type Deleter[T any] struct {
	table string
	where []Predicate
	Builder
}

func (s *Deleter[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	sb := s.sb
	var err error
	s.model, err = ParseModel(new(T))
	if err != nil {
		return nil, err
	}
	// 处理空格问题
	sb.WriteString("DELETE FROM ")
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
		if err = s.buildPredicates(s.where); err != nil {
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

func (s *Deleter[T]) From(table string) *Deleter[T] {
	s.table = table
	return s
}

func (s *Deleter[T]) Where(p ...Predicate) *Deleter[T] {
	s.where = p
	return s
}

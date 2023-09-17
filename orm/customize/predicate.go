// @Author: zqy
// @File: predicate.go
// @Date: 2023/6/5 21:07
// @Description todo

package customize

// string的衍生类型
// 可以清晰地表达该类型是代表某种操作或运算符
type op string

// 类型别名 重新创建了已有的类型
//type t = string

func (o op) String() string {
	return string(o)
}

const (
	opEq  op = "="
	opGt  op = ">"
	opGe  op = ">="
	opLt  op = "<"
	opLe  op = "=<"
	opAnd op = "AND"
	opNot op = "NOT"
	opOr  op = "OR"
)

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (p Predicate) expr() {
}

// Expression 用于标记是一个表达式
// 标记接口（Marker Interface）是一个不包含任何方法的空接口，它没有提供任何方法实现，
// 仅用于标识一个类型是否符合特定的约定或属性。
// 然后将value,column注册为Expression，后续便于统一管理
type Expression interface {
	expr()
}
type Column struct {
	name string
}

// Value 标记参数
type Value struct {
	val any
}

func (v Value) expr() {
}

func (c Column) expr() {
}

func C(name string) Column {
	return Column{name: name}
}

// Eq 用法 Eq(sub.id,12)
//func Eq(column string,arg any) Predicate {
//	return Predicate{
//		Column: column,
//		Op:     "=",
//		Arg:    arg,
//	}
//}

// Eq  用法：C("id").Eq(12)
// 方便后续子查询和join查询
// where id = 12
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: Value{val: arg},
	}
}

// Not Not(C("id").Eq(12))
// where not (id = 12)
func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}

// And C("id").Eq(12).And(C("id").Eq(12))
// where (id = 12) and (id = 12)
func (p Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAnd,
		right: right,
	}
}

func (p Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOr,
		right: right,
	}
}

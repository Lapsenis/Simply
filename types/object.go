package types

import (
	"Simply/ast"
	"fmt"
	"strconv"
)

var (
	TRUE  = &Bool{Value: true}
	FALSE = &Bool{Value: false}
	NULL  = &Null{}
)

type Object interface {
	String() string
}

type Error struct {
	Value string
}

func (e *Error) String() string { return e.Value }

type Int struct {
	Value int64
}

func (i *Int) String() string { return strconv.FormatInt(i.Value, 10) }

type String struct {
	Value string
}

func (s *String) String() string { return s.Value }

type Bool struct {
	Value bool
}

func (b *Bool) String() string { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) String() string { return "null" }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.CodeBlock
	Ctx        *Context
}

func (f *Function) String() string { return "TODO types.Function" }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) String() string { return r.Value.String() }

type InternalCallFunc func(args ...Object) Object
type InternalCall struct {
	Fn InternalCallFunc
}

func (i *InternalCall) String() string { return "Internal call" }

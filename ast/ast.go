package ast

import (
	"fmt"
	"strconv"
	"strings"
)

type Node interface {
	String() string
}

type Program struct {
	Statements []Node
}

func (p *Program) String() string {
	var sb strings.Builder

	for _, v := range p.Statements {
		sb.WriteString(v.String())
	}

	return sb.String()
}

type DeclarativeStatement struct {
	Name  Identifier
	Value Node //Expression
}

func (d *DeclarativeStatement) String() string {
	return fmt.Sprintf("%s : %s", d.Name.String(), d.Value.String())
}

type ReturnStatement struct {
	Value Node
}

func (r *ReturnStatement) String() string { return r.Value.String() }

type ExpressionStatement struct {
	Expression Node
}

func (e *ExpressionStatement) String() string { return e.Expression.String() }

type PrefixExpression struct {
	Prefix     string
	Expression Node
}

func (p *PrefixExpression) String() string { return "TODO PrefixExpression" }

type ConditionalExpression struct {
	Condition Node
	True      *CodeBlock
	False     *CodeBlock
}

func (c *ConditionalExpression) String() string { return "TODO ConditionalExpression" }

type CallExpression struct {
	Arguments []Node
	Function  Node
}

func (c *CallExpression) String() string { return "TODO CallExpression" }

type Identifier struct {
	Value string
}

func (i *Identifier) String() string { return i.Value }

type FunctionLiteral struct {
	Parameters []*Identifier
	Body       *CodeBlock
}

func (f *FunctionLiteral) String() string { return "TODO FunctionLiteral" }

type CodeBlock struct {
	Statements []Node
}

func (c *CodeBlock) String() string {
	var sb strings.Builder

	for _, v := range c.Statements {
		sb.WriteString(v.String())
	}

	return sb.String()
}

type IntLiteral struct {
	Value int64
}

func (il *IntLiteral) String() string { return strconv.FormatInt(il.Value, 0) }

type BoolLiteral struct {
	Value bool
}

func (b *BoolLiteral) String() string { return strconv.FormatBool(b.Value) }

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) String() string { return s.Value }

type InfixExpression struct {
	Left     Node
	Operator string
	Right    Node
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("%s %s %s", i.Left.String(), i.Operator, i.Right.String())
}

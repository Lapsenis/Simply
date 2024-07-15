package evaluator

import (
	"Simply/ast"
	"Simply/types"
	"fmt"
)

func Eval(n ast.Node, ctx *types.Context) types.Object {

	switch node := n.(type) {
	case *ast.Program:
		return evalProgram(node, ctx)
	case *ast.DeclarativeStatement:
		return evalDeclarativeStatement(node, ctx)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, ctx)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, ctx)
	case *ast.CallExpression:
		return evalCallExpression(node, ctx)
	case *ast.InfixExpression:
		return evalInfixExpression(node, ctx)
	case *ast.ConditionalExpression:
		return evalConditionalExpression(node, ctx)
	case *ast.Identifier:
		return evalIdentifier(node, ctx)
	case *ast.IntLiteral:
		return &types.Int{Value: node.Value}
	case *ast.BoolLiteral:
		return getBoolType(node.Value)
	case *ast.StringLiteral:
		return &types.String{Value: node.Value}
	case *ast.FunctionLiteral:
		return &types.Function{Parameters: node.Parameters, Body: node.Body, Ctx: ctx}
	case *ast.ReturnStatement:
		return evalReturnStatement(node, ctx)
	case *ast.CodeBlock:
		return evalCodeBlock(node, ctx)
	}

	return newError("Failed execute node %t", n)
}

func evalProgram(p *ast.Program, ctx *types.Context) types.Object {
	var result types.Object

	for _, v := range p.Statements {
		result = Eval(v, ctx)

		switch result := result.(type) {
		case *types.ReturnValue:
			return result.Value
		case *types.Error:
			return result
		}
	}

	return result
}

func evalDeclarativeStatement(d *ast.DeclarativeStatement, ctx *types.Context) types.Object {
	result := Eval(d.Value, ctx)

	if isError(result) {
		return result
	}

	ctx.Set(d.Name.Value, result)

	return nil //Good job? Here is nothing :P
}

func evalIdentifier(node *ast.Identifier, ctx *types.Context) types.Object {
	if val, ok := ctx.Get(node.Value); ok {
		return val
	}

	if builtin, ok := internalCalls[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalPrefixExpression(node *ast.PrefixExpression, ctx *types.Context) types.Object {
	exp := Eval(node.Expression, ctx)
	if isError(exp) {
		return exp
	}

	switch node.Prefix {
	case "!":
		if !isType[*types.Bool](exp) {
			return getBoolType(false)
		}

		return getBoolType(!exp.(*types.Bool).Value)
	case "-":
		if !isType[*types.Int](exp) {
			return newError("unknown operator: -%s", exp.String())
		}
		return &types.Int{Value: -(exp.(*types.Int).Value)}
	default:
		return newError("unknown prefix: %s%s", node.Prefix, exp.String())
	}
}

func evalCallExpression(node *ast.CallExpression, ctx *types.Context) types.Object {
	f := Eval(node.Function, ctx)
	if isError(f) {
		return f
	}

	var args []types.Object
	for _, e := range node.Arguments {
		evaluated := Eval(e, ctx)
		if isError(evaluated) {
			return evaluated
		}
		args = append(args, evaluated)
	}

	return executeFunction(f, args)
}

func executeFunction(fn types.Object, args []types.Object) types.Object {
	switch fn := fn.(type) {
	case *types.Function:
		newCtx := createFuncCtx(fn, args)
		evaluated := Eval(fn.Body, newCtx)
		return unwrapReturnValue(evaluated)
	case *types.InternalCall:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.String())
	}
}

func unwrapReturnValue(obj types.Object) types.Object {
	if returnValue, ok := obj.(*types.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func createFuncCtx(fn *types.Function, args []types.Object) *types.Context {
	env := types.NewContext(fn.Ctx)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func evalInfixExpression(node *ast.InfixExpression, ctx *types.Context) types.Object {
	left := Eval(node.Left, ctx)

	if isError(left) {
		return left
	}

	right := Eval(node.Right, ctx)

	if isError(right) {
		return right
	}

	return evalInfixOperation(node.Operator, left, right)
}

func evalInfixOperation(op string, left, right types.Object) types.Object {
	if isTypeEqual[*types.Int](left, right) {
		return evalIntInfixExpression(op, left, right)
	}

	return newError("Unknown inflix operation %s %s %s", left.String(), op, right.String())
}

func evalIntInfixExpression(op string, left, right types.Object) types.Object {
	leftValue := left.(*types.Int).Value
	rightValue := right.(*types.Int).Value

	switch op {
	case "+":
		return &types.Int{Value: leftValue + rightValue}
	case "-":
		return &types.Int{Value: leftValue - rightValue}
	case "*":
		return &types.Int{Value: leftValue * rightValue}
	case "/":
		return &types.Int{Value: leftValue + rightValue}
	case "<":
		return getBoolType(leftValue < rightValue)
	case ">":
		return getBoolType(leftValue > rightValue)
	case "==":
		return getBoolType(leftValue == rightValue)
	case "!=":
		return getBoolType(leftValue != rightValue)
	default:
		return newError("operator %s not supported for int", op)
	}
}

func getBoolType(b bool) *types.Bool {
	if b {
		return types.TRUE
	} else {
		return types.FALSE
	}
}

func evalConditionalExpression(node *ast.ConditionalExpression, ctx *types.Context) types.Object {
	condition := Eval(node.Condition, ctx)

	if isError(condition) {
		return condition
	}

	if isConditionTrue(condition) {
		return Eval(node.True, ctx)
	} else if node.False != nil {
		return Eval(node.False, ctx)
	} else {
		return types.NULL
	}
}

func isConditionTrue(o types.Object) bool {
	switch o {
	case types.NULL:
		return false
	case types.TRUE:
		return true
	case types.FALSE:
		return false
	default:
		return true
	}
}

func evalReturnStatement(node *ast.ReturnStatement, ctx *types.Context) types.Object {
	val := Eval(node.Value, ctx)
	if isError(val) {
		return val
	}
	return &types.ReturnValue{Value: val}
}

func evalCodeBlock(block *ast.CodeBlock, env *types.Context) types.Object {
	var result types.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			if isType[*types.ReturnValue](result) || isType[*types.Error](result) {
				return result
			}
		}
	}
	return result
}

func isError(obj types.Object) bool {
	if _, ok := obj.(*types.Error); ok {
		return true
	} else {
		return false
	}
}

func newError(format string, a ...interface{}) types.Object {
	return &types.Error{Value: fmt.Sprintf(format, a...)}
}

func isType[T any](o types.Object) bool {
	_, ok := o.(T)
	return ok
}

func isTypeEqual[T any](left, right types.Object) bool {
	_, leftOk := left.(T)
	_, rightOk := right.(T)

	return leftOk && rightOk
}

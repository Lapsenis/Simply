package evaluator

import (
	"Simply/types"
	"bufio"
	"fmt"
	"os"
)

var internalCalls = map[string]*types.InternalCall{
	"len":     {Fn: internal_len},
	"println": {Fn: internal_println},
	"print":   {Fn: internal_print},
	"input":   {Fn: internal_input},
}

func internal_len(args ...types.Object) types.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *types.String:
		return &types.Int{Value: int64(len(arg.Value))}
	default:
		return newError("argument to `len` not supported")
	}
}

func internal_println(args ...types.Object) types.Object {
	for _, a := range args {
		fmt.Println(a)
	}

	return types.NULL
}

func internal_print(args ...types.Object) types.Object {
	for _, a := range args {
		fmt.Print(a)
	}

	return types.NULL
}

func internal_input(args ...types.Object) types.Object {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return &types.String{Value: scanner.Text()}
}

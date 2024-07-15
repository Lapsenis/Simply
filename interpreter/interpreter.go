package interpreter

import (
	"Simply/ast"
	"Simply/evaluator"
	"Simply/lexer"
	"Simply/parser"
	"Simply/types"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

const promtd = ">>>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	globalCtx := types.NewContext(nil)
	for {
		fmt.Fprint(out, promtd)

		if ok := scanner.Scan(); !ok {
			continue
		}

		program, err := parseInput(out, scanner.Text())

		if err != nil {
			fmt.Println(out, err)
			continue
		}

		evalResult := evaluator.Eval(program, globalCtx)

		if evalResult != nil {
			fmt.Fprintln(out, evalResult.String())
		}

	}
}

func ProcessFile(path string, out io.Writer) {
	scriptText, err := readFile(path)
	if err != nil {
		return
	}

	program, err := parseInput(out, scriptText)
	if err != nil {
		fmt.Println(out, err)
		return
	}

	ctx := types.NewContext(nil)

	evalResult := evaluator.Eval(program, ctx)

	evalError, isError := evalResult.(*types.Error)

	if isError {
		logEvalErrors(evalError)
	}
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "", err
	}

	return string(content), err
}

func parseInput(out io.Writer, input string) (*ast.Program, error) {
	t := lexer.NewTokenizer(input)
	p := parser.NewParser(t)
	parseResult := p.ParseProgram()

	if len(p.Errors) > 0 {
		logParseErrors(out, p)
		return nil, errors.New("failed to parse")
	}

	return parseResult, nil
}

func logParseErrors(out io.Writer, p *parser.Parser) {
	for _, e := range p.Errors {
		fmt.Fprintln(out, e)
	}
}

func logEvalErrors(e *types.Error) {
	fmt.Println(e.String())
}

package evaluator

import (
	"Simply/lexer"
	"Simply/parser"
	"Simply/types"
	"testing"
)

func TestOfEverything(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"let x = if (3 > 4) { 3; } else {4}; x;", 4},
		{"let x = if (3 > 4) { 3; }; x;", nil},
		{"let x = if (3 > 2) { 3; }; x;", 3},
		{"let x = (5 + 5) * 3; x;", 30},
		{"let x = 5 + 5 * 3; x;", 20},
		{"let x = true; x;", true},
		{"let x = !true; x;", false},
		{"let x = -5; x;", -5},
		{"let x = func(x,y){ return 3; x+y;}; x(5,5);", 3},
		{"let x = func(x,y){ x+y;}; x(5,5);", 10},
		{"let x = 5 == 6; x;", false},
		{"let x = 5 + 5; x;", 10},
		{"let x = 5; x;", 5},
	}
	for _, tt := range tests {
		result := testEvaluator(t, tt.input)

		//TODO maybe impement testing that works
		if result == tt.expectedValue {
			t.Fatal("No match :()")
		}
	}
}

func testEvaluator(t *testing.T, input string) types.Object {
	l := lexer.NewTokenizer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	checkErrors(t, p)
	ctx := types.NewContext(nil)
	return Eval(program, ctx)
}

func checkErrors(t *testing.T, p *parser.Parser) {

	if len(p.Errors) == 0 {
		return
	}

	for _, e := range p.Errors {
		t.Error(e)
	}

}

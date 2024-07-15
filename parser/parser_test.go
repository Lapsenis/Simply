package parser

import (
	"Simply/lexer"
	"testing"
)

func TestDeclarativeStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5 + 5;", "x", "5 + 5"},
		{"let x = \"bazuka\";", "x", "bazuka"},
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}
	for _, tt := range tests {
		tokenizer := lexer.NewTokenizer(tt.input)
		p := NewParser(tokenizer)
		program := p.ParseProgram()
		checkErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatal("program statement is 0")
		}
	}
}

func checkErrors(t *testing.T, p *Parser) {

	if len(p.Errors) == 0 {
		return
	}

	for _, e := range p.Errors {
		t.Error(e)
	}

}

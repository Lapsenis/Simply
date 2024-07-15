package lexer

type Tokenizer struct {
	input      string
	currentPos int
	ch         byte
	//row, col   int TODO implement row/col trackin for debuging
}

func NewTokenizer(input string) *Tokenizer {
	t := &Tokenizer{input: input}
	t.readChar()
	return t
}

func (t *Tokenizer) NextToken() Token {
	var tok Token

	t.skipToNextCh()

	switch t.ch {
	//String
	case '"':
		t.readChar()
		tok.Type = STRING
		tok.Literal = t.readValue(isString)
	//Operators
	case '+':
		tok = Token{Type: PLUS, Literal: string(t.ch)}
	case '-':
		tok = Token{Type: MINUS, Literal: string(t.ch)}
	case '/':
		tok = Token{Type: SLASH, Literal: string(t.ch)}
	case '*':
		tok = Token{Type: ASTERISK, Literal: string(t.ch)}
	case '<':
		tok = Token{Type: LT, Literal: string(t.ch)}
	case '>':
		tok = Token{Type: GT, Literal: string(t.ch)}
	//Two byte operators
	case '=':
		if t.peekChar() == '=' {
			ch := t.ch
			t.readChar()
			literal := string(ch) + string(t.ch)
			tok = Token{Type: EQ, Literal: literal}
		} else {
			tok = Token{Type: ASSIGN, Literal: string(t.ch)}
		}
	case '!':
		if t.peekChar() == '=' {
			ch := t.ch
			t.readChar()
			literal := string(ch) + string(t.ch)
			tok = Token{Type: NOT_EQ, Literal: literal}
		} else {
			tok = Token{Type: BANG, Literal: string(t.ch)}
		}
	//Delimiters
	case ';':
		tok = Token{Type: SEMICOLON, Literal: string(t.ch)}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(t.ch)}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(t.ch)}
	case ',':
		tok = Token{Type: COMMA, Literal: string(t.ch)}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(t.ch)}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(t.ch)}
	//End of file
	case 0:
		tok = Token{Type: EOF, Literal: ""}
	default:
		if isLetter(t.ch) {
			tok.Literal = t.readValue(isLetter)
			tok.Type = GetType(tok.Literal)
			return tok
		} else if isDigit(t.ch) {
			tok.Type = INT
			tok.Literal = t.readValue(isDigit)
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(t.ch)}
		}
	}

	t.readChar()
	return tok
}

func (t *Tokenizer) skipToNextCh() {
	for t.ch == ' ' || t.ch == '\t' || t.ch == '\n' || t.ch == '\r' {
		t.readChar()
	}
}

func (t *Tokenizer) readChar() {
	if t.currentPos >= len(t.input) {
		t.ch = 0
	} else {
		t.ch = t.input[t.currentPos]
	}

	t.currentPos++
}

func (t *Tokenizer) peekChar() byte {
	if t.currentPos >= len(t.input) {
		return 0
	} else {
		return t.input[t.currentPos]
	}
}

func (t *Tokenizer) readValue(fn func(byte) bool) string {
	startPos := t.currentPos - 1
	for fn(t.ch) {
		t.readChar()
	}

	return t.input[startPos : t.currentPos-1]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isString(ch byte) bool {
	return ch != '"'
}

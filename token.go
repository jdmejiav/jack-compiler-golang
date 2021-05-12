package main

import (
	"bufio"
	"io"
	"strconv"
	"unicode"
)

type Token struct {
	tokenType TokenType
	token     string
}

func NewToken(tipo TokenType, token string) *Token {
	return &Token{
		tokenType: tipo,
		token:     token,
	}
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

type TokenType int

const (
	EOF = iota
	ILLEGAL
	IDENTIFIER
	KEYWORD
	SYMBOL
	STRINGCONST
	NUMBER
	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	LPARENT // (
	RPARENT // )
	LKEY    // {
	RKEY    // }
	LCORCH  // [
	RCORCH  // ]
	DOT     // .
	SLASH   // /

	COMA       // ,
	DOTANDCOMA // ;
	ASTERISCO  // *
	AND        // &
	OR         // |
	MINOR      // <
	MAJOR      // >

	CLASS
	CONSTRUCTOR
	FUNCTION
	METHOD
	FIELD
	STATIC
	VAR
	INT
	CHAR
	BOOLEAN
	VOID
	TRUE
	FALSE
	NULL
	THIS
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN

	ASSIGN // =

)

var tokens = []string{
	EOF:         "EOF",
	ILLEGAL:     "ILLEGAL",
	IDENTIFIER:  "IDENTIFIER",
	KEYWORD:     "KEYWORD",
	SYMBOL:      "SYMBOL",
	STRINGCONST: "STRINGCONST",
	NUMBER:      "NUMBER",
	// Infix ops
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	LKEY: "{",
	RKEY: "}",

	LPARENT: "(",
	RPARENT: ")",

	LCORCH: "[",
	RCORCH: "]",

	SLASH: "/",

	DOT: ".",

	COMA:       ",",
	DOTANDCOMA: ";",
	ASTERISCO:  "*",
	AND:        "&",
	OR:         "|",
	MINOR:      "<",
	MAJOR:      ">",

	ASSIGN: "=",
}

var tokenKeyword = []string{
	CLASS:       "class",
	CONSTRUCTOR: "constructor",
	FUNCTION:    "function",
	METHOD:      "method",
	FIELD:       "field",
	STATIC:      "static",
	VAR:         "var",
	INT:         "int",
	CHAR:        "char",
	BOOLEAN:     "boolean",
	VOID:        "void",
	TRUE:        "true",
	FALSE:       "false",
	NULL:        "null",
	THIS:        "this",
	LET:         "let",
	DO:          "do",
	IF:          "if",
	ELSE:        "else",
	WHILE:       "while",
	RETURN:      "return",
}

func (t TokenType) String() string {
	return tokens[t]
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

// Lex scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.
func (l *Lexer) tokenize() *Token {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return NewToken(EOF, "")
				//return l.pos, EOF, ""
			}

			// at this point there isn't much we can do, and the compiler
			// should just return the raw error to the user
			panic(err)
		}

		l.pos.column++

		switch r {

		case '\n':
			l.resetPosition()

		case '+':
			return NewToken(SYMBOL, "+")
			//return l.pos, SYMBOL, "+"
		case '-':
			return NewToken(SYMBOL, "-")
		case '*':
			return NewToken(SYMBOL, "*")
		case '/':
			return NewToken(SYMBOL, "/")
		case '=':
			return NewToken(SYMBOL, "=")
		case '{':
			return NewToken(SYMBOL, "{")
		case '}':
			return NewToken(SYMBOL, "}")
		case '(':
			return NewToken(SYMBOL, "(")
		case ')':
			return NewToken(SYMBOL, ")")
		case '[':
			return NewToken(SYMBOL, "[")
		case ']':
			return NewToken(SYMBOL, "]")
		case '.':
			return NewToken(SYMBOL, ".")
		case ',':
			return NewToken(SYMBOL, ",")
		case ';':
			return NewToken(SYMBOL, ";")
		case '&':
			return NewToken(SYMBOL, "&")
		case '|':
			return NewToken(SYMBOL, "|")
		case '<':
			return NewToken(SYMBOL, "<")
		case '>':
			return NewToken(SYMBOL, ">")
		case '~':
			return NewToken(SYMBOL, "~")
		case '"':
			return NewToken(STRINGCONST, l.lexStringConst())
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				//startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return NewToken(NUMBER, lit)
			} else if unicode.IsLetter(r) {
				// backup and let lexIdent rescan the beginning of the ident
				//startPos := l.pos
				l.backup()
				lit := l.lexString()

				for _, i := range tokenKeyword {
					if i == lit {
						return NewToken(KEYWORD, lit)
					}
				}

				return NewToken(IDENTIFIER, lit)
			} else {
				return NewToken(ILLEGAL, string(r))
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func (l *Lexer) lexStringConst() string {
	var lit string
	l.pos.column++
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				panic("There is one \" missing at " + strconv.Itoa(l.pos.column))

			}
		}

		l.pos.column++

		if r == 34 {
			l.pos.column++
			return lit
		} else {
			lit += string(r)

		}
	}
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexString() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}

		l.pos.column++

		if unicode.IsLetter(r) || (len(lit) > 0 && unicode.IsDigit(r) || (len(lit) > 0 && r == 95)) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}

// lexInt scans the input until the end of an integer and then returns the
// literal.
func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return lit
			}
		}
		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the integer
			l.backup()
			return lit
		}
	}
}

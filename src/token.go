package main

import (
	"bufio"
	"io"
	"strconv"
	"unicode"
//	"fmt"
)

type Token struct {
	tokenType TokenType
	token     string
	line      int
	column    int
}

func NewToken(tipo TokenType, token string, line int, column int) *Token {
	return &Token{
		tokenType: tipo,
		token:     token,
		line:      line,
		column:    column,
	}
}

func (t *Token) getToken() string {
	return t.token
}
func (t *Token) getTokenType() TokenType {
	return t.tokenType
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

type Copy struct {
	pos    Position
	reader *bufio.Reader
}

type TokenType int

const (
	EOF             = iota //0
	ILLEGAL                //0
	IDENTIFIER             //2
	KEYWORD                //3
	SYMBOL                 //4
	INTEGERCONSTANT        //5
	STRINGCONST            //6

	AXIOMA //7

	CLASS           //8
	CLASSVARDEC     //9
	TYPE            //10
	SUBROUTINEDEC   //11
	PARAMETERLIST   //12
	SUBROUTINEBODY  //13
	VARDEC          //14
	CLASSNAME       //15
	SUBROUTINENAME  //16
	VARNAME         //17
	STATEMENTS      //18
	STATEMENT       //19
	LETSTATEMENT    //20
	IFSTATEMENT     //21
	WHILESTATEMENT  //22
	DOSTATEMENT     //23
	RETURNSTATEMENT //24
	EXPRESSION      //25
	TERM            //26
	SUBROUTINECALL  //27
	EXPRESSIONLIST  //28
	OP              //29
	UNARYOP         //30
	KEYWORDCONSTANT //31
	DOTANDCOMA      //32

	CONSTRUCTOR        //33
	FUNCTION           //34
	METHOD             //35
	FIELD              //36
	STATIC             //37
	VAR                //38
	INT                //39
	CHAR               //40
	BOOLEAN            //41
	VOID               //42
	TRUE               //43
	FALSE              //44
	NULL               //45
	THIS               //46
	LET                //47
	DO                 //48
	IF                 //49
	ELSE               //50
	WHILE              //51
	RETURN             //52
	LPARENT            //53
	RPARENT            //54
	COMA               //55
	ARRAY              //56
	ELSESTATEMENT      //57
	EXPRESSIONCOND     //58
	TERMCOND           //59
	TERMPROD           //60
	EXPRESSIONLISTCOND //61
	SUBROUTINEDECCOND  //62
)

var tokens = []string{
	EOF:             "EOF",
	ILLEGAL:         "ILLEGAL",
	IDENTIFIER:      "IDENTIFIER",
	KEYWORD:         "KEYWORD",
	SYMBOL:          "SYMBOL",
	STRINGCONST:     "STRINGCONST",
	INTEGERCONSTANT: "NUMBER",
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
				return NewToken(EOF, "", l.pos.line, l.pos.column)
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
			return NewToken(OP, "+", l.pos.line, l.pos.column)
		case '-':
			return NewToken(OP, "-", l.pos.line, l.pos.column)
		case '*':
			return NewToken(OP, "*", l.pos.line, l.pos.column)
		case '/':
			return NewToken(OP, "/", l.pos.line, l.pos.column)
		case '=':
			return NewToken(OP, "=", l.pos.line, l.pos.column)
		case '{':
			return NewToken(SYMBOL, "{", l.pos.line, l.pos.column)
		case '}':
			return NewToken(SYMBOL, "}", l.pos.line, l.pos.column)
		case '(':
			return NewToken(SYMBOL, "(", l.pos.line, l.pos.column)
		case ')':
			return NewToken(SYMBOL, ")", l.pos.line, l.pos.column)
		case '[':
			return NewToken(SYMBOL, "[", l.pos.line, l.pos.column)
		case ']':
			return NewToken(SYMBOL, "]", l.pos.line, l.pos.column)
		case '.':
			return NewToken(SYMBOL, ".", l.pos.line, l.pos.column)
		case ',':
			return NewToken(SYMBOL, ",", l.pos.line, l.pos.column)
		case ';':
			return NewToken(SYMBOL, ";", l.pos.line, l.pos.column)
		case '&':
			return NewToken(OP, "&", l.pos.line, l.pos.column)
		case '|':
			return NewToken(OP, "|", l.pos.line, l.pos.column)
		case '<':
			return NewToken(OP, "<", l.pos.line, l.pos.column)
		case '>':
			return NewToken(OP, ">", l.pos.line, l.pos.column)
		case '~':
			return NewToken(UNARYOP, "~", l.pos.line, l.pos.column)
		case '!':
			return NewToken(UNARYOP, "~", l.pos.line, l.pos.column)
		case '"':
			return NewToken(STRINGCONST, l.lexStringConst(), l.pos.line, l.pos.column)
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				//startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return NewToken(INTEGERCONSTANT, lit, l.pos.line, l.pos.column)
			} else if unicode.IsLetter(r) {
				// backup and let lexIdent rescan the beginning of the ident
				//startPos := l.pos
				l.backup()
				lit := l.lexString()

				for _, i := range tokenKeyword {
					if i == lit {
						return NewToken(KEYWORD, lit, l.pos.line, l.pos.column)
					}
				}

				return NewToken(IDENTIFIER, lit, l.pos.line, l.pos.column)
			} else {
				return NewToken(ILLEGAL, string(r), l.pos.line, l.pos.column)
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


func (l *Lexer) removeComent(){
	for {

		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
					panic("")
			}
		}
		if r=='\n'{
			l.pos.column=0
			l.pos.line++
			return
		}
	}
}

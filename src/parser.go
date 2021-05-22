package main

import (
	"fmt"
	"strconv"
)

type TokenAnalysis struct {
	terminalToken bool
	tokenType     TokenType
	stringLiteral string
}

func NewTokenAnalysis(terminalToken bool, tokenType TokenType, stringLiteral string) TokenAnalysis {
	return TokenAnalysis{
		terminalToken: terminalToken,
		tokenType:     tokenType,
		stringLiteral: stringLiteral,
	}
}

func (t *TokenAnalysis) ToString() {
	fmt.Printf("%s \t\t%d\n", t.stringLiteral, t.tokenType)
}

type Analyzer struct {
	tokens []Token
	stack  Stack
}

func NewAnalyzer(tokens []Token) *Analyzer {
	return &Analyzer{
		tokens: tokens,
		stack:  *NewStack(),
	}
}

func (a *Analyzer) Analyze() {
	a.stack.push(NewTokenAnalysis(false, CLASS, ""))
	i := 0
	a.identifyToken(i)
}

func (a *Analyzer) identifyToken(i int) {
	//a.stack.ToString()

	//fmt.Print(a.stack.pos(i).stringLiteral)
	//fmt.Print("\t\t")
	//fmt.Printf("%d",a.stack.pos(i).tokenType)
	if a.stack.pos(i).terminalToken {
		if a.tokens[i].token == a.stack.pos(i).stringLiteral {
			i++
			a.identifyToken(i)
		} else {
			panic("There was an error at line " + strconv.Itoa(a.tokens[i].line) + " symbol <" + a.tokens[i].token + "> not found")
		}
	} else {
		switch a.stack.pos(i).tokenType {
		case CLASS:
			a.axioma(i)
			break

		case CLASSNAME:
			a.className(i)
			break

		case CLASSVARDEC:
			a.classVarDec(i)
			break

		case TYPE:
			a.typeId(i)
			break

		case VARNAME:
			a.varName(i)
			break

		case SUBROUTINEDEC:
			a.subRoutineDec(i)
			break

		case SUBROUTINENAME:
			a.subRoutineName(i)
			break

		case PARAMETERLIST:
			a.parameterList(i)
			break

		case SUBROUTINEBODY:
			a.subRoutineBody(i)
			break

		case VARDEC:
			a.varDec(i)
			break

		case STATEMENTS:
			a.statements(i)
			break

		case LETSTATEMENT:
			a.letStatement(i)
			break

		case IFSTATEMENT:
			a.ifStatement(i)
			break
		case WHILESTATEMENT:
			a.whileStatement(i)
			break

		case DOSTATEMENT:
			a.doStatement(i)
			break

		case RETURNSTATEMENT:
			a.returnStatement(i)
			break

		case ELSESTATEMENT:
			a.elseStatement(i)
			break

		case EXPRESSION:
			a.expression(i)
			break

		case ARRAY:
			a.array(i)
			break

		case TERM:
			a.term(i)
			break

		case TERMPROD:
			a.termProd(i) //este es el de las nuevas reglas en term
			break

		case TERMCOND:
			a.termCond(i) //este es el de <term>(<op><term>)*
			break

		case EXPRESSIONCOND:
			a.expressionCond(i) // Este es el del return
			break

		case SUBROUTINECALL:
			a.subRoutineCall(i)
			break
		case EXPRESSIONLIST:
			a.expressionList(i) // Este es el del return
			break
		case EXPRESSIONLISTCOND:
				a.expressionListCond(i) // Este es el del return
				break
		case SUBROUTINEDECCOND:
			a.subRoutineDecCond(i)
			break
		}


	}
}

func (a *Analyzer) axioma(i int) {
	arr := []TokenAnalysis{
		NewTokenAnalysis(true, CLASS, "class"),
		NewTokenAnalysis(false, CLASSNAME, ""),
		NewTokenAnalysis(true, SYMBOL, "{"),
		NewTokenAnalysis(false, CLASSVARDEC, ""),
		NewTokenAnalysis(false, SUBROUTINEDEC, ""),
		NewTokenAnalysis(true, SYMBOL, "}"),
	}
	a.stack.pushAtPos(i, arr)
	a.identifyToken(i)
}

func (a *Analyzer) className(i int) {
	arr := []TokenAnalysis{
		NewTokenAnalysis(true, IDENTIFIER, a.tokens[i].token),
	}
	a.stack.pushAtPos(i, arr)
	a.identifyToken(i)

}

func (a *Analyzer) classVarDec(i int) {
	var arr []TokenAnalysis
	if "static" == a.tokens[i].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, STATIC, "static"),
			NewTokenAnalysis(false, TYPE, ""),
			NewTokenAnalysis(false, VARNAME, ""),
			NewTokenAnalysis(true, DOTANDCOMA, ";"),
			NewTokenAnalysis(false, CLASSVARDEC, ""),
		}
		a.stack.pushAtPos(i, arr)
		//a.stack.ToString()
		a.identifyToken(i)

	} else if "field" == a.tokens[i].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, STATIC, "field"),
			NewTokenAnalysis(false, TYPE, ""),
			NewTokenAnalysis(false, VARNAME, ""),
			NewTokenAnalysis(true, DOTANDCOMA, ";"),
			NewTokenAnalysis(false, CLASSVARDEC, ""),
		}
	}
	a.stack.pushAtPos(i, arr)
	//a.stack.ToString()
	a.identifyToken(i)

}

func (a *Analyzer) typeId(i int) {
	var arr []TokenAnalysis
	if "int" == a.tokens[i].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, INT, "int"),
		}
	} else if "char" == a.tokens[i].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, CHAR, "char"),
		}
	} else if "boolean" == a.tokens[i].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, CHAR, "boolean"),
		}
	} else if IDENTIFIER == a.tokens[i].tokenType {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, IDENTIFIER, a.tokens[i].token),
		}
	} else {

		panic("Invalid data type " + a.tokens[i].token + " at line " + strconv.Itoa(a.tokens[i].line) + ":" + strconv.Itoa(a.tokens[i].column))
	}
	a.stack.pushAtPos(i, arr)
	//a.stack.ToString()
	a.identifyToken(i)
}

func (a *Analyzer) varName(i int) {
	arr := []TokenAnalysis{
		NewTokenAnalysis(true, IDENTIFIER, a.tokens[i].token),
	}
	a.stack.pushAtPos(i, arr)
	//i++
	a.identifyToken(i)
}

func (a *Analyzer) subRoutineDec(i int) {


	var arr []TokenAnalysis
	if "constructor" == a.tokens[i].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, CONSTRUCTOR, "constructor"),
		}
	} else if "function" == a.tokens[i].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, FUNCTION, "function"),
		}
	} else if "method" == a.tokens[i].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, METHOD, "method"),
		}
	} else {
		panic("Unknown subroutine name " + a.tokens[i].token + "at line " + strconv.Itoa(a.tokens[i].line) + ":" + strconv.Itoa(a.tokens[i].column))
	}

	if "void" == a.tokens[i+1].token {
		arr = append(arr, []TokenAnalysis{
			NewTokenAnalysis(true, VOID, "void"),
			NewTokenAnalysis(false, SUBROUTINENAME, ""),
			NewTokenAnalysis(true, LPARENT, "("),
			NewTokenAnalysis(false, PARAMETERLIST, ""),
			NewTokenAnalysis(true, RPARENT, ")"),
			NewTokenAnalysis(false, SUBROUTINEBODY, ""),
			NewTokenAnalysis(false, SUBROUTINEDECCOND, ""),
		}...)
	} else {
		arr = append(arr, []TokenAnalysis{
			NewTokenAnalysis(false, TYPE, ""),
			NewTokenAnalysis(false, SUBROUTINENAME, ""),
			NewTokenAnalysis(true, LPARENT, "("),
			NewTokenAnalysis(false, PARAMETERLIST, ""),
			NewTokenAnalysis(true, RPARENT, ")"),
			NewTokenAnalysis(false, SUBROUTINEBODY, ""),
			NewTokenAnalysis(false, SUBROUTINEDECCOND, ""),
		}...)
	}
	a.stack.pushAtPos(i, arr)

	//a.stack.ToString()
	a.identifyToken(i)
}

func (a *Analyzer) subRoutineName(i int) {
	arr := []TokenAnalysis{
		NewTokenAnalysis(true, IDENTIFIER, a.tokens[i].token),
	}
	a.stack.pushAtPos(i, arr)

	//a.stack.ToString()
	//fmt.Println("sale del subRoutineName")
	a.identifyToken(i)
}

func (a *Analyzer) parameterList(i int) {
	var arr []TokenAnalysis
	//fmt.Println(a.tokens[i+2].token)

	if ")" == a.tokens[i].token {
		a.stack.deletePos(i)

		a.identifyToken(i)
		return
	} else if ")" == a.tokens[i+2].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, TYPE, ""),
			NewTokenAnalysis(false, VARNAME, ""),
		}
	} else if "," == a.tokens[i+2].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, TYPE, ""),
			NewTokenAnalysis(false, VARNAME, ""),
			NewTokenAnalysis(true, COMA, ","),
			NewTokenAnalysis(false, PARAMETERLIST, ""),
		}
	}
	a.stack.pushAtPos(i, arr)
	//i++
	a.identifyToken(i)
}

func (a *Analyzer) subRoutineBody(i int) {
	//fmt.Println("acá se inserta el subroutinebody")
	arr := []TokenAnalysis{
		NewTokenAnalysis(true, SYMBOL, "{"),
		NewTokenAnalysis(false, VARDEC, ""),
		NewTokenAnalysis(false, STATEMENTS, ""),
		NewTokenAnalysis(true, SYMBOL, "}"),
	}
	a.stack.pushAtPos(i, arr)
	a.identifyToken(i)
}

func (a *Analyzer) varDec(i int) {
	var arr []TokenAnalysis
	if a.tokens[i].token == "var" {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, VAR, "var"),
			NewTokenAnalysis(false, TYPE, ""),
			NewTokenAnalysis(false, VARNAME, ""),
			NewTokenAnalysis(true, VAR, ";"),
			NewTokenAnalysis(false, VARDEC, ""),
		}
		a.stack.pushAtPos(i, arr)
	} else {
		a.stack.deletePos(i)
	}
	a.identifyToken(i)
}

func (a *Analyzer) statements(i int) {
	var arr []TokenAnalysis

	if a.tokens[i].token == "let" {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, LETSTATEMENT, ""),
			NewTokenAnalysis(false, STATEMENTS, ""),
		}
	} else if a.tokens[i].token == "if" {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, IFSTATEMENT, ""),
			NewTokenAnalysis(false, STATEMENTS, ""),
		}
	} else if a.tokens[i].token == "while" {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, WHILESTATEMENT, ""),
			NewTokenAnalysis(false, STATEMENTS, ""),
		}
	} else if a.tokens[i].token == "do" {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, DOSTATEMENT, ""),
			NewTokenAnalysis(false, STATEMENTS, ""),
		}
	} else if a.tokens[i].token == "return" {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, RETURNSTATEMENT, ""),
			NewTokenAnalysis(false, STATEMENTS, ""),
		}
	} else {
		a.stack.deletePos(i)
		a.identifyToken(i)
		return
	}
	a.stack.pushAtPos(i, arr)
	a.identifyToken(i)
}
func (a *Analyzer) array(i int) {
	if a.tokens[i].token == "[" {
		arr := []TokenAnalysis{
			NewTokenAnalysis(true, SYMBOL, "["),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(true, SYMBOL, "]"),
		}
		a.stack.pushAtPos(i, arr)
		a.identifyToken(i)
	} else {
		a.stack.deletePos(i)
		a.identifyToken(i)
	}
}

func (a *Analyzer) letStatement(i int) {
	var arr []TokenAnalysis
	if a.tokens[i].token == "let" {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, KEYWORD, "let"),
			NewTokenAnalysis(false, VARNAME, ""),
			NewTokenAnalysis(false, ARRAY, ""),
			NewTokenAnalysis(true, SYMBOL, "="),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(true, DOTANDCOMA, ";"),
		}
	} else {
		panic("There is an error at line " + strconv.Itoa(a.tokens[i].line) + ":" + strconv.Itoa(a.tokens[i].column))
	}
	a.stack.pushAtPos(i, arr)
	a.identifyToken(i)
}

func (a *Analyzer) ifStatement(i int) {
	//fmt.Println("Entra al if")
	if a.tokens[i].token == "if" {

		arr := []TokenAnalysis{
			NewTokenAnalysis(true, KEYWORD, "if"),
			NewTokenAnalysis(true, LPARENT, "("),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(true, RPARENT, ")"),
			NewTokenAnalysis(true, SYMBOL, "{"),
			NewTokenAnalysis(false, STATEMENTS, ""),
			NewTokenAnalysis(true, SYMBOL, "}"),
			NewTokenAnalysis(false, ELSESTATEMENT, ""),
		}
		a.stack.pushAtPos(i, arr)
		a.identifyToken(i)
	} else {

		panic("There is an error at line " + strconv.Itoa(a.tokens[i].line) + ":" + strconv.Itoa(a.tokens[i].column));
	}
}

func (a *Analyzer) elseStatement(i int) {
	if a.tokens[i].token == "else" {
		arr := []TokenAnalysis{
			NewTokenAnalysis(true, KEYWORD, "else"),
			NewTokenAnalysis(true, SYMBOL, "{"),
			NewTokenAnalysis(false, STATEMENTS, ""),
			NewTokenAnalysis(true, SYMBOL, "}"),
		}
		a.stack.pushAtPos(i, arr)
		a.identifyToken(i)
	} else {
		a.stack.deletePos(i)
		a.identifyToken(i)
	}
}

func (a *Analyzer) whileStatement(i int) {
	if a.tokens[i].token == "while" {
		arr := []TokenAnalysis{
			NewTokenAnalysis(true, KEYWORD, "while"),
			NewTokenAnalysis(true, LPARENT, "("),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(true, RPARENT, ")"),
			NewTokenAnalysis(true, SYMBOL, "{"),
			NewTokenAnalysis(false, STATEMENTS, ""),
			NewTokenAnalysis(true, SYMBOL, "}"),
			NewTokenAnalysis(false, ELSESTATEMENT, ""),
		}
		a.stack.pushAtPos(i, arr)
		a.identifyToken(i)
	} else {
		panic("There is an error at line " + strconv.Itoa(a.tokens[i].line) + ":" + strconv.Itoa(a.tokens[i].column))
	}
}

func (a *Analyzer) doStatement(i int) {

	if a.tokens[i].token == "do" {
		arr := []TokenAnalysis{
			NewTokenAnalysis(true, KEYWORD, "do"),
			NewTokenAnalysis(false, SUBROUTINECALL, ""),
			NewTokenAnalysis(true, DOTANDCOMA, ";"),
		}
		a.stack.pushAtPos(i, arr)
		a.identifyToken(i)
	} else {
		panic("There is an error at line " + strconv.Itoa(a.tokens[i].line) + ":" + strconv.Itoa(a.tokens[i].column))
	}

}

func (a *Analyzer) returnStatement(i int) {

	if a.tokens[i].token == "return" {
		arr := []TokenAnalysis{
			NewTokenAnalysis(true, KEYWORD, "return"),
			NewTokenAnalysis(false, EXPRESSIONCOND, ""),
			NewTokenAnalysis(true, SYMBOL, ";"),
		}
		a.stack.pushAtPos(i, arr)
		a.identifyToken(i)
	} else {
		panic("There is an error at line " + strconv.Itoa(a.tokens[i].line) + ":" + strconv.Itoa(a.tokens[i].column))
	}
}

func (a *Analyzer) expression(i int) {

	arr := []TokenAnalysis{
		NewTokenAnalysis(false, TERM, ""),
		NewTokenAnalysis(false, TERMCOND, ""),
	}
	a.stack.pushAtPos(i, arr)
	a.identifyToken(i)
}

func (a *Analyzer) term(i int) {

	var arr []TokenAnalysis
	if a.tokens[i].tokenType == INTEGERCONSTANT {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, INTEGERCONSTANT, a.tokens[i].token),
		}
	} else if a.tokens[i].tokenType == STRINGCONST {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, STRINGCONST, a.tokens[i].token),
		}
	} else if a.tokens[i].tokenType == KEYWORD {
		if a.tokens[i].token == "true" {
			arr = []TokenAnalysis{
				NewTokenAnalysis(true, KEYWORD, "true"),
			}
		} else if a.tokens[i].token == "false" {
			arr = []TokenAnalysis{
				NewTokenAnalysis(true, KEYWORD, "false"),
			}
		} else if a.tokens[i].token == "null" {
			arr = []TokenAnalysis{
				NewTokenAnalysis(true, KEYWORD, "null"),
			}
		} else if a.tokens[i].token == "this" {
			arr = []TokenAnalysis{
				NewTokenAnalysis(true, KEYWORD, "this"),
			}
		}
	} else if a.tokens[i].tokenType == OP || a.tokens[i].tokenType == UNARYOP {
		if a.tokens[i].token == "-" {
			arr = []TokenAnalysis{
				NewTokenAnalysis(true, UNARYOP, "-"),
				NewTokenAnalysis(false, TERM, ""),
			}
		} else if a.tokens[i].token == "~" || a.tokens[i].token == "!" {
			arr = []TokenAnalysis{
				NewTokenAnalysis(true, UNARYOP, "!"),
				NewTokenAnalysis(false, TERM, ""),
			}
		}
	} else {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, TERMPROD, ""),
		}
	}
	a.stack.pushAtPos(i, arr)
	a.identifyToken(i)
}

func (a *Analyzer) termProd(i int) {
	//fmt.Println("Entra a termProd")
	//fmt.Println(a.tokens[i+1].token)
	var arr []TokenAnalysis
	if a.tokens[i].token == "(" {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, LPARENT, "("),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(true, RPARENT, ")"),
		}
	} else if a.tokens[i+1].token == "(" || a.tokens[i+1].token == "." {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, SUBROUTINECALL, ""),
		}
	} else if a.tokens[i].tokenType == IDENTIFIER {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, VARNAME, ""),
			NewTokenAnalysis(false, ARRAY, ""),
		}
	} else {
		panic("Unknown token <" + a.tokens[i].token + "> at line " + strconv.Itoa(a.tokens[i].line) + ":" + strconv.Itoa(a.tokens[i].column))
	}
	a.stack.pushAtPos(i, arr)
	a.identifyToken(i)
}

func (a *Analyzer) termCond(i int) {


	if a.tokens[i].tokenType == OP {
		arr := []TokenAnalysis{
			NewTokenAnalysis(true, OP, a.tokens[i].token),
			NewTokenAnalysis(false, TERM,""),
			NewTokenAnalysis(false, TERMCOND, ""),
		}
		a.stack.pushAtPos(i, arr)
		a.identifyToken(i)
	} else {
		a.stack.deletePos(i)
		a.identifyToken(i)
	}
}

func (a *Analyzer) expressionCond(i int) {
	if a.tokens[i].token == ";" {
		a.stack.deletePos(i)
		a.identifyToken(i)
	} else {
		arr := []TokenAnalysis{
			NewTokenAnalysis(false, EXPRESSION, ""),
		}
		a.stack.pushAtPos(i, arr)
		a.identifyToken(i)
	}
}

func (a *Analyzer) subRoutineCall(i int) {
	var arr []TokenAnalysis
	if a.tokens[i+1].token == "." {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false,VARNAME,""),
			NewTokenAnalysis(true,SYMBOL,"."),
			NewTokenAnalysis(false,SUBROUTINENAME,""),
			NewTokenAnalysis(true,LPARENT,"("),
			NewTokenAnalysis(false,EXPRESSIONLIST,""),
			NewTokenAnalysis(true,RPARENT,")"),
		}
	} else {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false,SUBROUTINENAME,""),
			NewTokenAnalysis(true,LPARENT,"("),
			NewTokenAnalysis(false,EXPRESSIONLIST,""),
			NewTokenAnalysis(true,RPARENT,")"),
		}
	}
	a.stack.pushAtPos(i, arr)
	a.identifyToken(i)
}


func (a *Analyzer) expressionList(i int) {

	if a.tokens[i].token==")"{
		a.stack.deletePos(i)
		a.identifyToken(i)
	}else{
		arr:= []TokenAnalysis{
			NewTokenAnalysis(false,EXPRESSION,""),
			NewTokenAnalysis(false,EXPRESSIONLISTCOND,""),
		}
		a.stack.pushAtPos(i, arr)
		a.identifyToken(i)
	}
}



func (a *Analyzer) expressionListCond(i int) {
	//fmt.Println("entra 2")
	if a.tokens[i].token==")"{
		a.stack.deletePos(i)
		a.identifyToken(i)
	}else if a.tokens[i].token==","{
		arr:= []TokenAnalysis{
			NewTokenAnalysis(true,SYMBOL,","),
			NewTokenAnalysis(false,EXPRESSION,""),
			NewTokenAnalysis(false,EXPRESSIONLISTCOND,""),
		}
		a.stack.pushAtPos(i, arr)
		a.identifyToken(i)
	}
}

func (a *Analyzer) subRoutineDecCond(i int){

	if a.tokens[i].token =="}"{
		a.stack.deletePos(i)
		//a.stack.identifyToken(i)
	}else{
		arr:=[]TokenAnalysis{
			NewTokenAnalysis(false,SUBROUTINEDEC,""),
		}
		a.stack.pushAtPos(i,arr)
		a.identifyToken(i)
	}
}

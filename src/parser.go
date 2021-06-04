package main

import (
	"fmt"
	"strconv"
)


var className string


var currentScope string
var backupScope	string

var countStatic int
var countField int
var countLocal int
var countArgument int

var tokensPos int
var subroutineTemp string


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
	listaPrecedencia map[string]*Node
	text string
}

func NewAnalyzer(tokens []Token) *Analyzer {
	return &Analyzer{
		tokens: tokens,
		stack:  *NewStack(),
		listaPrecedencia: make(map[string]*Node),
		text: "",
	}
}

func (a *Analyzer) Analyze() {
	a.stack.push(NewTokenAnalysis(false, CLASS, ""))
	i := 0
	tokensPos = 0
	a.identifyToken(i)
}

func (a *Analyzer) identifyToken(i int) {
	//fmt.Println(a.text)
	if a.stack.pos(tokensPos).terminalToken {
		if a.tokens[i].token == a.stack.pos(tokensPos).stringLiteral {
			i++
			tokensPos++
			a.identifyToken(i)
		} else {
			panic("There was an error at line " + strconv.Itoa(a.tokens[i].line) + " symbol <" + a.tokens[i].token + "> not found")
		}
	} else {
		switch a.stack.pos(tokensPos).tokenType {
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
		case INIEXP:
			//i++
			tokensPos++
			a.identifyToken(i)
			break
		case FINEXP:
			//i++
			tokensPos++
			a.identifyToken(i)
			break
		case CURRENTSCOPE:
			backupScope=currentScope

			currentScope=a.stack.pos(tokensPos-1).stringLiteral
			a.listaPrecedencia[a.stack.pos(tokensPos-1).stringLiteral]= NewNode([][]string{})
			tokensPos++
			a.identifyToken(i)
			break
		case BACKUPSCOPE:
			a.listaPrecedencia[currentScope].next=a.listaPrecedencia[a.stack.pos(tokensPos).stringLiteral]

			currentScope=a.stack.pos(tokensPos).stringLiteral
			//a.listaPrecedencia[currentScope].imprimirLista()
			//i++
			tokensPos++
			a.identifyToken(i)
			break

		case COMPILEFUNCTION:
			//fmt.Println(a.listaPrecedencia[currentScope].data)
			a.text+="function "+className+"."+currentScope+" "+strconv.Itoa(countLocal)+"\n"
			a.stack.deletePos(tokensPos)
			a.identifyToken(i)
			break

		case CONSTRUCTOR:
			a.text+="push constant "+strconv.Itoa(countField)+"\n"
			a.text+="call Memory.alloc 1\n"
			a.text+="pop pointer 0\n"
			a.stack.deletePos(tokensPos)
			a.identifyToken(i)
			break
		case METHOD:
			a.text+="push argument 0\n"
			a.text+="pop pointer 0\n"
			a.stack.deletePos(tokensPos)
			a.identifyToken(i)
			break

		case COMPILESUBROUTINECALL:
			a.text+="call "+subroutineTemp+" "+strconv.Itoa(countArgument)+"\n"
			a.stack.deletePos(tokensPos)
			subroutineTemp=""
			a.identifyToken(i)
			break
		}


	}
}


func (a *Analyzer) getText() string {

	return a.text

}



func (a *Analyzer) axioma(i int) {
	arr := []TokenAnalysis{
		NewTokenAnalysis(true, CLASS, "class"),
		NewTokenAnalysis(false, CLASSNAME, ""),
		NewTokenAnalysis(false, CURRENTSCOPE, a.tokens[i+1].token),
		NewTokenAnalysis(true, SYMBOL, "{"),
		NewTokenAnalysis(false, CLASSVARDEC, ""),
		NewTokenAnalysis(false, SUBROUTINEDEC, ""),
		NewTokenAnalysis(false, BACKUPSCOPE, currentScope),
		NewTokenAnalysis(true, SYMBOL, "}"),
	}
	a.stack.pushAtPos(tokensPos, arr)
	a.identifyToken(i)
}

func (a *Analyzer) className(i int) {
	arr := []TokenAnalysis{
		NewTokenAnalysis(true, IDENTIFIER, a.tokens[i].token),
	}
	className= a.tokens[i].token
	currentScope = a.tokens[i].token
	a.stack.pushAtPos(tokensPos, arr)
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
		tempMat := [][]string{{a.tokens[i+2].token,a.tokens[i+1].token,"static",strconv.Itoa(countStatic)}}
		countStatic++

		if a.listaPrecedencia[currentScope]==nil{
			a.listaPrecedencia[currentScope] = NewNode(tempMat)
		}else{
			a.listaPrecedencia[currentScope].data = append( a.listaPrecedencia[currentScope].data,tempMat...)
		}

	} else if "field" == a.tokens[i].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, FIELD, "field"),
			NewTokenAnalysis(false, TYPE, ""),
			NewTokenAnalysis(false, VARNAME, ""),
			NewTokenAnalysis(true, DOTANDCOMA, ";"),
			NewTokenAnalysis(false, CLASSVARDEC, ""),
		}

		tempMat := [][]string{{a.tokens[i+2].token,a.tokens[+1].token,"field",strconv.Itoa(countField)}}
		countField++
		if a.listaPrecedencia[currentScope]==nil{
			a.listaPrecedencia[currentScope] = NewNode(tempMat)
		}else{

			a.listaPrecedencia[currentScope].data = append( a.listaPrecedencia[currentScope].data,tempMat...)
		}
	}
	a.stack.pushAtPos(tokensPos, arr)
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
	a.stack.pushAtPos(tokensPos, arr)
	//a.stack.ToString()
	a.identifyToken(i)
}

func (a *Analyzer) varName(i int) {
	arr := []TokenAnalysis{
		NewTokenAnalysis(true, IDENTIFIER, a.tokens[i].token),
	}

	a.stack.pushAtPos(tokensPos, arr)
	//i++
	a.identifyToken(i)
}

func (a *Analyzer) subRoutineDec(i int) {
	countLocal = 0
	countArgument = 0
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
			NewTokenAnalysis(false, CURRENTSCOPE, a.tokens[i+1].token),
			NewTokenAnalysis(true, LPARENT, "("),
			NewTokenAnalysis(false, PARAMETERLIST, ""),
			NewTokenAnalysis(true, RPARENT, ")"),
			NewTokenAnalysis(false, SUBROUTINEBODY, ""),
			NewTokenAnalysis(false, BACKUPSCOPE, currentScope),
			NewTokenAnalysis(false, SUBROUTINEDECCOND, ""),

		}...)
	} else {

		if "constructor" == a.tokens[i].token{
			arr = append(arr, []TokenAnalysis{
				NewTokenAnalysis(false, TYPE, ""),
				NewTokenAnalysis(false, SUBROUTINENAME, ""),
				NewTokenAnalysis(false, CURRENTSCOPE, a.tokens[i+1].token),
				NewTokenAnalysis(true, LPARENT, "("),
				NewTokenAnalysis(false, PARAMETERLIST, ""),
				NewTokenAnalysis(true, RPARENT, ")"),
				NewTokenAnalysis(false, SUBROUTINEBODY, ""),
				NewTokenAnalysis(false, CONSTRUCTOR, ""),
				NewTokenAnalysis(false, BACKUPSCOPE, currentScope),
				NewTokenAnalysis(false, SUBROUTINEDECCOND, ""),
			}...)
		}else if "method" == a.tokens[i].token{
			arr = append(arr, []TokenAnalysis{
				NewTokenAnalysis(false, TYPE, ""),
				NewTokenAnalysis(false, SUBROUTINENAME, ""),
				NewTokenAnalysis(false, CURRENTSCOPE, a.tokens[i+1].token),
				NewTokenAnalysis(true, LPARENT, "("),
				NewTokenAnalysis(false, PARAMETERLIST, ""),
				NewTokenAnalysis(true, RPARENT, ")"),
				NewTokenAnalysis(false, SUBROUTINEBODY, ""),
				NewTokenAnalysis(false, METHOD, ""),
				NewTokenAnalysis(false, BACKUPSCOPE, currentScope),
				NewTokenAnalysis(false, SUBROUTINEDECCOND, ""),
			}...)
		}
	}
	a.stack.pushAtPos(tokensPos, arr)

	//a.stack.ToString()
	a.identifyToken(i)
}

func (a *Analyzer) subRoutineName(i int) {
	arr := []TokenAnalysis{
		NewTokenAnalysis(true, IDENTIFIER, a.tokens[i].token),
	}
	a.stack.pushAtPos(tokensPos, arr)
	a.identifyToken(i)
}

func (a *Analyzer) parameterList(i int) {
	var arr []TokenAnalysis


	if ")" == a.tokens[i].token {
		a.stack.deletePos(tokensPos)

		a.identifyToken(i)
		return
	} else if ")" == a.tokens[i+2].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, TYPE, ""),
			NewTokenAnalysis(false, VARNAME, ""),
		}

		tempMat := [][]string{{a.tokens[i+1].token,a.tokens[i].token,"argument",strconv.Itoa(countArgument)}}
		countArgument++
		if a.listaPrecedencia[currentScope]==nil{
			a.listaPrecedencia[currentScope] = NewNode(tempMat)
		}else{
			a.listaPrecedencia[currentScope].data = append( a.listaPrecedencia[currentScope].data,tempMat...)
		}
	} else if "," == a.tokens[i+2].token {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false, TYPE, ""),
			NewTokenAnalysis(false, VARNAME, ""),
			NewTokenAnalysis(true, COMA, ","),
			NewTokenAnalysis(false, PARAMETERLIST, ""),
		}
		tempMat := [][]string{{a.tokens[i+1].token,a.tokens[i].token,"argument",strconv.Itoa(countArgument)}}
		countArgument++
		if a.listaPrecedencia[currentScope]==nil{
			a.listaPrecedencia[currentScope] = NewNode(tempMat)
		}else{
			a.listaPrecedencia[currentScope].data = append( a.listaPrecedencia[currentScope].data,tempMat...)
		}
	}
	a.stack.pushAtPos(tokensPos, arr)
	//i++
	a.identifyToken(i)
}

func (a *Analyzer) subRoutineBody(i int) {
	//fmt.Println("acá se inserta el subroutinebody")
	arr := []TokenAnalysis{
		NewTokenAnalysis(true, SYMBOL, "{"),
		NewTokenAnalysis(false, VARDEC, ""),
		NewTokenAnalysis(false,COMPILEFUNCTION,""),
		NewTokenAnalysis(false, STATEMENTS, ""),
		NewTokenAnalysis(true, SYMBOL, "}"),
	}
	a.stack.pushAtPos(tokensPos, arr)
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
		tempMat := [][]string{{a.tokens[i+2].token,a.tokens[i+1].token,"local",strconv.Itoa(countLocal)}}
		countLocal++
		if a.listaPrecedencia[currentScope]==nil{
			a.listaPrecedencia[currentScope] = NewNode(tempMat)
		}else{
			a.listaPrecedencia[currentScope].data = append( a.listaPrecedencia[currentScope].data,tempMat...)

		}


		a.stack.pushAtPos(tokensPos, arr)
	} else {
		a.stack.deletePos(tokensPos)
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
		a.stack.deletePos(tokensPos)
		a.identifyToken(i)
		return
	}
	a.stack.pushAtPos(tokensPos, arr)
	a.identifyToken(i)
}
func (a *Analyzer) array(i int) {
	if a.tokens[i].token == "[" {
		arr := []TokenAnalysis{
			NewTokenAnalysis(true, SYMBOL, "["),
			NewTokenAnalysis(false,INIEXP,""),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(false,FINEXP,""),
			NewTokenAnalysis(true, SYMBOL, "]"),
		}
		a.stack.pushAtPos(tokensPos, arr)
		a.identifyToken(i)
	} else {
		a.stack.deletePos(tokensPos)
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
			NewTokenAnalysis(false,INIEXP,""),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(false,FINEXP,""),
			NewTokenAnalysis(true, DOTANDCOMA, ";"),
		}
	} else {
		panic("There is an error at line " + strconv.Itoa(a.tokens[i].line) + ":" + strconv.Itoa(a.tokens[i].column))
	}
	a.stack.pushAtPos(tokensPos, arr)
	a.identifyToken(i)
}

func (a *Analyzer) ifStatement(i int) {
	//fmt.Println("Entra al if")
	if a.tokens[i].token == "if" {

		arr := []TokenAnalysis{
			NewTokenAnalysis(true, KEYWORD, "if"),
			NewTokenAnalysis(true, LPARENT, "("),
			NewTokenAnalysis(false,INIEXP,""),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(false,FINEXP,""),
			NewTokenAnalysis(true, RPARENT, ")"),
			NewTokenAnalysis(true, SYMBOL, "{"),
			NewTokenAnalysis(false, STATEMENTS, ""),
			NewTokenAnalysis(true, SYMBOL, "}"),
			NewTokenAnalysis(false, ELSESTATEMENT, ""),
		}
		a.stack.pushAtPos(tokensPos, arr)
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
		a.stack.pushAtPos(tokensPos, arr)
		a.identifyToken(i)
	} else {
		a.stack.deletePos(tokensPos)
		a.identifyToken(i)
	}
}

func (a *Analyzer) whileStatement(i int) {
	if a.tokens[i].token == "while" {
		arr := []TokenAnalysis{
			NewTokenAnalysis(true, KEYWORD, "while"),
			NewTokenAnalysis(true, LPARENT, "("),
			NewTokenAnalysis(false,INIEXP,""),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(false,FINEXP,""),
			NewTokenAnalysis(true, RPARENT, ")"),
			NewTokenAnalysis(true, SYMBOL, "{"),
			NewTokenAnalysis(false, STATEMENTS, ""),
			NewTokenAnalysis(true, SYMBOL, "}"),
			NewTokenAnalysis(false, ELSESTATEMENT, ""),
		}
		a.stack.pushAtPos(tokensPos, arr)
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
			NewTokenAnalysis(false, COMPILESUBROUTINECALL, ""),
			NewTokenAnalysis(true, DOTANDCOMA, ";"),
		}
		a.stack.pushAtPos(tokensPos, arr)
		a.identifyToken(i)
	} else {
		panic("There is an error at line " + strconv.Itoa(a.tokens[i].line) + ":" + strconv.Itoa(a.tokens[i].column))
	}

}

func (a *Analyzer) returnStatement(i int) {

	if a.tokens[i].token == "return" {
		arr := []TokenAnalysis{
			NewTokenAnalysis(true, KEYWORD, "return"),
			NewTokenAnalysis(false,INIEXP,""),
			NewTokenAnalysis(false, EXPRESSIONCOND, ""),
			NewTokenAnalysis(false,FINEXP,""),
			NewTokenAnalysis(true, SYMBOL, ";"),
		}
		a.stack.pushAtPos(tokensPos, arr)
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
	a.stack.pushAtPos(tokensPos, arr)
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
	a.stack.pushAtPos(tokensPos, arr)
	a.identifyToken(i)
}

func (a *Analyzer) termProd(i int) {
	//fmt.Println("Entra a termProd")
	//fmt.Println(a.tokens[i+1].token)
	var arr []TokenAnalysis
	if a.tokens[i].token == "(" {
		arr = []TokenAnalysis{
			NewTokenAnalysis(true, LPARENT, "("),
			NewTokenAnalysis(false,INIEXP,""),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(false,FINEXP,""),
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
	a.stack.pushAtPos(tokensPos, arr)
	a.identifyToken(i)
}

func (a *Analyzer) termCond(i int) {


	if a.tokens[i].tokenType == OP {
		arr := []TokenAnalysis{
			NewTokenAnalysis(true, OP, a.tokens[i].token),
			NewTokenAnalysis(false, TERM,""),
			NewTokenAnalysis(false, TERMCOND, ""),
		}
		a.stack.pushAtPos(tokensPos, arr)
		a.identifyToken(i)
	} else {
		a.stack.deletePos(tokensPos)
		a.identifyToken(i)
	}
}

func (a *Analyzer) expressionCond(i int) {
	if a.tokens[i].token == ";" {
		a.stack.deletePos(tokensPos)
		a.identifyToken(i)
	} else {
		arr := []TokenAnalysis{
			NewTokenAnalysis(false,INIEXP,""),
			NewTokenAnalysis(false, EXPRESSION, ""),
			NewTokenAnalysis(false,FINEXP,""),
		}
		a.stack.pushAtPos(tokensPos, arr)
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
		tempSub:=""

		for _,j := range a.listaPrecedencia[currentScope].data{
			if j[0]==a.tokens[i].token{
				tempSub = j[1]
				fmt.Println(tempSub)
			}
		}
		if tempSub==""{
			if a.listaPrecedencia[currentScope].next != nil{
				fmt.Println("Entra a la precedencia")
				for _,j := range a.listaPrecedencia[currentScope].next.data{
					if j[0]==a.tokens[i].token{
						tempSub = j[1]
						fmt.Println(tempSub)
					}
				}
			}else{
				tempSub=a.tokens[i].token
			}
		}
		subroutineTemp = tempSub + "." + a.tokens[i+2].token
	} else {
		arr = []TokenAnalysis{
			NewTokenAnalysis(false,SUBROUTINENAME,""),
			NewTokenAnalysis(true,LPARENT,"("),
			NewTokenAnalysis(false,EXPRESSIONLIST,""),
			NewTokenAnalysis(true,RPARENT,")"),
		}
		fmt.Println("se hace la vuelta")
		fmt.Println(className)

		subroutineTemp = className +"."+ a.tokens[i].token
		fmt.Println(subroutineTemp)
	}
	a.stack.pushAtPos(tokensPos, arr)
	a.identifyToken(i)
}


func (a *Analyzer) expressionList(i int) {

	if a.tokens[i].token==")"{
		a.stack.deletePos(tokensPos)
		a.identifyToken(i)
	}else{
		arr:= []TokenAnalysis{
			NewTokenAnalysis(false,INIEXP,""),
			NewTokenAnalysis(false,EXPRESSION,""),
			NewTokenAnalysis(false,FINEXP,""),
			NewTokenAnalysis(false,EXPRESSIONLISTCOND,""),
		}
		a.stack.pushAtPos(tokensPos, arr)
		a.identifyToken(i)
	}
}



func (a *Analyzer) expressionListCond(i int) {
	//fmt.Println("entra 2")
	if a.tokens[i].token==")"{
		a.stack.deletePos(tokensPos)
		a.identifyToken(i)
	}else if a.tokens[i].token==","{
		arr:= []TokenAnalysis{
			NewTokenAnalysis(true,SYMBOL,","),
			NewTokenAnalysis(false,INIEXP,""),
			NewTokenAnalysis(false,EXPRESSION,""),
			NewTokenAnalysis(false,FINEXP,""),
			NewTokenAnalysis(false,EXPRESSIONLISTCOND,""),
		}
		a.stack.pushAtPos(tokensPos, arr)
		a.identifyToken(i)
	}
}

func (a *Analyzer) subRoutineDecCond(i int){

	if a.tokens[i].token =="}"{
		a.stack.deletePos(tokensPos)
		//a.stack.identifyToken(i)
	}else{
		arr:=[]TokenAnalysis{
			NewTokenAnalysis(false,SUBROUTINEDEC,""),
		}
		a.stack.pushAtPos(tokensPos,arr)
		a.identifyToken(i)
	}
}

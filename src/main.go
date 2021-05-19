package main

import (
	"fmt"
	"os"
)

func main() {
	files := os.Args[1:]
	for _, i := range files {
		file, err := os.Open(i)
		if err != nil {
			panic(err)
		}
		var tokens []Token = []Token{}
		fmt.Printf("<%v>\n", CLASSVARDEC)
		lexer := NewLexer(file)
		for {
			tokenTemp := lexer.tokenize()
			if tokenTemp.tokenType == EOF {
				break
			}
			tokens = append(tokens, *tokenTemp)

			fmt.Printf("<%v> %s </%v>\n", tokenTemp.tokenType, tokenTemp.token, tokenTemp.tokenType)
		}
		analyzer := NewAnalyzer(tokens)
		analyzer.Analyze()
	}
}

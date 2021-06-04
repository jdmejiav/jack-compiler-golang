package main

import (
	"fmt"
	"io/ioutil"
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

		//fmt.Printf("<%v>\n", CLASSVARDEC)
		lexer := NewLexer(file)
		for {
			tokenTemp := lexer.tokenize()
			if tokenTemp.tokenType == EOF {
				break
			}
			if len(tokens)>0 {

				if (tokenTemp.token=="/" && tokens[len(tokens)-1].token=="/"){
					//fmt.Println("Llega acá")
					tokens=tokens[:len(tokens)-1]
					lexer.removeComent()
					continue
				}else{
					tokens = append(tokens, *tokenTemp)
				}
			}else{
				tokens = append(tokens, *tokenTemp)
			}
			//fmt.Println(tokenTemp.token)
			//fmt.Printf("<%v> %s </%v>\n", tokenTemp.tokenType, tokenTemp.token, tokenTemp.tokenType)

		}

		analyzer := NewAnalyzer(tokens)
		analyzer.Analyze()

		fmt.Printf("\n\n%s Grammar OK\n\n", i)
		fmt.Println(i[:len(i)-4]+"vm")


		f, _ := os.Create(i[:len(i)-4]+"vm")
    defer f.Close()

		b := []byte(analyzer.getText())
		ioutil.WriteFile(i[:len(i)-4]+"vm", b, 0644)
		
	}
}

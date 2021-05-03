package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	argsWithProg := os.Args

	if len(argsWithProg) != 1 {
		fmt.Println(argsWithProg[1])
		file, err := ioutil.ReadFile(argsWithProg[1])
		if err == nil {

			exp := string(file)

		} else {
			panic(err)
		}
	} else {
		fmt.Println("You must pass the route of the file .jack as argument")
		return
	}

}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	filename := "README.md"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(os.Stderr, err)
		return
	}
	fmt.Println(Lex(data))
}

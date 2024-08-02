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
	Parse(data, Lex(data))
	// res := Render(data, Parse(data, Lex(data)))
	// fmt.Println(res)
	// os.WriteFile("rendered.html", []byte(res), 0666)
}

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func readFile(filename string) ([]rune, error) {
	data, err := ioutil.ReadFile(filename)
	return bytes.Runes(data), err
}

func main() {
	filename := "README.md"
	data, err := readFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		tks []Token
		ast *Node
		res *string = new(string)
	)

	// defer debugInfo(data, res)()

	tks = Lex(data)
	ast = Parse(data, tks)
	*res = Render(data, ast)

	os.WriteFile("rendered.html", []byte(*res), 0666)
}

func debugInfo(data []rune, res *string) func() {
	for _, t := range Lex(data) {
		t.Print(data)
	}
	return func() {
		fmt.Println("Render result:")
		fmt.Println(*res)
	}
}

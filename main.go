package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	return

	filename := "README.md"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		tks []Token
		ast *Node
		res *string = new(string)
	)

	defer debugInfo(data, res)()

	tks = Lex(data)
	ast = Parse(data, tks)
	*res = Render(data, ast)

	os.WriteFile("rendered.html", []byte(*res), 0666)
}

func debugInfo(data []byte, res *string) func() {
	for _, t := range Lex(data) {
		t.Print(data)
	}
	return func() {
		fmt.Println("Render result:")
		fmt.Println(*res)
	}
}

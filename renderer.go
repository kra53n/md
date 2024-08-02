package main

import (
	"fmt"
	// "unsafe"
)

func Render(d []byte, tokens []Token) string {
	var res string
	var deque []Token
	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Type {
		case TokenH1, TokenH2, TokenH3, TokenH4, TokenH5, TokenH6:
			for ; i < len(tokens) && tokens[i].Type != TokenNewL; i++ {
				deque = append(deque, tokens[i])
			}
			res += load(d, &deque)
			res += upload(d, &deque)
		}
	}
	return res
}

func load(d []byte, tokens *[]Token) string {
	for _, t := range *tokens {
		fmt.Println(t.Type)
	}
	fmt.Println()
	var res string
	for _, t := range *tokens {
		switch t.Type {
		case TokenH1:
			res += "<h1>"
		case TokenH2:
			res += "<h2>"
		case TokenH3:
			res += "<h3>"
		case TokenBacktick:
			res += "<code>"
		case TokenSpace:
			res += " "
		case TokenPlainText:
			res += string(d[t.Start:t.End])
		}
	}
	return res
}

func upload(d []byte, tokens *[]Token) string {
	var res string
	for i := len(*tokens) - 1; i >= 0; i-- {
		switch (*tokens)[i].Type {
		case TokenH1:
			res += "</h1>\n"
		case TokenH2:
			res += "</h2>\n"
		case TokenH3:
			res += "</h3>\n"
		case TokenBacktick:
			res += "</code>"
		}
	}
	*tokens = (*tokens)[:0]
	return res
}

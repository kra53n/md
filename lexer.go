package main

// import "fmt"

type Lexer struct {
	Data []byte
	Pos  int
	Bol  int
	Row  int
}

func Lex(d []byte) []Token {
	l := Lexer{}
	l.Data = d
	tokens := l.tokenize()

	/* TODO(kra53n):
	 * Maybe this stage we can call analysis.
	 *
	 * Define `*` acceptance, it can be:
	 *   1) unordered list
	 *   2) italic
	 *   3) bold
	 *   4) italic and bold
	 *
	 * Group '`' if it can be group.
	 *
	 * Tables.
	 *
	 * New lines.
	 *
	 * Code inserts.
	 */
	tokens = analyse(tokens)

	return tokens
}

func analyse(tokens []Token) []Token {
	return tokens
}

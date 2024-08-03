package main

import "fmt"

type Lexer struct {
	Data []byte
	Pos  int
	Bol  int
	Row  int
}

type Token struct {
	Type  TokenType
	Start int
	End   int
	Row   int
}

type TokenType int

const (
	TokenNil TokenType = iota
	TokenH1
	TokenH2
	TokenH3
	TokenH4
	TokenH5
	TokenH6
	TokenNewL
	TokenSpace
	TokenAsterisk
	TokenBacktick
	TokenQuote
	TokenUnderscore
	TokenTilde
	TokenPlainText
	TokenLink
	TokenImg
	TokenUnorderedList
	TokenOrderedList
)

func Lex(d []byte) []Token {
	l := Lexer{}
	l.Data = d

	var tokens []Token
	var t Token
	for l.Pos = 0; l.Pos < len(l.Data); l.Pos++ {
		switch l.Data[l.Pos] {
		case '\r':
			t = l.newL()
		case ' ':
			t = l.space()
		case '#':
			t = l.header()
		case '*':
			t = l.single(TokenAsterisk)
		case '`':
			t = l.single(TokenBacktick)
		case '>':
			t = l.single(TokenQuote)
		case '_':
			t = l.single(TokenUnderscore)
		case '~':
			t = l.single(TokenTilde)
		case '!':
			t = l.exclamationMark()
		case '[':
			t = l.openBrac()
		case '-':
			t = l.unorderedList('-')
		case '+':
			t = l.unorderedList('+')
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t = l.digit()
		default:
			t = l.plainText()
		}

		if t.Type == TokenNil {
			continue
		}
		if shouldSkipDueNewLRepetitions(tokens, &t) {
			continue
		}
		if hasExcessSapce(tokens, &t) {
			tokens = tokens[:len(tokens)-1]
			continue
		}
		tokens = append(tokens, t)
	}
	i := len(tokens) - 1
	for ; i > 0 && tokens[i].Type == TokenNewL; i-- {
	}
	tokens = tokens[:i+1]
	return tokens
}

func (l *Lexer) newL() Token {
	t := Token{}
	if l.Pos < len(l.Data)-1 && l.Data[l.Pos+1] == '\n' {
		t.Type = TokenNewL
		t.Start = l.Pos
		t.End = l.Pos + 2
		t.Row = l.Row
		l.Pos++
		l.Row++
		return t
	}
	return t
}

func (l *Lexer) space() Token {
	i := l.Pos
	for ; i < len(l.Data) && l.Data[i] == ' '; i++ {
	}
	t := Token{
		Type:  TokenSpace,
		Start: l.Pos,
		End:   i,
		Row:   l.Row,
	}
	l.Pos = i - 1
	return t
}

func (l *Lexer) header() Token {
	t := Token{Row: l.Row}
	i := l.Pos
	for ; i < len(l.Data); i++ {
		switch {
		case l.Data[i] == '#' && i-l.Pos == 6:
			return t
		case l.Data[i] == '#':
			continue
		case l.Data[i] == '\r':
			return t
		case l.Data[i] == ' ':
			goto CheckContent
		default:
			return l.plainText()
		}
	}
CheckContent:
	for j := i; j < len(l.Data); j++ {
		switch {
		case chopChar(l.Data[j]):
			continue
		case l.Data[j] == '\r':
			return t
		default:
			goto Ok
		}
	}
Ok:
	t.Start = l.Pos
	t.End = i
	t.Type = TokenType(int(TokenH1) + i - l.Pos - 1)
	l.Pos = i
	return t
}

func (l *Lexer) single(tp TokenType) Token {
	return Token{
		Type:  tp,
		Start: l.Pos,
		End:   l.Pos + 1,
		Row:   l.Row,
	}
}

func (l *Lexer) exclamationMark() Token {
	if l.Pos+1 < len(l.Data) && l.Data[l.Pos+1] != '[' {
		return l.single(TokenPlainText)
	}
	openBracs := 1
	i := l.Pos + 2
	for ; i < len(l.Data) && openBracs != 0 && l.Data[i] != '\r'; i++ {
		switch l.Data[i] {
		case '[':
			openBracs++
		case ']':
			openBracs--
		}
	}
	if openBracs != 0 || i >= len(l.Data) || l.Data[i] != '(' {
		return l.single(TokenPlainText)
	}
	i++
	openBracs = 0
	for ; i < len(l.Data) && l.Data[i] != '\r'; i++ {
		switch l.Data[i] {
		case '(':
			openBracs++
		case ')':
			if openBracs != 0 {
				openBracs--
				continue
			}
			t := Token{
				Type:  TokenImg,
				Start: l.Pos,
				End:   i + 1,
				Row:   l.Row,
			}
			l.Pos = i
			return t
		}
	}
	return l.single(TokenPlainText)
}

func (l *Lexer) openBrac() Token {
	openBracs := 1
	i := l.Pos + 1
	for ; i < len(l.Data) && openBracs != 0 && l.Data[i] != '\r'; i++ {
		switch l.Data[i] {
		case '[':
			openBracs++
		case ']':
			openBracs--
		}
	}
	if openBracs != 0 || l.Data[i] != '(' {
		return l.single(TokenPlainText)
	}
	i++
	openBracs = 0
	for ; i < len(l.Data) && l.Data[i] != '\r'; i++ {
		switch l.Data[i] {
		case '(':
			return l.single(TokenPlainText)
		case ')':
			t := Token{
				Type:  TokenLink,
				Start: l.Pos,
				End:   i + 1,
				Row:   l.Row,
			}
			l.Pos = i
			return t
		}
	}
	return l.single(TokenPlainText)
}

func (l *Lexer) unorderedList(b byte) Token {
	t := Token{
		Type:  TokenUnorderedList,
		Start: l.Pos,
		End:   l.Pos + 1,
		Row:   l.Row,
	}
	if l.Pos == 0 {
		return t
	}
	i := l.Pos - 1
	for ; i > 0 && l.Data[i] == ' '; i-- {
	}
	if i == l.Pos-1 {
		i++
	}
	if l.Data[i] == b {
		return t
	}
	t.Type = TokenPlainText
	return t
}

func (l *Lexer) digit() Token {
	t := Token{
		Type:  TokenOrderedList,
		Start: l.Pos,
		End:   l.Pos + 1,
		Row:   l.Row,
	}
	i := l.Pos
	for ; i > 0 && l.Data[i] == ' '; i-- {
	}
	if !(i == l.Pos && i == 0 || l.Data[i] != '\r') {
		t.Type = TokenPlainText
		return t
	}
	i = l.Pos
	for ; i < len(l.Data); i++ {
		c := l.Data[i]
		if '0' <= c && c <= '9' {
			continue
		}
		switch c {
		case ' ', '\r', 0:
			t.Type = TokenPlainText
			t.Start = l.Pos
			t.End = i
			l.Pos = i - 1
			return t
		case '.', ')':
			t.Start = l.Pos
			t.End = i + 1
			l.Pos = i
			return t
		}
	}
	t.Type = TokenPlainText
	return t
}

func (l *Lexer) plainText() Token {
	i := l.Pos
	for ; i < len(l.Data); i++ {
		switch l.Data[i] {
		case '\r', ' ', '*', '`', '_', '~':
			goto End
		}
	}
End:
	t := Token{
		Type:  TokenPlainText,
		Start: l.Pos,
		End:   i,
		Row:   l.Row,
	}
	l.Pos = i - 1
	return t
}

func chopChar(c byte) bool {
	switch c {
	case ' ', '\t':
		return true
	}
	return false
}

func shouldSkipDueNewLRepetitions(tokens []Token, cur *Token) bool {
	if len(tokens) < 2 {
		return false
	}
	var prv1, prv2 TokenType
	prv1 = tokens[len(tokens)-1].Type
	prv2 = tokens[len(tokens)-2].Type
	return prv1 == prv2 && prv1 == cur.Type && cur.Type == TokenNewL
}

func hasExcessSapce(tokens []Token, cur *Token) bool {
	if len(tokens) < 1 {
		return false
	}
	return tokens[len(tokens)-1].Type == TokenSpace && cur.Type == TokenNewL
}

func (t *Token) Print(d []byte) {
	fmt.Printf("Type: %d Val: %s(%d, %d) Row: %d\n", t.Type, d[t.Start:t.End], t.Start, t.End, t.Row)
}

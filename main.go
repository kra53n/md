package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

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
)

func main() {
	filename := "README.md"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(os.Stderr, err)
		return
	}
	Lex(data)
}

func Lex(d []byte) []Token {
	l := Lexer{}
	l.Data = d

	var t Token
	var tokens []Token
	for l.Pos = 0; l.Pos < len(l.Data); l.Pos++ {
		switch l.Data[l.Pos] {
		case '\r':
			t = l.newL()
		case ' ':
			t = l.space()
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
		default:
			t = l.plainText()
		}

		if t.Type == TokenNil {
			continue
		}

		t.print(&l)
	}

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

/* NOTE(kra53n): wait for parsing
const (
	noteLowerS = "note"
	noteUpperS = "NOTE"

	tipLowerS = "tip"
	tipUpperS = "TIP"

	importantLowerS = "important"
	importantUpperS = "IMPORTANT"

	warningLowerS = "warning"
	warningUpperS = "WARNING"

	cautionLowerS = "caution"
	cautionUpperS = "CAUTION"
)
*/

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

func (t *Token) print(l *Lexer) {
	fmt.Printf("Type: %d Val: %s(%d, %d) Row: %d\n", t.Type, l.Data[t.Start:t.End], t.Start, t.End, t.Row)
}

func chopChar(c byte) bool {
	switch c {
	case ' ', '\t':
		return true
	}
	return false

}

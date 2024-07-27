package main

import (
	"errors"
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
	TokenErr TokenType = iota
	TokenNewL
	TokenH1
	TokenH2
	TokenH3
	TokenH4
	TokenH5
	TokenH6
	TokenSpace
	TokenCodeLine
	TokenPlainText
	TokenAsterisk
	TokenBacktick
	TokenUnderscoreL
	TokenUnderscoreR
	TokenTildeL
	TokenTildeR
	TokenQuote
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

func Lex(d []byte) {
	l := Lexer{}
	l.Data = d

	var t *Token
	var err error
	for l.Pos = 0; l.Pos < len(l.Data); l.Pos++ {
		switch l.Data[l.Pos] {
		case '\r':
			t, err = l.newL()
		case ' ':
			if t = l.space(); t == nil {
				continue
			}
		case '#':
			t = l.header()
		case '*':
			t = l.asterisk()
		case '`':
			t = l.asterisk()
			t.Type = TokenBacktick
		case '_':
			t = l.underscore()
		case '~':
			if t = l.underscore(); t != nil {
				switch t.Type {
				case TokenUnderscoreL:
					t.Type = TokenTildeL
				case TokenUnderscoreR:
					t.Type = TokenTildeR
				}
			}
		case '>':
			t = l.gt()
		default:
			t = l.plainText()
		}

		if t == nil {
			continue
		}

		t.print(&l)
		if t.Type == TokenErr {
			goto Err
		}
	}

	return
Err:
	fmt.Println("Error while lexing", *t)
	if err != nil {
		fmt.Println(err)
	}
}

func (l *Lexer) newL() (*Token, error) {
	t := Token{}
	if l.Pos < len(l.Data)-1 && l.Data[l.Pos+1] == '\n' {
		t.Type = TokenNewL
		t.Start = l.Pos
		t.End = l.Pos + 2
		t.Row = l.Row
		l.Pos++
		l.Row++
		return &t, nil
	}
	return &t, errors.New("new line identation error")
}

func (l *Lexer) space() *Token {
	if !l.beginning() {
		return nil
	}
	i := l.Pos
	for ; i < len(l.Data); i++ {
		switch l.Data[i] {
		case ' ':
			continue
		default:
			if i-l.Pos >= 4 && l.Data[i] != '\r' {
				l.Pos = i
				return l.codeLine()
			}
			return nil
		}
	}
	return nil
}

func (l *Lexer) beginning() bool {
	return l.Pos == 0 || l.Data[l.Pos-1] == '\n'
}

func (l *Lexer) codeLine() *Token {
	i := l.Pos
	for ; i < len(l.Data) && l.Data[i] != '\r'; i++ {
	}
	t := Token{
		Type: TokenCodeLine,
		Start: l.Pos,
		End: i,
		Row: l.Row,
	}
	l.Pos = i-1
	return &t
}

func (l *Lexer) header() *Token {
	t := Token{Type: TokenErr, Row: l.Row}
	i := l.Pos
	for ; i < len(l.Data); i++ {
		switch {
		case l.Data[i] == '#' && i-l.Pos == 6:
			return &t
		case l.Data[i] == '#':
			continue
		case l.Data[i] == '\r':
			return &t
		case l.Data[i] == ' ':
			goto CheckContent
		default:
			return l.plainText()
		}
	}
CheckContent:
	for j := i; j < len(l.data); j++ {
		switch {
		case chopchar(l.data[j]):
			continue
		case l.data[j] == '\r':
			return &t
		default:
			goto ok
		}
	}
Ok:
	t.Start = l.Pos
	t.End = i
	t.Type = TokenH1
	l.Pos = i
	return &t
}

func (l *Lexer) asterisk() *Token {
	return &Token{
		Type:  TokenAsterisk,
		Start: l.Pos,
		End:   l.Pos + 1,
		Row:   l.Row,
	}
}

func (l *Lexer) underscore() *Token {
	t := Token{
		Type:  TokenErr,
		Start: l.Pos,
		End:   l.Pos + 1,
		Row:   l.Row,
	}
	var lChoping, rChoping bool
	if len(l.Data) == 1 {
		goto PlainText
	}

	rChoping = chopChar(l.Data[l.Pos+1])
	if l.Pos == 0 && !rChoping {
		goto Left
	}
	if l.Pos-1 < 0 {
		goto PlainText
	}
	if l.Data[l.Pos-1] == '\n' && !rChoping {
		goto Left
	}

	lChoping = chopChar(l.Data[l.Pos-1])
	if l.Pos == len(l.Data)-1 && !lChoping || l.Data[l.Pos+1] == '\r' && !lChoping {
		goto Right
	}
	lChoping = chopChar(l.Data[l.Pos-1])
	rChoping = chopChar(l.Data[l.Pos+1])
	if l.Data[l.Pos-1] != '\n' && l.Data[l.Pos+1] != '\r' && lChoping != rChoping {
		switch {
		case lChoping:
			goto Left
		case rChoping:
			goto Right
		}
	}
PlainText:
	t.Type = TokenPlainText
	return &t
Left:
	t.Type = TokenUnderscoreL
	return &t
Right:
	t.Type = TokenUnderscoreR
	return &t
}

func (l *Lexer) gt() *Token {
	i := l.Pos-1
	for ; i >= 0 && l.Data[i] != '\n' && l.Data[i] != '>'; i-- {
	}
	if i < 0 || l.Pos-i < 5 {
		return &Token{
			Type: TokenQuote,
			Start: l.Pos,
			End: l.Pos+1,
			Row: l.Row,
		}
	}
	return l.codeLine()
}


func (l *Lexer) plainText() *Token {
	i := l.Pos
	for ; i < len(l.Data); i++ {
		switch l.Data[i] {
		case '\r', '*', '`', '_', '~':
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
	return &t
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

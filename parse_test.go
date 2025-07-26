package main

import (
	"fmt"
	"testing"
)

type parserTestSuite []parserTestCase

type parserTestCase struct {
	src []Token
	dst *Node
}

func parserHeaders() parserTestSuite {
	testSuite := parserTestSuite{}

	{
		src := []Token{
			Token{Type: TokenH1},
			Token{Type: TokenPlainText},
			Token{Type: TokenNewL},
			Token{Type: TokenH2},
			Token{Type: TokenPlainText},
		}

		dst := &Node{}

		h1 := &Node{T: Token{Type: TokenH1}}
		h1PlainText := &Node{T: Token{Type: TokenPlainText}}
		h1.FstChd = h1PlainText

		h2 := &Node{T: Token{Type: TokenH2}}
		h2PlainText := &Node{T: Token{Type: TokenPlainText}}
		h2.FstChd = h2PlainText

		dst.FstChd = h1
		h1.Nxt = h2

		testSuite = append(testSuite, parserTestCase{src, dst})
	}

	return testSuite
}

func TestParse(t *testing.T) {
	runParserTestSuite(t, parserHeaders())
}

func runParserTestSuite(t *testing.T, testSuite parserTestSuite) {
	for i, testCase := range testSuite {
		t.Run(fmt.Sprintf("%d", i), func(st *testing.T) {
			src := Parse([]rune{}, testCase.src)
			if ok := astEquals(src, testCase.dst); !ok {
				st.Errorf("src and dst trees are not equals:\nsrc:\n%v\ndst:\n%v\n", rootStringRepr(src), rootStringRepr(testCase.dst))
			}
		})
	}
}

func astEquals(n1, n2 *Node) bool {
	if n1 == nil || n2 == nil {
		return n1 == nil && n2 == nil
	}
	if n1.FstChd != nil || n2.FstChd != nil {
		if n1.FstChd == nil || n2.FstChd == nil {
			return false
		}
		if !astEquals(n1.FstChd, n2.FstChd) {
			return false
		}
	}
	if n1.Nxt != nil || n2.Nxt != nil {
		if n1.Nxt == nil || n2.Nxt == nil {
			return false
		}
		if !astEquals(n1.Nxt, n2.Nxt) {
			return false
		}
	}
	return n1.T.Type == n2.T.Type
}

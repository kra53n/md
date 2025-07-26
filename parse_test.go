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

var testParserHeaders = parserTestSuite{
	parserTestCase{
		src: []Token{
			Token{Type: TokenH1},
			Token{Type: TokenSpace},
			Token{Type: TokenPlainText},
			Token{Type: TokenNewL},
			Token{Type: TokenH2},
			Token{Type: TokenSpace},
			Token{Type: TokenPlainText},
		},
		dst: &Node{
			FstChd: &Node{
				T: Token{Type: TokenH1},
				FstChd: &Node{
					T: Token{Type: TokenSpace},
				},
				Nxt: &Node{
					T: Token{Type: TokenPlainText},
				},
			},
		},
	},
}

func TestParse(t *testing.T) {
	runParserTestSuite(t, testParserHeaders)
}

func runParserTestSuite(t *testing.T, testSuite parserTestSuite) {
	for i, testCase := range testSuite {
		t.Run(fmt.Sprintf("%d", i), func(st *testing.T) {
			src := Parse([]rune{}, testCase.src)
			yes := astEquals(src, testCase.dst)
			if !yes {
				st.Errorf("src and dst trees are not equals:\nsrc: %v\ndst: %v\n", src, testCase.dst)
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

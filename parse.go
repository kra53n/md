package main

import (
	"fmt"
	_ "unsafe"
)

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

type Node struct {
	T      Token
	Prt    *Node
	Nxt    *Node
	Prv    *Node
	FstChd *Node
	LstChd *Node
}

type Parser struct {
	d      []rune
	tokens []Token
	pos    int
	root   *Node
	cur    *Node
}

func Parse(d []rune, tokens []Token) *Node {
	p := Parser{
		d:      d,
		tokens: tokens,
	}
	return p.parse()
}

func (p *Parser) parse() *Node {
	p.root = new(Node)
	p.cur = p.root

	defer func() {
		fmt.Println("nani")
		println()
		println("Parse tree before processing:")
		println(rootStringRepr(p.root))
		println()

		p.root = processTree(p.root)

		println()
		println("Parse tree after processing:")
		println(rootStringRepr(p.root))
		println()
	}()

	for {
		if p.pos >= len(p.tokens) {
			break
		}
		switch p.tokens[p.pos].Type {
		case TokenH1,
			TokenH2,
			TokenH3,
			TokenH4,
			TokenH5,
			TokenH6,
			TokenUnorderedListType1,
			TokenUnorderedListType2,
			TokenUnorderedListType3,
			TokenOrderedListType1,
			TokenOrderedListType2:
			p.becomeNewChd()

		case TokenNewL:
			p.tokenNewL()

		case TokenPlainText,
			TokenSpace,
			TokenUnderscore,
			TokenBacktick,
			TokenCodeBlock:
			p.addChd()

		case TokenBoldStart,
			TokenItalicStart:
			p.addChd()
			p.cur = p.cur.LstChd
		case TokenBoldEnd,
			TokenItalicEnd:
			p.cur = p.cur.Prt

		case TokenTableStart, TokenTableHeaderStart, TokenTableBodyStart:
			p.addChd()
			p.cur = p.cur.LstChd
		case TokenTableHeaderEnd:
			for p.cur.T.Type != TokenTableHeaderStart {
				p.cur = p.cur.Prt
			}
			p.cur = p.cur.Prt
		case TokenTableBodyEnd:
			for p.cur.T.Type != TokenTableBodyStart {
				p.cur = p.cur.Prt
			}
			p.cur = p.cur.Prt
		case TokenTableLeftAlign, TokenTableCenterAlign, TokenTableRightAlign:
			switch p.cur.T.Type {
			case TokenTableLeftAlign, TokenTableCenterAlign, TokenTableRightAlign:
				p.cur = p.cur.Prt
			}
			p.addChd()
			p.cur = p.cur.LstChd
		case TokenTableRow:
			switch p.cur.T.Type {
			case TokenTableRow:
				p.cur = p.cur.Prt
			case TokenTableCol:
				p.cur = p.cur.Prt.Prt
			}
			p.addChd()
			p.cur = p.cur.LstChd
		case TokenTableCol:
			if p.cur.T.Type == TokenTableCol {
				p.cur = p.cur.Prt
			}
			p.addChd()
			p.cur = p.cur.LstChd
		case TokenTableEnd:
			p.cur = p.root

		}

		p.pos++
	}

	return p.root
}

func processTree(root *Node) *Node {
	var cur *Node

	cur = root.FstChd
	if cur == nil {
		return root
	}

	for cur != nil {
		switch cur.T.Type {
		case TokenUnorderedListType1,
			TokenUnorderedListType2,
			TokenUnorderedListType3,
			TokenOrderedListType1,
			TokenOrderedListType2:

			var ulNodeTokenType TokenType
			switch cur.T.Type {
			case TokenUnorderedListType1,
				TokenUnorderedListType2,
				TokenUnorderedListType3:
				ulNodeTokenType = TokenUnorderedList
			case TokenOrderedListType1,
				TokenOrderedListType2:
				ulNodeTokenType = TokenOrderedList
			}
			ulNode := &Node{
				T:      Token{Type: ulNodeTokenType},
				Prt:    cur.Prt,
				FstChd: cur,
			}
			if cur.Prv != nil {
				cur.Prv.Nxt = ulNode
			} else {
				cur.Prt.FstChd = ulNode
			}
			end := cur
			for end.Nxt != nil && end.Nxt.T.Type == cur.T.Type {
				end.Prt = ulNode
				end = end.Nxt
			}
			ulNode.Nxt = end.Nxt
			end.Nxt = nil
			ulNode.LstChd = end
			if ulNode.Nxt == nil {
				return root
			}
			cur = ulNode.Nxt
			cur.Prv = ulNode
		default:
			cur = cur.Nxt
		}
	}

	return root
}

func rootStringRepr(root *Node) string {
	return rootStringReprWithSpacesParam(root, 2)
}

func rootStringReprWithSpacesParam(root *Node, spaces int) string {
	var res, s string
	if root == nil {
		return res
	}
	for range spaces {
		s += " "
	}
	for i := root.FstChd; i != nil; i = i.Nxt {
		res += fmt.Sprintf("%s%d\n", s, i.T.Type)
		res += rootStringReprWithSpacesParam(i, spaces+2)
	}
	return res
}

func (p *Parser) addChd() {
	chd := &Node{T: p.tokens[p.pos]}
	if p.cur.FstChd == nil {
		p.cur.FstChd = chd
		p.cur.LstChd = chd
		chd.Prt = p.cur
	} else {
		p.cur.LstChd.Nxt = chd
		chd.Prv = p.cur.LstChd
		p.cur.LstChd = chd
		p.cur.LstChd.Prt = p.cur
	}
}

func (cur *Node) addNxt(nxt *Node) {
	nxt.Prt = cur.Nxt
	if cur.Nxt == nil {
		cur.Nxt = nxt
		nxt.Prt.LstChd = nxt
	} else {
		cur.Nxt.Prv = nxt
		nxt.Prv = cur
		nxt.Nxt = cur.Nxt
		cur.Nxt = nxt
	}
}

func (p *Parser) tokenNewL() {
	p.cur = p.root
	// for node := p.cur; node != nil; node = node.Prt {
	// 	switch node.T.Type {
	// 	case TokenH1, TokenH2, TokenH3, TokenH4, TokenH5, TokenH6:
	// 		p.cur = p.root
	// 	}
	// }
	// if p.tokens[p.pos+1].Type == TokenNewL {
	// 	p.cur = p.root
	// } else {
	// 	// TODO(kra53n):
	// 	// p.cur.addChd(&Node{T: Token{Type: TokenSpace}})
	// }
}

func (p *Parser) becomeNewChd() {
	p.addChd()
	p.cur = p.cur.LstChd
}

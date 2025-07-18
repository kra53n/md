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

/* Operations:
 *    2) addchild
 *    3) getroot
 */

// type Parser struct {
// 	Cur  *Node
// 	Prvs []*Node
// }

func Parse(d []rune, tokens []Token) *Node {
	/* TODO(kra53n):
	 * Maybe this stage we can call analysis.

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

	var root, cur *Node
	root = new(Node)
	cur = root

	// NOTE(kra53n): for more clarity make separate functions with names
	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Type {
		case TokenH1,
			TokenH2,
			TokenH3,
			TokenH4,
			TokenH5,
			TokenH6,
			TokenUnorderedListElem1,
			TokenUnorderedListElem2,
			TokenUnorderedListElem3,
			TokenOrderedListElem1,
			TokenOrderedListElem2:
			root.addChd(&Node{T: tokens[i]})
			cur = root.LstChd

		case TokenNewL:
			for node := cur; node != nil; node = node.Prt {
				switch node.T.Type {
				case TokenH1, TokenH2, TokenH3, TokenH4, TokenH5, TokenH6:
					cur = root
				}
			}
			if tokens[i+1].Type == TokenNewL {
				cur = root
			} else {
				cur.addChd(&Node{T: Token{Type: TokenSpace}})
			}

		case TokenPlainText,
			TokenSpace,
			TokenUnderscore,
			TokenBacktick,
			TokenCodeBlock:
			cur.addChd(&Node{T: tokens[i]})

		case TokenBoldStart,
			TokenItalicStart:
			cur.addChd(&Node{T: tokens[i]})
			cur = cur.LstChd
		case TokenBoldEnd,
			TokenItalicEnd:
			cur = cur.Prt

		case TokenTableStart, TokenTableHeaderStart, TokenTableBodyStart:
			cur.addChd(&Node{T: tokens[i]})
			cur = cur.LstChd
		case TokenTableHeaderEnd:
			for cur.T.Type != TokenTableHeaderStart {
				cur = cur.Prt
			}
			cur = cur.Prt
		case TokenTableBodyEnd:
			for cur.T.Type != TokenTableBodyStart {
				cur = cur.Prt
			}
			cur = cur.Prt
		case TokenTableLeftAlign, TokenTableCenterAlign, TokenTableRightAlign:
			switch cur.T.Type {
			case TokenTableLeftAlign, TokenTableCenterAlign, TokenTableRightAlign:
				cur = cur.Prt
			}
			cur.addChd(&Node{T: tokens[i]})
			cur = cur.LstChd
		case TokenTableRow:
			switch cur.T.Type {
			case TokenTableRow:
				cur = cur.Prt
			case TokenTableCol:
				cur = cur.Prt.Prt
			}
			cur.addChd(&Node{T: tokens[i]})
			cur = cur.LstChd
		case TokenTableCol:
			if cur.T.Type == TokenTableCol {
				cur = cur.Prt
			}
			cur.addChd(&Node{T: tokens[i]})
			cur = cur.LstChd
		case TokenTableEnd:
			cur = root

		}
	}

	// println()
	// println("Parse tree before processing:")
	// printRoot(root, 2)
	// println()

	// root = processTree(root)

	// println()
	// println("Parse tree after processing:")
	// printRoot(root, 2)
	// println()

	return root
}

func processTree(root *Node) *Node {
	var cur *Node

	cur = root.FstChd
	if cur == nil {
		return root
	}

	for cur != nil {
		switch cur.T.Type {
		case TokenUnorderedListElem1,
			TokenUnorderedListElem2,
			TokenUnorderedListElem3,
			TokenOrderedListElem1,
			TokenOrderedListElem2:

			var ulNodeTokenType TokenType
			switch cur.T.Type {
			case TokenUnorderedListElem1,
				TokenUnorderedListElem2,
				TokenUnorderedListElem3:
				ulNodeTokenType = TokenUnorderedList
			case TokenOrderedListElem1,
				TokenOrderedListElem2:
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

func printRoot(root *Node, spaces int) {
	if root == nil {
		return
	}
	var s string
	for i := 0; i < spaces; i++ {
		s += " "
	}
	for i := root.FstChd; i != nil; i = i.Nxt {
		fmt.Printf("%s%d\n", s, i.T.Type)
		printRoot(i, spaces+2)
	}
}

func (whose *Node) addChd(what *Node) {
	if whose.FstChd == nil {
		whose.FstChd = what
		whose.LstChd = what
		what.Prt = whose
	} else {
		whose.LstChd.Nxt = what
		what.Prv = whose.LstChd
		whose.LstChd = what
		whose.LstChd.Prt = whose
	}
}

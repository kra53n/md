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

func Parse(d []byte, tokens []Token) *Node {
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

	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Type {
		case TokenH1, TokenH2, TokenH3, TokenH4, TokenH5, TokenH6:
			root.addChd(&Node{T: tokens[i]})
			cur = root.LstChd

		// case TokenBacktick:
		// 	if cur.T.Type == TokenBacktick {
		// 		cur = cur.Prt
		// 	} else {
		// 		cur.addChd(&Node{T: tokens[i]})
		// 		cur = cur.LstChd
		// 	}

		// case TokenAsterisk:
		// 	if cur.T.Type == TokenAsterisk {
		// 		cur = cur.Prt
		// 	} else {
		// 		cur.addChd(&Node{T: tokens[i]})
		// 		cur = cur.LstChd
		// 	}

		case TokenNewL:
			tmp := cur
			for ; tmp != nil; tmp = tmp.Prt {
				switch tmp.T.Type {
				case TokenH1, TokenH2, TokenH3, TokenH4, TokenH5, TokenH6:
					cur = root
				}
			}
			if tokens[i+1].Type == TokenNewL {
				cur = root
			}

		case TokenPlainText:
			cur.addChd(&Node{T: tokens[i]})

		// case TokenTableStart:
		// 	cur.addChd(&Node{T: tokens[i]})
		// 	cur = cur.LstChd
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
			fmt.Println(cur.T.Type, cur.Prt.T.Type)
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

	printRoot(root, 2)

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

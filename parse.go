package main

import "fmt"

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

type Elem struct {
	Elems []Elem
	T     Token
}

type Parser struct {
	Cur  *Elem
	Prvs []*Elem
}

func Parse(d []byte, tokens []Token) {
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

	root := Elem{}
	p := Parser{Cur: &root}

	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Type {
		case TokenH1, TokenH2, TokenH3, TokenH4, TokenH5, TokenH6:
			p.Cur.append(tokens[i])
			p.Prvs = append(p.Prvs, p.Cur)
			p.Cur = &p.Cur.Elems[0]
			j := i + 1
			if tokens[j].Type == TokenSpace {
				j++
			}
			for ; j < len(tokens) && tokens[j].Type != TokenNewL; j++ {
				p.Cur.append(tokens[i])
			}
			if tokens[j].Type == TokenNewL {
				p.Cur = p.Prvs[len(p.Prvs)-1]
				p.Prvs = p.Prvs[:len(p.Prvs)-1]
				continue
			}

		}
	}

	for _, v := range root.Elems {
		fmt.Println("-->", v.T.Type)
	}
}

func (e *Elem) append(t Token) {
	e.Elems = append(e.Elems, Elem{T: t})
}

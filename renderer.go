package main

/* NOTE(kra53n): here we can implement render of the AST using 2 methods:
 * 1) using the stack
 * 2) using the recursion
 * We must test both variants and choose the best
 */
func Render(d []byte, root *Node) string {
	var res *string
	res = new(string)

	recursiveTraversal(res, d, root, new([]Token))
	return *res
}

func recursiveTraversal(res *string, d []byte, n *Node, ptr *[]Token) {
	if n == nil {
		return
	}
	var tag string
	for i := n.FstChd; i != nil; i = i.Nxt {
		*ptr = append(*ptr, i.T)
		_, tag = getOpenedTag(d, &i.T)
		*res = *res + tag
		recursiveTraversal(res, d, i, ptr)
	}
	if len(*ptr) > 0 {
		_, tag = getClosedTag(&(*ptr)[len(*ptr)-1])
		*res += tag
		*ptr = (*ptr)[:len(*ptr)-1]
	}
}

var tagNames map[TokenType]string = map[TokenType]string{
	TokenH1:               "h1",
	TokenH2:               "h2",
	TokenH3:               "h3",
	TokenH4:               "h4",
	TokenH5:               "h5",
	TokenH6:               "h6",
	TokenUnorderedList:    "ul",
	TokenUnorderedList1:   "li",
	TokenUnorderedList2:   "li",
	TokenUnorderedList3:   "li",
	TokenBoldStart:        "strong",
	TokenItalicStart:      "em",
	TokenTableStart:       "table",
	TokenTableHeaderStart: "thead",
	TokenTableCenterAlign: "th",
	TokenTableRow:         "tr",
	TokenTableCol:         "td",
}

func getOpenedTag(d []byte, t *Token) (int, string) {
	switch t.Type {
	case TokenPlainText:
		s := string(d[t.Start:t.End])
		return len(s), s
	case TokenUnderscore,
		TokenAsterisk,
		TokenBacktick:
		return 1, string(d[t.Start : t.Start+1])
	case TokenSpace:
		return 1, " "
	default:
		var tagString string = tagNames[t.Type]
		if len(tagString) > 0 {
			return len(tagString) + 2, ("<" + tagString + ">")
		}
	}
	return 0, ""
}

func getClosedTag(t *Token) (int, string) {
	var tagString string = tagNames[t.Type]
	if len(tagString) > 0 {
		return len(tagString) + 3, "</" + tagString + ">"
	}
	return 0, ""
}

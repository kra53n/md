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
	for i := n.FstChd; i != nil; i = i.Nxt {
		*ptr = append(*ptr, i.T)
		*res = *res + getOpenedTag(d, &i.T)
		recursiveTraversal(res, d, i, ptr)
	}
	if len(*ptr) > 0 {
		*res += getClosedTag(&(*ptr)[len(*ptr)-1])
		*ptr = (*ptr)[:len(*ptr)-1]
	}
}

func getOpenedTag(d []byte, t *Token) string {
	switch t.Type {
	case TokenH1:
		return "<h1>"
	case TokenH2:
		return "<h2>"
	case TokenH3:
		return "<h3>"
	case TokenBacktick:
		return "<code>"
	case TokenAsterisk:
		return "<i>"
	case TokenSpace:
		return " "
	case TokenPlainText:
		return string(d[t.Start:t.End])
	}
	return ""
}

func getClosedTag(t *Token) string {
	switch t.Type {
	case TokenH1:
		return "</h1>\n"
	case TokenH2:
		return "</h2>\n"
	case TokenH3:
		return "</h3>\n"
	case TokenAsterisk:
		return "</i>"
	case TokenBacktick:
		return "</code>"
	}
	return ""
}

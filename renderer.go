package main

/* NOTE(kra53n): here we can implement render of the AST using 2 methods:
 * 1) using the stack
 * 2) using the recursion
 * We must test both variants and choose the best
 */
func Render(d []rune, root *Node) string {
	var res *string
	res = new(string)

	recursiveTraversal(res, d, root, new([]Token))
	return *res
}

func recursiveTraversal(res *string, d []rune, n *Node, ptr *[]Token) {
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

var tagNames map[TokenType][]string = map[TokenType][]string{
	TokenH1:                 []string{"h1"},
	TokenH2:                 []string{"h2"},
	TokenH3:                 []string{"h3"},
	TokenH4:                 []string{"h4"},
	TokenH5:                 []string{"h5"},
	TokenH6:                 []string{"h6"},
	TokenUnorderedList:      []string{"ul"},
	TokenUnorderedListType1: []string{"li"},
	TokenUnorderedListType2: []string{"li"},
	TokenUnorderedListType3: []string{"li"},
	TokenOrderedList:        []string{"ol"},
	TokenOrderedListType1:   []string{"li"},
	TokenOrderedListType2:   []string{"li"},
	TokenCodeBlock:          []string{"pre", "code"},
	TokenBoldStart:          []string{"strong"},
	TokenItalicStart:        []string{"em"},
	TokenTableStart:         []string{"table"},
	TokenTableHeaderStart:   []string{"thead"},
	TokenTableCenterAlign:   []string{"th"},
	TokenTableRow:           []string{"tr"},
	TokenTableCol:           []string{"td"},
}

func getOpenedTag(d []rune, t *Token) (int, string) {
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
		var l int
		var res string
		for _, tagString := range tagNames[t.Type] {
			res += "<" + tagString + ">"
			l += len(tagString) + 2
		}
		switch t.Type {
		case TokenCodeBlock:
			res += string(d[t.Start:t.End])
			l += t.End - t.Start
		}
		return l, res
	}
	return 0, ""
}

func getClosedTag(t *Token) (int, string) {
	var l int
	var res string
	for i := len(tagNames[t.Type]) - 1; i >= 0; i-- {
		tagString := tagNames[t.Type][i]
		l += len(tagString) + 3
		res += "</" + tagString + ">"
	}
	return l, res
}

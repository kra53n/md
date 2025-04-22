package main

import "fmt"

type Lexer struct {
	Data []byte
	Pos  int
}

type Token struct {
	Type  TokenType
	Start int
	End   int
}

type TokenType int

// NOTE(kra53n):
// We can use asterisk sign for bold notation and for unordered lists, so
// we should do something with it. For example have a function or maybe we
// already have a solution for that, we must check it.
//
// Solution is that we are tokenize the chars and then process them in next
// stage (parsing).

const (
	TokenNil TokenType = iota
	TokenH1
	TokenH2
	TokenH3
	TokenH4
	TokenH5
	TokenH6
	TokenNewL
	TokenSpace
	TokenAsterisk
	TokenBacktick
	TokenDash
	TokenPlus
	TokenQuote
	TokenUnderscore
	TokenTilde
	TokenPlainText
	TokenLink
	TokenImg
	TokenUnorderedList
	TokenUnorderedList1 // TODO rename it to TokenUnorderedListElemX where X - digits
	TokenUnorderedList2
	TokenUnorderedList3
	TokenOrderedList // TODO make 2 ordered lists: with `.` notation and `)`
	TokenTableStart
	TokenTableHeaderStart
	TokenTableHeaderEnd
	TokenTableBodyStart
	TokenTableBodyEnd
	TokenTableLeftAlign
	TokenTableCenterAlign
	TokenTableRightAlign
	TokenTableRow
	TokenTableCol
	TokenTableEnd
	TokenCodeLine
	TokenCodeBlock
	TokenBoldStart
	TokenBoldEnd
	TokenItalicStart
	TokenItalicEnd
	TokenStrikeThrough
)

func Lex(d []byte) []Token {
	l := Lexer{}
	l.Data = d

	var tokens []Token
	var t Token
	var isTable bool
	for l.Pos = 0; l.Pos < len(l.Data); l.Pos++ {
		tokens, isTable = l.table(tokens)
		if isTable && l.tableStartsWithPipe() {
			continue
		}
		t = l.single()
		if t.Type == TokenNil {
			t = l.plainText()
		}

		if shouldSkipDueNewLRepetitions(tokens, &t) {
			continue
		}
		if hasExcessSapce(tokens, &t) {
			tokens = tokens[:len(tokens)-1]
			continue
		}
		tokens = append(tokens, t)
	}
	return analyze(d, tokens)
}

func (l *Lexer) table(tokens []Token) ([]Token, bool) {
	start := l.lineBeginning()

	if !l.isTable(start) {
		return tokens, false
	}

	tokens = append(tokens, Token{Type: TokenTableStart})
	{
		var headerPipes, i, tableDataPipes int
		i = start
		headerPipes, _ = l.tableGetPipesNum(start)
		tokens = append(tokens, Token{Type: TokenTableHeaderStart})
		{
			tokens = l.tableAppendHeaders(tokens, headerPipes, i)
		}
		tokens = append(tokens, Token{Type: TokenTableHeaderEnd})

		tokens = append(tokens, Token{Type: TokenTableBodyStart})
		{
			i = l.skipLine(start)
			i = l.skipLine(i)
			for tableDataPipes, _ = l.tableGetPipesNum(i); tableDataPipes == headerPipes; {
				tokens = l.tableAppendData(tokens, headerPipes, i)
				i = l.skipLine(i)
				tableDataPipes, _ = l.tableGetPipesNum(i)
			}
			l.Pos = i - 1
		}
		tokens = append(tokens, Token{Type: TokenTableBodyEnd})
	}
	tokens = append(tokens, Token{Type: TokenTableEnd})

	return tokens, true
}

func (l *Lexer) single() Token {
	switch l.Data[l.Pos] {
	case '\r':
		return l.newL()
	case ' ':
		return l.space()
	case '#':
		return l.header()
	case '*':
		return l.charToken(TokenAsterisk)
	case '`':
		return l.charToken(TokenBacktick)
	case '>':
		return l.charToken(TokenQuote)
	case '_':
		return l.charToken(TokenUnderscore)
	case '~':
		return l.charToken(TokenTilde)
	case '!':
		return l.exclamationMark()
	case '[':
		return l.openBrac()
	case '-':
		return l.charToken(TokenDash)
	case '+':
		return l.charToken(TokenPlus)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.digit()
	default:
		return Token{}
	}
}

func (l *Lexer) lineBeginning() int {
	i := l.Pos
	for i != 0 {
		if l.Data[i] == '\n' {
			i++
			break
		}
		i--
	}
	return i
}

func (l *Lexer) eof(pos int) bool {
	return pos >= len(l.Data)
}

func (l *Lexer) eol(pos int) bool {
	if l.eof(pos) {
		return true
	}
	switch l.Data[pos] {
	case '\r', '\n', 0:
		return true
	}
	return false
}

func (l *Lexer) skipLine(pos int) int {
	i := pos
	for i < len(l.Data) {
		if l.Data[i] == '\r' {
			return i + 2
		}
		i++
	}
	return i
}

func (l *Lexer) lTrim(beg int) int {
	i := beg
	for i < len(l.Data) && l.skipChar(i) {
		i++
	}
	return i
}

func (l *Lexer) rTrim(end int) int {
	i := end
	for i > 0 {
		if !l.skipChar(i) {
			return i + 1
		}
		i--
	}
	return i
}

func (l *Lexer) skipChar(pos int) bool {
	switch l.Data[pos] {
	case ' ', '\t':
		return true
	}
	return false
}

func (t *Token) Print(d []byte) {
	fmt.Printf("Type: %d Val: %s(%d, %d)\n", t.Type, d[t.Start:t.End], t.Start, t.End)
}

func (l *Lexer) isTable(start int) bool {
	var i, headerPipes, alignPipes int
	i = start
	headerPipes, i = l.tableGetPipesNum(i)
	if headerPipes == 0 {
		return false
	}
	i += 2
	alignPipes, _ = l.tableGetPipesNum(i)
	if alignPipes == 0 || headerPipes != alignPipes {
		return false
	}
	return l.isTableAlignsCorrect(i)
}

func (l *Lexer) tableStartsWithPipe() bool {
	return l.Data[l.Pos] == '|' && !l.tableIgnorePipe(l.Pos)
}

func (l *Lexer) tableGetPipesNum(start int) (int, int) {
	var i, pipes int
	i = start
	for !l.eol(i) {
		if l.Data[i] == '|' && !l.tableIgnorePipe(i) {
			pipes++
		}
		i++
	}
	return pipes, i
}

func (l *Lexer) tableIgnorePipe(pos int) bool {
	if pos-1 < 0 ||
		(pos == 0 || l.Data[pos-1] == '\n') && pos+1 < len(l.Data) && l.Data[pos+1] == ' ' ||
		pos+1 < len(l.Data) && l.Data[pos+1] == '\r' && l.skipChar(pos-1) {
		return true
	}
	return l.Data[pos-1] == '\\'
}

func (l *Lexer) isTableAlignsCorrect(start int) bool {
	for i := start; i < len(l.Data) && l.Data[i] != '\r'; i++ {
		if l.tableIgnorePipe(i) || l.skipChar(i) || l.Data[i] == '-' {
			continue
		}
		if l.Data[i] == ':' &&
			!(i-1 < 0 || l.Data[i-1] == '\n' || l.skipChar(i-1) ||
				i+1 == len(l.Data) || l.Data[i+1] == '\r' || l.skipChar(i+1)) {
			return false
		}
		if i+1 != len(l.Data) && l.skipChar(i+1) && l.Data[i] != '|' {
			i++
			for ; i < len(l.Data) && l.Data[i] != '-' && l.Data[i] != '|' && l.Data[i] != '\r'; i++ {
			}
			if l.Data[i] == '-' {
				return false
			}
		}
	}
	return true
}

func (l *Lexer) tableAppendHeaders(tokens []Token, pipes int, pos int) []Token {
	var iHeaders, iAligns int
	iHeaders = pos
	iAligns = l.skipLine(pos)
	tokens = append(tokens, Token{Type: TokenTableRow})
	for ; pipes >= 0; pipes-- {
		var startHeader, startAlign int

		startHeader = iHeaders
		startAlign = iAligns

		iHeaders = l.tableNxtPipe(iHeaders) - 1
		iAligns = l.tableNxtPipe(iAligns) - 1

		startHeader = l.tableLTrim(startHeader)
		fmt.Println("nnnn", startHeader)
		iHeaders = l.tableRTrim(iHeaders)

		tokens = append(tokens, Token{
			Type:  l.tableGetAlignType(l.tableLTrim(startAlign), l.tableRTrim(iAligns)),
			Start: startHeader,
			End:   startHeader,
		})

		var t Token
		for l.Pos = startHeader; l.Pos < iHeaders; l.Pos++ {
			t = l.single()
			if t.Type == TokenNil {
				t = l.plainText()
			}
			tokens = append(tokens, t)
		}

		iHeaders = l.tableNxtPipe(iHeaders) + 1
		iAligns = l.tableNxtPipe(iAligns) + 1
	}
	return tokens
}

func (l *Lexer) tableNxtPipe(pos int) int {
	var i int = pos
	for ; (l.Data[i] != '|' && !l.eol(i)) || l.tableIgnorePipe(i); i++ {
	}
	return i
}

func (l *Lexer) tableLTrim(beg int) int {
	if l.Data[beg] == '|' && l.tableIgnorePipe(beg) {
		beg++
	}
	fmt.Println("wassup", beg)
	beg = l.lTrim(beg)
	return beg
}

func (l *Lexer) tableRTrim(end int) int {
	if l.Data[end] == '\r' {
		end--
	}
	if l.tableIgnorePipe(end) {
		end--
	}
	end = l.rTrim(end)
	return end
}

func (l *Lexer) tableGetAlignType(beg int, end int) TokenType {
	if l.Data[beg] == ':' && l.Data[end-1] == ':' {
		return TokenTableCenterAlign
	} else if l.Data[beg] == ':' {
		return TokenTableLeftAlign
	} else if l.Data[end-1] == ':' {
		return TokenTableRightAlign
	} else {
		return TokenTableCenterAlign
	}
}

func (l *Lexer) tableAppendData(tokens []Token, pipes int, pos int) []Token {
	i := pos
	tokens = append(tokens, Token{
		Type:  TokenTableRow,
		Start: i,
		End:   i,
	})
	for ; pipes >= 0; pipes-- {
		var start int = i

		i = l.tableNxtPipe(i)

		start = l.tableLTrim(start)

		i = l.tableRTrim(i - 1)

		tokens = append(tokens, Token{
			Type:  TokenTableCol,
			Start: start,
			End:   start,
		})
		var t Token
		for l.Pos = start; l.Pos < i; l.Pos++ {
			t = l.single()
			if t.Type == TokenNil {
				t = l.plainText()
			}
			tokens = append(tokens, t)
		}
		i = l.tableNxtPipe(i) + 1
	}
	return tokens
}

func (l *Lexer) newL() Token {
	t := Token{}
	if l.Pos < len(l.Data)-1 && l.Data[l.Pos+1] == '\n' {
		t.Type = TokenNewL
		t.Start = l.Pos
		t.End = l.Pos + 2
		l.Pos++
		return t
	}
	return t
}

func (l *Lexer) space() Token {
	i := l.Pos
	for i < len(l.Data) && l.Data[i] == ' ' {
		i++
	}
	t := Token{
		Type:  TokenSpace,
		Start: l.Pos,
		End:   i,
	}
	l.Pos = i - 1
	return t
}

func (l *Lexer) header() Token {
	t := Token{}
	i := l.Pos
	for ; i < len(l.Data); i++ {
		switch {
		case l.Data[i] == '#' && i-l.Pos == 6:
			return t
		case l.Data[i] == '#':
			continue
		case l.Data[i] == '\r':
			return t
		case l.Data[i] == ' ':
			goto CheckContent
		default:
			return l.plainText()
		}
	}
CheckContent:
	for j := i; j < len(l.Data); j++ {
		switch {
		case l.skipChar(j):
			continue
		case l.Data[j] == '\r':
			return t
		default:
			goto Ok
		}
	}
Ok:
	t.Start = l.Pos
	t.End = i
	t.Type = TokenType(int(TokenH1) + i - l.Pos - 1)
	l.Pos = i
	return t
}

func (l *Lexer) charToken(tp TokenType) Token {
	return Token{
		Type:  tp,
		Start: l.Pos,
		End:   l.Pos + 1,
	}
}

func (l *Lexer) exclamationMark() Token {
	if l.Pos+1 < len(l.Data) && l.Data[l.Pos+1] != '[' {
		return l.charToken(TokenPlainText)
	}
	openBracs := 1
	i := l.Pos + 2
	for ; i < len(l.Data) && openBracs != 0 && l.Data[i] != '\r'; i++ {
		switch l.Data[i] {
		case '[':
			openBracs++
		case ']':
			openBracs--
		}
	}
	if openBracs != 0 || i >= len(l.Data) || l.Data[i] != '(' {
		return l.charToken(TokenPlainText)
	}
	i++
	openBracs = 0
	for ; i < len(l.Data) && l.Data[i] != '\r'; i++ {
		switch l.Data[i] {
		case '(':
			openBracs++
		case ')':
			if openBracs != 0 {
				openBracs--
				continue
			}
			t := Token{
				Type:  TokenImg,
				Start: l.Pos,
				End:   i + 1,
			}
			l.Pos = i
			return t
		}
	}
	return l.charToken(TokenPlainText)
}

func (l *Lexer) openBrac() Token {
	openBracs := 1
	i := l.Pos + 1
	for ; i < len(l.Data) && openBracs != 0 && l.Data[i] != '\r'; i++ {
		switch l.Data[i] {
		case '[':
			openBracs++
		case ']':
			openBracs--
		}
	}
	if openBracs != 0 || l.Data[i] != '(' {
		return l.charToken(TokenPlainText)
	}
	i++
	openBracs = 0
	for ; i < len(l.Data) && l.Data[i] != '\r'; i++ {
		switch l.Data[i] {
		case '(':
			return l.charToken(TokenPlainText)
		case ')':
			t := Token{
				Type:  TokenLink,
				Start: l.Pos,
				End:   i + 1,
			}
			l.Pos = i
			return t
		}
	}
	return l.charToken(TokenPlainText)
}

// func (l *Lexer) unorderedList(b byte) Token {
// 	// TODO define here what tokenunorderedlist type (1 or 2)
// 	var ttype TokenType
// 	switch b {
// 	case '-':
// 		ttype = TokenUnorderedList1
// 	case '*':
// 		ttype = TokenUnorderedList2
// 	case '+':
// 		ttype = TokenUnorderedList3
// 	default:
// 		panic("in markdown there is only 3 notations (`-`, `*`, `+`) for declaring the unordered list")
// 	}
// 	t := Token{
// 		Type:  ttype,
// 		Start: l.Pos,
// 		End:   l.Pos + 1,
// 	}
// 	if l.Pos == 0 {
// 		return t
// 	}
// 	i := l.Pos - 1
// 	for i > 0 && l.Data[i] == ' ' {
// 		i--
// 	}
// 	if i == l.Pos-1 {
// 		i++
// 	}
// 	if l.Data[i] == b {
// 		return t
// 	}
// 	t.Type = TokenPlainText
// 	return t
// }

func (l *Lexer) digit() Token {
	t := Token{
		Type:  TokenOrderedList,
		Start: l.Pos,
		End:   l.Pos + 1,
	}
	i := l.Pos
	for i > 0 && l.Data[i] == ' ' {
		i--
	}
	if !(i == l.Pos && i == 0 || l.Data[i] != '\r') {
		t.Type = TokenPlainText
		return t
	}
	i = l.Pos
	for ; i < len(l.Data); i++ {
		c := l.Data[i]
		if '0' <= c && c <= '9' {
			continue
		}
		switch c {
		// TODO(kra53n): check one different cases
		// case ' ', '\r', 0:
		// 	t.Type = TokenPlainText
		// 	t.Start = l.Pos
		// 	t.End = i
		// 	l.Pos = i - 1
		// 	return t
		case '.', ')':
			t.Start = l.Pos
			t.End = i + 1
			l.Pos = i
			return t
		default:
			t.Type = TokenPlainText
			t.Start = l.Pos
			t.End = i
			l.Pos = i - 1
			return t
		}
	}
	t.Type = TokenPlainText
	return t
}

func (l *Lexer) plainText() Token {
	i := l.Pos
	for ; i < len(l.Data); i++ {
		switch l.Data[i] {
		case '\r', ' ', '*', '`', '_', '~':
			goto End
		}
	}
End:
	t := Token{
		Type:  TokenPlainText,
		Start: l.Pos,
		End:   i,
	}
	l.Pos = i - 1
	return t
}

func shouldSkipDueNewLRepetitions(tokens []Token, cur *Token) bool {
	if len(tokens) < 2 {
		return false
	}
	var prv1, prv2 TokenType
	prv1 = tokens[len(tokens)-1].Type
	prv2 = tokens[len(tokens)-2].Type
	return prv1 == prv2 && prv1 == cur.Type && cur.Type == TokenNewL
}

func hasExcessSapce(tokens []Token, cur *Token) bool {
	if len(tokens) < 1 {
		return false
	}
	return tokens[len(tokens)-1].Type == TokenSpace && cur.Type == TokenNewL
}

func analyze(d []byte, tokens []Token) []Token {
	tokens = delExtraNewLinesAtTheEnd(tokens)
	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Type {
		// TODO(kra53n): look at the spaces in TokenUnorderedList, TokenOrderedList
		case TokenSpace:
			tokens = analyzeSpace(tokens, i)
		case TokenAsterisk, TokenDash, TokenPlus:
			tokens = analyzeOnUnorderedList(tokens, i)
		case TokenUnderscore:
			tokens = analyzeAsteriskAndUnderscore(tokens, i)
		case TokenBacktick:
			tokens = analyzeBacktick(tokens, i)
		}
	}
	return tokens
}

func delExtraNewLinesAtTheEnd(tokens []Token) []Token {
	i := len(tokens) - 1
	for i > 0 && tokens[i].Type == TokenNewL {
		i--
	}
	return tokens[:i+1]
}

func shiftTokens(tokens []Token, beg int, end int) []Token {
	i := beg
	for i+end-beg < len(tokens) {
		tokens[i] = tokens[i+end-beg]
		i++
	}
	return tokens[:len(tokens)-(end-beg)]
}

func analyzeSpace(tokens []Token, pos int) []Token {
	cur := tokens[pos]
	if (cur.End-cur.Start >= 4) && (pos == 0 || tokens[pos-1].Type == TokenNewL) {
		i := pos + 1
		for i < len(tokens) && tokens[i].Type != TokenNewL {
			i++
		}
		tokens[pos] = Token{
			Type:  TokenCodeBlock,
			Start: cur.Start,
			End:   tokens[i].End,
		}
		tokens = shiftTokens(tokens, pos+1, i+1)
	}
	return tokens
}

func analyzeOnUnorderedList(tokens []Token, pos int) []Token {
	t := &tokens[pos]
	pT := prvToken(tokens, pos)
	ppT := prvToken(tokens, pos-1)
	// ppT vals: TokenNil, TokenNewL, SomeOtherToken
	//  pT vals: TokenNil, TokenSpace, TokenNewL, SomeOtherToken
	if pT.Type == TokenNewL || pT.Type == TokenSpace && ppT.Type == TokenNewL || pT.Type == TokenNil {
		switch t.Type {
		case TokenAsterisk:
			t.Type = TokenUnorderedList1
		case TokenDash:
			t.Type = TokenUnorderedList2
		case TokenPlus:
			t.Type = TokenUnorderedList3
		}
	}
	return tokens
}

func prvToken(tokens []Token, pos int) Token {
	if pos <= 0 {
		return Token{Type: TokenNil}
	}
	return tokens[pos-1]
}

func analyzeAsteriskAndUnderscore(tokens []Token, pos int) []Token {
	var t TokenType = tokens[pos].Type
	var i, countBeg, countEnd, count int

	i = pos + 1
	if i >= len(tokens) {
		return tokens
	}

	i, countBeg = matchBoldOrItalicBeg(tokens, i, t)
	if countBeg == 0 {
		return tokens
	}
Loop:
	for ; i < len(tokens); i++ {
		switch tokens[i].Type {
		case t:
			break Loop
		case TokenTableCol:
			return tokens
		}
	}
	i, countEnd = matchBoldOrItalicEnd(tokens, i, t)
	if countEnd == 0 {
		return tokens
	}

	count = min(countBeg, countEnd)
	tokens = replaceWithEmphasis(tokens, pos, i-countEnd, count)

	return tokens
}

func matchBoldOrItalicBeg(tokens []Token, pos int, t TokenType) (int, int) {
	i := pos
	count := 1
	for ; i < len(tokens); i++ {
		if tokens[i].Type == t {
			count++
			continue
		} else if tokens[i].Type == TokenAsterisk || tokens[i].Type == TokenUnderscore {
			break
		} else if tokens[i].Type == TokenSpace || tokens[i].Type == TokenNewL {
			return 0, 0
		} else {
			break
		}
	}
	if count > 2 {
		count = 2
	}
	return i, count
}

func matchBoldOrItalicEnd(tokens []Token, pos int, t TokenType) (int, int) {
	i := pos
	count := 0
	for ; i < len(tokens); i++ {
		if tokens[i].Type == t {
			count++
			if count >= 2 {
				count = 2
				return i + 1, count
			}
		} else if count == 1 {
			// NOTE(kra53n): may be have some problems due having here
			// the TokenNewL
			return i, count
		} else if tokens[i].Type == TokenAsterisk || tokens[i].Type == TokenUnderscore {
			break
		} else if tokens[i].Type == TokenNewL {
			if i+1 >= len(tokens) || tokens[i+1].Type == TokenNewL {
				return 0, 0
			}
			continue
		} else {
			continue
		}
	}
	return i, count
}

func replaceWithEmphasis(tokens []Token, beg int, end int, length int) []Token {
	if end >= len(tokens) {
		return tokens
	}
	var tStart, tEnd TokenType
	tStart = TokenItalicStart
	tEnd = TokenItalicEnd
	if length == 2 {
		tStart = TokenBoldStart
		tEnd = TokenBoldEnd
	}
	tokens[beg] = Token{
		Type:  tStart,
		Start: tokens[beg].Start,
		End:   tokens[beg].Start + length,
	}
	tokens[end] = Token{
		Type:  tEnd,
		Start: tokens[end].Start,
		End:   tokens[end].Start + length,
	}
	if length == 1 {
		return tokens
	}
	tokens = shiftTokens(tokens, end+1, end+2)
	tokens = shiftTokens(tokens, beg+1, beg+2)
	return tokens
}

func analyzeBacktick(tokens []Token, pos int) []Token {
	if tokens[pos+1].Type != TokenBacktick && tokens[pos+1].Type != TokenNewL {
		return matchBacktickAsCodeLine(tokens, pos)
	}
	if pos+2 < len(tokens) && tokens[pos+1].Type == TokenBacktick && tokens[pos+2].Type == TokenBacktick {
		return matchBacktickAsCodeBlock(tokens, pos)
	}
	return tokens
}

func matchBacktickAsCodeLine(tokens []Token, pos int) []Token {
	i := pos + 1
	for ; i < len(tokens); i++ {
		switch tokens[i].Type {
		case TokenBacktick:
			tokens[pos] = Token{
				Type:  TokenCodeLine,
				Start: tokens[pos].Start + 1,
				End:   tokens[i].Start,
			}
			return shiftTokens(tokens, pos+1, i+1)
		case TokenNewL:
			return tokens
		default:
			continue
		}
	}
	return tokens
}

func matchBacktickAsCodeBlock(tokens []Token, pos int) []Token {
	var newL Token = tokens[pos]
	newL.Start += 1
	var i int

	for ; i < len(tokens); i++ {
		if tokens[i].Type == TokenNewL {
			newL = tokens[i]
			break
		}
	}

	i = pos + 3
	for ; i+2 < len(tokens); i++ {
		if tokens[i].Type == TokenBacktick && tokens[i].Type == tokens[i+1].Type && tokens[i].Type == tokens[i+2].Type {
			tokens[pos] = Token{
				Type:  TokenCodeBlock,
				Start: newL.Start + 2,
				End:   tokens[i].Start,
			}
			return shiftTokens(tokens, pos+1, i+1)
		}
	}
	return tokens
}

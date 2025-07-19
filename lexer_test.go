package main

import (
	"fmt"
	"testing"
)

var buf string

var tokenNames map[TokenType]string = map[TokenType]string{
	TokenNil:                "Nil",
	TokenH1:                 "H1",
	TokenH2:                 "H2",
	TokenH3:                 "H3",
	TokenH4:                 "H4",
	TokenH5:                 "H5",
	TokenH6:                 "H6",
	TokenNewL:               "NewL",
	TokenSpace:              "Space",
	TokenAsterisk:           "Asterisk",
	TokenBacktick:           "Backtick",
	TokenDash:               "Dash",
	TokenPlus:               "Plus",
	TokenQuote:              "Quote",
	TokenUnderscore:         "Underscore",
	TokenTilde:              "Tilde",
	TokenPlainText:          "PlainText",
	TokenLink:               "Link",
	TokenImg:                "Img",
	TokenUnorderedList:      "UnorderedList",
	TokenUnorderedListType1: "UnorderedListType1",
	TokenUnorderedListType2: "UnorderedListType2",
	TokenUnorderedListType3: "UnorderedListType3",
	TokenOrderedList:        "OrderedList",
	TokenOrderedListType1:   "OrderedListType1",
	TokenOrderedListType2:   "OrderedListType2",
	TokenTableStart:         "TableStart",
	TokenTableHeaderStart:   "TableHeaderStart",
	TokenTableHeaderEnd:     "TableHeaderEnd",
	TokenTableBodyStart:     "TableBodyStart",
	TokenTableBodyEnd:       "TableBodyEnd",
	TokenTableLeftAlign:     "TableLeftAlign",
	TokenTableCenterAlign:   "TableCenterAlign",
	TokenTableRightAlign:    "TableRightAlign",
	TokenTableRow:           "TableRow",
	TokenTableCol:           "TableCol",
	TokenTableEnd:           "TableEnd",
	TokenCodeLine:           "CodeLine",
	TokenCodeBlock:          "CodeBlock",
	TokenBoldStart:          "BoldStart",
	TokenBoldEnd:            "BoldEnd",
	TokenItalicStart:        "ItalicStart",
	TokenItalicEnd:          "ItalicEnd",
	TokenStrikeThrough:      "StrikeThrough",
}

type lexerTestSuite struct {
	md     string
	expect []Token
}

var lexerTestHeaders []lexerTestSuite = []lexerTestSuite{
	{"# header 1", []Token{Token{Type: TokenH1}, Token{Type: TokenPlainText}, Token{Type: TokenSpace}, Token{Type: TokenPlainText}}},
	{"## header 2", []Token{Token{Type: TokenH2}, Token{Type: TokenPlainText}, Token{Type: TokenSpace}, Token{Type: TokenPlainText}}},
	{"### header 3", []Token{Token{Type: TokenH3}, Token{Type: TokenPlainText}, Token{Type: TokenSpace}, Token{Type: TokenPlainText}}},
	{"#### header 4", []Token{Token{Type: TokenH4}, Token{Type: TokenPlainText}, Token{Type: TokenSpace}, Token{Type: TokenPlainText}}},
	{"##### header 5", []Token{Token{Type: TokenH5}, Token{Type: TokenPlainText}, Token{Type: TokenSpace}, Token{Type: TokenPlainText}}},
	{"###### header 6", []Token{Token{Type: TokenH6}, Token{Type: TokenPlainText}, Token{Type: TokenSpace}, Token{Type: TokenPlainText}}},
}

func (t Token) String() string {
	// start := strconv.Itoa(t.Start)
	// end := strconv.Itoa(t.End)
	// return tokenNames[t.Type] + "(" + buf[t.Start:t.End] + "){" + start + ", " + end + "}"
	return tokenNames[t.Type]
}

func tksEqualsByTypes(a, b []Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i].Type != b[i].Type {
			return false
		}
	}
	return true
}

func TestLexer(t *testing.T) {
	runLexerTestSuite(t, lexerTestHeaders)
}

func runLexerTestSuite(t *testing.T, testSuite []lexerTestSuite) {
	for i, testCase := range testSuite {
		t.Run(fmt.Sprintf("%d", i), func(st *testing.T) {
			buf = testCase.md
			tks := Lex([]rune(buf))
			if !tksEqualsByTypes(tks, testCase.expect) {
				st.Errorf("\n     got: %s\nexpected: %s\n", tks, testCase.expect)
			}
		})
	}
}

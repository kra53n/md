package main

type Token struct {
	Type  TokenType
	Start int
	End   int
}

type TokenType int

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
	TokenTab
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
	TokenUnorderedListType1
	TokenUnorderedListType2
	TokenUnorderedListType3
	TokenOrderedList
	TokenOrderedListType1
	TokenOrderedListType2
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

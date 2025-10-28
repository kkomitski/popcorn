package tokens

type TokenType int

const (
	// Literal types
	Number TokenType = iota
	Identifier

	// Keywords
	Let
	Const
	Fn
	Pop

	// Grouping & Operators
	Equals
	Comma
	Dot
	Colon
	Semicolon
	OpenParen    // (
	CloseParen   // )
	OpenBrace    // {
	CloseBrace   // }
	OpenBracket  // [
	CloseBracket // ]
	BinaryOperator

	// End of File
	EOF
)

type Token struct {
	Value     string
	TokenType TokenType
}

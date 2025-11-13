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
	NewLine      // \n
	BinaryOperator

	// Strings
	Quotes // "

	// Comparison operators
	Equal        // ==
	NotEqual     // !=
	Less         // <
	Greater      // >
	LessEqual    // <=
	GreaterEqual // >=

	// Logical operators
	And // &&
	Or  // ||
	Not // !

	// Booleans/null
	Null
	True
	False

	// End of File
	EOF
)

type Token struct {
	Value     string
	TokenType TokenType
}

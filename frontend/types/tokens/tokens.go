package tokens

import "fmt"

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
    OpenParen      // (
    CloseParen     // )
    OpenBrace      // {
    CloseBrace     // }
    OpenBracket    // [
    CloseBracket   // ]
    NewLine        // \n
    BinaryOperator
    UnaryOperator // !

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

func (t TokenType) String() string {
	switch t {
	case Number:
		return "Number"
	case Identifier:
		return "Identifier"
	case Let:
		return "Let"
	case Const:
		return "Const"
	case Fn:
		return "Fn"
	case Pop:
		return "Pop"
	case Equals:
		return "Equals"
	case Comma:
		return "Comma"
	case Dot:
		return "Dot"
	case Colon:
		return "Colon"
	case Semicolon:
		return "Semicolon"
	case OpenParen:
		return "OpenParen"
	case CloseParen:
		return "CloseParen"
	case OpenBrace:
		return "OpenBrace"
	case CloseBrace:
		return "CloseBrace"
	case OpenBracket:
		return "OpenBracket"
	case CloseBracket:
		return "CloseBracket"
	case NewLine:
		return "NewLine"
	case BinaryOperator:
		return "BinaryOperator"
	case Quotes:
		return "Quotes"
	case Equal:
		return "Equal"
	case NotEqual:
		return "NotEqual"
	case Less:
		return "Less"
	case Greater:
		return "Greater"
	case LessEqual:
		return "LessEqual"
	case GreaterEqual:
		return "GreaterEqual"
	case And:
		return "And"
	case Or:
		return "Or"
	case Not:
		return "Not"
	case Null:
		return "Null"
	case True:
		return "True"
	case False:
		return "False"
	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("UnknownTokenType(%d)", int(t))
	}
}

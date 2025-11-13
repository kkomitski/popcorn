package test_frontend

import (
	"os"
	FE "pop/frontend"
	"pop/frontend/types/tokens"
	"testing"
)

const TOKENS_FILE = "../mocks/all-tokens.pop"

func TestLexer(t *testing.T) {
	content, err := os.ReadFile(TOKENS_FILE)
	if err != nil {
		t.Fatalf("Failed to read file %v", err)
	}

	tokensOut := FE.Tokenize(string(content))

	expected := []tokens.Token{
		// "// Literals" comment
		{Value: "\n", TokenType: tokens.NewLine},

		// Literals
		{Value: "123", TokenType: tokens.Number},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "identifier", TokenType: tokens.Identifier},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "\"", TokenType: tokens.Quotes},
		{Value: "string", TokenType: tokens.Identifier},
		{Value: "\"", TokenType: tokens.Quotes},
		{Value: "\n", TokenType: tokens.NewLine},

		// Blank line
		{Value: "\n", TokenType: tokens.NewLine},

		// "// Keywords" comment
		{Value: "\n", TokenType: tokens.NewLine},

		// Keywords
		{Value: "let", TokenType: tokens.Let},
		{Value: "x", TokenType: tokens.Identifier},
		{Value: "=", TokenType: tokens.Equals},
		{Value: "1", TokenType: tokens.Number},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "const", TokenType: tokens.Const},
		{Value: "y", TokenType: tokens.Identifier},
		{Value: "=", TokenType: tokens.Equals},
		{Value: "2", TokenType: tokens.Number},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "fn", TokenType: tokens.Fn},
		{Value: "add", TokenType: tokens.Identifier},
		{Value: "(", TokenType: tokens.OpenParen},
		{Value: "a", TokenType: tokens.Identifier},
		{Value: ",", TokenType: tokens.Comma},
		{Value: "b", TokenType: tokens.Identifier},
		{Value: ")", TokenType: tokens.CloseParen},
		{Value: "{", TokenType: tokens.OpenBrace},
		{Value: "pop", TokenType: tokens.Pop},
		{Value: "a", TokenType: tokens.Identifier},
		{Value: "+", TokenType: tokens.BinaryOperator},
		{Value: "b", TokenType: tokens.Identifier},
		{Value: "}", TokenType: tokens.CloseBrace},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "pop", TokenType: tokens.Pop},
		{Value: "x", TokenType: tokens.Identifier},
		{Value: "\n", TokenType: tokens.NewLine},

		// Blank line
		{Value: "\n", TokenType: tokens.NewLine},

		// "// Grouping & Operators" comment
		{Value: "\n", TokenType: tokens.NewLine},

		// Grouping & Operators
		{Value: "x", TokenType: tokens.Identifier},
		{Value: "=", TokenType: tokens.Equals},
		{Value: "(", TokenType: tokens.OpenParen},
		{Value: "1", TokenType: tokens.Number},
		{Value: "+", TokenType: tokens.BinaryOperator},
		{Value: "2", TokenType: tokens.Number},
		{Value: ")", TokenType: tokens.CloseParen},
		{Value: "*", TokenType: tokens.BinaryOperator},
		{Value: "[", TokenType: tokens.OpenBracket},
		{Value: "3", TokenType: tokens.Number},
		{Value: ",", TokenType: tokens.Comma},
		{Value: "4", TokenType: tokens.Number},
		{Value: "]", TokenType: tokens.CloseBracket},
		{Value: "/", TokenType: tokens.BinaryOperator},
		{Value: "{", TokenType: tokens.OpenBrace},
		{Value: "a", TokenType: tokens.Identifier},
		{Value: ":", TokenType: tokens.Colon},
		{Value: "5", TokenType: tokens.Number},
		{Value: "}", TokenType: tokens.CloseBrace},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "x", TokenType: tokens.Identifier},
		{Value: ",", TokenType: tokens.Comma},
		{Value: "y", TokenType: tokens.Identifier},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "obj", TokenType: tokens.Identifier},
		{Value: ".", TokenType: tokens.Dot},
		{Value: "prop", TokenType: tokens.Identifier},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "arr", TokenType: tokens.Identifier},
		{Value: "[", TokenType: tokens.OpenBracket},
		{Value: "0", TokenType: tokens.Number},
		{Value: "]", TokenType: tokens.CloseBracket},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "{", TokenType: tokens.OpenBrace},
		{Value: "}", TokenType: tokens.CloseBrace},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "[", TokenType: tokens.OpenBracket},
		{Value: "]", TokenType: tokens.CloseBracket},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "(", TokenType: tokens.OpenParen},
		{Value: ")", TokenType: tokens.CloseParen},
		{Value: "\n", TokenType: tokens.NewLine},

		// Blank line
		{Value: "\n", TokenType: tokens.NewLine},

		// "// Comparison operators" comment
		{Value: "\n", TokenType: tokens.NewLine},

		// Comparison operators
		{Value: "a", TokenType: tokens.Identifier},
		{Value: "==", TokenType: tokens.Equal},
		{Value: "b", TokenType: tokens.Identifier},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "a", TokenType: tokens.Identifier},
		{Value: "!=", TokenType: tokens.NotEqual},
		{Value: "b", TokenType: tokens.Identifier},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "a", TokenType: tokens.Identifier},
		{Value: "<", TokenType: tokens.Less},
		{Value: "b", TokenType: tokens.Identifier},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "a", TokenType: tokens.Identifier},
		{Value: ">", TokenType: tokens.Greater},
		{Value: "b", TokenType: tokens.Identifier},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "a", TokenType: tokens.Identifier},
		{Value: "<=", TokenType: tokens.LessEqual},
		{Value: "b", TokenType: tokens.Identifier},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "a", TokenType: tokens.Identifier},
		{Value: ">=", TokenType: tokens.GreaterEqual},
		{Value: "b", TokenType: tokens.Identifier},
		{Value: "\n", TokenType: tokens.NewLine},

		// Blank line
		{Value: "\n", TokenType: tokens.NewLine},

		// "// Logical operators" comment
		{Value: "\n", TokenType: tokens.NewLine},

		// Logical operators
		{Value: "&&", TokenType: tokens.And},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "||", TokenType: tokens.Or},
		{Value: "\n", TokenType: tokens.NewLine},

		// Blank line
		{Value: "\n", TokenType: tokens.NewLine},

		// "// Boolean and null" comment
		{Value: "\n", TokenType: tokens.NewLine},

		// Boolean and null
		{Value: "true", TokenType: tokens.True},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "false", TokenType: tokens.False},
		{Value: "\n", TokenType: tokens.NewLine},
		{Value: "null", TokenType: tokens.Null},
		// Removed the trailing newline - file doesn't end with one

		// End of file
		{Value: "EndOfFile", TokenType: tokens.EOF},
	}

	if len(tokensOut) != len(expected) {
		t.Fatalf("Token count mismatch: got %d, want %d", len(tokensOut), len(expected))
	}

	for i, exp := range expected {
		got := tokensOut[i]
		if got.TokenType != exp.TokenType || got.Value != exp.Value {
			t.Errorf("Token %d: got {%v, %v}, want {%v, %v}", i, got.Value, got.TokenType, exp.Value, exp.TokenType)
		}
	}
}

package frontend

import (
	"log"
	"pop/frontend/types/tokens"
	utils "pop/lib"
)

func Tokenize(sourceCode string) []tokens.Token {
	chars := []rune(sourceCode)
	tokensList := make([]tokens.Token, 0, len(chars))

	singleCharTokens := map[rune]tokens.TokenType{
		'+': tokens.BinaryOperator,
		'-': tokens.BinaryOperator,
		'/': tokens.BinaryOperator,
		'*': tokens.BinaryOperator,
		'%': tokens.BinaryOperator,
		'(': tokens.OpenParen,
		')': tokens.CloseParen,
		'{': tokens.OpenBrace,
		'}': tokens.CloseBrace,
		'[': tokens.OpenBracket,
		']': tokens.CloseBracket,
		'=': tokens.Equals,
		';': tokens.Semicolon,
		':': tokens.Colon,
		',': tokens.Comma,
		'.': tokens.Dot,
		'"': tokens.Quotes,
		'<': tokens.Less,
		'>': tokens.Greater,
	}

	keywords := map[string]tokens.TokenType{
		"let":   tokens.Let,
		"const": tokens.Const,
		"fn":    tokens.Fn,
		"pop":   tokens.Pop,
	}

	comparers := map[string]tokens.TokenType{
		"==": tokens.Equal,
		"!=": tokens.NotEqual,
		"<=": tokens.LessEqual,
		">=": tokens.GreaterEqual,
	}

	i := 0
	for i < len(chars) {
		c := chars[i]

		if i+1 < len(chars) && utils.IsComparer(string(c)+string(chars[i+1])) {
			comparer := string(c) + string(chars[i+1])

			tokenType, _ := comparers[comparer]
			tokensList = append(tokensList, tokens.Token{Value: string(c) + string(chars[i+1]), TokenType: tokenType})

			i++
			i++
		} else if tokenType, ok := singleCharTokens[c]; ok {
			tokensList = append(tokensList, tokens.Token{Value: string(c), TokenType: tokenType})
			i++
		} else if utils.IsDigit(c) {
			start := i
			for i < len(chars) && utils.IsDigit(chars[i]) {
				i++
			}
			digit := string(chars[start:i])
			tokensList = append(tokensList, tokens.Token{Value: digit, TokenType: tokens.Number})
		} else if utils.IsAlphabetical(c) {
			start := i
			for i < len(chars) && (utils.IsAlphabetical(chars[i]) || utils.IsDigit(chars[i])) {
				i++
			}
			word := string(chars[start:i])
			keyword, ok := keywords[word]
			if ok {
				tokensList = append(tokensList, tokens.Token{Value: word, TokenType: keyword})
			} else {
				tokensList = append(tokensList, tokens.Token{Value: word, TokenType: tokens.Identifier})
			}
		} else if utils.IsSkippable(c) {
			i++
		} else {
			log.Fatalf("Token of type '%s' is not yet processable", string(c))
			log.Fatalf("Failed at: %s", string(chars[i:i+30]))
		}
	}

	tokensList = append(tokensList, tokens.Token{Value: "EndOfFile", TokenType: tokens.EOF})

	return tokensList
}

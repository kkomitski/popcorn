package lib

func IsDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func IsAlphabetical(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func IsSkippable(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IsComparer(s string) bool {
	return s == "==" || s == "!=" || s == "<=" || s == ">="
}

func IsLogical(s string) bool {
	return s == "&&" || s == "||"
}

func IsComment(s string ) bool {
	return s == "//"
}

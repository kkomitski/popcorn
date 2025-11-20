package backend

import (
	"bufio"
	"fmt"
	"os"
	FE "pop/frontend"
	T "pop/frontend/types/tokens"
	"strings"
)

func RunFile(filePath string) error {
	// Read the file contents
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Tokenize and parse
	tokens := FE.Tokenize(string(content))
	ast := FE.ProduceAST(tokens)

	// Create environment and evaluate
	env := MakeEnvironment()
	result := Evaluate(ast, env)

	// Print the final result
	fmt.Printf("%+v\n", result)

	return nil
}

// ANSI color codes
const (
	colorReset      = "\033[0m"
	colorKeyword    = "\033[1;35m" // Magenta
	colorNumber     = "\033[1;33m" // Yellow
	colorOperator   = "\033[1;36m" // Cyan
	colorString     = "\033[1;32m" // Green
	colorIdentifier = "\033[0;37m" // White
)

// highlightSyntax applies syntax highlighting to the input line
func highlightSyntax(line string) string {
	if strings.TrimSpace(line) == "" {
		return line
	}

	tokensList := FE.Tokenize(line)
	var highlighted strings.Builder
	lastEnd := 0

	for _, token := range tokensList {
		if token.TokenType == T.EOF {
			break
		}

		// Find the token in the original line
		pos := strings.Index(line[lastEnd:], token.Value)
		if pos >= 0 {
			// Add any whitespace before the token
			highlighted.WriteString(line[lastEnd : lastEnd+pos])

			// Add the colorized token
			switch token.TokenType {
			case T.Let, T.Const, T.Fn:
				highlighted.WriteString(colorKeyword + token.Value + colorReset)
			case T.Number:
				highlighted.WriteString(colorNumber + token.Value + colorReset)
			case T.BinaryOperator, T.Equals:
				highlighted.WriteString(colorOperator + token.Value + colorReset)
			case T.Identifier:
				highlighted.WriteString(colorIdentifier + token.Value + colorReset)
			default:
				highlighted.WriteString(token.Value)
			}

			lastEnd += pos + len(token.Value)
		}
	}

	// Add any remaining characters
	if lastEnd < len(line) {
		highlighted.WriteString(line[lastEnd:])
	}

	return highlighted.String()
}

func printHeader() {
    fmt.Println("\033[1;35m╔══════════════════════════════════════════════════════════╗\033[0m")
    fmt.Println("\033[1;35m║\033[0m          🍿 \033[1;36mPopcorn Language REPL\033[0m 🍿           \033[1;35m║\033[0m")
    fmt.Println("\033[1;35m╠══════════════════════════════════════════════════════════╣\033[0m")
    fmt.Println("\033[1;35m║\033[0m  \033[0;33mType 'exit' to quit\033[0m                           \033[1;35m║\033[0m")
    fmt.Println("\033[1;35m║\033[0m  \033[0;33mType 'clear' to clear screen\033[0m                  \033[1;35m║\033[0m")
    fmt.Println("\033[1;35m║\033[0m  \033[0;33mType 'verbose' for multi-line input (finish with ':send')\033[0m \033[1;35m║\033[0m")
    fmt.Println("\033[1;35m╚══════════════════════════════════════════════════════════╝\033[0m")
    fmt.Println()
}

func Repl() RuntimeVal {
	scanner := bufio.NewScanner(os.Stdin)

	printHeader()

	env := MakeEnvironment()
	verboseMode := false
	var verboseBuffer strings.Builder

	for {
		if !verboseMode {
			// Cyan prompt
			fmt.Print("🍿 \033[1;36m>>\033[0m ")
		} else {
			fmt.Print("🍿 \033[1;33m(v)\033[0m ")
		}

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()

		if !verboseMode && line == "exit" {
			fmt.Println("\033[1;31m👋 Exiting REPL... Enjoy your popcorn!\033[0m")
			break
		}

		if !verboseMode && line == "clear" {
			fmt.Print("\033[2J\033[H")
			printHeader()
			continue
		}

		if !verboseMode && line == "verbose" {
			fmt.Println("\033[1;33m[Verbose mode: Enter multiple lines. Type ':send' on a new line to execute.]\033[0m")
			verboseMode = true
			verboseBuffer.Reset()
			continue
		}

		if verboseMode {
			if line == ":send" {
				// Process the collected input
				code := verboseBuffer.String()
				if strings.TrimSpace(code) != "" {
					highlighted := highlightSyntax(code)
					fmt.Printf("   \033[2m→\033[0m %s\n", highlighted)

					tokensSlice := FE.Tokenize(code)
					ast := FE.ProduceAST(tokensSlice)
					res := Evaluate(ast, env)
					fmt.Printf("   \033[1;32m←\033[0m \033[1;32m%+v\033[0m\n\n", res)
				}
				verboseMode = false
				continue
			}
			verboseBuffer.WriteString(line + "\n")
			continue
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		// Add a newline so the REPL doesn't throw "VariableDeclarations must end with a newline error"
		if !strings.HasSuffix(line, "\n") {
			line += "\n"
		}

		// Show syntax highlighted version
		highlighted := highlightSyntax(line)
		fmt.Printf("   \033[2m→\033[0m %s\n", highlighted)

		tokensSlice := FE.Tokenize(line)
		ast := FE.ProduceAST(tokensSlice)

		res := Evaluate(ast, env)

		// Green output with arrow
		fmt.Printf("   \033[1;32m←\033[0m \033[1;32m%+v\033[0m\n\n", res)
	}

	return Null
}

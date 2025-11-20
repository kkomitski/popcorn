package frontend

import (
	"fmt"
	"log"
	"os"
	"pop/frontend/types/ast"
	"pop/frontend/types/tokens"
	"strconv"
)

type Parser struct {
	Tokens []tokens.Token
	Pos    int
}

// * ========= UTILS ========= * \\

func (p *Parser) notEOF() bool {
	return p.Pos < len(p.Tokens) && p.Tokens[p.Pos].TokenType != tokens.EOF
}

func (p *Parser) eat() tokens.Token {
	curr := p.Tokens[p.Pos]
	p.Pos++

	return curr
}

func (p *Parser) at() tokens.Token {
	if p.Pos < len(p.Tokens) {
		return p.Tokens[p.Pos]
	}
	return tokens.Token{TokenType: tokens.EOF}
}

func (p *Parser) expect(tokenType tokens.TokenType, err string) tokens.Token {
	prev := p.eat()
	if prev.TokenType != tokenType {
		log.Fatalf("Parser error: %s\nExpected: '%v', but got: '%v'.", err, tokenType.String(), prev.TokenType.String())
	}
	return prev
}

func (p *Parser) skipNewlines() {
	for p.at().TokenType == tokens.NewLine {
		p.eat()
	}
}

// * ======== STATEMENTS ======== * \\

func (p *Parser) parseStatement() ast.ASTNode {
	switch p.at().TokenType {
	case tokens.Let, tokens.Const:
		return p.parseVarDeclaration()
	case tokens.Fn:
		return p.parseFnDeclaration()
	case tokens.Pop:
		return p.parseFnReturn()
	default:
		node := p.parseExpr()

		if p.at().TokenType == tokens.NewLine {
			p.eat()
		} else if p.at().TokenType != tokens.EOF {
			log.Fatalf("Expected newline or EOF after statement, got: %v", p.at())
		}

		return node
	}
}

func (p *Parser) parseFnReturn() ast.ASTNode {
	p.eat() // Eat the `pop` keyword

	// Pop statements need to complete with a semicolon or an expression

	// Is the next token an expression
	isNextExpr := p.at().TokenType == tokens.Identifier ||
		p.at().TokenType == tokens.Number ||
		p.at().TokenType == tokens.OpenParen

	// Is the next token an expression?
	if isNextExpr {
		val := p.parseExpr()
		// Optionally: check for newline, block end, or EOF here
		return ast.ReturnStatementNode{Value: val}
	} else if p.at().TokenType == tokens.CloseBrace ||
		p.at().TokenType == tokens.EOF {
		// No value, just pop (return)
		return ast.ReturnStatementNode{Value: nil}
	} else {
		log.Fatalf("Expected an expression or end of statement after 'pop', got: %v", p.at())
		return nil
	}
}

func (p *Parser) parseVarDeclaration() ast.ASTNode {
	isConstant := p.eat().TokenType == tokens.Const

	identifier := p.expect(tokens.Identifier, "Expected identifier name following 'let' | 'const' keywords").Value

	if p.at().TokenType == tokens.NewLine {
		p.eat()
		if isConstant {
			log.Fatalf("Must assign value to constant expression. No value provided.")
		}
		return ast.VariableDeclarationNode{
			Identifier: identifier,
			Constant:   isConstant,
		}
	}

	p.expect(tokens.Equals, "Expected equals token following identifier in variable declaration.")

	declaration := ast.VariableDeclarationNode{
		Constant:   isConstant,
		Identifier: identifier,
		Value:      p.parseExpr(),
	}

	p.expect(tokens.NewLine, "Variable declaration statement must end with a new line")

	return declaration
}

func (p *Parser) parseFnDeclaration() ast.ASTNode {
	p.eat() // Eat the 'fn' keyword

	name := p.expect(tokens.Identifier, "Expected a function name following the 'fn' keyword.").Value

	args := p.parseArgs()
	params := []string{}

	for _, arg := range args {
		identifier, ok := arg.(ast.IdentifierExprNode)
		if !ok {
			log.Fatalf("Inside function declaration expected parameters to be of type 'Identifier'. Got: %v", arg)
		}
		params = append(params, identifier.Symbol)
	}

	p.expect(tokens.OpenBrace, "Expected fn body following a declaration")

	body := []ast.ASTNode{}

	for p.notEOF() {
		p.skipNewlines()
		if p.at().TokenType == tokens.CloseBrace {
			break
		}
		body = append(body, p.parseStatement())
	}

	p.expect(tokens.CloseBrace, "Closing bracket expected inside function declaration")

	// Consume trailing newLine
	if p.at().TokenType == tokens.NewLine {
		p.eat()
	}

	return ast.FunctionDeclarationNode{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

// * ======= EXPRESSIONS ======= * \\

func (p *Parser) parseExpr() ast.ASTNode {
	return p.parseAssignmentExpr()
}

func (p *Parser) parseAssignmentExpr() ast.ASTNode {
	left := p.parseLogicalExpr()

	if p.at().TokenType == tokens.Equals {
		p.eat() // Advance past equals
		value := p.parseAssignmentExpr()
		return ast.AssignmentExprNode{
			Value:    value,
			Assignee: left,
		}
	}

	return left
}

func (p *Parser) parseLogicalExpr() ast.ASTNode {
	left := p.parseComparisonExpr()

	for {
		tk := p.at().TokenType
		if tk == tokens.And || tk == tokens.Or {
			operator := p.eat().Value
			right := p.parseComparisonExpr()

			left = ast.LogicalExprNode{
				Left:     left,
				Right:    right,
				Operator: ast.BinaryOperatorKind(operator),
			}
		} else {
			break
		}
	}

	return left
}

func (p *Parser) parseComparisonExpr() ast.ASTNode {
	left := p.parseObjectExpr()

	for {
		op := p.at().Value
		if op == "==" || op == "!=" || op == "<" || op == ">" || op == "<=" || op == ">=" {
			operator := p.eat().Value
			right := p.parseObjectExpr()
			left = ast.BinaryExprNode{
				Left:     left,
				Right:    right,
				Operator: ast.BinaryOperatorKind(operator),
			}
		} else {
			break
		}
	}

	return left
}

func (p *Parser) parseObjectExpr() ast.ASTNode {
	if p.at().TokenType != tokens.OpenBrace {
		return p.parseAdditiveExpr()
	}

	p.eat() // advance past the open brace

	properties := []ast.PropertyNode{}

	// eat the new lines before object member assignments
	for p.notEOF() && p.at().TokenType != tokens.CloseBrace {
		if p.at().TokenType == tokens.NewLine {
			p.eat()
		}

		key := p.expect(tokens.Identifier, "Object literal key expected!").Value

		// Shorthand property: { key }
		if p.at().TokenType == tokens.Comma {
			p.eat()
			properties = append(properties, ast.PropertyNode{
				Key:   key,
				Value: nil,
			})
			continue
		} else if p.at().TokenType == tokens.CloseBrace {
			properties = append(properties, ast.PropertyNode{
				Key:   key,
				Value: nil,
			})
			continue
		}

		// Full property: { key: value }
		p.expect(tokens.Colon, "Missing colon following identifier in ObjectExpression")
		value := p.parseExpr()

		properties = append(properties, ast.PropertyNode{
			Key:   key,
			Value: value,
		})

		// eat the new lines after object member assignments
		if p.at().TokenType == tokens.NewLine {
			p.eat()
		}

		if p.at().TokenType != tokens.CloseBrace {
			p.expect(tokens.Comma, "Expected comma or closing bracket following property")
		}
	}

	p.expect(tokens.CloseBrace, "Object literal missing closing brace.")

	return ast.ObjectLiteralExprNode{
		Properties: properties,
	}
}

func (p *Parser) parseAdditiveExpr() ast.ASTNode {
	left := p.parseMultiplicativeExpr()

	for p.at().Value == "+" || p.at().Value == "-" {
		operator := p.eat().Value
		right := p.parseMultiplicativeExpr()
		left = ast.BinaryExprNode{
			Left:     left,
			Right:    right,
			Operator: ast.BinaryOperatorKind(operator),
		}
	}

	return left
}

func (p *Parser) parseMultiplicativeExpr() ast.ASTNode {
	left := p.parseUnaryExpr()

	for p.at().Value == "/" || p.at().Value == "*" || p.at().Value == "%" {
		operator := p.eat().Value
		right := p.parseUnaryExpr()
		left = ast.BinaryExprNode{
			Left:     left,
			Right:    right,
			Operator: ast.BinaryOperatorKind(operator),
		}
	}

	return left
}

func (p *Parser) parseUnaryExpr() ast.ASTNode {
	tk := p.at()
	// Check for unary minus or logical not
	if (tk.TokenType == tokens.BinaryOperator && (tk.Value == "-" || tk.Value == "+")) ||
		(tk.TokenType == tokens.UnaryOperator && tk.Value == "!") {
		operator := p.eat().Value
		operand := p.parseUnaryExpr()
		if operator == "+" {
			// Unary plus: just return the operand as-is
			return operand
		}
		return ast.UnaryExprNode{
			Operator: ast.UnaryOperatorKind(operator),
			Operand:  operand,
		}
	}
	return p.parseCallMemberExpr()
}

// * ======= CALL & MEMBER EXPRESSIONS ======= * \\

func (p *Parser) parseCallMemberExpr() ast.ASTNode {
	member := p.parseMemberExpr()

	if p.at().TokenType == tokens.OpenParen {
		return p.parseCallExpr(member)
	}

	return member
}

func (p *Parser) parseCallExpr(caller ast.ASTNode) ast.ASTNode {
	callExpr := ast.CallExprNode{
		Caller: caller,
		Args:   p.parseArgs(),
	}

	if p.at().TokenType == tokens.OpenParen {
		callExpr = p.parseCallExpr(callExpr).(ast.CallExprNode)
	}

	return callExpr
}

func (p *Parser) parseArgs() []ast.ASTNode {
	p.expect(tokens.OpenParen, "Expected open parenthesis")

	var args []ast.ASTNode
	if p.at().TokenType == tokens.CloseParen {
		args = []ast.ASTNode{}
	} else {
		args = p.parseArgumentsList()
	}

	p.expect(tokens.CloseParen, "Missing closing parenthesis")
	return args
}

func (p *Parser) parseArgumentsList() []ast.ASTNode {
	args := []ast.ASTNode{p.parseAssignmentExpr()}

	for p.at().TokenType == tokens.Comma {
		p.eat()
		args = append(args, p.parseAssignmentExpr())
	}

	return args
}

func (p *Parser) parseMemberExpr() ast.ASTNode {
	object := p.parsePrimaryExpr()

	for p.at().TokenType == tokens.Dot || p.at().TokenType == tokens.OpenBracket {
		operator := p.eat()
		var property ast.ASTNode
		var computed bool

		if operator.TokenType == tokens.Dot {
			computed = false
			property = p.parsePrimaryExpr()

			if ast.GetNodeKind(property) != ast.IdentifierExpr {
				log.Fatalf("Cannot use dot operator without right hand side being an identifier")
			}
		} else {
			computed = true
			property = p.parseExpr()
			p.expect(tokens.CloseBracket, "Missing closing bracket in computed value.")
		}

		object = ast.MemberExprNode{
			Object:   object,
			Property: property,
			Computed: computed,
		}
	}

	return object
}

// * ======= PRIMARY EXPRESSIONS ======= * \\

func (p *Parser) parsePrimaryExpr() ast.ASTNode {
	tk := p.at().TokenType

	switch tk {
	case tokens.Identifier:
		return ast.IdentifierExprNode{
			Symbol: p.eat().Value,
		}
	case tokens.Number:
		value, err := strconv.ParseFloat(p.eat().Value, 64)
		if err != nil {
			log.Fatalf("Failed to parse number: %v", err)
		}
		return ast.NumericLiteralExprNode{
			Value: value,
		}
	case tokens.OpenParen:
		p.eat() // Eat the opening paren
		value := p.parseExpr()
		p.expect(tokens.CloseParen, "Unexpected token found inside parenthesised expression. Expected closing parenthesis.")
		return value
	case tokens.OpenBracket:
		p.eat() // Eat the opening bracket

		// Look-ahead to count elements for optimal pre-allocation
		elementCount := 0
		if p.at().TokenType != tokens.CloseBracket {
			tempPos := p.Pos
			depth := 0
			for tempPos < len(p.Tokens) {
				tk := p.Tokens[tempPos]
				if tk.TokenType == tokens.OpenBracket {
					depth++
				} else if tk.TokenType == tokens.CloseBracket {
					if depth == 0 {
						break
					}
					depth--
				} else if tk.TokenType == tokens.Comma && depth == 0 {
					elementCount++
				}
				tempPos++
			}
			elementCount++ // Add 1 for the last element (no trailing comma)
		}

		// Pre-allocate with exact capacity
		elements := make([]ast.ASTNode, 0, elementCount)

		if p.at().TokenType != tokens.CloseBracket {
			for {
				elements = append(elements, p.parseExpr())
				if p.at().TokenType == tokens.CloseBracket {
					break
				}
				// Require a comma between elements
				p.expect(tokens.Comma, "Array elements should be separated with commas.")
				// Check for trailing comma (optional)
				if p.at().TokenType == tokens.CloseBracket {
					break
				}
			}
		}

		p.expect(tokens.CloseBracket, "Expected closing bracket for array literal.")
		return ast.ArrayLiteralExprNode{
			Elements: elements,
			Size:     int64(len(elements)),
		}
	case tokens.OpenBrace:
		return p.parseObjectExpr()
	case tokens.Quotes:
		p.eat() // Eat the opening quote

		val := p.eat().Value

		p.expect(tokens.Quotes, "String literals should end with a closing quote.")
		return ast.StringLiteralExprNode{
			Value: val,
		}
	case tokens.True, tokens.False:
		val := p.eat().Value

		fmt.Println(val)

		if val == "true" {
			return ast.BooleanLiteralExprNode{
				Value: true,
			}
		}

		return ast.BooleanLiteralExprNode{
			Value: false,
		}

	default:
		log.Fatalf("Unexpected token found during parsing: %v", p.at())
		return nil
	}
}

// Helpers
// WrapASTWithKind recursively wraps all ASTNodes with JSONNode for JSON marshaling
func WrapASTWithKind(node ast.ASTNode) ast.JSONNode {
	if node == nil {
		return ast.JSONNode{Kind: "ERR_UNKNOWN", Data: nil}
	}

	switch n := node.(type) {
	case ast.Program:
		body := make([]ast.ASTNode, len(n.Body))
		for i, child := range n.Body {
			body[i] = WrapASTWithKind(child)
		}
		return ast.JSONNode{Data: ast.Program{Body: body}}
	case ast.VariableDeclarationNode:
		return ast.JSONNode{Data: ast.VariableDeclarationNode{
			Constant:   n.Constant,
			Identifier: n.Identifier,
			Value:      WrapASTWithKind(n.Value),
		}}
	case ast.FunctionDeclarationNode:
		body := make([]ast.ASTNode, len(n.Body))
		for i, child := range n.Body {
			body[i] = WrapASTWithKind(child)
		}
		return ast.JSONNode{Data: ast.FunctionDeclarationNode{
			Name:   n.Name,
			Params: n.Params,
			Body:   body,
		}}
	case ast.AssignmentExprNode:
		return ast.JSONNode{Data: ast.AssignmentExprNode{
			Assignee: WrapASTWithKind(n.Assignee),
			Value:    WrapASTWithKind(n.Value),
		}}
	case ast.BinaryExprNode:
		return ast.JSONNode{Data: ast.BinaryExprNode{
			Left:     WrapASTWithKind(n.Left),
			Right:    WrapASTWithKind(n.Right),
			Operator: n.Operator,
		}}
	case ast.LogicalExprNode:
		return ast.JSONNode{Data: ast.LogicalExprNode{
			Left:     WrapASTWithKind(n.Left),
			Right:    WrapASTWithKind(n.Right),
			Operator: n.Operator,
		}}
	case ast.MemberExprNode:
		return ast.JSONNode{Data: ast.MemberExprNode{
			Object:   WrapASTWithKind(n.Object),
			Property: WrapASTWithKind(n.Property),
			Computed: n.Computed,
		}}
	case ast.CallExprNode:
		args := make([]ast.ASTNode, len(n.Args))
		for i, arg := range n.Args {
			args[i] = WrapASTWithKind(arg)
		}
		return ast.JSONNode{Data: ast.CallExprNode{
			Caller: WrapASTWithKind(n.Caller),
			Args:   args,
		}}
	case ast.ArrayLiteralExprNode:
		elements := make([]ast.ASTNode, len(n.Elements))
		for i, elem := range n.Elements {
			elements[i] = WrapASTWithKind(elem)
		}
		return ast.JSONNode{Data: ast.ArrayLiteralExprNode{
			Elements: elements,
			Size:     n.Size,
		}}
	case ast.ObjectLiteralExprNode:
		props := make([]ast.PropertyNode, len(n.Properties))
		for i, prop := range n.Properties {
			props[i] = ast.PropertyNode{
				Key:   prop.Key,
				Value: WrapASTWithKind(prop.Value),
			}
		}
		return ast.JSONNode{Data: ast.ObjectLiteralExprNode{
			Properties: props,
		}}
	case ast.ReturnStatementNode:
		return ast.JSONNode{Data: ast.ReturnStatementNode{
			Value: WrapASTWithKind(n.Value),
		}}
	// Leaf nodes
	case ast.IdentifierExprNode, ast.NumericLiteralExprNode, ast.StringLiteralExprNode,
		ast.BooleanLiteralExprNode, ast.NullLiteralExprNode:
		return ast.JSONNode{Data: n}
	default:
		return ast.JSONNode{Data: n}
	}
}

// * ======= PUBLIC API ======= * \\

func ProduceAST(tokens []tokens.Token, verbose ...bool) ast.Program {
	parser := Parser{
		Tokens: tokens,
		Pos:    0,
	}

	program := ast.Program{
		Body: []ast.ASTNode{},
	}

	for parser.notEOF() {
		parser.skipNewlines()

		program.Body = append(program.Body, parser.parseStatement())
	}

	isVerbose := true
	if len(verbose) > 0 {
		isVerbose = verbose[0]
	}

	if isVerbose {
		wrappedAST := WrapASTWithKind(program)
		jsonBytes, err := wrappedAST.MarshalJSON()
		if err != nil {
			log.Fatalf("Failed to marshal AST to JSON: %v", err)
		}

		os.WriteFile("current_ast.json", jsonBytes, 0777)
	}

	return program
}

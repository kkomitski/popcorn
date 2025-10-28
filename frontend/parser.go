package frontend

import (
	"log"
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
		log.Fatalf("Parser error: %s\nExpected: %v, but got: %v", err, tokenType, prev.TokenType)
	}
	return prev
}

// * ======== STATEMENTS ======== * \\

func (p *Parser) parseStatement() ast.ASTNode {
	switch p.at().TokenType {
	case tokens.Let, tokens.Const:
		return p.parseVarDeclaration()
	case tokens.Fn:
		return p.parseFnDeclaration()
	default:
		expr := p.parseExpr()
		// Eat the extra semicolon at the end of expression statements
		if p.at().TokenType == tokens.Semicolon {
			p.eat()
		}
		return expr
	}
}

func (p *Parser) parseVarDeclaration() ast.ASTNode {
	isConstant := p.eat().TokenType == tokens.Const

	identifier := p.expect(tokens.Identifier, "Expected identifier name following 'let' | 'const' keywords").Value

	if p.at().TokenType == tokens.Semicolon {
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

	p.expect(tokens.Semicolon, "Variable declaration statement must end with semicolon")

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

	for p.notEOF() && p.at().TokenType != tokens.CloseBrace {
		body = append(body, p.parseStatement())
	}

	p.expect(tokens.CloseBrace, "Closing bracket expected inside function declaration")

	// Consume trailing semicolon
	if p.at().TokenType == tokens.Semicolon {
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
	left := p.parseObjectExpr()

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

func (p *Parser) parseObjectExpr() ast.ASTNode {
	if p.at().TokenType != tokens.OpenBrace {
		return p.parseAdditiveExpr()
	}

	p.eat() // advance past the open brace

	properties := []ast.PropertyNode{}

	for p.notEOF() && p.at().TokenType != tokens.CloseBrace {
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
	left := p.parseCallMemberExpr()

	for p.at().Value == "/" || p.at().Value == "*" || p.at().Value == "%" {
		operator := p.eat().Value
		right := p.parseCallMemberExpr()
		left = ast.BinaryExprNode{
			Left:     left,
			Right:    right,
			Operator: ast.BinaryOperatorKind(operator),
		}
	}

	return left
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
	default:
		log.Fatalf("Unexpected token found during parsing: %v", p.at())
		return nil
	}
}

// * ======= PUBLIC API ======= * \\

func ProduceAST(tokens []tokens.Token) ast.Program {
	parser := Parser{
		Tokens: tokens,
		Pos:    0,
	}

	program := ast.Program{
		Body: []ast.ASTNode{},
	}

	for parser.notEOF() {
		program.Body = append(program.Body, parser.parseStatement())
	}

	return program
}

package test_frontend

import (
	"os"
	FE "pop/frontend"
	"pop/frontend/types/ast"
	"testing"
)

const PARSER_FILE = "../mocks/parser-mock.pop"

func TestParser(t *testing.T) {
	return
	content, err := os.ReadFile(PARSER_FILE)
	if err != nil {
		t.Fatalf("Failed to read file %v", err)
	}

	tokensOut := FE.Tokenize(string(content))
	astOut := FE.ProduceAST(tokensOut)

	// We expect multiple statements in our mock file
	expectedStatementsCount := 17 // Update based on mock file content
	if len(astOut.Body) != expectedStatementsCount {
		t.Fatalf("Expected %d top-level statements, got %d", expectedStatementsCount, len(astOut.Body))
	}

	// Test 1: Variable declaration with let
	t.Run("VariableDeclaration_Let", func(t *testing.T) {
		node := astOut.Body[0]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}
		if varDecl.Identifier != "x" {
			t.Errorf("Expected identifier 'x', got '%s'", varDecl.Identifier)
		}
		if varDecl.Constant {
			t.Errorf("Expected non-constant declaration")
		}
		numLiteral, ok := varDecl.Value.(ast.NumericLiteralExprNode)
		if !ok {
			t.Fatalf("Expected NumericLiteralExprNode, got %T", varDecl.Value)
		}
		if numLiteral.Value != 42 {
			t.Errorf("Expected value 42, got %f", numLiteral.Value)
		}
	})

	// Test 2: Variable declaration with const
	t.Run("VariableDeclaration_Const", func(t *testing.T) {
		node := astOut.Body[1]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}
		if varDecl.Identifier != "y" {
			t.Errorf("Expected identifier 'y', got '%s'", varDecl.Identifier)
		}
		if !varDecl.Constant {
			t.Errorf("Expected constant declaration")
		}
		numLiteral, ok := varDecl.Value.(ast.NumericLiteralExprNode)
		if !ok {
			t.Fatalf("Expected NumericLiteralExprNode, got %T", varDecl.Value)
		}
		if numLiteral.Value != 100 {
			t.Errorf("Expected value 100, got %f", numLiteral.Value)
		}
	})

	// Test 3: Binary expression (addition)
	t.Run("BinaryExpression_Addition", func(t *testing.T) {
		node := astOut.Body[2]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}
		if varDecl.Identifier != "sum" {
			t.Errorf("Expected identifier 'sum', got '%s'", varDecl.Identifier)
		}

		binaryExpr, ok := varDecl.Value.(ast.BinaryExprNode)
		if !ok {
			t.Fatalf("Expected BinaryExprNode, got %T", varDecl.Value)
		}
		if binaryExpr.Operator != ast.Add {
			t.Errorf("Expected operator '+', got '%s'", binaryExpr.Operator)
		}

		left, ok := binaryExpr.Left.(ast.NumericLiteralExprNode)
		if !ok || left.Value != 10 {
			t.Errorf("Expected left operand 10")
		}
		right, ok := binaryExpr.Right.(ast.NumericLiteralExprNode)
		if !ok || right.Value != 20 {
			t.Errorf("Expected right operand 20")
		}
	})

	// Test 4: Binary expression (multiplication)
	t.Run("BinaryExpression_Multiplication", func(t *testing.T) {
		node := astOut.Body[3]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		binaryExpr, ok := varDecl.Value.(ast.BinaryExprNode)
		if !ok {
			t.Fatalf("Expected BinaryExprNode, got %T", varDecl.Value)
		}
		if binaryExpr.Operator != ast.Multiply {
			t.Errorf("Expected operator '*', got '%s'", binaryExpr.Operator)
		}
	})

	// Test 5: Comparison expression
	t.Run("ComparisonExpression_LessThan", func(t *testing.T) {
		node := astOut.Body[4]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		binaryExpr, ok := varDecl.Value.(ast.BinaryExprNode)
		if !ok {
			t.Fatalf("Expected BinaryExprNode, got %T", varDecl.Value)
		}
		if binaryExpr.Operator != ast.LessThan {
			t.Errorf("Expected operator '<', got '%s'", binaryExpr.Operator)
		}
	})

	// Test 6: Assignment expression
	t.Run("AssignmentExpression", func(t *testing.T) {
		node := astOut.Body[5]
		assignExpr, ok := node.(ast.AssignmentExprNode)
		if !ok {
			t.Fatalf("Expected AssignmentExprNode, got %T", node)
		}

		assignee, ok := assignExpr.Assignee.(ast.IdentifierExprNode)
		if !ok {
			t.Fatalf("Expected IdentifierExprNode assignee, got %T", assignExpr.Assignee)
		}
		if assignee.Symbol != "x" {
			t.Errorf("Expected assignee 'x', got '%s'", assignee.Symbol)
		}

		value, ok := assignExpr.Value.(ast.NumericLiteralExprNode)
		if !ok {
			t.Fatalf("Expected NumericLiteralExprNode value, got %T", assignExpr.Value)
		}
		if value.Value != 50 {
			t.Errorf("Expected value 50, got %f", value.Value)
		}
	})

	// Test 7: Function declaration
	t.Run("FunctionDeclaration", func(t *testing.T) {
		node := astOut.Body[6]
		fnDecl, ok := node.(ast.FunctionDeclarationNode)
		if !ok {
			t.Fatalf("Expected FunctionDeclarationNode, got %T", node)
		}
		if fnDecl.Name != "add" {
			t.Errorf("Expected function name 'add', got '%s'", fnDecl.Name)
		}
		if len(fnDecl.Params) != 2 {
			t.Fatalf("Expected 2 parameters, got %d", len(fnDecl.Params))
		}
		if fnDecl.Params[0] != "a" || fnDecl.Params[1] != "b" {
			t.Errorf("Expected params ['a', 'b'], got %v", fnDecl.Params)
		}
		if len(fnDecl.Body) != 1 {
			t.Fatalf("Expected 1 statement in body, got %d", len(fnDecl.Body))
		}

		// Check the return statement
		retStmt, ok := fnDecl.Body[0].(ast.ReturnStatementNode)
		if !ok {
			t.Fatalf("Expected ReturnStatementNode, got %T", fnDecl.Body[0])
		}

		binaryExpr, ok := retStmt.Value.(ast.BinaryExprNode)
		if !ok {
			t.Fatalf("Expected BinaryExprNode in return, got %T", retStmt.Value)
		}
		if binaryExpr.Operator != ast.Add {
			t.Errorf("Expected operator '+', got '%s'", binaryExpr.Operator)
		}
	})

	// Test 8: Function call
	t.Run("FunctionCall", func(t *testing.T) {
		node := astOut.Body[7]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		callExpr, ok := varDecl.Value.(ast.CallExprNode)
		if !ok {
			t.Fatalf("Expected CallExprNode, got %T", varDecl.Value)
		}

		caller, ok := callExpr.Caller.(ast.IdentifierExprNode)
		if !ok {
			t.Fatalf("Expected IdentifierExprNode caller, got %T", callExpr.Caller)
		}
		if caller.Symbol != "add" {
			t.Errorf("Expected caller 'add', got '%s'", caller.Symbol)
		}

		if len(callExpr.Args) != 2 {
			t.Fatalf("Expected 2 arguments, got %d", len(callExpr.Args))
		}
	})

	// Test 9: String literal
	t.Run("StringLiteral", func(t *testing.T) {
		node := astOut.Body[8]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		strLiteral, ok := varDecl.Value.(ast.StringLiteralExprNode)
		if !ok {
			t.Fatalf("Expected StringLiteralExprNode, got %T", varDecl.Value)
		}
		if strLiteral.Value != "hello" {
			t.Errorf("Expected string 'hello', got '%s'", strLiteral.Value)
		}
	})

	// Test 10: Array literal
	t.Run("ArrayLiteral", func(t *testing.T) {
		node := astOut.Body[9]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		arrLiteral, ok := varDecl.Value.(ast.ArrayLiteralExprNode)
		if !ok {
			t.Fatalf("Expected ArrayLiteralExprNode, got %T", varDecl.Value)
		}
		if len(arrLiteral.Elements) != 5 {
			t.Fatalf("Expected 5 elements, got %d", len(arrLiteral.Elements))
		}
		if arrLiteral.Size != 5 {
			t.Errorf("Expected size 5, got %d", arrLiteral.Size)
		}

		// Check first element
		firstElem, ok := arrLiteral.Elements[0].(ast.NumericLiteralExprNode)
		if !ok || firstElem.Value != 1 {
			t.Errorf("Expected first element to be 1")
		}
	})

	// Test 11: Object literal
	t.Run("ObjectLiteral", func(t *testing.T) {
		node := astOut.Body[10]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		objLiteral, ok := varDecl.Value.(ast.ObjectLiteralExprNode)
		if !ok {
			t.Fatalf("Expected ObjectLiteralExprNode, got %T", varDecl.Value)
		}
		if len(objLiteral.Properties) != 2 {
			t.Fatalf("Expected 2 properties, got %d", len(objLiteral.Properties))
		}

		// Check first property
		if objLiteral.Properties[0].Key != "name" {
			t.Errorf("Expected first property key 'name', got '%s'", objLiteral.Properties[0].Key)
		}

		strValue, ok := objLiteral.Properties[0].Value.(ast.StringLiteralExprNode)
		if !ok || strValue.Value != "John" {
			t.Errorf("Expected first property value 'John'")
		}

		// Check second property
		if objLiteral.Properties[1].Key != "age" {
			t.Errorf("Expected second property key 'age', got '%s'", objLiteral.Properties[1].Key)
		}
	})

	// Test 12: Member expression (dot notation)
	t.Run("MemberExpression_Dot", func(t *testing.T) {
		node := astOut.Body[11]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		memberExpr, ok := varDecl.Value.(ast.MemberExprNode)
		if !ok {
			t.Fatalf("Expected MemberExprNode, got %T", varDecl.Value)
		}

		if memberExpr.Computed {
			t.Errorf("Expected non-computed (dot notation) access")
		}

		object, ok := memberExpr.Object.(ast.IdentifierExprNode)
		if !ok || object.Symbol != "person" {
			t.Errorf("Expected object 'person'")
		}

		property, ok := memberExpr.Property.(ast.IdentifierExprNode)
		if !ok || property.Symbol != "name" {
			t.Errorf("Expected property 'name'")
		}
	})

	// Test 13: Member expression (bracket notation)
	t.Run("MemberExpression_Bracket", func(t *testing.T) {
		node := astOut.Body[12]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		memberExpr, ok := varDecl.Value.(ast.MemberExprNode)
		if !ok {
			t.Fatalf("Expected MemberExprNode, got %T", varDecl.Value)
		}

		if !memberExpr.Computed {
			t.Errorf("Expected computed (bracket notation) access")
		}

		object, ok := memberExpr.Object.(ast.IdentifierExprNode)
		if !ok || object.Symbol != "numbers" {
			t.Errorf("Expected object 'numbers'")
		}

		property, ok := memberExpr.Property.(ast.NumericLiteralExprNode)
		if !ok || property.Value != 0 {
			t.Errorf("Expected property index 0")
		}
	})

	// Test 14: Comparison operators (==)
	t.Run("ComparisonOperator_Equal", func(t *testing.T) {
		node := astOut.Body[13]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		binaryExpr, ok := varDecl.Value.(ast.BinaryExprNode)
		if !ok {
			t.Fatalf("Expected BinaryExprNode, got %T", varDecl.Value)
		}
		if binaryExpr.Operator != ast.Equal {
			t.Errorf("Expected operator '==', got '%s'", binaryExpr.Operator)
		}
	})

	// Test 15: Comparison operators (!=)
	t.Run("ComparisonOperator_NotEqual", func(t *testing.T) {
		node := astOut.Body[14]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		binaryExpr, ok := varDecl.Value.(ast.BinaryExprNode)
		if !ok {
			t.Fatalf("Expected BinaryExprNode, got %T", varDecl.Value)
		}
		if binaryExpr.Operator != ast.NotEqual {
			t.Errorf("Expected operator '!=', got '%s'", binaryExpr.Operator)
		}
	})

	// Test 16: Comparison operators (<)
	t.Run("ComparisonOperator_LessThan", func(t *testing.T) {
		node := astOut.Body[15]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		binaryExpr, ok := varDecl.Value.(ast.BinaryExprNode)
		if !ok {
			t.Fatalf("Expected BinaryExprNode, got %T", varDecl.Value)
		}
		if binaryExpr.Operator != ast.LessThan {
			t.Errorf("Expected operator '<', got '%s'", binaryExpr.Operator)
		}
	})

	// Test 17: Comparison operators (>)
	t.Run("ComparisonOperator_GreaterThan", func(t *testing.T) {
		node := astOut.Body[16]
		varDecl, ok := node.(ast.VariableDeclarationNode)
		if !ok {
			t.Fatalf("Expected VariableDeclarationNode, got %T", node)
		}

		binaryExpr, ok := varDecl.Value.(ast.BinaryExprNode)
		if !ok {
			t.Fatalf("Expected BinaryExprNode, got %T", varDecl.Value)
		}
		if binaryExpr.Operator != ast.GreaterThan {
			t.Errorf("Expected operator '>', got '%s'", binaryExpr.Operator)
		}
	})
}

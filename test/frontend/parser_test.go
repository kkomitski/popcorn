package test_frontend

import (
	"os"
	FE "pop/frontend"
	"pop/frontend/types/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const PARSER_FILE = "../mocks/parser-mock.pop"

// Update based on mock file content
const EXPECTED_STATEMENTS_COUNT = 20

func assertVariableDeclaration(t *testing.T, node ast.ASTNode, expectedIdentifier string, shouldBeConstant bool) ast.VariableDeclarationNode {
	varDecl, ok := node.(ast.VariableDeclarationNode)
	require.True(t, ok, "Expected VariableDeclarationNode, got %T", node)

	assert.Equalf(t, expectedIdentifier, varDecl.Identifier, "Expected identifier to be '%s', instead got '%s'", expectedIdentifier, varDecl.Identifier)

	if shouldBeConstant {
		assert.True(t, varDecl.Constant, "Should be constant")
	} else {
		assert.False(t, varDecl.Constant, "Should not be constant")
	}

	return varDecl
}

func TestParser(t *testing.T) {
	content, err := os.ReadFile(PARSER_FILE)
	if err != nil {
		t.Fatalf("Failed to read file %v", err)
	}

	tokensOut := FE.Tokenize(string(content))
	astOut := FE.ProduceAST(tokensOut)

	t.Run("Have enough statements in the mock file", func(t *testing.T) {
		require.Equal(t, EXPECTED_STATEMENTS_COUNT, len(astOut.Body))
	})

	// Test 1: Variable declaration with let
	t.Run("VariableDeclaration_Let", func(t *testing.T) {
		node := astOut.Body[0]

		varDecl := assertVariableDeclaration(t, node, "x", false)

		// Use require again for type assertion
		numLiteral, ok := varDecl.Value.(ast.NumericLiteralExprNode)
		require.True(t, ok, "Expected NumericLiteralExprNode, got %T", varDecl.Value)
		assert.Equal(t, float64(42), numLiteral.Value, "Value should be 42")
	})

	// Test 2: Variable declaration with const
	t.Run("VariableDeclaration_Const", func(t *testing.T) {
		node := astOut.Body[1]

		varDecl := assertVariableDeclaration(t, node, "y", true)

		numLiteral, ok := varDecl.Value.(ast.NumericLiteralExprNode)
		require.True(t, ok, "Expected NumericLiteralExprNode, got %T", node)
		assert.Equal(t, float64(100), numLiteral.Value, "Value should be 100")
	})

	// Test 3: Binary expression (addition)
	t.Run("BinaryExpression_Addition", func(t *testing.T) {
		node := astOut.Body[2]

		varDecl := assertVariableDeclaration(t, node, "sum", false)

		binExpr, ok := varDecl.Value.(ast.BinaryExprNode)
		require.True(t, ok, "Expected 'BinaryExprNode', got %T", node)

		left, isLeftCorrectType := binExpr.Left.(ast.NumericLiteralExprNode)
		right, isRightCorrectType := binExpr.Right.(ast.NumericLiteralExprNode)

		require.True(t, isLeftCorrectType, "Expected left side of binary expression to be NumericLiteralExprNode")
		require.True(t, isRightCorrectType, "Expected right side of binary expression to be NumericLiteralExprNode")

		assert.Equal(t, binExpr.Operator, ast.BinaryOperatorKind("+"), "Expected operator for addition binary expressions to be '+', instead got %s", binExpr.Operator)

		assert.Equal(t, left.Value, float64(10), "Expected left side of binary expression to be '10', got %v", left.Value)
		assert.Equal(t, right.Value, float64(20), "Expected right side of binary expression to be '10', got %v", right.Value)
	})

	// Test 4: Binary expression (multiplication)
	t.Run("BinaryExpression_Multiplication", func(t *testing.T) {
		node := astOut.Body[3]

		varDecl := assertVariableDeclaration(t, node, "product", false)

		binExpr, ok := varDecl.Value.(ast.BinaryExprNode)
		require.True(t, ok, "Expected 'BinaryExprNode', got %T", node)

		left, isLeftCorrectType := binExpr.Left.(ast.NumericLiteralExprNode)
		right, isRightCorrectType := binExpr.Right.(ast.NumericLiteralExprNode)

		require.True(t, isLeftCorrectType, "Expected left side of binary expression to be NumericLiteralExprNode")
		require.True(t, isRightCorrectType, "Expected right side of binary expression to be NumericLiteralExprNode")

		assert.Equal(t, binExpr.Operator, ast.BinaryOperatorKind("*"), "Expected operator for multiplication binary expressions to be '*', instead got %s", binExpr.Operator)

		assert.Equal(t, left.Value, float64(5), "Expected left side of binary expression to be '5', got %v", left.Value)
		assert.Equal(t, right.Value, float64(6), "Expected right side of binary expression to be '6', got %v", right.Value)
	})

	// Test 5: Comparison expression
	t.Run("ComparisonOperators", func(t *testing.T) {
		tests := []struct {
			name       string
			nodeIndex  int
			identifier string
			operator   ast.BinaryOperatorKind
			leftValue  float64
			rightValue float64
		}{
			{"Equal", 4, "isEqual", ast.Equal, 10, 10},
			{"NotEqual", 5, "isNotEqual", ast.NotEqual, 5, 10},
			{"LessThan", 6, "isLessThan", ast.LessThan, 3, 5},
			{"GreaterThan", 7, "isGreaterThan", ast.GreaterThan, 10, 5},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				node := astOut.Body[tt.nodeIndex]

				varDecl := assertVariableDeclaration(t, node, tt.identifier, false)

				binExpr, ok := varDecl.Value.(ast.BinaryExprNode)
				require.True(t, ok, "Expected BinaryExprNode, got %T", varDecl.Value)

				left, ok := binExpr.Left.(ast.NumericLiteralExprNode)
				require.True(t, ok, "Expected left side to be NumericLiteralExprNode, got %T", binExpr.Left)

				right, ok := binExpr.Right.(ast.NumericLiteralExprNode)
				require.True(t, ok, "Expected right side to be NumericLiteralExprNode, got %T", binExpr.Right)

				assert.Equal(t, tt.operator, binExpr.Operator, "Expected operator '%s'", tt.operator)
				assert.Equal(t, tt.leftValue, left.Value, "Expected left value to be %v", tt.leftValue)
				assert.Equal(t, tt.rightValue, right.Value, "Expected right value to be %v", tt.rightValue)
			})
		}
	})

	// Test 5: Assignment expression
	t.Run("AssignmentExpression", func(t *testing.T) {
		node := astOut.Body[8]

		varAssignment, ok := node.(ast.AssignmentExprNode)
		require.True(t, ok, "Expected AssignmentExprNode, got %T", node)

		assignee, ok := varAssignment.Assignee.(ast.IdentifierExprNode)
		require.True(t, ok, "Expected assignee to be of type 'IdentifierExprNode', got %T", node)
		assert.Equal(t, assignee.Symbol, "x", "Expected assignment variable name to be 'x', got %s", assignee.Symbol)

		val, ok := varAssignment.Value.(ast.NumericLiteralExprNode)
		require.True(t, ok, "Expected Value property to be of type 'NumericLiteralExprNode', got %T", node)
		assert.Equal(t, val.Value, float64(50), "Expected value to be '50'")
	})

	// Test 7: Function declaration
	t.Run("FunctionDeclaration", func(t *testing.T) {
		node := astOut.Body[9]

		fnDecl, ok := node.(ast.FunctionDeclarationNode)
		require.True(t, ok, "Expected FunctionDeclarationNode, got %T", node)

		assert.Equal(t, "add", fnDecl.Name, "Expected function name to be 'add'")
		require.Len(t, fnDecl.Params, 2, "Expected 2 parameters")
		assert.Equal(t, "a", fnDecl.Params[0], "Expected first parameter to be 'a'")
		assert.Equal(t, "b", fnDecl.Params[1], "Expected second parameter to be 'b'")
		require.Len(t, fnDecl.Body, 1, "Expected 1 statement in body")

		retStmt, ok := fnDecl.Body[0].(ast.ReturnStatementNode)
		require.True(t, ok, "Expected ReturnStatementNode, got %T", fnDecl.Body[0])

		binaryExpr, ok := retStmt.Value.(ast.BinaryExprNode)
		require.True(t, ok, "Expected BinaryExprNode in return, got %T", retStmt.Value)
		assert.Equal(t, ast.Add, binaryExpr.Operator, "Expected operator to be '+'")
	})

	// Test 8: Function call
	t.Run("FunctionCall", func(t *testing.T) {
		node := astOut.Body[10]

		varDecl := assertVariableDeclaration(t, node, "result", false)

		callExpr, ok := varDecl.Value.(ast.CallExprNode)
		require.True(t, ok, "Expected CallExprNode, got %T", varDecl.Value)

		caller, ok := callExpr.Caller.(ast.IdentifierExprNode)
		require.True(t, ok, "Expected IdentifierExprNode caller, got %T", callExpr.Caller)
		assert.Equal(t, "add", caller.Symbol, "Expected caller to be 'add'")

		require.Len(t, callExpr.Args, 2, "Expected 2 arguments")
	})

	// Test 9: String literal
	t.Run("StringLiteral", func(t *testing.T) {
		node := astOut.Body[11]

		varDecl := assertVariableDeclaration(t, node, "greeting", false)

		strLiteral, ok := varDecl.Value.(ast.StringLiteralExprNode)
		require.True(t, ok, "Expected StringLiteralExprNode, got %T", varDecl.Value)
		assert.Equal(t, "hello", strLiteral.Value, "Expected string to be 'hello'")
	})

	// Test 10: Array literal
	t.Run("ArrayLiteral", func(t *testing.T) {
		node := astOut.Body[12]

		varDecl := assertVariableDeclaration(t, node, "numbers", false)

		arrLiteral, ok := varDecl.Value.(ast.ArrayLiteralExprNode)
		require.True(t, ok, "Expected ArrayLiteralExprNode, got %T", varDecl.Value)
		require.Len(t, arrLiteral.Elements, 5, "Expected 5 elements")
		assert.Equal(t, int64(5), arrLiteral.Size, "Expected size to be 5")

		firstElem, ok := arrLiteral.Elements[0].(ast.NumericLiteralExprNode)
		require.True(t, ok, "Expected first element to be NumericLiteralExprNode")
		assert.Equal(t, float64(1), firstElem.Value, "Expected first element to be 1")
	})

	// Test 11: Object literal
	t.Run("ObjectLiteral", func(t *testing.T) {
		node := astOut.Body[13]

		varDecl := assertVariableDeclaration(t, node, "person", false)

		objLiteral, ok := varDecl.Value.(ast.ObjectLiteralExprNode)
		require.True(t, ok, "Expected ObjectLiteralExprNode, got %T", varDecl.Value)
		require.Len(t, objLiteral.Properties, 2, "Expected 2 properties")

		assert.Equal(t, "name", objLiteral.Properties[0].Key, "Expected first property key to be 'name'")

		strValue, ok := objLiteral.Properties[0].Value.(ast.StringLiteralExprNode)
		require.True(t, ok, "Expected first property value to be StringLiteralExprNode")
		assert.Equal(t, "John", strValue.Value, "Expected first property value to be 'John'")

		assert.Equal(t, "age", objLiteral.Properties[1].Key, "Expected second property key to be 'age'")
	})

	// Test 12: Member expression (dot notation)
	t.Run("MemberExpression_Dot", func(t *testing.T) {
		node := astOut.Body[14]

		varDecl := assertVariableDeclaration(t, node, "personName", false)

		memberExpr, ok := varDecl.Value.(ast.MemberExprNode)
		require.True(t, ok, "Expected MemberExprNode, got %T", varDecl.Value)

		assert.False(t, memberExpr.Computed, "Should be non-computed (dot notation) access")

		object, ok := memberExpr.Object.(ast.IdentifierExprNode)
		require.True(t, ok, "Expected object to be IdentifierExprNode")
		assert.Equal(t, "person", object.Symbol, "Expected object to be 'person'")

		property, ok := memberExpr.Property.(ast.IdentifierExprNode)
		require.True(t, ok, "Expected property to be IdentifierExprNode")
		assert.Equal(t, "name", property.Symbol, "Expected property to be 'name'")
	})

	// Test 13: Member expression (bracket notation)
	t.Run("MemberExpression_Bracket", func(t *testing.T) {
		node := astOut.Body[15]

		varDecl := assertVariableDeclaration(t, node, "firstNumber", false)

		memberExpr, ok := varDecl.Value.(ast.MemberExprNode)
		require.True(t, ok, "Expected MemberExprNode, got %T", varDecl.Value)

		assert.True(t, memberExpr.Computed, "Should be computed (bracket notation) access")

		object, ok := memberExpr.Object.(ast.IdentifierExprNode)
		require.True(t, ok, "Expected object to be IdentifierExprNode")
		assert.Equal(t, "numbers", object.Symbol, "Expected object to be 'numbers'")

		property, ok := memberExpr.Property.(ast.NumericLiteralExprNode)
		require.True(t, ok, "Expected property to be NumericLiteralExprNode")
		assert.Equal(t, float64(0), property.Value, "Expected property index to be 0")
	})

	// Test 14: Logical operators
	t.Run("LogicalOperators", func(t *testing.T) {
		tests := []struct {
			name       string
			nodeIndex  int
			identifier string
			operator   ast.BinaryOperatorKind
			leftValue  bool
			rightValue bool
		}{
			{"And_True", 16, "isTrue", ast.And, true, true},
			{"And_False", 17, "isFalse", ast.And, true, false},
			{"Or_True", 18, "isTrue", ast.Or, true, false},
			{"Or_False", 19, "isFalse", ast.Or, false, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				node := astOut.Body[tt.nodeIndex]

				varDecl := assertVariableDeclaration(t, node, tt.identifier, false)

				logicalExpr, ok := varDecl.Value.(ast.LogicalExprNode)
				require.True(t, ok, "Expected LogicalExprNode, got %T", varDecl.Value)

				left, ok := logicalExpr.Left.(ast.BooleanLiteralExprNode)
				require.True(t, ok, "Expected left side to be BooleanLiteralExprNode, got %T", logicalExpr.Left)

				right, ok := logicalExpr.Right.(ast.BooleanLiteralExprNode)
				require.True(t, ok, "Expected right side to be BooleanLiteralExprNode, got %T", logicalExpr.Right)

				assert.Equal(t, tt.operator, logicalExpr.Operator, "Expected operator '%s'", tt.operator)
				assert.Equal(t, tt.leftValue, left.Value, "Expected left value to be %v", tt.leftValue)
				assert.Equal(t, tt.rightValue, right.Value, "Expected right value to be %v", tt.rightValue)
			})
		}
	})
}

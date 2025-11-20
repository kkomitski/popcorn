package ast

import (
	"encoding/json"
	"fmt"
)

type NodeKind int

const (
	// * ==================== Statements ==================== *

	/* Root node */
	ProgramStatement NodeKind = iota

	/* For `let`, `const` declarations */
	VariableDeclaration

	/* For function definitions */
	FunctionDeclaration

	/* For `if` statements */
	IfStatement

	/* For `while` loops */
	WhileStatement

	/* For `for` loops */
	ForStatement

	/* For `return` statements */
	ReturnStatement

	/* For blocks of statements enclosed in braces */
	BlockStatement

	// * ==================== Expressions ==================== *

	/* For assignment expressions (e.g., a = b) */
	AssignmentExpr

	/* For identifiers (variable and function names) */
	IdentifierExpr

	/* For binary operations (e.g., a + b) */
	BinaryExpr

	/* For logical operations (e.g., a && b, a || b) */
	LogicalExpr

	/* For function calls (e.g., foo()) */
	CallExpr

	/* For unary operations (e.g., -a, !a) */
	UnaryExpr

	/* For member access (e.g., obj.prop) */
	MemberExpr

	/* For index access (e.g., arr[0]) */
	IndexExpr

	/* For conditional/ternary expressions (e.g., a ? b : c) */
	ConditionalExpr

	// * ==================== Literals ==================== *

	/* For numeric literals (e.g., 42, 3.14) */
	NumericLiteral

	/* For string literals (e.g., "hello") */
	StringLiteral

	/* For boolean literals (e.g., true, false) */
	BooleanLiteral

	/* For null literals */
	NullLiteral

	/* For array literals (e.g., [1, 2, 3]) */
	ArrayLiteral

	/* For object literals (e.g., { key: value }) */
	ObjectLiteral

	/* For object properties */
	Property
)

type BinaryOperatorKind string

const (
	Add                BinaryOperatorKind = "+"
	Subtract           BinaryOperatorKind = "-"
	Multiply           BinaryOperatorKind = "*"
	Divide             BinaryOperatorKind = "/"
	Modulo             BinaryOperatorKind = "%"
	Equal              BinaryOperatorKind = "=="
	NotEqual           BinaryOperatorKind = "!="
	LessThan           BinaryOperatorKind = "<"
	GreaterThan        BinaryOperatorKind = ">"
	LessThanOrEqual    BinaryOperatorKind = "<="
	GreaterThanOrEqual BinaryOperatorKind = ">="
	And                BinaryOperatorKind = "&&"
	Or                 BinaryOperatorKind = "||"
)

// ASTNode can be any AST node type
type ASTNode any

// GetNodeKind returns the NodeKind for any ASTNode using type switch
func GetNodeKind(node ASTNode) NodeKind {
	switch node.(type) {
	case Program, *Program:
		return ProgramStatement
	case VariableDeclarationNode, *VariableDeclarationNode:
		return VariableDeclaration
	case FunctionDeclarationNode, *FunctionDeclarationNode:
		return FunctionDeclaration
	case AssignmentExprNode, *AssignmentExprNode:
		return AssignmentExpr
	case BinaryExprNode, *BinaryExprNode:
		return BinaryExpr
	case MemberExprNode, *MemberExprNode:
		return MemberExpr
	case CallExprNode, *CallExprNode:
		return CallExpr
	case IdentifierExprNode, *IdentifierExprNode:
		return IdentifierExpr
	case NumericLiteralExprNode, *NumericLiteralExprNode:
		return NumericLiteral
	case StringLiteralExprNode, *StringLiteralExprNode:
		return StringLiteral
	case BooleanLiteralExprNode, *BooleanLiteralExprNode:
		return BooleanLiteral
	case NullLiteralExprNode, *NullLiteralExprNode:
		return NullLiteral
	case ArrayLiteralExprNode, *ArrayLiteralExprNode:
		return ArrayLiteral
	case PropertyNode, *PropertyNode:
		return Property
	case ObjectLiteralExprNode, *ObjectLiteralExprNode:
		return ObjectLiteral
	case UnaryExprNode, *UnaryExprNode:
		return UnaryExpr
	case LogicalExprNode, *LogicalExprNode:
		return LogicalExpr
	case ConditionalExprNode, *ConditionalExprNode:
		return ConditionalExpr
	case IndexExprNode, *IndexExprNode:
		return IndexExpr
	case IfStatementNode, *IfStatementNode:
		return IfStatement
	case WhileStatementNode, *WhileStatementNode:
		return WhileStatement
	case ForStatementNode, *ForStatementNode:
		return ForStatement
	case ReturnStatementNode, *ReturnStatementNode:
		return ReturnStatement
	case BlockStatementNode, *BlockStatementNode:
		return BlockStatement
	default:
		return -1
	}
}

// GetNodeKind returns the NodeKind for any ASTNode using type switch
func GetNodeKindAsString(node ASTNode) string {
	switch node.(type) {
	case Program, *Program:
		return "ProgramStatement"
	case VariableDeclarationNode, *VariableDeclarationNode:
		return "VariableDeclaration"
	case FunctionDeclarationNode, *FunctionDeclarationNode:
		return "FunctionDeclaration"
	case AssignmentExprNode, *AssignmentExprNode:
		return "AssignmentExpr"
	case BinaryExprNode, *BinaryExprNode:
		return "BinaryExpr"
	case MemberExprNode, *MemberExprNode:
		return "MemberExpr"
	case CallExprNode, *CallExprNode:
		return "CallExpr"
	case IdentifierExprNode, *IdentifierExprNode:
		return "IdentifierExpr"
	case NumericLiteralExprNode, *NumericLiteralExprNode:
		return "NumericLiteral"
	case StringLiteralExprNode, *StringLiteralExprNode:
		return "StringLiteral"
	case BooleanLiteralExprNode, *BooleanLiteralExprNode:
		return "BooleanLiteral"
	case NullLiteralExprNode, *NullLiteralExprNode:
		return "NullLiteral"
	case ArrayLiteralExprNode, *ArrayLiteralExprNode:
		return "ArrayLiteral"
	case PropertyNode, *PropertyNode:
		return "Property"
	case ObjectLiteralExprNode, *ObjectLiteralExprNode:
		return "ObjectLiteral"
	case UnaryExprNode, *UnaryExprNode:
		return "UnaryExpr"
	case LogicalExprNode, *LogicalExprNode:
		return "LogicalExpr"
	case ConditionalExprNode, *ConditionalExprNode:
		return "ConditionalExpr"
	case IndexExprNode, *IndexExprNode:
		return "IndexExpr"
	case IfStatementNode, *IfStatementNode:
		return "IfStatement"
	case WhileStatementNode, *WhileStatementNode:
		return "WhileStatement"
	case ForStatementNode, *ForStatementNode:
		return "ForStatement"
	case ReturnStatementNode, *ReturnStatementNode:
		return "ReturnStatement"
	case BlockStatementNode, *BlockStatementNode:
		return "BlockStatement"
	default:
		return "ERR_UNKNOWN"
	}
}

// JSONNode wraps an ASTNode with its kind for JSON marshalling
type JSONNode struct {
	Kind string  `json:"kind"`
	Data ASTNode `json:"data"`
}

// MarshalJSON custom marshaller for JSONNode
func (n JSONNode) MarshalJSON() ([]byte, error) {
	// First, marshal the data node to get its fields
	dataBytes, err := json.Marshal(n.Data)
	if err != nil {
		return nil, err
	}

	// Unmarshal into a map to manipulate
	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, err
	}

	// Add the kind field to the map
	dataMap["kind"] = GetNodeKindAsString(n.Data)

	// Marshal the combined map
	return json.Marshal(dataMap)
}

// UnmarshalJSON custom unmarshaller for JSONNode
func (n *JSONNode) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a map to extract the kind
	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return err
	}

	// Extract and remove the kind field
	kindStr, ok := dataMap["kind"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'kind' field")
	}
	n.Kind = kindStr
	delete(dataMap, "kind")

	// Re-marshal the remaining data (without kind)
	remainingData, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	// Create the appropriate node type based on kind
	var node ASTNode
	switch kindStr {
	case "ProgramStatement":
		node = &Program{}
	case "VariableDeclaration":
		node = &VariableDeclarationNode{}
	case "FunctionDeclaration":
		node = &FunctionDeclarationNode{}
	case "AssignmentExpr":
		node = &AssignmentExprNode{}
	case "BinaryExpr":
		node = &BinaryExprNode{}
	case "MemberExpr":
		node = &MemberExprNode{}
	case "CallExpr":
		node = &CallExprNode{}
	case "IdentifierExpr":
		node = &IdentifierExprNode{}
	case "NumericLiteral":
		node = &NumericLiteralExprNode{}
	case "StringLiteral":
		node = &StringLiteralExprNode{}
	case "BooleanLiteral":
		node = &BooleanLiteralExprNode{}
	case "NullLiteral":
		node = &NullLiteralExprNode{}
	case "ArrayLiteral":
		node = &ArrayLiteralExprNode{}
	case "Property":
		node = &PropertyNode{}
	case "ObjectLiteral":
		node = &ObjectLiteralExprNode{}
	case "UnaryExpr":
		node = &UnaryExprNode{}
	case "LogicalExpr":
		node = &LogicalExprNode{}
	case "ConditionalExpr":
		node = &ConditionalExprNode{}
	case "IndexExpr":
		node = &IndexExprNode{}
	case "IfStatement":
		node = &IfStatementNode{}
	case "WhileStatement":
		node = &WhileStatementNode{}
	case "ForStatement":
		node = &ForStatementNode{}
	case "ReturnStatement":
		node = &ReturnStatementNode{}
	case "BlockStatement":
		node = &BlockStatementNode{}
	default:
		return fmt.Errorf("unknown node kind: %s", kindStr)
	}

	// Unmarshal the remaining data into the node
	if err := json.Unmarshal(remainingData, node); err != nil {
		return err
	}
	n.Data = node
	return nil
}

// Program represents the root node of the AST.
// It contains all top-level statements in the program.
type Program struct {
	// Body contains all top-level statements
	Body []ASTNode
}

// VariableDeclarationNode represents a variable declaration statement in the AST.
// It handles both `let` and `const` declarations.
type VariableDeclarationNode struct {
	// Constant is true for `const` declarations, false for `let`
	Constant bool
	// Identifier is the variable name being declared
	Identifier string
	// Value is the initial value assigned to the variable
	Value ASTNode
}

// FunctionDeclarationNode represents a function declaration statement in the AST.
// It includes the function name, parameters, and body.
type FunctionDeclarationNode struct {
	// Params contains the parameter names for the function
	Params []string
	// Name is the function identifier
	Name string
	// Body contains the statements within the function
	Body []ASTNode
}

// AssignmentExprNode represents an assignment expression in the AST.
// It handles expressions like `x = 5` or `obj.prop = value`.
type AssignmentExprNode struct {
	// Assignee is the target being assigned to (e.g., x, obj.prop)
	Assignee ASTNode
	// Value is the expression being assigned
	Value ASTNode
}

// BinaryExprNode represents a binary operation expression in the AST.
// It handles operations like addition, subtraction, comparison, etc.
type BinaryExprNode struct {
	// Left is the left-hand operand
	Left ASTNode
	// Right is the right-hand operand
	Right ASTNode
	// Operator specifies the binary operation (e.g., +, -, ==)
	Operator BinaryOperatorKind
}

// MemberExprNode represents a member access expression in the AST.
// It supports both dot notation (obj.prop) and bracket notation (obj["prop"]).
type MemberExprNode struct {
	// Object is the expression being accessed (e.g., obj)
	Object ASTNode
	// Property is the property/key being accessed (e.g., prop or "prop")
	Property ASTNode
	// Computed is true if bracket notation is used, false for dot notation
	Computed bool
}

// CallExprNode represents a function call expression in the AST.
// It handles function invocations like `foo()` or `obj.method(arg1, arg2)`.
type CallExprNode struct {
	// Callee is the function being called
	Caller ASTNode
	// Args contains the arguments passed to the function
	Args []ASTNode
}

// IdentifierExprNode represents an identifier (variable or function name) in the AST.
type IdentifierExprNode struct {
	// Symbol is the identifier's name
	Symbol string
}

// NumericLiteralExprNode represents a numeric literal value in the AST.
type NumericLiteralExprNode struct {
	// Value is the numeric value
	Value float64
}

// StringLiteralExprNode represents a string literal value in the AST.
type StringLiteralExprNode struct {
	// Value is the string content
	Value string
}

// BooleanLiteralExprNode represents a boolean literal value in the AST.
type BooleanLiteralExprNode struct {
	// Value is the boolean value (true or false)
	Value bool
}

// NullLiteralExprNode represents a null literal value in the AST.
type NullLiteralExprNode struct{}

// ArrayLiteralExprNode represents an array literal expression in the AST.
// It contains a list of expressions enclosed in brackets (e.g., [1, 2, 3]).
type ArrayLiteralExprNode struct {
	// Elements contains all expressions in the array
	Elements []ASTNode
	Size     int64
}

// PropertyNode represents a key-value pair property inside an object literal in the AST.
// Key is the property's name, and Value is the expression assigned to that property.
type PropertyNode struct {
	// Key is the property's name (e.g., "foo" in {foo: 42})
	Key string
	// Value is the expression assigned to the property (e.g., 42 in {foo: 42})
	Value ASTNode
}

// ObjectLiteralExprNode represents an object literal expression in the AST.
// It contains key-value pairs defined within braces (e.g., {foo: 42, bar: "hello"}).
type ObjectLiteralExprNode struct {
	// Properties contains all key-value pairs in the object
	Properties []PropertyNode
}

// UnaryExprNode represents a unary operation expression in the AST.
// It handles operations like negation (-x) or logical NOT (!x).
type UnaryExprNode struct {
	// Operator specifies the unary operation (e.g., -, !, +)
	Operator string
	// Operand is the expression being operated on
	Operand ASTNode
}

// LogicalExprNode represents a logical operation expression in the AST.
// It handles logical AND (&&) and OR (||) operations.
type LogicalExprNode struct {
	// Left is the left-hand operand
	Left ASTNode
	// Right is the right-hand operand
	Right ASTNode
	// Operator specifies the logical operation ("&&" or "||")
	Operator BinaryOperatorKind
}

// ConditionalExprNode represents a ternary conditional expression in the AST.
// It handles expressions like `condition ? ifTrue : ifFalse`.
type ConditionalExprNode struct {
	// Condition is the test expression
	Condition ASTNode
	// Consequent is the expression evaluated if condition is true
	Consequent ASTNode
	// Alternate is the expression evaluated if condition is false
	Alternate ASTNode
}

// IndexExprNode represents an index access expression in the AST.
// It handles array/object indexing like `arr[0]` or `obj[key]`.
type IndexExprNode struct {
	// Object is the expression being indexed
	Object ASTNode
	// Index is the index expression
	Index ASTNode
}

// IfStatementNode represents an if statement in the AST.
// It handles conditional execution with optional else branches.
type IfStatementNode struct {
	// Condition is the test expression
	Condition ASTNode
	// Consequent is the statement/block executed if condition is true
	Consequent ASTNode
	// Alternate is the optional else statement/block
	Alternate ASTNode
}

// WhileStatementNode represents a while loop in the AST.
type WhileStatementNode struct {
	// Condition is the loop test expression
	Condition ASTNode
	// Body is the statement/block executed while condition is true
	Body ASTNode
}

// ForStatementNode represents a for loop in the AST.
type ForStatementNode struct {
	// Init is the optional initialization statement
	Init ASTNode
	// Condition is the optional loop test expression
	Condition ASTNode
	// Update is the optional update expression
	Update ASTNode
	// Body is the statement/block executed in each iteration
	Body ASTNode
}

// ReturnStatementNode represents a return statement in the AST.
type ReturnStatementNode struct {
	// Value is the optional expression being returned
	Value ASTNode
}

// BlockStatementNode represents a block of statements enclosed in braces in the AST.
type BlockStatementNode struct {
	// Body contains all statements within the block
	Body []ASTNode
}

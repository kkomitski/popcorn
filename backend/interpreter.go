package backend

import (
	"fmt"
	"log"
	"pop/frontend/types/ast"
)

func evalAssignment(node ast.AssignmentExprNode, env *Environment) RuntimeVal {
	// Only allow assignment to identifiers for now
	ident, ok := node.Assignee.(ast.IdentifierExprNode)
	if !ok {
		log.Fatalf("Invalid LHS in assignment: %+v", node.Assignee)
	}
	val := Evaluate(node.Value, env)
	return env.AssignVar(ident.Symbol, val)
}

func evalObjectLiteral(node ast.ObjectLiteralExprNode, env *Environment) RuntimeVal {
	obj := ObjectVal{Properties: make(map[string]RuntimeVal)}
	for _, prop := range node.Properties {
		var val RuntimeVal
		if prop.Value == nil {
			val = env.GetVar(prop.Key)
		} else {
			val = Evaluate(prop.Value, env)
		}
		obj.Properties[prop.Key] = val
	}
	return obj
}

func evalCallExpression(node ast.CallExprNode, env *Environment) RuntimeVal {
	callee := Evaluate(node.Caller, env)
	args := make([]RuntimeVal, len(node.Args))
	for i, arg := range node.Args {
		args[i] = Evaluate(arg, env)
	}

	switch fn := callee.(type) {
	case NativeFunctionVal:
		return fn.Call(args, env)
	case FunctionVal:
		scope := MakeEnvironment()
		scope.Parent = fn.DeclarationEnv
		for i, param := range fn.Params {
			if i < len(args) {
				scope.DeclareVar(param, false, args[i])
			} else {
				scope.DeclareVar(param, false, Null)
			}
		}
		var result RuntimeVal = Null
		for _, stmt := range fn.Body {
			result = Evaluate(stmt, scope)
			// Check if a return statement was executed
			if retVal, isReturn := result.(ReturnVal); isReturn {
				return retVal.Value
			}
		}
		return result
	default:
		log.Fatalf("Cannot call value that is not a function: %+v", callee)
	}
	return Null
}

func evalVarDeclaration(node ast.VariableDeclarationNode, env *Environment) RuntimeVal {
	var val RuntimeVal

	if node.Value != nil {
		val = Evaluate(node.Value, env)
	} else {
		val = Null
	}

	env.DeclareVar(node.Identifier, node.Constant, val)
	return val
}

func evalFnDeclaration(node ast.FunctionDeclarationNode, env *Environment) RuntimeVal {
	fn := FunctionVal{
		Name:           node.Name,
		Params:         node.Params,
		DeclarationEnv: env,
		Body:           node.Body,
	}
	env.DeclareVar(node.Name, true, fn)
	return fn
}

func evalProgram(node ast.Program, env *Environment) RuntimeVal {
	var final RuntimeVal = Null

	for _, stmt := range node.Body {
		final = Evaluate(stmt, env)
	}

	return final
}

func evalNumber(node ast.NumericLiteralExprNode, env *Environment) RuntimeVal {
	return NumberVal{Value: node.Value}
}

func evalString(node ast.StringLiteralExprNode, env *Environment) RuntimeVal {
	return StringVal{
		Value: node.Value,
	}
}

func evalBool(node ast.BooleanLiteralExprNode, env *Environment) RuntimeVal {
	return BoolValue{Value: node.Value}
}

func evalLogicalExpr(node ast.LogicalExprNode, env *Environment) RuntimeVal {
	left := Evaluate(node.Left, env)
	leftBool, isLeftBool := left.(BoolValue)

	if !isLeftBool {
		log.Fatalf("Logical operators require boolean operands, got: %v", left)
	}

	// Short-circuit evaluation
	if node.Operator == "&&" {
		if !leftBool.Value {
			return BoolValue{Value: false}
		}
	} else if node.Operator == "||" {
		if leftBool.Value {
			return BoolValue{Value: true}
		}
	}

	right := Evaluate(node.Right, env)
	rightBool, isRightBool := right.(BoolValue)

	if !isRightBool {
		log.Fatalf("Logical operators require boolean operands, got: %v", right)
	}

	if node.Operator == "&&" {
		return BoolValue{Value: rightBool.Value}
	} else {
		return BoolValue{Value: rightBool.Value}
	}
}

func evalBinaryOp(node ast.BinaryExprNode, env *Environment) RuntimeVal {
	left := Evaluate(node.Left, env)
	right := Evaluate(node.Right, env)

	switch node.Operator {
	case "+", "-", "*", "/", "%":
		leftNum, leftIsNum := left.(NumberVal)
		rightNum, rightIsNum := right.(NumberVal)
		if !leftIsNum || !rightIsNum {
			log.Fatalf("Cannot perform arithmetic operation on non-number values: %v, %v", left, right)
		}
		switch node.Operator {
		case "+":
			return NumberVal{Value: leftNum.Value + rightNum.Value}
		case "-":
			return NumberVal{Value: leftNum.Value - rightNum.Value}
		case "*":
			return NumberVal{Value: leftNum.Value * rightNum.Value}
		case "/":
			return NumberVal{Value: leftNum.Value / rightNum.Value}
		case "%":
			return NumberVal{Value: float64(int(leftNum.Value) % int(rightNum.Value))}
		}
	case "==", "!=":
		// Allow equality/inequality for numbers, strings, booleans, and null
		switch l := left.(type) {
		case NumberVal:
			if r, ok := right.(NumberVal); ok {
				if node.Operator == "==" {
					return BoolValue{Value: l.Value == r.Value}
				} else {
					return BoolValue{Value: l.Value != r.Value}
				}
			}
		case StringVal:
			if r, ok := right.(StringVal); ok {
				if node.Operator == "==" {
					return BoolValue{Value: l.Value == r.Value}
				} else {
					return BoolValue{Value: l.Value != r.Value}
				}
			}
		case BoolValue:
			if r, ok := right.(BoolValue); ok {
				if node.Operator == "==" {
					return BoolValue{Value: l.Value == r.Value}
				} else {
					return BoolValue{Value: l.Value != r.Value}
				}
			}
			// Null equality: compare with Null singleton value
		}
		// Null equality: if left is Null
		if left == Null {
			if right == Null {
				if node.Operator == "==" {
					return BoolValue{Value: true}
				} else {
					return BoolValue{Value: false}
				}
			} else {
				if node.Operator == "==" {
					return BoolValue{Value: false}
				} else {
					return BoolValue{Value: true}
				}
			}
		}
		// If types don't match, not equal
		if node.Operator == "==" {
			return BoolValue{Value: false}
		} else {
			return BoolValue{Value: true}
		}
	case "<", ">", "<=", ">=":
		leftNum, leftIsNum := left.(NumberVal)
		rightNum, rightIsNum := right.(NumberVal)
		if !leftIsNum || !rightIsNum {
			log.Fatalf("Cannot perform comparison operation on non-number values: %v, %v", left, right)
		}
		switch node.Operator {
		case "<":
			return BoolValue{Value: leftNum.Value < rightNum.Value}
		case ">":
			return BoolValue{Value: leftNum.Value > rightNum.Value}
		case "<=":
			return BoolValue{Value: leftNum.Value <= rightNum.Value}
		case ">=":
			return BoolValue{Value: leftNum.Value >= rightNum.Value}
		}
	default:
		log.Fatalf("Unknown binary operator: %s", node.Operator)
	}
	return Null
}

func evalUnaryOp(node ast.UnaryExprNode, env *Environment) RuntimeVal {
	right := Evaluate(node.Operand, env)

	switch node.Operator {
	case "!":
		rightBool, isRightBool := right.(BoolValue)
		if !isRightBool {
			log.Fatalf("Cannot negate a non-bool value! %v", right)
		}
		return BoolValue{Value: !rightBool.Value}
	case "-":
		rightNum, isRightNum := right.(NumberVal)
		if !isRightNum {
			log.Fatalf("Cannot negate a non-number value! %v", right)
		}
		return NumberVal{Value: -rightNum.Value}
	default:
		log.Fatalf("Unknown unary operator: %v", node.Operator)
		return Null // unreachable, but keeps compiler happy
	}
}

func evalVarLookup(node ast.IdentifierExprNode, env *Environment) RuntimeVal {
	return env.GetVar(node.Symbol)
}

func evalReturnStatement(node ast.ReturnStatementNode, env *Environment) RuntimeVal {
	var value RuntimeVal = Null
	if node.Value != nil {
		value = Evaluate(node.Value, env)
	}
	return ReturnVal{Value: value}
}

func evalArray(node ast.ArrayLiteralExprNode, env *Environment) RuntimeVal {
	// Pre-allocate slice with exact capacity needed
	elements := make([]RuntimeVal, len(node.Elements))

	// Evaluate each element
	for i, elem := range node.Elements {
		elements[i] = Evaluate(elem, env)
	}

	return ArrayVal{Elements: elements}
}

func evalMember(node ast.MemberExprNode, env *Environment) RuntimeVal {
	object := Evaluate(node.Object, env)

	// Computed access: obj[expr] or array[index]
	if node.Computed {
		property := Evaluate(node.Property, env)

		// Array access
		if arr, isArray := object.(ArrayVal); isArray {
			index, isNum := property.(NumberVal)
			if !isNum {
				log.Fatalf("Array index must be a number, got: %+v", property)
			}
			idx := int(index.Value)
			if idx < 0 || idx >= len(arr.Elements) {
				log.Fatalf("Array index out of bounds: %d (length: %d)", idx, len(arr.Elements))
			}
			return arr.Elements[idx]
		}

		// Object computed access: obj[key]
		if obj, isObj := object.(ObjectVal); isObj {
			// Convert property to string key
			var key string
			if num, isNum := property.(NumberVal); isNum {
				key = fmt.Sprintf("%v", num.Value)
			} else {
				log.Fatalf("Object key must be string or number, got: %+v", property)
			}
			if val, exists := obj.Properties[key]; exists {
				return val
			}
			return Null
		}

		log.Fatalf("Cannot use computed access on non-object/array: %+v", object)
	}

	// Dot access: obj.property
	if obj, isObj := object.(ObjectVal); isObj {
		// Property should be an identifier
		ident, ok := node.Property.(ast.IdentifierExprNode)
		if !ok {
			log.Fatalf("Property in dot notation must be identifier, got: %+v", node.Property)
		}
		if val, exists := obj.Properties[ident.Symbol]; exists {
			return val
		}
		return Null
	}

	log.Fatalf("Cannot access property on non-object: %+v", object)
	return Null
}

func evalForLoop(node ast.ForStatementNode, env *Environment) RuntimeVal {
	// New scope for the for loop body
	loopEnv := MakeEnvironment()
	loopEnv.Parent = env

	// Evaluate the init to load it into env
	if node.Init != nil {
		Evaluate(node.Init, loopEnv)
	}

	for {
		// Re-evaluate the condition
		conditionVal := Evaluate(node.Condition, loopEnv)
		condition, isBoolCondition := conditionVal.(BoolValue)
		if !isBoolCondition {
			log.Fatalf("For loop condition does not evaluate to a boolean value: %v", condition)
		}

		// Condition is false
		if !condition.Value {
			break
		}

		// Evaluate the body
		Evaluate(node.Body, loopEnv)

		if node.Update != nil {
			Evaluate(node.Update, loopEnv)
		}
	}

	// For loops are statements so they don't resolve to a value
	return Null
}

func evalWhileLoop(node ast.WhileStatementNode, env *Environment) RuntimeVal {
	// New scope for the for loop body
	loopEnv := MakeEnvironment()
	loopEnv.Parent = env

	for {
		// Re-evaluate the condition
		conditionVal := Evaluate(node.Condition, loopEnv)
		condition, isBoolCondition := conditionVal.(BoolValue)
		if !isBoolCondition {
			log.Fatalf("For loop condition does not evaluate to a boolean value: %v", condition)
		}

		// Condition is false
		if !condition.Value {
			break
		}

		// Evaluate the body
		Evaluate(node.Body, loopEnv)
	}

	// While loops are statements so they don't resolve to a value
	return Null
}


func evalIfStatement(node ast.IfStatementNode, env *Environment) RuntimeVal {
	// New scope for the for loop body
	ifBlockEnv := MakeEnvironment()
	ifBlockEnv.Parent = env

	condition := Evaluate(node.Condition, ifBlockEnv)

	conditionVal, isConditionBool := condition.(BoolValue)

	if !isConditionBool {
		log.Fatalf("If statement condition must evaluate to a boolean: %v", conditionVal)
	}

	if !conditionVal.Value {
		return Null
	}

	// consequent := Evaluate(node.Consequent, ifBlockEnv)

	// We will only support block statements as consequent
	consequentVal, isConsequentBlock := node.Consequent.(ast.BlockStatementNode)

	if !isConsequentBlock {
		log.Fatalf("If statement must have a block statement as the body %v: ", consequentVal)
	}

	Evaluate(consequentVal, ifBlockEnv)

	// While loops are statements so they don't resolve to a value
	return Null
}

func evalBlockStatement(node ast.BlockStatementNode, env *Environment) RuntimeVal {
	for _, stmt := range node.Body {
		Evaluate(stmt, env)
	}

	return Null
}

func Evaluate(astNode ast.ASTNode, env *Environment) RuntimeVal {
	switch node := astNode.(type) {
	case ast.AssignmentExprNode:
		return evalAssignment(node, env)
	case ast.ObjectLiteralExprNode:
		return evalObjectLiteral(node, env)
	case ast.CallExprNode:
		return evalCallExpression(node, env)
	case ast.Program:
		return evalProgram(node, env)
	case ast.VariableDeclarationNode:
		return evalVarDeclaration(node, env)
	case ast.FunctionDeclarationNode:
		return evalFnDeclaration(node, env)
	case ast.ReturnStatementNode:
		return evalReturnStatement(node, env)
	case ast.NumericLiteralExprNode:
		return evalNumber(node, env)
	case ast.BinaryExprNode:
		return evalBinaryOp(node, env)
	case ast.IdentifierExprNode:
		return evalVarLookup(node, env)
	case ast.ArrayLiteralExprNode:
		return evalArray(node, env)
	case ast.MemberExprNode:
		return evalMember(node, env)
	case ast.StringLiteralExprNode:
		return evalString(node, env)
	case ast.BooleanLiteralExprNode:
		return evalBool(node, env)
	case ast.LogicalExprNode:
		return evalLogicalExpr(node, env)
	case ast.UnaryExprNode:
		return evalUnaryOp(node, env)
	case ast.ForStatementNode:
		return evalForLoop(node, env)
	case ast.WhileStatementNode:
		return evalWhileLoop(node, env)
	case ast.BlockStatementNode:
		return evalBlockStatement(node, env)
	case ast.IfStatementNode:
		return evalIfStatement(node, env)
	default:
		log.Fatalf("Node of type '%s' is not setup for evaluation.", ast.GetNodeKindAsString(node))
	}

	return Null
}

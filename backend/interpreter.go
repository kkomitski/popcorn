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

func evalBinaryOp(node ast.BinaryExprNode, env *Environment) RuntimeVal {
	left := Evaluate(node.Left, env)
	right := Evaluate(node.Right, env)

	leftNum, leftIsNum := left.(NumberVal)
	rightNum, rightIsNum := right.(NumberVal)

	if !leftIsNum || !rightIsNum {
		log.Fatalf("Cannot perform a binary operation on values not of type num. %v, %v", left, right)
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
	default:
		log.Fatalf("Unknown binary operator: %s", node.Operator)
	}

	return Null
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
	default:
		log.Fatalf("Node of type '%s' is not setup for evaluation.", node)
	}

	return Null
}

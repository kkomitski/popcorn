package backend

import (
	"pop/frontend/types/ast"
)

type ValueType = int

const (
	NullType = iota
	NumberType
	BooleanType
	ObjectType
	NativeFunctionType
	FunctionType
	ReturnType
)

type RuntimeVal any

var Null = NullValue{Value: nil}

func GetValType(val RuntimeVal) ValueType {
	switch val.(type) {
	case NullValue, *NullValue:
		return NullType
	case BoolValue, *BoolValue:
		return BooleanType
	case NumberVal, *NumberVal:
		return NumberType
	case ObjectVal, *ObjectVal:
		return ObjectType
	case NativeFunctionVal, *NativeFunctionVal:
		return NativeFunctionType
	case FunctionVal, *FunctionVal:
		return FunctionType
	case ReturnVal, *ReturnVal:
		return ReturnType
	default:
		return -1
	}
}

type NullValue struct {
	Value interface{}
}

type BoolValue struct {
	Value bool
}

type NumberVal struct {
	Value float64
}

type ObjectVal struct {
	Properties map[string]RuntimeVal
}

type FunctionCall func(args []RuntimeVal, env *Environment) RuntimeVal

type NativeFunctionVal struct {
	Call FunctionCall
}

type FunctionVal struct {
	Name           string
	Params         []string
	DeclarationEnv *Environment
	Body           []ast.ASTNode
}

type ReturnVal struct {
	Value RuntimeVal
}

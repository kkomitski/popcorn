package backend

import (
	"log"
)

type Environment struct {
	Parent    *Environment
	Variables map[string]RuntimeVal
	Constants map[string]struct{}
}

func (e *Environment) resolveEnv(varName string) *Environment {
	if _, ok := e.Variables[varName]; ok {
		return e
	}

	if e.Parent == nil {
		log.Fatalf("Cannot resolve variable '%s' !", varName)
	}

	return e.Parent.resolveEnv(varName)
}

func (e *Environment) GetVar(varName string) RuntimeVal {
	env := e.resolveEnv(varName)

	if val, ok := env.Variables[varName]; ok {
		return val
	}

	return Null
}

func (e *Environment) AssignVar(varName string, value RuntimeVal) RuntimeVal {
	env := e.resolveEnv(varName)

	if _, ok := env.Variables[varName]; !ok {
		log.Fatalf("No variable with the '%s' identifier found.", varName)
	}

	if _, ok := env.Constants[varName]; ok {
		log.Fatalf("Cannot reassign constant variable '%s'", varName)
	}

	env.Variables[varName] = value

	return value
}

func (e *Environment) DeclareVar(varName string, isConstant bool, value RuntimeVal) RuntimeVal {
	if _, ok := e.Variables[varName]; ok {
		log.Fatalf("Cannot declare variable '%s' as its already present in the current scope.", varName)
	}

	var val RuntimeVal
	if value == nil {
		val = Null
	} else {
		val = value
	}

	if isConstant {
		if val == Null {
			log.Fatalf("Cannot declare a constant variable '%s' without a value.", varName)
		}

		e.Constants[varName] = struct{}{}
	}

	e.Variables[varName] = val
	return val
}

func MakeEnvironment() *Environment {
	env := &Environment{
		Parent:    nil,
		Variables: make(map[string]RuntimeVal),
		Constants: map[string]struct{}{},
	}

	return env
}

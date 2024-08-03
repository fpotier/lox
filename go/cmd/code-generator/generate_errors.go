package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type Field struct {
	Name string
	Type string
}

type ErrorType struct {
	Name    string
	Fields  []Field
	Message string
}

type Data struct {
	Package   string
	ErrorKind string
	Imports   []string
	Types     []ErrorType
}

var lexingErrors = Data{
	Package:   "lexer",
	ErrorKind: "LexingError",
	Imports:   []string{"fmt"},
	Types: []ErrorType{
		{
			Name: "UnexpectedCharacter",
			Fields: []Field{
				{Name: "character", Type: "byte"},
			},
			Message: "Unexpected character '%c'",
		},
		{
			Name:    "UnterminatedString",
			Fields:  []Field{},
			Message: "Unterminated string literal",
		},
		{
			Name: "InvalidFloat",
			Fields: []Field{
				{Name: "text", Type: "string"},
			},
			Message: "Error converting %s to float",
		},
	},
}

var runtimeErrors = Data{
	Package:   "runtime",
	ErrorKind: "RuntimeError",
	Imports:   []string{"fmt"},
	Types: []ErrorType{
		{
			Name: "UndefinedVariable",
			Fields: []Field{
				{Name: "identifier", Type: "string"},
			},
			Message: "Undefined variable '%s'",
		},
		{
			Name: "UndefinedProperty",
			Fields: []Field{
				{Name: "propertyName", Type: "string"},
				{Name: "className", Type: "string"},
			},
			Message: "Undefined property '%s' for class '%s'",
		},
		{
			Name: "InvalidSuper",
			Fields: []Field{
				{Name: "location", Type: "string"},
			},
			Message: "Can't use 'super' %s",
		},
		{
			Name: "InvalidThis",
			Fields: []Field{
				{Name: "location", Type: "string"},
			},
			Message: "Can't use 'this' %s",
		},
		{
			Name:    "UninitializedRead",
			Fields:  []Field{},
			Message: "Can't read local variable in its own initializer",
		},
		{
			Name: "InvalidInheritance",
			Fields: []Field{
				{Name: "errorMsg", Type: "string"},
			},
			Message: "%s",
		},
		{
			Name: "InvalidReturn",
			Fields: []Field{
				{Name: "location", Type: "string"},
			},
			Message: "Can't return from %s",
		},
		{
			Name: "VariableRedeclaration",
			Fields: []Field{
				{Name: "identifier", Type: "string"},
			},
			Message: "Variable '%s' is already declared in this scope",
		},
		{
			Name: "UnsupportedBinaryOperation",
			Fields: []Field{
				{Name: "operator", Type: "string"},
				{Name: "lhs", Type: "string"},
				{Name: "rhs", Type: "string"},
			},
			Message: "Operator '%s': incompatible types '%v' and '%v'",
		},
		{
			Name: "UnsupportedUnaryOperation",
			Fields: []Field{
				{Name: "operator", Type: "string"},
				{Name: "rhs", Type: "string"},
			},
			Message: "Operator '%s': incompatible type '%v'",
		},
		{
			Name:    "InvalidSetGet",
			Fields:  []Field{},
			Message: "Only class instances have properties",
		},
		{
			Name:    "NotCallable",
			Fields:  []Field{},
			Message: "Only classes and functions are callable",
		},
		{
			Name: "BadArity",
			Fields: []Field{
				{Name: "functionName", Type: "string"},
				{Name: "expected", Type: "int"},
				{Name: "got", Type: "int"},
			},
			Message: "Function '%s' expected %d arguments but got %d",
		},
	},
}

const (
	nbArgsRequired = 2
	filePerm       = 0644
)

func main() {
	if len(os.Args) != nbArgsRequired {
		fmt.Fprintf(os.Stderr, "Missing argument")
		return
	}

	var data Data
	prefix := os.Args[1]
	switch prefix {
	case "lexer":
		data = lexingErrors
	case "runtime":
		data = runtimeErrors
	default:
		fmt.Fprintf(os.Stderr, "Unknown argument %s", os.Args[1])
		return
	}

	tmpl, err := template.ParseFiles("../../cmd/code-generator/errors_template.go.tmpl")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	filename := "generated-errors" + ".go"
	err = os.WriteFile(filename, buf.Bytes(), filePerm)
	if err != nil {
		panic(err)
	}
}

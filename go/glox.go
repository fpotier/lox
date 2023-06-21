package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fpotier/crafting-interpreters/go/ast"
	"github.com/fpotier/crafting-interpreters/go/ast/visitor"
	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/sean-/sysexits"
)

var hadError = false

func main() {
	expr := ast.NewBinaryExpression(
		ast.NewUnaryExpression(
			*lexer.NewToken(lexer.DASH, "-", nil, 1),
			ast.NewLiteralExpression(&lexer.NumberLiteral{Value: 123})),
		*lexer.NewToken(lexer.STAR, "*", nil, 1),
		ast.NewGroupingExpression(ast.NewLiteralExpression(&lexer.NumberLiteral{Value: 45.67})))

	printer := &visitor.LispPrinter{}
	fmt.Println(printer.String(expr))

	expr2 := ast.NewBinaryExpression(
		ast.NewGroupingExpression(ast.NewBinaryExpression(
			ast.NewLiteralExpression(&lexer.NumberLiteral{Value: 1}),
			*lexer.NewToken(lexer.PLUS, "+", nil, 1),
			ast.NewLiteralExpression(&lexer.NumberLiteral{Value: 2}))),
		*lexer.NewToken(lexer.STAR, "*", nil, 1),
		ast.NewGroupingExpression(ast.NewBinaryExpression(
			ast.NewLiteralExpression(&lexer.NumberLiteral{Value: 4}),
			*lexer.NewToken(lexer.DASH, "-", nil, 1),
			ast.NewLiteralExpression(&lexer.NumberLiteral{Value: 3}))))

	rpnPrinter := &visitor.RPNPrinter{}
	fmt.Println(rpnPrinter.String(expr2))

	nbArgs := len(os.Args)
	if nbArgs > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(sysexits.Usage)
	} else if nbArgs == 2 {
		fmt.Printf("Run file %v\n", os.Args[1])
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		if input := scanner.Text(); input == "exit" || len(input) == 0 {
			break
		}
		run(scanner.Text())
		hadError = false
	}
}

func runFile(filepath string) {
	sourceCode, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	run(string(sourceCode))
	if hadError {
		os.Exit(sysexits.DataErr)
	}
}

func run(source_code string) {
	lexer := lexer.NewLexer(source_code)
	tokens := lexer.Tokens()
	parser := ast.NewParser(tokens)
	expr := parser.Parse()

	// TODO error handling

	fmt.Println((&visitor.LispPrinter{}).String(expr))
}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line %v] Error %v : %v\n", line, where, message)
	hadError = true
}

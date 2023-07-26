package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fpotier/crafting-interpreters/go/ast"
	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
	"github.com/sean-/sysexits"
)

// Maybe get rid of this global variable
var interpreter = ast.NewInterpreter()

func RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("lox> ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		if input := scanner.Text(); input == "exit" || len(input) == 0 {
			break
		}
		run(scanner.Text())
		loxerror.HadError = false
	}
}

func RunFile(filepath string) {
	sourceCode, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	run(string(sourceCode))
	if loxerror.HadError {
		os.Exit(sysexits.DataErr)
	}
}

func run(sourceCode string) {
	lexer := lexer.NewLexer(sourceCode)
	tokens := lexer.Tokens()

	parser := ast.NewParser(tokens)
	statements := parser.Parse()
	if loxerror.HadError {
		return
	}

	resolver := ast.NewResolver(interpreter)
	resolver.ResolveProgram(statements)
	if loxerror.HadError {
		return
	}

	interpreter.Eval(statements)
}

func main() {
	const maxArgs = 2
	nbArgs := len(os.Args)
	switch {
	case nbArgs > maxArgs:
		fmt.Println("Usage: glox [script]")
		os.Exit(sysexits.Usage)
	case nbArgs == maxArgs:
		fmt.Printf("Run file %v\n", os.Args[1])
		RunFile(os.Args[1])
	default:
		RunPrompt()
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fpotier/crafting-interpreters/go/ast"
	"github.com/fpotier/crafting-interpreters/go/ast/visitor"
	"github.com/fpotier/crafting-interpreters/go/lexer"
	"github.com/fpotier/crafting-interpreters/go/loxerror"
	"github.com/sean-/sysexits"
)

var interpreter = visitor.NewInterpreter()

func RunPrompt() {
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
	statements, err := parser.Parse()

	if loxerror.HadError || err != nil {
		fmt.Println(err)
	}

	interpreter.Eval(statements)
}

func main() {
	const maxArgs = 2
	if nbArgs := len(os.Args); nbArgs > maxArgs {
		fmt.Println("Usage: glox [script]")
		os.Exit(sysexits.Usage)
	} else if nbArgs == maxArgs {
		fmt.Printf("Run file %v\n", os.Args[1])
		RunFile(os.Args[1])
	} else {
		RunPrompt()
	}
}

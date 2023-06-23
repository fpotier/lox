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

func run(source_code string) {
	lexer := lexer.NewLexer(source_code)
	tokens := lexer.Tokens()
	parser := ast.NewParser(tokens)
	expr, err := parser.Parse()

	if loxerror.HadError || err != nil {
		return
	}

	fmt.Println((&visitor.LispPrinter{}).String(expr))
	fmt.Println((&visitor.RPNPrinter{}).String(expr))
	fmt.Println((&visitor.Interpreter{}).Eval(expr))
}

func main() {
	nbArgs := len(os.Args)
	if nbArgs > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(sysexits.Usage)
	} else if nbArgs == 2 {
		fmt.Printf("Run file %v\n", os.Args[1])
		RunFile(os.Args[1])
	} else {
		RunPrompt()
	}
}

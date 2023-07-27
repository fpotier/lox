package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fpotier/lox/go/ast"
	"github.com/fpotier/lox/go/lexer"
	"github.com/sean-/sysexits"
)

type Lox struct {
	hadError    bool
	lexer       *lexer.Lexer
	parser      *ast.Parser
	resolver    *ast.Resolver
	interpreter *ast.Interpreter
	stdout      io.Writer
}

func NewLox(stdout io.Writer) *Lox {
	lox := Lox{
		hadError:    false,
		lexer:       nil,
		parser:      nil,
		resolver:    nil,
		interpreter: nil,
		stdout:      os.Stdout,
	}
	lox.interpreter = ast.NewInterpreter(stdout)

	return &lox
}

func (l *Lox) RunPrompt() {
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
		l.run(scanner.Text())
		l.hadError = false
	}
}

func (l *Lox) RunFile(filepath string) {
	sourceCode, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	l.run(string(sourceCode))
	if l.hadError {
		os.Exit(sysexits.DataErr)
	}
}

func (l *Lox) run(sourceCode string) {
	// TODO: avoid to recreate all components each time
	l.lexer = lexer.NewLexer(l, sourceCode)
	tokens := l.lexer.Tokens()

	l.parser = ast.NewParser(l, tokens)
	statements := l.parser.Parse()
	if l.hadError {
		return
	}

	l.resolver = ast.NewResolver(l, l.interpreter)
	l.resolver.ResolveProgram(statements)
	if l.hadError {
		return
	}

	l.interpreter.Eval(statements)
}

func (l *Lox) Error(line int, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, message)
	l.hadError = true
}

func main() {
	const maxArgs = 2
	nbArgs := len(os.Args)
	lox := NewLox(os.Stdout)
	switch {
	case nbArgs > maxArgs:
		fmt.Println("Usage: glox [script]")
		os.Exit(sysexits.Usage)
	case nbArgs == maxArgs:
		lox.RunFile(os.Args[1])
	default:
		lox.RunPrompt()
	}
}

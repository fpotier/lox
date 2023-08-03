package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fpotier/lox/go/pkg/lexer"
	"github.com/fpotier/lox/go/pkg/parser"
	"github.com/fpotier/lox/go/pkg/runtime"
	"github.com/sean-/sysexits"
)

type Lox struct {
	hadError    bool
	lexer       *lexer.Lexer
	parser      *parser.Parser
	resolver    *runtime.Resolver
	interpreter *runtime.Interpreter
	stdout      io.Writer
	stderr      io.Writer
}

func NewLox(fds ...io.Writer) *Lox {
	lox := Lox{
		hadError:    false,
		lexer:       nil,
		parser:      nil,
		resolver:    nil,
		interpreter: nil,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
	}
	for i, fd := range fds {
		switch i {
		case 0:
			lox.stdout = fd
		case 1:
			lox.stderr = fd
		}
	}
	lox.interpreter = runtime.NewInterpreter(lox.stdout, &lox)

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

func (l *Lox) RunFile(filepath string) int {
	sourceCode, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	l.run(string(sourceCode))

	if l.hadError {
		return sysexits.DataErr
	}

	return 0
}

func (l *Lox) run(sourceCode string) {
	// TODO: avoid to recreate all components each time
	l.lexer = lexer.NewLexer(l, sourceCode)
	tokens := l.lexer.Tokens()

	l.parser = parser.NewParser(l, tokens)
	statements := l.parser.Parse()
	if l.hadError {
		return
	}

	l.resolver = runtime.NewResolver(l, l.interpreter)
	l.resolver.ResolveProgram(statements)
	if l.hadError {
		return
	}

	l.interpreter.Eval(statements)
}

func (l *Lox) Error(line int, message string) {
	fmt.Fprintf(l.stderr, "[line %d] Error: %s\n", line, message)
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
		os.Exit(lox.RunFile(os.Args[1]))
	default:
		lox.RunPrompt()
	}
}

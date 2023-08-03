package runtime

import (
	"os"
	"time"

	"github.com/fpotier/lox/go/pkg/ast"
)

type NativeFunction struct {
	name  string
	arity int
	code  CallableCode
}

func (f NativeFunction) Kind() ast.Kind             { return ast.NativeFunc }
func (f NativeFunction) IsTruthy() bool             { return true }
func (f NativeFunction) String() string             { return "<native fn>" }
func (f NativeFunction) Equals(_ ast.LoxValue) bool { return false }
func (f NativeFunction) Call(i *Interpreter, arguments []ast.LoxValue) ast.LoxValue {
	return f.code(i, arguments)
}
func (f NativeFunction) Arity() int { return f.arity }

var builtinNativeFunctions = []NativeFunction{
	{
		name:  "clock",
		arity: 0,
		code: func(*Interpreter, []ast.LoxValue) ast.LoxValue {
			return ast.NewNumberValue(float64(time.Now().Unix()))
		},
	},
	{
		name:  "getchar",
		arity: 0,
		code: func(*Interpreter, []ast.LoxValue) ast.LoxValue {
			// FIXME: terminates the REPL mode
			var buffer [1]byte
			bytesRead, err := os.Stdin.Read(buffer[:])
			if bytesRead == 0 || err != nil {
				return ast.NewNumberValue(-1)
			}

			return ast.NewNumberValue(float64(int(buffer[0])))
		},
	},
}

package runtime

import (
	"os"
	"time"

	"github.com/fpotier/lox/go/pkg/value"
)

type NativeFunction struct {
	name  string
	arity int
	code  CallableCode
}

func (*NativeFunction) TypeName() string { return "native function" }
func (f *NativeFunction) String() string { return "<native fn>" }
func (f *NativeFunction) Name() string   { return f.name }
func (f *NativeFunction) Call(i *Interpreter, arguments []value.Value) value.Value {
	return f.code(i, arguments)
}
func (f *NativeFunction) Arity() int { return f.arity }

// Pointers, so that == on Lox values is identity (a struct holding a func is not comparable).
var builtinNativeFunctions = []*NativeFunction{
	{
		name:  "clock",
		arity: 0,
		code: func(*Interpreter, []value.Value) value.Value {
			return value.Number(float64(time.Now().Unix()))
		},
	},
	{
		name:  "getchar",
		arity: 0,
		code: func(*Interpreter, []value.Value) value.Value {
			// FIXME: terminates the REPL mode
			var buffer [1]byte
			bytesRead, err := os.Stdin.Read(buffer[:])
			if bytesRead == 0 || err != nil {
				return value.Number(-1)
			}

			return value.Number(float64(int(buffer[0])))
		},
	},
}

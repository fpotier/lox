package runtime

import (
	"fmt"

	"github.com/fpotier/lox/go/pkg/ast"
)

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		name:       name,
		superclass: superclass,
		methods:    methods,
	}
}

func (c *LoxClass) Kind() ast.Kind             { return ast.Class }
func (c *LoxClass) IsTruthy() bool             { return true }
func (c *LoxClass) String() string             { return c.name }
func (c *LoxClass) Name() string               { return fmt.Sprintf("%s::%s", c.name, c.name) }
func (c *LoxClass) Equals(v ast.LoxValue) bool { return c == v }

func (c *LoxClass) Call(i *Interpreter, arguments []ast.LoxValue) ast.LoxValue {
	instance := NewLoxInstance(*c)
	if constructor, ok := c.findMethod("init"); ok {
		constructor.Bind(instance).Call(i, arguments)
	}

	return instance
}

func (c *LoxClass) Arity() int {
	if constructor, ok := c.findMethod("init"); ok {
		return constructor.Arity()
	}
	return 0
}

func (c *LoxClass) findMethod(name string) (*LoxFunction, bool) {
	if method, ok := c.methods[name]; ok {
		return method, true
	}

	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}

	return nil, false
}

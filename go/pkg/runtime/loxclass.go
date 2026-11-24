package runtime

import (
	"fmt"

	"github.com/fpotier/lox/go/pkg/value"
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

func (*LoxClass) TypeName() string { return "class" }
func (c *LoxClass) String() string { return c.name }
func (c *LoxClass) Name() string   { return fmt.Sprintf("%s::%s", c.name, c.name) }

func (c *LoxClass) Call(i *Interpreter, arguments []value.Value) value.Value {
	instance := NewLoxInstance(c)
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

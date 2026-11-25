package ast

var Limits = struct {
	MaxArgs int
}{
	MaxArgs: 255, //nolint:mnd // defining the limit is the point of this file
}

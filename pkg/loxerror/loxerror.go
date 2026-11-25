package loxerror

type LoxError interface {
	Line() int
	Kind() string
	Message() string
}

package loxerror

type ErrorFormatter interface {
	PushError(loxError LoxError)
	PopError() (LoxError, error)
	Format(loxerror LoxError) string
	HasErrors() bool
	Reset()
}

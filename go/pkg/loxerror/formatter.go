package loxerror

type ErrorFormatter interface {
	PushError(LoxError)
	PopError() (LoxError, error)
	Format(LoxError) string
	HasErrors() bool
	Reset()
}

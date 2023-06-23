package loxerror

type RuntimeError struct {
	Message string
}

func (e RuntimeError) Error() string {
	return e.Message
}

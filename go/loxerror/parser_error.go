package loxerror

type ParserError struct {
	Message string
}

func (e ParserError) Error() string {
	return e.Message
}

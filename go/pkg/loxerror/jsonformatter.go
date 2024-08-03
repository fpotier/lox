package loxerror

import (
	"bytes"
	"encoding/json"
	"errors"
)

type JsonErrorFormatter struct {
	errors []LoxError
}

func NewJsonErrorFormatter() *JsonErrorFormatter {
	return &JsonErrorFormatter{
		errors: make([]LoxError, 0),
	}
}

func (f *JsonErrorFormatter) PushError(e LoxError) {
	f.errors = append(f.errors, e)
}

func (f *JsonErrorFormatter) PopError() (LoxError, error) {
	if !f.HasErrors() {
		return nil, errors.New("empty error stack")
	}

	err := f.errors[0]
	f.errors = f.errors[1:]
	return err, nil
}

func (f *JsonErrorFormatter) HasErrors() bool {
	return len(f.errors) > 0
}

func (f *JsonErrorFormatter) Format(e LoxError) string {
	raw_string, err := MarshalJSON(e)
	if err != nil {
		panic(err)
	}

	return string(raw_string)
}

func (f *JsonErrorFormatter) Errors() []LoxError {
	return f.errors
}

func (f *JsonErrorFormatter) Reset() {
	f.errors = f.errors[:0]
}

func MarshalJSON[T LoxError](e T) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(map[string]any{
		"line":    e.Line(),
		"type":    e.Kind(),
		"message": e.Message(),
	}); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

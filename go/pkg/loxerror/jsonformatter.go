package loxerror

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var ErrEmptyErrorStack = errors.New("empty error stack")

type JSONErrorFormatter struct {
	errors []LoxError
}

func NewJSONErrorFormatter() *JSONErrorFormatter {
	return &JSONErrorFormatter{
		errors: make([]LoxError, 0),
	}
}

func (f *JSONErrorFormatter) PushError(e LoxError) {
	f.errors = append(f.errors, e)
}

func (f *JSONErrorFormatter) PopError() (LoxError, error) {
	if !f.HasErrors() {
		return nil, ErrEmptyErrorStack
	}

	err := f.errors[0]
	f.errors = f.errors[1:]
	return err, nil
}

func (f *JSONErrorFormatter) HasErrors() bool {
	return len(f.errors) > 0
}

func (f *JSONErrorFormatter) Format(e LoxError) string {
	rawString, err := MarshalJSON(e)
	if err != nil {
		panic(err)
	}

	return string(rawString)
}

func (f *JSONErrorFormatter) Errors() []LoxError {
	return f.errors
}

func (f *JSONErrorFormatter) Reset() {
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
		return nil, fmt.Errorf("failed to encode the error message: %w", err)
	}

	return buffer.Bytes(), nil
}

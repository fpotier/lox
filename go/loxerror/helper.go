package loxerror

import "fmt"

var HadError = false

func Error(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Printf("[line %v] Error %v : %v\n", line, where, message)
	HadError = true
}

package main

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pkg/diff"
)

const TestDirectory = "../../../test/official_tests"

var testedDirectories = [...]string{
	".",
	"assignment",
	"block",
	"bool",
	"call",
	"class",
	"closure",
	"comments",
	"constructor",
	"field",
	"for",
	"function",
	"if",
	"inheritance",
	// "limit",
	"logical_operator",
	"method",
	"nil",
	"number",
	"operator",
	"print",
	"regression",
	"return",
	"string",
	"super",
	"this",
	"variable",
	"while",
}

func loxFilesInDir(path string) ([]string, error) {
	return filepath.Glob(filepath.Join(path + "/*.lox"))
}

func TestRunFile(t *testing.T) {
	for _, dir := range testedDirectories {
		absolutePath, err := filepath.Abs(TestDirectory + "/" + dir)
		if err != nil {
			t.Fatal(err.Error())
		}

		t.Run(dir, func(t *testing.T) {
			runFilesInDir(t, absolutePath)
		})
	}

}

func runFilesInDir(t *testing.T, dirPath string) {
	loxFiles, err := loxFilesInDir(dirPath)
	if err != nil || len(loxFiles) == 0 {
		t.Logf("No .lox files in %s", dirPath)
		t.Skip()
	}

	for _, file := range loxFiles {
		t.Run(file, func(t *testing.T) {
			diffReport := strings.Builder{}
			err := checkRunOutput(file, &diffReport)
			if diffReport.Len() > 0 {
				t.Fatal(diffReport.String())
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func checkRunOutput(filename string, diffOutput io.Writer) error {
	var (
		outputPattern = regexp.MustCompile("^.*expect: (.*)$")
		errorPattern  = regexp.MustCompile("^.*error: (.*)$")
		stdoutBuilder = strings.Builder{}
		stderrBuilder = strings.Builder{}
	)

	lox := NewLox(&stdoutBuilder, &stderrBuilder)
	lox.RunFile(filename)
	programOutput := stdoutBuilder.String()
	programErr := stderrBuilder.String()

	expectedBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	expectedStdoutBuilder := strings.Builder{}
	for _, line := range strings.Split(string(expectedBytes), "\n") {
		match := outputPattern.FindStringSubmatch(line)
		if len(match) > 1 {
			expectedStdoutBuilder.WriteString(match[1])
			expectedStdoutBuilder.WriteByte('\n')
		}
	}
	expectedOutput := expectedStdoutBuilder.String()
	if expectedOutput != programOutput {
		err := diff.Text(filename, filename+".expected", programOutput, expectedOutput, diffOutput)
		if err != nil {
			return err
		}
	}

	expectedStderrBuilder := strings.Builder{}
	for _, line := range strings.Split(string(expectedBytes), "\n") {
		match := errorPattern.FindStringSubmatch(line)
		if len(match) > 1 {
			expectedStderrBuilder.WriteString(match[1])
			expectedStderrBuilder.WriteByte('\n')
		}
	}
	expectedError := expectedStderrBuilder.String()
	if expectedError != programErr {
		err := diff.Text(filename, filename+".expected", programErr, expectedError, diffOutput)
		if err != nil {
			return err
		}
	}

	return nil
}

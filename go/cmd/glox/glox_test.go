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
const BenchmarkDirectory = "../../../benchmark/official_benchmarks"

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
	return filepath.Glob(filepath.Join(path, "/*.lox"))
}

func TestRunFile(t *testing.T) {
	t.Parallel()
	for _, dir := range testedDirectories {
		absolutePath, err := filepath.Abs(TestDirectory + "/" + dir)
		if err != nil {
			t.Fatal(err.Error())
		}

		t.Run(dir, func(t *testing.T) {
			t.Parallel()
			runFilesInDir(t, absolutePath)
		})
	}
}

func BenchmarkRunFile(b *testing.B) {
	runFilesInDirBench(b, BenchmarkDirectory)
}

func runFilesInDir(t *testing.T, dirPath string) {
	t.Helper()
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

func runFilesInDirBench(b *testing.B, dirPath string) {
	b.Helper()
	loxFiles, err := loxFilesInDir(dirPath)
	if err != nil || len(loxFiles) == 0 {
		b.Logf("No .lox files in %s", dirPath)
		b.Skip()
	}

	for _, file := range loxFiles {
		b.Run(file, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var (
					stdoutBuilder = strings.Builder{}
					stderrBuilder = strings.Builder{}
				)
				lox := NewLox(&stdoutBuilder, &stderrBuilder)
				lox.RunFile(file)
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

	rawFileContent, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	// Required when files use CRLF or CR instead of LF (Go doesn't convert when reading)
	fileContent := strings.ReplaceAll(string(rawFileContent), "\r", "")
	expectedStdoutBuilder := strings.Builder{}
	for _, line := range strings.Split(fileContent, "\n") {
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
	for _, line := range strings.Split(fileContent, "\n") {
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

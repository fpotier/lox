package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pkg/diff"
)

const TestDirectory = "../../../test/"

func TestRunFile(t *testing.T) {
	t.Parallel()

	err := filepath.WalkDir(TestDirectory, func(path string, _ fs.DirEntry, _ error) error {
		pattern := regexp.MustCompile("^.*expect: (.*)$")

		t.Run(path, func(t *testing.T) {
			if filepath.Ext(path) != ".lox" {
				return
			}

			if filepath.Base(filepath.Dir(path)) == "limit" {
				t.Skip()
			}
			if filepath.Base(filepath.Dir(path)) == "scanning" {
				t.Skip()
			}

			stdoutBuilder := strings.Builder{}
			stderrBuilder := strings.Builder{}
			lox := NewLox(&stdoutBuilder, &stderrBuilder)
			lox.RunFile(path)
			programOutput := stdoutBuilder.String()
			// programErr := stderrBuilder.String()

			expectedBytes, err := os.ReadFile(path)
			if err != nil {
				t.Fatal("Failed to read ", path)
			}
			expectedStdoutBuilder := strings.Builder{}
			for _, line := range strings.Split(string(expectedBytes), "\n") {
				match := pattern.FindStringSubmatch(line)
				if len(match) > 1 {
					expectedStdoutBuilder.WriteString(match[1])
					expectedStdoutBuilder.WriteByte('\n')
				}
			}
			expectedOutput := expectedStdoutBuilder.String()

			if expectedOutput != programOutput {
				err := diff.Text(filepath.Base(path), path+".expected", programOutput, expectedOutput, os.Stdout)
				if err != nil {
					fmt.Println(err)
				}
				t.Fail()
			}
		})

		return nil
	})

	if err != nil {
		t.Fail()
	}
}

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

const TestDirectory = "../test/"

func TestRunFile(t *testing.T) {
	t.Parallel()

	err := filepath.WalkDir(TestDirectory, func(path string, d fs.DirEntry, err error) error {
		pattern := regexp.MustCompile("^.*expect: (.*)$")

		t.Run(path, func(t *testing.T) {
			if filepath.Ext(path) != ".lox" {
				return
			}

			if filepath.Base(filepath.Dir(path)) == "limit" {
				t.Skip()
			}

			builder := strings.Builder{}
			lox := NewLox(&builder, &builder)
			lox.RunFile(path)
			programOutput := builder.String()

			expectedBytes, err := os.ReadFile(path)
			if err != nil {
				t.Fatal("Failed to read ", path)
			}
			builder = strings.Builder{}
			for _, line := range strings.Split(string(expectedBytes), "\n") {
				match := pattern.FindStringSubmatch(line)
				if len(match) > 1 {
					builder.WriteString(match[1])
					builder.WriteByte('\n')
				}
			}
			expectedOutput := builder.String()

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

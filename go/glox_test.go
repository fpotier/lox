package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/diff"
)

const TestDirectory = "../test/"

func TestRunFile(t *testing.T) {
	filepath.WalkDir(TestDirectory, func(path string, d fs.DirEntry, err error) error {
		t.Run(path, func(t *testing.T) {
			if filepath.Ext(path) != ".lox" {
				return
			}

			expectedPath := path + ".expected"

			_, err := os.Stat(expectedPath)
			if err != nil {
				t.FailNow()
			}

			builder := strings.Builder{}
			interpreter.OutputStream = &builder
			RunFile(path)
			programOutput := builder.String()

			expectedBytes, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatal("Failed to read " + expectedPath)
			}
			expectedOutput := string(expectedBytes)

			if string(expectedOutput) != programOutput {
				err := diff.Text(filepath.Base(path), filepath.Base(expectedPath), programOutput, expectedOutput, os.Stdout)
				if err != nil {
					fmt.Println(err)
				}
				t.Fail()
			}
		})

		return nil
	})
}

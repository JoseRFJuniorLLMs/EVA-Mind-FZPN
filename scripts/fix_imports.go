// +build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Fix accidental replacement of "brainstem" -> "cortex/brainstem"
	// Because "brain" is a substring of "brainstem"
	target := "\"eva-mind/internal/brainstem"
	replacement := "\"eva-mind/internal/brainstem"

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		newContent := strings.ReplaceAll(string(content), target, replacement)

		if newContent != string(content) {
			err = os.WriteFile(path, []byte(newContent), info.Mode())
			if err != nil {
				fmt.Printf("Error writing %s: %v\n", path, err)
			} else {
				fmt.Printf("Fixed %s\n", path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Walk error: %v\n", err)
	}
}

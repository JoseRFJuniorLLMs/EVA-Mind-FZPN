// +build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	replacements := map[string]string{
		// Cortex
		"eva-mind/internal/cortex/brain":       "eva-mind/internal/cortex/brain",
		"eva-mind/internal/cortex/personality": "eva-mind/internal/cortex/personality",
		"eva-mind/internal/cortex/lacan":       "eva-mind/internal/cortex/lacan",
		"eva-mind/internal/cortex/transnar":    "eva-mind/internal/cortex/transnar",
		"eva-mind/internal/cortex/medgemma":    "eva-mind/internal/cortex/medgemma",
		"eva-mind/internal/cortex/veracity":    "eva-mind/internal/cortex/veracity",
		"eva-mind/internal/cortex/llm":         "eva-mind/internal/cortex/llm",
		"eva-mind/internal/cortex/gemini":      "eva-mind/internal/cortex/gemini",

		// Hippocampus
		"eva-mind/internal/hippocampus/memory":    "eva-mind/internal/hippocampus/memory",
		"eva-mind/internal/hippocampus/knowledge": "eva-mind/internal/hippocampus/knowledge",
		"eva-mind/internal/hippocampus/stories":   "eva-mind/internal/hippocampus/stories",

		// Senses
		"eva-mind/internal/senses/signaling": "eva-mind/internal/senses/signaling",
		"eva-mind/internal/senses/api":       "eva-mind/internal/senses/api",
		"eva-mind/internal/senses/voice":     "eva-mind/internal/senses/voice",
		"eva-mind/internal/senses/telemetry": "eva-mind/internal/senses/telemetry",

		// Motor
		"eva-mind/internal/motor/scheduler":   "eva-mind/internal/motor/scheduler",
		"eva-mind/internal/motor/workers":     "eva-mind/internal/motor/workers",
		"eva-mind/internal/motor/calendar":    "eva-mind/internal/motor/calendar",
		"eva-mind/internal/motor/drive":       "eva-mind/internal/motor/drive",
		"eva-mind/internal/motor/sheets":      "eva-mind/internal/motor/sheets",
		"eva-mind/internal/motor/docs":        "eva-mind/internal/motor/docs",
		"eva-mind/internal/motor/email":       "eva-mind/internal/motor/email",
		"eva-mind/internal/motor/gmail":       "eva-mind/internal/motor/gmail",
		"eva-mind/internal/motor/spotify":     "eva-mind/internal/motor/spotify",
		"eva-mind/internal/motor/youtube":     "eva-mind/internal/motor/youtube",
		"eva-mind/internal/motor/uber":        "eva-mind/internal/motor/uber",
		"eva-mind/internal/motor/maps":        "eva-mind/internal/motor/maps",
		"eva-mind/internal/motor/whatsapp":    "eva-mind/internal/motor/whatsapp",
		"eva-mind/internal/motor/actions":     "eva-mind/internal/motor/actions",
		"eva-mind/internal/motor/computeruse": "eva-mind/internal/motor/computeruse",
		"eva-mind/internal/motor/googlefit":   "eva-mind/internal/motor/googlefit",

		// Brainstem
		"eva-mind/internal/brainstem/config":         "eva-mind/internal/brainstem/config",
		"eva-mind/internal/brainstem/database":       "eva-mind/internal/brainstem/database",
		"eva-mind/internal/brainstem/infrastructure": "eva-mind/internal/brainstem/infrastructure",
		"eva-mind/internal/brainstem/logger":         "eva-mind/internal/brainstem/logger",
		"eva-mind/internal/brainstem/auth":           "eva-mind/internal/brainstem/auth",
		"eva-mind/internal/brainstem/oauth":          "eva-mind/internal/brainstem/oauth",
		"eva-mind/internal/brainstem/push":           "eva-mind/internal/brainstem/push",
		"eva-mind/internal/brainstem/middleware":     "eva-mind/internal/brainstem/middleware",
	}

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

		// Read file
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", path, err)
			return nil
		}

		newContent := string(content)
		modified := false

		for oldImport, newImport := range replacements {
			// Replace "eva-mind/internal/x" with "eva-mind/internal/region/x"
			// Check for exact quote match to safely replace
			quotedOld := "\"" + oldImport
			quotedNew := "\"" + newImport

			// Also check for subpackages e.g. "eva-mind/internal/brainstem/infrastructure/vector"
			// Code above won't catch it if I don't check prefixes.
			// Better strategy: Replace longest matches first? Or just strict checking.
			// The map keys are full package paths.
			// Let's replace simple string matches but we must be careful.
			// Actually replace `eva-mind/internal/infrastructure` will also replace `eva-mind/internal/infrastructure/vector` correctly
			// IF we just do string replacement.
			// BUT we must do it carefully to not double replace?
			// Since target paths contain source paths?? No.
			// internal/brain -> internal/cortex/brain (contains internal/brain?? No)
			// internal/brainstem/database vs internal/database.
			// internal/database is "old". internal/brainstem/database is "new".

			if strings.Contains(newContent, quotedOld) {
				newContent = strings.ReplaceAll(newContent, quotedOld, quotedNew)
				modified = true
			}

			// Handle subpackages manually roughly?
			// Example: eva-mind/internal/infrastructure/vector
			// My map has "eva-mind/internal/brainstem/infrastructure".
			// Replacing that string works for subpackages too.
			// "eva-mind/internal/brainstem/infrastructure/vector" -> "eva-mind/internal/brainstem/infrastructure/vector"
		}

		if modified {
			err = os.WriteFile(path, []byte(newContent), info.Mode())
			if err != nil {
				fmt.Printf("Error writing %s: %v\n", path, err)
			} else {
				fmt.Printf("Updated %s\n", path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Walk error: %v\n", err)
	}
}

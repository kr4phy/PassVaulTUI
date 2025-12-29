package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// dataFilePath returns the data file path located beside the running binary.
// When running via `go run`, the binary lives in a temporary directory, so we
// instead use the current working directory to keep the file with the sources.
func dataFilePath() string {
	exePath, err := os.Executable()

	if err != nil {
		log.Printf("Warning: Failed to get executable path: %v. Using default path 'data.bin'", err)
		return "data.bin"
	}

	exePath, err = filepath.EvalSymlinks(exePath)

	if err != nil {
		log.Printf("Warning: Failed to evaluate symlinks for %s: %v. Using default path 'data.bin'", exePath, err)
		return "data.bin"
	}

	tmpDir := filepath.Clean(os.TempDir())

	if strings.HasPrefix(exePath, tmpDir) {
		cwd, cwdErr := os.Getwd()

		if cwdErr != nil {
			log.Printf("Warning: Failed to get working directory: %v. Using default path 'data.bin'", cwdErr)
			return "data.bin"
		}

		return filepath.Join(cwd, "data.bin")
	}

	return filepath.Join(filepath.Dir(exePath), "data.bin")
}

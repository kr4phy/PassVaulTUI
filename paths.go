package main

import (
	"log"
	"os"
	"path/filepath"
)

// dataFilePath returns the data file path located beside the running binary.
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
	return filepath.Join(filepath.Dir(exePath), "data.bin")
}

package utils

import (
	"fmt"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"os"
	"path/filepath"
	"runtime"
)

func FindDirectory(dirName string) (string, error) {
	// Get the executable path
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exePath := filepath.Dir(ex)

	// Get the source file directory (for when running with 'go run')
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get source file path")
	}
	srcPath := filepath.Dir(filename)

	// Define possible base directories
	baseDirs := []string{
		exePath,
		filepath.Join(exePath, ".."),
		srcPath,
		filepath.Join(srcPath, ".."),
		filepath.Join(srcPath, "../.."),
		".", // Current working directory
	}

	// Check each possible location
	for _, baseDir := range baseDirs {
		dirPath := filepath.Join(baseDir, dirName)
		if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
			return dirPath, nil
		}
	}

	return "", errorx.New(errorx.ErrCodeNotFound, "directory '%s' not found in any of the searched locations", errorx.ErrInternal)
}

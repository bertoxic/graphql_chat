package config

import (
	"fmt"
	"golang.org/x/sys/windows"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DataBaseINFO  *database
	JWTSecret     string
	Port          string
	UserCache     string
	TemplateCache map[string]*template.Template
}

func getLongPathName(shortPath string) (string, error) {
	longPath := make([]uint16, windows.MAX_LONG_PATH)
	n, err := windows.GetLongPathName(&windows.StringToUTF16(shortPath)[0], &longPath[0], windows.MAX_LONG_PATH)
	if err != nil {
		return "", err
	}
	return windows.UTF16ToString(longPath[:n]), nil
}

//nolint:funlen
func NewConfig(jwtSecret, port string) (*AppConfig, error) {
	ex, err := os.Executable()
	if err != nil {
		//log.Fatalf("unable to get executable path: %v", err)
	}
	exePath := filepath.Dir(ex)

	// Convert short path to long path
	longExePath, err := getLongPathName(exePath)
	if err != nil {
		log.Printf("Warning: Unable to get long path name: %v", err)
		longExePath = exePath
	}

	// Get the source file directory (for when running with 'go run')
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("unable to get source file path")
	}
	srcPath := filepath.Dir(filename)
	// Try to load .env from multiple possible locations
	envPaths := []string{
		filepath.Join(longExePath, ".env"),
		filepath.Join(longExePath, "../.env"),
		filepath.Join(srcPath, ".env"),
		filepath.Join(srcPath, "../.env"),
		filepath.Join(srcPath, "../../.env"),
		".env", // Try current working directory as well
	}

	envLoaded := false
	for _, path := range envPaths {
		absPath, _ := filepath.Abs(path)
		//fmt.Printf("Trying to load .env from: %s\n", absPath)
		err = godotenv.Load(path)
		if err == nil {
			fmt.Printf("Successfully loaded .env from: %s\n", absPath)
			envLoaded = true
			break
		} else {
			//fmt.Printf("Failed to load .env from %s: %v\n", absPath, err)
		}
	}

	if !envLoaded {
		log.Fatal("unable to load .env file from any location 3")
	}

	if err != nil {
		log.Fatal("unable to load config file")
	}
	if err != nil {
		return &AppConfig{}, err
	}
	return &AppConfig{
		DataBaseINFO: &database{URL: os.Getenv("DATABASE_URL")},
		JWTSecret:    jwtSecret,
		Port:         port,
	}, nil
}

type database struct {
	URL string
}

package config

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DataBaseINFO  *database
	RedisDBinfo   *database
	InProduction  bool
	JWT           JWT
	Port          string
	UserCache     string
	TemplateCache map[string]*template.Template
}

type JWT struct {
	Secret []byte
	Issuer string
}

func NewConfig(fileName, port string) (*AppConfig, error) {
	// Check if the environment variables are already set, to skip .env file loading
	if os.Getenv("DATABASE_URL") != "" {
		fmt.Println("Environment variables already set, skipping .env file loading")

		return &AppConfig{
			DataBaseINFO: &database{URL: os.Getenv("DATABASE_URL")},
			RedisDBinfo:  &database{URL: os.Getenv("REDIS_DNS")},
			JWT: JWT{
				Secret: []byte(os.Getenv("JWT_SECRET")),
				Issuer: os.Getenv("ISSUER"),
			},
			Port: port,
		}, nil
	}

	exePath, err := getExecutablePath()
	if err != nil {
		log.Printf("Warning: Unable to get executable path: %v", err)
	}

	// Get the source file directory (for when running with 'go run')
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("unable to get source file path")
	}
	srcPath := filepath.Dir(filename)

	// Trying to load .env from multiple possible locations
	envPaths := []string{
		filepath.Join(exePath, fmt.Sprintf("%s", fileName)),
		filepath.Join(exePath, fmt.Sprintf("../%s", fileName)),
		filepath.Join(srcPath, fmt.Sprintf("%s", fileName)),
		filepath.Join(srcPath, fmt.Sprintf("../%s", fileName)),
		filepath.Join(srcPath, fmt.Sprintf("../../%s", fileName)),
		filepath.Join(exePath, fmt.Sprintf("../../%s", fileName)),
		fmt.Sprintf("%s", fileName),
		fileName, // Trying current working directory as well
	}

	envLoaded := false
	for _, path := range envPaths {
		absPath, _ := filepath.Abs(path)
		err = godotenv.Load(absPath)

		if err == nil {
			fmt.Printf("Successfully loaded .env from: %s\n", absPath)
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		return nil, fmt.Errorf("unable to load %s file from any location", fileName)
	}

	return &AppConfig{
		DataBaseINFO: &database{URL: os.Getenv("DATABASE_URL")},
		JWT: JWT{
			Secret: []byte(os.Getenv("JWT_SECRET")),
			Issuer: os.Getenv("ISSUER"),
		},
		Port: port,
	}, nil
}

func getExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

type database struct {
	URL string
}

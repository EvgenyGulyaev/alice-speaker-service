package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Env          map[string]string
	Port         string
	ServiceToken string
	DBPath       string
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		instance = load()
	})
	return instance
}

func load() *Config {
	startDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	env := map[string]string{}
	if envPath, err := findEnvPath(startDir); err == nil {
		_ = godotenv.Load(envPath)
		if loaded, readErr := godotenv.Read(envPath); readErr == nil {
			env = loaded
		}
	}

	return &Config{
		Env:          env,
		Port:         normalizePort(readEnv("PORT", env, "8090")),
		ServiceToken: readEnv("ALICE_SERVICE_TOKEN", env, ""),
		DBPath:       readEnv("ALICE_DB_PATH", env, "alice.db"),
	}
}

func readEnv(key string, env map[string]string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	if value := strings.TrimSpace(env[key]); value != "" {
		return value
	}
	return fallback
}

func normalizePort(value string) string {
	if value == "" {
		return ":8090"
	}
	if strings.HasPrefix(value, ":") {
		return value
	}
	return ":" + value
}

func findEnvPath(startDir string) (string, error) {
	current := startDir
	for {
		candidate := filepath.Join(current, ".env")
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return candidate, nil
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", os.ErrNotExist
		}
		current = parent
	}
}

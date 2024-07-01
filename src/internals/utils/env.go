package utils

import (
	"log/slog"
	"os"
)

func MustGetEnv(envVarName string) string {
	envVarValue := os.Getenv(envVarName)
	if envVarValue == "" {
		slog.Error("Missing Environment Variable", "environmentVariable", envVarName)
		os.Exit(1)
	}

	return envVarValue
}


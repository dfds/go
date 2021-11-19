package config

import (
	"os"
	"strings"
)

func GetEnvValue(key string, def string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return def
	}

	return val
}

func GetEnvBool(key string, def bool) bool {
	val := os.Getenv(key)

	if len(val) == 0 {
		return def
	}

	if strings.Compare("false", strings.ToLower(val)) == 0 {
		return false
	}

	if strings.Compare("true", strings.ToLower(val)) == 0 {
		return true
	}

	return def
}


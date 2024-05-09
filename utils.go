package loggy

import (
	"os"
)

func EnvOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func Map[From any, To any](mapper func(From) To, arr ...From) []To {
	var result []To
	for _, v := range arr {
		result = append(result, mapper(v))
	}
	return result
}

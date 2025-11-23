package utils

import (
	"os"
)

func Getenv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func MustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

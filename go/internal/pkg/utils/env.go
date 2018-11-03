package utils

import (
	"os"

	"github.com/renwick/go/internal/pkg/log"
)

func MustGetEnv(envKey string) string {
	val := os.Getenv(envKey)
	if val == "" {
		log.FatalF("utils.MustGetEnv - Key %s was empty, erroring out", envKey)
	}
	return val
}

package util

import (
	"log"
	"os"
)

func MustGetEnvVar(v string) string {
	var res string

	if val, ok := os.LookupEnv(v); !ok {
		log.Fatalf("%s must be set", v)
	} else {
		res = val
	}

	return res
}

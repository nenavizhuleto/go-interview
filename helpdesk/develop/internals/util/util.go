package util

import (
	"log"
	"math/rand"
	"os"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func MustGetEnvVar(v string) string {
	var res string

	if val, ok := os.LookupEnv(v); !ok {
		log.Fatalf("%s must be set", v)
	} else {
		res = val
	}

	return res
}

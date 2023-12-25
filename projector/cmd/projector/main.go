package main

import (
	application "projector/internals/applications/projector"
)

func main() {
	if _, err := application.New().Run(); err != nil {
		panic(err)
	}
}

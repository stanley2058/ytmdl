package main

import (
	"os"
)

var downloadPath string

func main() {
	homePath := os.Getenv("HOME")
	downloadPath = homePath + "/Downloads"

	createTui()
}

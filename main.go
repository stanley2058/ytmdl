package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var downloadPath string

func main() {
	log.SetOutput(io.Discard)

	downloadPath = os.Getenv("XDG_DOWNLOAD_DIR")
	if downloadPath == "" {
		homePath, err := os.UserHomeDir()
		if err == nil {
			downloadPath = fmt.Sprintf("%s%cDownloads", homePath, os.PathSeparator)
		}
	}

	createTui()
}

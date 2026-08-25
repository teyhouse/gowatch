package main

import (
	"fmt"
	"os"

	"github.com/teyhouse/gowatch/watcher"
)

const version = "1.2.2"

func main() {
	fmt.Printf("📁 GOWATCH - Version %s\n", version)
	if err := watcher.Watch(); err != nil {
		fmt.Fprintf(os.Stderr, "GOWATCH failed: %s\n", err)
		os.Exit(1)
	}
}

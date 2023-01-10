package main

import (
	"fmt"

	"github.com/teyhouse/gowatch/watcher"
)

const version = "1.1"

func main() {
	fmt.Printf("📁 GOWATCH - Version %s\n", version)
	watcher.Watch()
}

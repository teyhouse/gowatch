//go:build !windows

package logger

import (
	"log"
	"log/slog"
	"log/syslog"
	"sync"
)

var once sync.Once

func setup() {
	// Setup Logger: change from STDOUT to SYSLOG (once, not per call —
	// every syslog.New opens a connection that would otherwise leak)
	once.Do(func() {
		logwriter, err := syslog.New(syslog.LOG_NOTICE, "[GOWATCH]")
		if err != nil {
			slog.Warn("syslog unavailable, logging to stderr", "err", err)
			return
		}
		log.SetOutput(logwriter)
	})
}

func Log(message string) {
	setup()
	log.Println(message)
}

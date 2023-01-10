//go:build !windows
// +build !windows

package logger

import (
	"log"
	"log/syslog"
)

func setup() {
	//Setup Logger: change from STDOUT to SYSLOG
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "[GOWATCH]")
	if e == nil {
		log.SetOutput(logwriter)
	}
}

func Log(message string) {
	setup()
	log.Println(message)
}

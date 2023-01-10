//go:build windows

package logger

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func amAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func Log(message string) {
	if !amAdmin() {
		runMeElevated()
	}
	time.Sleep(1 * time.Second)

	command := "EventCreate"
	args := []string{"/T", "INFORMATION", "/ID", "777", "/L", "APPLICATION",
		"/SO", "Go", "/D", message}
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

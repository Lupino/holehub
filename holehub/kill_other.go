// +build !windows

package main

import "syscall"

func killApp(pid int) {
	syscall.Kill(pid, syscall.SIGINT)
}

// +build windows

package main

import "fmt"

func killApp(pid int) {
	fmt.Printf("Warning: kill not support on windows.\n")
}

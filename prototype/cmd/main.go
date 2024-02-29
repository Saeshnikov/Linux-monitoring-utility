package main

import (
	"abcd/internal/config"
	"fmt"
)

func main() {
	var bpftrace_time int
	var syscalls []string

	bpftrace_time, syscalls = config.ConfigRead()
	fmt.Print(bpftrace_time)
	fmt.Print(syscalls)
}

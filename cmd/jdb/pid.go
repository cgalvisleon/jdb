package main

import (
	"fmt"
	"os"
)

const pidFile = "./tmp/myservice.pid"

func writePID() error {
	f, err := os.Create(pidFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%d", os.Getpid()))
	return err
}

func readPID() (int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}
	var pid int
	fmt.Sscanf(string(data), "%d", &pid)
	return pid, nil
}

package utils

import (
	"fmt"

	"github.com/mitchellh/go-ps"
)

func FindProcessByName(processName string) (int, error) {
	processes, err := ps.Processes()
	if err != nil {
		return -1, err
	}

	for _, process := range processes {
		if process.Executable() == processName {
			return process.Pid(), nil
		}
	}

	return -1, fmt.Errorf("process not found")
}

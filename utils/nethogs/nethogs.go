package nethogs

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"
)

// PID USER     PROGRAM  DEV         SENT      RECEIVED
type NethogsIOInfo struct {
	PID  int
	SENT float64
	RECV float64
}

func GetNetworkIOInfo(pid int, netIOInfo chan NethogsIOInfo, errs chan error) {
	cmd := exec.Command("nethogs", "-t", "eth0")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errs <- err

		return
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		errs <- err

		return
	}

	r := bufio.NewReader(stdout)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			errs <- err

			return
		}

		info := strings.Fields(line)
		if len(info) != 3 {
			continue
		}

		processData := strings.Split(info[0], "/")

		if len(processData) < 2 {
			continue
		}

		processID, err := strconv.Atoi(processData[len(processData)-2])
		if err != nil {
			continue
		}

		if processID != pid {
			continue
		}

		recv, err := strconv.ParseFloat(info[1], 64)
		if err != nil {
			continue
		}

		sent, err := strconv.ParseFloat(info[2], 64)
		if err != nil {
			continue
		}

		netIOInfo <- NethogsIOInfo{
			PID:  pid,
			SENT: sent,
			RECV: recv,
		}
	}
}

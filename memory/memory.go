package memory

import (
	"strconv"
	"time"

	"profiling/utils"
	"profiling/utils/pidstat"

	"github.com/shirou/gopsutil/process"
)

type Memory struct {
	pid              int
	memoryMonitoring chan pidstat.PidstatMemoryInfo
	memoryErrors     chan error
}

const (
	proccessName              = "publisher"
	duration                  = 2 * time.Minute
	interval                  = 1 * time.Second
	intervalPidstatMonitoring = 1
)

func New() (*Memory, error) {
	pid, err := utils.FindProcessByName(proccessName)
	if err != nil {
		return nil, err
	}

	return &Memory{
		pid:              pid,
		memoryMonitoring: make(chan pidstat.PidstatMemoryInfo, 50),
		memoryErrors:     make(chan error, 1),
	}, nil
}

func (c *Memory) GetPidstatMonitoringChannel() chan pidstat.PidstatMemoryInfo {
	return c.memoryMonitoring
}

func (c *Memory) GetPidstatMonitoringErrsChannel() chan error {
	return c.memoryErrors
}

func (c *Memory) RunMemoryMonitoring() {
	pidstat.GetPidstatMemoryMonitoring(
		strconv.Itoa(c.pid), strconv.Itoa(intervalPidstatMonitoring), c.memoryMonitoring, c.memoryErrors)
}

func (c *Memory) CalculateAverageMemoryUsage() (float32, error) {
	process, err := process.NewProcess(int32(c.pid))
	if err != nil {
		return 0, err
	}

	var totalMemoryUsage float32

	iteration := int(duration / interval)

	for i := 0; i < iteration; i++ {
		cpuUsage, err := process.MemoryPercent()
		if err != nil {
			return 0, err
		}

		totalMemoryUsage += cpuUsage

		time.Sleep(interval)
	}

	return totalMemoryUsage / float32(iteration), nil
}

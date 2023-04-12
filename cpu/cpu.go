package cpu

import (
	"strconv"
	"time"

	"profiling/utils"
	"profiling/utils/pidstat"

	"github.com/shirou/gopsutil/process"
)

type CPU struct {
	pid               int
	cpuMonitoring     chan pidstat.PidstatCPUInfo
	cpuMonitoringErrs chan error
}

const (
	proccessName              = "publisher"
	duration                  = 2 * time.Minute
	interval                  = 1 * time.Second
	intervalPidstatMonitoring = 1
)

func New() (*CPU, error) {
	pid, err := utils.FindProcessByName(proccessName)
	if err != nil {
		return nil, err
	}

	return &CPU{
		pid:               pid,
		cpuMonitoring:     make(chan pidstat.PidstatCPUInfo, 50),
		cpuMonitoringErrs: make(chan error, 1),
	}, nil
}

func (c *CPU) GetPidstatMonitoringChannel() chan pidstat.PidstatCPUInfo {
	return c.cpuMonitoring
}

func (c *CPU) GetPidstatMonitoringErrsChannel() chan error {
	return c.cpuMonitoringErrs
}

func (c *CPU) CalculateAverageCPUUsage() (float64, error) {
	process, err := process.NewProcess(int32(c.pid))
	if err != nil {
		return 0, err
	}

	var totalCPUUsage float64

	iteration := int(duration / interval)

	for i := 0; i < iteration; i++ {
		cpuUsage, err := process.CPUPercent()
		if err != nil {
			return 0, err
		}

		totalCPUUsage += cpuUsage

		time.Sleep(interval)
	}

	return totalCPUUsage / float64(iteration), nil
}

func (c *CPU) RunCPUMonitoring() {
	pidstat.GetPidstatCPUMonitoring(
		strconv.Itoa(c.pid), strconv.Itoa(intervalPidstatMonitoring), c.cpuMonitoring, c.cpuMonitoringErrs)
}

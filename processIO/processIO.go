package processio

import (
	"fmt"

	"profiling/utils"
	"profiling/utils/nethogs"

	"github.com/shirou/gopsutil/process"
)

type ProcessIO struct {
	pid               int
	networkMonitoring chan nethogs.NethogsIOInfo
	networkErrors     chan error
}

type ProcessIOInfo struct {
	ReadCount  uint64
	WriteCount uint64
	ReadMb     float64
	WriteMb    float64
}

const (
	proccessName = "publisher"
)

func New() (*ProcessIO, error) {
	pid, err := utils.FindProcessByName(proccessName)
	if err != nil {
		return nil, err
	}

	return &ProcessIO{
		pid:               pid,
		networkMonitoring: make(chan nethogs.NethogsIOInfo, 50),
		networkErrors:     make(chan error, 1),
	}, nil
}

func (c *ProcessIO) GetProcessIOInfo() (*ProcessIOInfo, error) {
	proc, err := process.NewProcess(int32(c.pid))
	if err != nil {
		return nil, err
	}

	ioCounters, err := proc.IOCounters()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return &ProcessIOInfo{
		ReadCount:  ioCounters.ReadCount,
		WriteCount: ioCounters.WriteCount,
		ReadMb:     float64(ioCounters.ReadBytes) / 1024 / 1024,
		WriteMb:    float64(ioCounters.WriteBytes) / 1024 / 1024,
	}, nil
}

func (c *ProcessIO) GetNetworkMonitoringChannel() chan nethogs.NethogsIOInfo {
	return c.networkMonitoring
}

func (c *ProcessIO) GetNetworkMonitoringErrsChannel() chan error {
	return c.networkErrors
}

func (c *ProcessIO) RunNetworkIOMonitoring() {
	nethogs.GetNetworkIOInfo(c.pid, c.networkMonitoring, c.networkErrors)
}

package main

import (
	"os"
	"sync"
	"time"

	"profiling/cpu"
	"profiling/memory"
	processio "profiling/processIO"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02 15:04:05.000",
		FullTimestamp:    true,
	})
	log.SetOutput(os.Stdout)

	log.Info("Start profiling")

	cpu, err := cpu.New()
	if err != nil {
		log.Error(err)

		return
	}

	memory, err := memory.New()
	if err != nil {
		log.Error(err)

		return
	}

	processIO, err := processio.New()
	if err != nil {
		log.Error(err)

		return
	}

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		cpuUsage, err := cpu.CalculateAverageCPUUsage()
		if err != nil {
			log.Error(err)
		}

		log.Info("===============CPU average usage:=================")
		log.Infof("%f", cpuUsage)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		memoryUsage, err := memory.CalculateAverageMemoryUsage()
		if err != nil {
			log.Error(err)
		}

		log.Info("===============Memory average usage:=================")
		log.Infof("%f", memoryUsage)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			processIOInfo, err := processIO.GetProcessIOInfo()
			if err != nil {
				log.Error(err)

				return
			}

			log.Info("===============Process IO memory  Info:=================")
			log.Infof(
				"read count %d, write count %d, read MB %f, write MB %f",
				processIOInfo.ReadCount, processIOInfo.WriteCount, processIOInfo.ReadMb, processIOInfo.WriteMb)

			time.Sleep(1 * time.Second)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case cpuInfo := <-cpu.GetPidstatMonitoringChannel():
				log.Info("===============CPU Info:=================")
				log.Infof(
					"%%usr %f, %%system %f, %%guest %f, %%wait %f, %%CPU %f, CPU %d",
					cpuInfo.Usr, cpuInfo.System, cpuInfo.Guest, cpuInfo.Wait, cpuInfo.CPU, cpuInfo.CPUID)

			case err := <-cpu.GetPidstatMonitoringErrsChannel():
				log.Error(err)

			case memoryInfo := <-memory.GetPidstatMonitoringChannel():
				log.Info("===============Memory Info:=================")
				log.Infof(
					"minflt/s %f, majflt/s %f, VSZ %d, RSS %d, %%MEM %f",
					memoryInfo.Minflt, memoryInfo.Majflt, memoryInfo.VSZ, memoryInfo.RSS, memoryInfo.MEM)

			case err := <-memory.GetPidstatMonitoringErrsChannel():
				log.Error(err)

			case networkIOInfo := <-processIO.GetNetworkMonitoringChannel():
				log.Info("===============Network Info:=================")
				log.Infof(
					"PID %d, SENT %f, RECV %f",
					networkIOInfo.PID, networkIOInfo.SENT, networkIOInfo.RECV)

			case err := <-processIO.GetNetworkMonitoringErrsChannel():
				log.Error(err)
			}
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		cpu.RunCPUMonitoring()
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		memory.RunMemoryMonitoring()
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		processIO.RunNetworkIOMonitoring()
	}()

	wg.Wait()
}

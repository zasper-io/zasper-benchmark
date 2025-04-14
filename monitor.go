package main

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/process"
)

func monitorProcessByPID(pid int32) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		log.Printf("Failed to get process with PID %d: %v", pid, err)
		return
	}

	for {
		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			log.Printf("Error getting CPU usage for PID %d: %v", pid, err)
			continue
		}

		memInfo, err := proc.MemoryInfo()
		if err != nil {
			log.Printf("Error getting memory info for PID %d: %v", pid, err)
			continue
		}

		benchmarkMutex.Lock()
		benchmarkResults = append(benchmarkResults, BenchmarkData{
			Timestamp:     time.Now().Format(time.RFC3339),
			CPUUsage:      cpuPercent,
			MemoryUsageMB: float64(memInfo.RSS) / (1024 * 1024),
		})
		benchmarkMutex.Unlock()

		time.Sleep(5 * time.Second)
	}
}

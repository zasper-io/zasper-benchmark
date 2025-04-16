package monitor

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"
	"github.com/zasper-io/bechmark/models"
)

var (
	benchmarkMutex sync.Mutex

	benchmarkResults []models.BenchmarkData
)

func InitializeBenchmarkResults() {
	benchmarkResults = make([]models.BenchmarkData, 0)
}

func MonitorProcessByPID(pid int32) {
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
		benchmarkResults = append(benchmarkResults, models.BenchmarkData{
			Timestamp:     time.Now().Format(time.RFC3339),
			CPUUsage:      cpuPercent,
			MemoryUsageMB: float64(memInfo.RSS) / (1024 * 1024),
		})
		benchmarkMutex.Unlock()

		time.Sleep(5 * time.Second)
	}
}

func WriteBenchmarkResultsPeriodically(filename string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C

		benchmarkMutex.Lock()
		dataCopy := make([]models.BenchmarkData, len(benchmarkResults))
		copy(dataCopy, benchmarkResults)
		benchmarkMutex.Unlock()

		file, err := os.Create(filename)
		if err != nil {
			log.Printf("Failed to create benchmark file: %v", err)
			continue
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(dataCopy); err != nil {
			log.Printf("Failed to write benchmark data: %v", err)
		} else {
			log.Printf("Benchmark data written to %s", filename)
		}
	}
}

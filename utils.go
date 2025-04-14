package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func writeBenchmarkResultsPeriodically(filename string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C

		benchmarkMutex.Lock()
		dataCopy := make([]BenchmarkData, len(benchmarkResults))
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

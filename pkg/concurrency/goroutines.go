package concurrency

import (
	"sync"
	"time"
)

const numGoroutines = 20000

// GoroutineOverheadBenchmark измеряет накладные расходы на создание горутин.
func GoroutineOverheadBenchmark() (float64, string) {
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	start := time.Now()
	for i := 0; i < numGoroutines; i++ {
		go func() {
			wg.Done()
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	// Возвращаем среднее время на горутину в наносекундах
	return float64(elapsed.Nanoseconds()) / float64(numGoroutines), "ns/op"
}

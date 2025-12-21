package concurrency

import (
	"sync"
	"sync/atomic"
	"time"
)

const syncOps = 2000000
const syncGoRoutines = 4

// MutexBenchmark тестирует производительность с использованием sync.Mutex.
func MutexBenchmark() (float64, string) {
	var counter int64
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(syncGoRoutines)

	start := time.Now()
	for i := 0; i < syncGoRoutines; i++ {
		go func() {
			for j := 0; j < syncOps/syncGoRoutines; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return float64(time.Since(start).Nanoseconds()) / float64(syncOps), "ns/op"
}

// AtomicBenchmark тестирует производительность с использованием sync/atomic.
func AtomicBenchmark() (float64, string) {
	var counter int64
	var wg sync.WaitGroup
	wg.Add(syncGoRoutines)

	start := time.Now()
	for i := 0; i < syncGoRoutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < syncOps/syncGoRoutines; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}
	wg.Wait()
	return float64(time.Since(start).Nanoseconds()) / float64(syncOps), "ns/op"
}

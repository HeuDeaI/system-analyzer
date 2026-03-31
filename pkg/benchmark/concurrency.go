package benchmark

import (
	"math"
	"sync"
	"time"
)

// ConcurrencyBenchmark measures multi-threaded Pi calculation performance.
func ConcurrencyBenchmark() (float64, string) {
	numPoints := 10000000
	numThreads := 8
	pointsPerThread := numPoints / numThreads

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			acc := 0.0
			for j := 0; j < pointsPerThread; j++ {
				x := float64(j) / float64(pointsPerThread)
				acc += math.Sqrt(1 - x*x)
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start).Seconds()

	return float64(numPoints) / elapsed / 1e6, "M-ops/s"
}

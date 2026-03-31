package benchmark

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// SortingBenchmark measures QuickSort performance on random integers.
func SortingBenchmark() (float64, string) {
	size := 1000000
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Int()
	}

	start := time.Now()
	sort.Ints(data)
	elapsed := time.Since(start).Seconds()

	return 1.0 / elapsed, "M-elements/s" // Normalized to 1M elements
}

// StringManipulationBenchmark measures string processing performance.
func StringManipulationBenchmark() (float64, string) {
	base := "The quick brown fox jumps over the lazy dog. "
	longStr := strings.Repeat(base, 10000) // ~450KB

	start := time.Now()
	for i := 0; i < 100; i++ {
		_ = strings.ReplaceAll(longStr, "fox", "cat")
		_ = strings.Split(longStr, " ")
		_ = strings.Contains(longStr, "lazy")
	}
	elapsed := time.Since(start).Seconds()

	return 100.0 / elapsed, "iterations/s"
}

package benchmark

import (
	"fmt"
	"math"
)

// TestResult содержит статистику по результатам тестов.
type Result struct {
	Min   float64
	Max   float64
	Avg   float64
	Unit  string
	Score int
}

func (r Result) String() string {
	return fmt.Sprintf("Результат: %.2f %s (мин: %.2f, макс: %.2f)", r.Avg, r.Unit, r.Min, r.Max)
}

// BenchmarkFunc представляет собой функцию, выполняющую бенчмарк и возвращающую результат и единицу измерения.
type BenchmarkFunc func() (float64, string)

// TestRunner выполняет бенчмарк несколько раз и собирает статистику.
func TestRunner(fn BenchmarkFunc, iterations int, progress chan<- float64) Result {
	results := make([]float64, iterations)
	var unit string
	var min, max, sum float64
	min = math.MaxFloat64
	max = 0.0

	for i := 0; i < iterations; i++ {
		val, u := fn()
		if i == 0 { // Захватываем единицу измерения только на первой итерации
			unit = u
		}
		results[i] = val
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
		sum += val
		progress <- float64(i+1) / float64(iterations)
	}

	return Result{
		Min:  min,
		Max:  max,
		Avg:  sum / float64(iterations),
		Unit: unit,
	}
}

package benchmark

import (
	"fmt"
)

// TestResult содержит статистику по результатам тестов.
type Result struct {
	Min   float64
	Max   float64
	Avg   float64
	Unit  string
	Score int
}

// BenchmarkStats содержит статистику по результатам выполнения бенчмарка.
type BenchmarkStats struct {
	Min  float64 // Минимальное значение
	Max  float64 // Максимальное значение
	Avg  float64 // Среднее значение
	Unit string  // Единица измерения
}

func (r Result) String() string {
	return fmt.Sprintf("Результат: %.2f %s (мин: %.2f, макс: %.2f)", r.Avg, r.Unit, r.Min, r.Max)
}

// BenchmarkFunc представляет собой функцию, выполняющую бенчмарк и возвращающую результат и единицу измерения.
// BenchmarkFunc определяет тип для функции бенчмарка, которая возвращает результат и единицу измерения.
type BenchmarkFunc func() (float64, string)

// TestRunner выполняет заданную функцию бенчмарка указанное количество раз и возвращает статистику.
// Он также может сообщать о прогрессе через опциональный канал.
func TestRunner(fn BenchmarkFunc, iterations int, progress chan<- float64) BenchmarkStats {
	results := make([]float64, iterations)
	var unit string

	// Выполняем тест iterations раз.
	for i := 0; i < iterations; i++ {
		val, u := fn()
		results[i] = val
		unit = u // Сохраняем единицу измерения (предполагается, что она одинакова для всех запусков).

		// Отправляем прогресс, если канал предоставлен.
		if progress != nil {
			progress <- float64(i+1) / float64(iterations)
		}
	}

	var sum, min, max float64
	min = results[0]
	max = results[0]

	// Агрегируем результаты.
	for _, v := range results {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	// Возвращаем собранную статистику.
	return BenchmarkStats{
		Min:  min,
		Max:  max,
		Avg:  sum / float64(iterations),
		Unit: unit,
	}
}

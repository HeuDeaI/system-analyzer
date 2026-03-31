package benchmark

import (
	"fmt"
	"math"
	"sort"
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
	Min      float64 // Минимальное значение
	Max      float64 // Максимальное значение
	Avg      float64 // Среднее значение
	Median   float64 // Медиана (устойчива к выбросам)
	Unit     string  // Единица измерения
	StdDev   float64 // Среднеквадратичное отклонение
	MAD      float64 // Median Absolute Deviation
	Variance float64 // Коэффициент вариации в процентах
}

func (r Result) String() string {
	return fmt.Sprintf("Результат: %.2f %s (мин: %.2f, макс: %.2f)", r.Avg, r.Unit, r.Min, r.Max)
}

// BenchmarkFunc представляет собой функцию, выполняющую бенчмарк и возвращающую результат и единицу измерения.
// BenchmarkFunc определяет тип для функции бенчмарка, которая возвращает результат и единицу измерения.
type BenchmarkFunc func() (float64, string)

// TestRunner выполняет заданную функцию бенчмарка указанное количество раз и возвращает статистику.
// Он также может сообщать о прогрессе через опциональный канал.
// Для повышения точности использует математические методы: медиану, MAD, IQR фильтрацию.
func TestRunner(fn BenchmarkFunc, iterations int, progress chan<- float64) BenchmarkStats {
	// Увеличиваем количество итераций для лучшей статистики
	if iterations < 10 {
		iterations = 10
	}

	// Разогрев - 2 запуска для инициализации кэшей и стабилизации CPU
	_, _ = fn()
	_, _ = fn()

	results := make([]float64, iterations)
	var unit string

	// Выполняем тест iterations раз.
	for i := 0; i < iterations; i++ {
		val, u := fn()
		results[i] = val
		unit = u

		if progress != nil {
			progress <- float64(i+1) / float64(iterations)
		}
	}

	// Сортируем результаты для статистического анализа
	sort.Float64s(results)

	// Вычисляем медиану (более устойчива к выбросам чем среднее)
	median := calculateMedian(results)

	// Median Absolute Deviation (MAD) - робастная оценка разброса
	mad := calculateMAD(results, median)

	// IQR фильтрация - удаляем выбросы за пределами 1.5*IQR
	filtered := filterOutliersIQR(results)

	// Вычисляем статистики на отфильтрованных данных
	n := float64(len(filtered))
	var sum, min, max float64
	min = filtered[0]
	max = filtered[len(filtered)-1]

	for _, v := range filtered {
		sum += v
	}
	avg := sum / n

	// Вычисляем медиану отфильтрованных данных
	filteredMedian := calculateMedian(filtered)

	// Вычисляем стандартное отклонение (на отфильтрованных данных)
	var varianceSum float64
	for _, v := range filtered {
		diff := v - avg
		varianceSum += diff * diff
	}
	stdDev := math.Sqrt(varianceSum / n)
	cv := (stdDev / avg) * 100 // Коэффициент вариации

	// Используем взвешенное среднее медианы и среднего для финального результата
	// Медиана имеет больший вес так как устойчива к выбросам
	finalResult := 0.6*filteredMedian + 0.4*avg

	return BenchmarkStats{
		Min:      min,
		Max:      max,
		Avg:      finalResult,
		Median:   filteredMedian,
		Unit:     unit,
		StdDev:   stdDev,
		MAD:      mad,
		Variance: cv,
	}
}

// calculateMedian вычисляет медиану отсортированного среза
func calculateMedian(data []float64) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	if n%2 == 0 {
		return (data[n/2-1] + data[n/2]) / 2
	}
	return data[n/2]
}

// calculateMAD вычисляет Median Absolute Deviation
func calculateMAD(data []float64, median float64) float64 {
	absDeviations := make([]float64, len(data))
	for i, v := range data {
		absDeviations[i] = math.Abs(v - median)
	}
	sort.Float64s(absDeviations)
	return calculateMedian(absDeviations)
}

// filterOutliersIQR удаляет выбросы используя Interquartile Range метод
func filterOutliersIQR(data []float64) []float64 {
	n := len(data)
	if n < 4 {
		return data
	}

	// Вычисляем квартили
	q1 := calculatePercentile(data, 25)
	q3 := calculatePercentile(data, 75)
	iqr := q3 - q1

	// Границы для определения выбросов
	lowerBound := q1 - 1.5*iqr
	upperBound := q3 + 1.5*iqr

	// Фильтруем данные
	var filtered []float64
	for _, v := range data {
		if v >= lowerBound && v <= upperBound {
			filtered = append(filtered, v)
		}
	}

	// Если все данные оказались выбросами (редкий случай), возвращаем исходные
	if len(filtered) == 0 {
		return data
	}

	return filtered
}

// calculatePercentile вычисляет указанный перцентиль
func calculatePercentile(data []float64, percentile float64) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	index := (percentile / 100) * float64(n-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))
	if lower == upper {
		return data[lower]
	}
	weight := index - float64(lower)
	return data[lower]*(1-weight) + data[upper]*weight
}

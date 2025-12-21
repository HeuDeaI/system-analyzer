package benchmark

import (
	"math/rand"
	"time"
)

const ops = 20000000 // 20 миллионов операций

// IntegerBenchmark выполняет тест целочисленной арифметики.
func IntegerBenchmark() (float64, string) {
	start := time.Now()
	var a, b, c, d int64 = 1, 2, 3, 4
	for i := 0; i < ops; i++ {
		a = a + b
		b = c - a
		c = d * b
		d = c / 2
	}
	elapsed := time.Since(start)
	return float64(ops) / elapsed.Seconds() / 1e9, "млрд оп/с"
}

// FloatBenchmark выполняет тест арифметики с плавающей запятой.
func FloatBenchmark() (float64, string) {
	start := time.Now()
	var a, b, c, d float64 = 1.1, 2.2, 3.3, 4.4
	for i := 0; i < ops; i++ {
		a = a + b
		b = c - a
		c = d * b
		d = c / 2.0
	}
	elapsed := time.Since(start)
	return float64(ops) / elapsed.Seconds() / 1e9, "млрд оп/с"
}

// ReadBandwidthBenchmark выполняет тест пропускной способности чтения памяти.
func ReadBandwidthBenchmark() (float64, string) {
	const bufferSize = 16 * 1024 * 1024 // 16MB
	data := make([]byte, bufferSize)
	for i := range data {
		data[i] = byte(i)
	}
	start := time.Now()
	var sum byte
	for i := 0; i < 10; i++ {
		for _, v := range data {
			sum += v
		}
	}
	elapsed := time.Since(start)
	return float64(bufferSize*10) / elapsed.Seconds() / (1024 * 1024 * 1024), "ГБ/с"
}

// WriteBandwidthBenchmark выполняет тест пропускной способности записи памяти.
func WriteBandwidthBenchmark() (float64, string) {
	const bufferSize = 16 * 1024 * 1024 // 16MB
	data := make([]byte, bufferSize)
	start := time.Now()
	for i := 0; i < 10; i++ {
		for j := range data {
			data[j] = byte(j)
		}
	}
	elapsed := time.Since(start)
	return float64(bufferSize*10) / elapsed.Seconds() / (1024 * 1024 * 1024), "ГБ/с"
}

// RandomBandwidthBenchmark выполняет тест пропускной способности случайного доступа к памяти.
func RandomBandwidthBenchmark() (float64, string) {
	const bufferSize = 16 * 1024 * 1024 // 16MB
	data := make([]byte, bufferSize)
	indices := make([]int, bufferSize)
	for i := range indices {
		indices[i] = rand.Intn(bufferSize)
	}
	start := time.Now()
	for i := 0; i < 10; i++ {
		for _, idx := range indices {
			data[idx]++
		}
	}
	elapsed := time.Since(start)
	return float64(bufferSize*10) / elapsed.Seconds() / (1024 * 1024 * 1024), "ГБ/с"
}

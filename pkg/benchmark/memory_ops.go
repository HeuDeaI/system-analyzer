package benchmark

import (
	"math/rand"
	"time"
)

const dataSize = 1 * 1024 * 1024 // 1 MB

// MemoryBandwidthBenchmark измеряет пропускную способность памяти.
func MemoryBandwidthBenchmark() (readSpeed, writeSpeed, randomSpeed float64) {
	data := make([]byte, dataSize)
	// Заполняем, чтобы избежать оптимизаций ОС с нулевыми страницами
	for i := range data {
		data[i] = byte(i)
	}

	// Тест на запись
	startWrite := time.Now()
	for i := 0; i < dataSize; i++ {
		data[i] = 1
	}
	durationWrite := time.Since(startWrite).Seconds()
	writeSpeed = float64(dataSize) / durationWrite / (1024 * 1024 * 1024)

	// Тест на чтение
	startRead := time.Now()
	var temp byte
	for i := 0; i < dataSize; i++ {
		temp = data[i]
	}
	durationRead := time.Since(startRead).Seconds()
	readSpeed = float64(dataSize) / durationRead / (1024 * 1024 * 1024)

	// Тест на случайный доступ
	indices := make([]int, dataSize)
	for i := range indices {
		indices[i] = i
	}
	rand.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })

	startRandom := time.Now()
	for i := 0; i < dataSize; i++ {
		temp = data[indices[i]]
	}
	durationRandom := time.Since(startRandom).Seconds()
	randomSpeed = float64(dataSize) / durationRandom / (1024 * 1024 * 1024)

	_ = temp // чтобы компилятор не ругался
	return
}

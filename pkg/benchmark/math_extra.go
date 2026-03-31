package benchmark

import (
	"crypto/rand"
	"math/big"
	mathrand "math/rand"
	"time"
)

// LargeIntArithmeticBenchmark измеряет скорость операций с большими числами (BigInt).
// Актуально для асимметричной криптографии и научных вычислений.
func LargeIntArithmeticBenchmark() (float64, string) {
	a, _ := rand.Int(rand.Reader, big.NewInt(1e18))
	b, _ := rand.Int(rand.Reader, big.NewInt(1e18))
	
	start := time.Now()
	for i := 0; i < 100000; i++ {
		c := new(big.Int).Mul(a, b)
		_ = new(big.Int).Div(c, a)
	}
	elapsed := time.Since(start).Seconds()
	
	return 100000.0 / elapsed, "оп/с"
}

// PrimeSearchBenchmark измеряет скорость поиска простых чисел.
// Критично для генерации ключей в RSA.
func PrimeSearchBenchmark() (float64, string) {
	start := time.Now()
	for i := 0; i < 10; i++ {
		_, _ = rand.Prime(rand.Reader, 1024)
	}
	elapsed := time.Since(start).Seconds()
	
	return 10.0 / elapsed, "ключей/с"
}

// MatrixMultiplicationBenchmark измеряет производительность умножения матриц.
func MatrixMultiplicationBenchmark() (float64, string) {
	size := 256
	a := make([][]float64, size)
	b := make([][]float64, size)
	c := make([][]float64, size)
	for i := range a {
		a[i] = make([]float64, size)
		b[i] = make([]float64, size)
		c[i] = make([]float64, size)
		for j := range a[i] {
			a[i][j] = mathrand.Float64()
			b[i][j] = mathrand.Float64()
		}
	}

	start := time.Now()
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	elapsed := time.Since(start).Seconds()
	
	ops := float64(size * size * size * 2) // mul + add
	return ops / elapsed / 1e6, "MFLOPS"
}

// FibonacciBenchmark измеряет скорость рекурсивного вычисления чисел Фибоначчи.
func FibonacciBenchmark() (float64, string) {
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}

	start := time.Now()
	_ = fib(35)
	elapsed := time.Since(start).Seconds()
	
	return 1.0 / elapsed, "оп/с"
}

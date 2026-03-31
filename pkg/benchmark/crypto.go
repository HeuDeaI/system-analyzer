package benchmark

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"math"
	mathrand "math/rand"
	"time"

	"golang.org/x/crypto/chacha20poly1305"
)

// CryptoSHA256Benchmark измеряет производительность хеширования SHA-256.
func CryptoSHA256Benchmark() (float64, string) {
	data := make([]byte, 1024*1024)
	_, _ = rand.Read(data)
	start := time.Now()
	for i := 0; i < 100; i++ {
		h := sha256.New()
		h.Write(data)
		h.Sum(nil)
	}
	elapsed := time.Since(start).Seconds()
	return 100.0 / elapsed, "МБ/с"
}

// CryptoSHA512Benchmark измеряет производительность SHA-512.
func CryptoSHA512Benchmark() (float64, string) {
	data := make([]byte, 1024*1024)
	_, _ = rand.Read(data)
	start := time.Now()
	for i := 0; i < 100; i++ {
		h := sha512.New()
		h.Write(data)
		h.Sum(nil)
	}
	elapsed := time.Since(start).Seconds()
	return 100.0 / elapsed, "МБ/с"
}

// CryptoAESBenchmark измеряет производительность AES-256-GCM.
func CryptoAESBenchmark() (float64, string) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	_, _ = rand.Read(nonce)
	data := make([]byte, 1024*1024)
	_, _ = rand.Read(data)
	start := time.Now()
	for i := 0; i < 100; i++ {
		gcm.Seal(nil, nonce, data, nil)
	}
	elapsed := time.Since(start).Seconds()
	return 100.0 / elapsed, "МБ/с"
}

// CryptoEd25519Benchmark измеряет скорость проверки Ed25519.
func CryptoEd25519Benchmark() (float64, string) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	message := []byte("benchmark")
	sig := ed25519.Sign(priv, message)
	start := time.Now()
	for i := 0; i < 10000; i++ {
		_ = ed25519.Verify(pub, message, sig)
	}
	elapsed := time.Since(start).Seconds()
	return 10000.0 / elapsed, "оп/с"
}

// CryptoAESDecryptionBenchmark имитирует время расшифрования AES-256.
func CryptoAESDecryptionBenchmark() (float64, string) {
	key := make([]byte, 32)
	rand.Read(key)
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	data := make([]byte, 1024*1024)
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	start := time.Now()
	for i := 0; i < 100; i++ {
		_, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			panic(err)
		}
	}
	elapsed := time.Since(start).Seconds()
	return 100.0 / elapsed, "МБ/с"
}

// CryptoRSADecryptionBenchmark имитирует время расшифрования RSA-2048.
func CryptoRSADecryptionBenchmark() (float64, string) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	message := []byte("benchmark")
	ciphertext, _ := rsa.EncryptPKCS1v15(rand.Reader, &priv.PublicKey, message)

	start := time.Now()
	iterations := 50
	for i := 0; i < iterations; i++ {
		_, _ = rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	}
	elapsed := time.Since(start).Seconds()
	return float64(iterations) / elapsed, "оп/с"
}

// CryptoRSAKeyGenBenchmark измеряет скорость генерации ключей RSA-2048.
func CryptoRSAKeyGenBenchmark() (float64, string) {
	start := time.Now()
	iterations := 5
	for i := 0; i < iterations; i++ {
		_, _ = rsa.GenerateKey(rand.Reader, 2048)
	}
	elapsed := time.Since(start).Seconds()
	return float64(iterations) / elapsed, "ключей/с"
}

// RSABruteForceEstimate оценивает теоретическое время полного перебора ключа RSA-2048.
// Возвращает время в годах для полного перебора всех возможных ключей.
func RSABruteForceEstimate() (float64, string) {
	// Генерируем тестовый ключ для измерения времени одной операции
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	message := []byte("test")
	ciphertext, _ := rsa.EncryptPKCS1v15(rand.Reader, &priv.PublicKey, message)

	// Измеряем время одной операции дешифрования
	start := time.Now()
	iterations := 10
	for i := 0; i < iterations; i++ {
		_, _ = rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	}
	elapsed := time.Since(start).Seconds()
	timePerOp := elapsed / float64(iterations)

	// RSA-2048 использует модуль 2048 бит (256 байт)
	// Теоретическое количество возможных ключей: 2^2048 (практически бесконечно)
	// Но для оценки используем энтропию закрытого ключа - примерно 2^128 операций
	// для атаки на основе решеток (lattice-based attacks) или 2^112 для brute-force
	// Для чистого brute-force: сложность ~ 2^2048 (невозможно)
	// Практическая оценка: время для подбора 2048-битного модуля
	const operationsNeeded = 1.15e77 // ~ 2^256 для эффективных атак на RSA
	secondsNeeded := operationsNeeded * timePerOp
	yearsNeeded := secondsNeeded / (365.25 * 24 * 3600)

	return yearsNeeded, "лет"
}

// CaesarCipherBruteForce оценивает время перебора шифра Цезаря (25 ключей).
func CaesarCipherBruteForce() (float64, string) {
	text := []byte("HELLO WORLD THIS IS A TEST MESSAGE")
	encrypted := make([]byte, len(text))
	// Шифруем со сдвигом 15
	shift := byte(15)
	for i, c := range text {
		if c >= 'A' && c <= 'Z' {
			encrypted[i] = byte((int(c-'A')+int(shift))%26) + 'A'
		} else {
			encrypted[i] = c
		}
	}

	start := time.Now()
	iterations := 10000
	for i := 0; i < iterations; i++ {
		// Перебираем все 25 возможных ключей
		for key := 1; key < 26; key++ {
			_ = key // Проверка ключа
			// Декодируем с текущим ключом
			for j, c := range encrypted {
				if c >= 'A' && c <= 'Z' {
					_ = byte((int(c-'A')-key+26)%26) + 'A'
				} else {
					_ = c
				}
				_ = j
			}
		}
	}
	elapsed := time.Since(start).Seconds()
	timePerAttempt := elapsed / float64(iterations*25)
	// Для взлома одного сообщения нужно 25 попыток
	bruteForceTime := timePerAttempt * 25
	return bruteForceTime, "сек"
}

// XOR8BitBruteForce оценивает время перебора 8-битного XOR ключа (256 вариантов).
func XOR8BitBruteForce() (float64, string) {
	data := make([]byte, 100)
	rand.Read(data)
	key := byte(0xAB)
	encrypted := make([]byte, len(data))
	for i, b := range data {
		encrypted[i] = b ^ key
	}

	start := time.Now()
	iterations := 5000
	for i := 0; i < iterations; i++ {
		// Перебираем все 256 возможных 8-битных ключей
		for testKey := 0; testKey < 256; testKey++ {
			k := byte(testKey)
			_ = k
			// Декодируем с текущим ключом
			for _, c := range encrypted {
				_ = c ^ k
			}
		}
	}
	elapsed := time.Since(start).Seconds()
	timePerAttempt := elapsed / float64(iterations*256)
	// Для взлома нужно 256 попыток
	bruteForceTime := timePerAttempt * 256
	return bruteForceTime, "сек"
}

// PIN4DigitBruteForce оценивает время перебора 4-значного PIN-кода (10000 комбинаций).
func PIN4DigitBruteForce() (float64, string) {
	targetPIN := 7384

	start := time.Now()
	iterations := 100
	for i := 0; i < iterations; i++ {
		// Перебираем все 10000 возможных PIN
		for pin := 0; pin < 10000; pin++ {
			if pin == targetPIN {
				_ = pin
				break
			}
		}
	}
	elapsed := time.Since(start).Seconds()
	timePerAttempt := elapsed / float64(iterations*10000)
	// В среднем нужно 5000 попыток (среднее от 1 до 10000)
	bruteForceTime := timePerAttempt * 5000
	return bruteForceTime, "сек"
}

// AES128BruteForceEstimate оценивает теоретическое время brute-force AES-128.
func AES128BruteForceEstimate() (float64, string) {
	// Измеряем время одной операции AES
	key := make([]byte, 16)
	rand.Read(key)
	block, _ := aes.NewCipher(key)
	data := make([]byte, 16)
	rand.Read(data)

	start := time.Now()
	iterations := 10000
	for i := 0; i < iterations; i++ {
		block.Encrypt(data, data)
	}
	elapsed := time.Since(start).Seconds()
	timePerOp := elapsed / float64(iterations)

	// AES-128 имеет 128-битный ключ = 2^128 возможных комбинаций
	const operationsNeeded = 3.4e38 // 2^128
	secondsNeeded := operationsNeeded * timePerOp
	yearsNeeded := secondsNeeded / (365.25 * 24 * 3600)

	return yearsNeeded, "лет"
}

// RAMCopyBenchmark measures memory copy performance.
func RAMCopyBenchmark() (float64, string) {
	const dataSize = 128 * 1024 * 1024 // 128 MB
	data := make([]byte, dataSize)

	// Fill with data to avoid OS zero-page optimizations
	for i := range data {
		data[i] = byte(i)
	}

	// Test write speed - byte by byte to prevent optimization
	startWrite := time.Now()
	for i := 0; i < dataSize; i++ {
		data[i] = 1
	}
	durationWrite := time.Since(startWrite).Seconds()
	writeSpeed := float64(dataSize) / durationWrite / (1024 * 1024 * 1024)

	// Test read speed - byte by byte to prevent optimization
	startRead := time.Now()
	var temp byte
	for i := 0; i < dataSize; i++ {
		temp = data[i]
	}
	durationRead := time.Since(startRead).Seconds()
	readSpeed := float64(dataSize) / durationRead / (1024 * 1024 * 1024)

	_ = temp // Prevent compiler from optimizing away the read

	// Return average of read and write speed
	avgSpeed := (readSpeed + writeSpeed) / 2
	return avgSpeed, "GB/s"
}

// RAMLatencyBenchmark measures memory latency.
func RAMLatencyBenchmark() (float64, string) {
	const size = 16 * 1024 * 1024 // 16MB to exceed L3 cache on many CPUs
	data := make([]int, size/8)
	
	// Create a random pointer chain to prevent prefetching
	indices := mathrand.Perm(len(data))
	for i := 0; i < len(indices)-1; i++ {
		data[indices[i]] = indices[i+1]
	}
	data[indices[len(indices)-1]] = indices[0]

	start := time.Now()
	curr := 0
	iterations := 1000000
	for i := 0; i < iterations; i++ {
		curr = data[curr]
	}
	elapsed := time.Since(start).Seconds()
	
	latencyNs := (elapsed / float64(iterations)) * 1e9
	return latencyNs, "ns"
}

// EstimateBruteForceTimeGeneric оценивает время brute-force для произвольного количества бит.
func EstimateBruteForceTimeGeneric(bits int) string {
	// Измеряем время одной базовой операции (например, SHA-256)
	data := make([]byte, 64)
	h := sha256.New()
	
	start := time.Now()
	iterations := 100000
	for i := 0; i < iterations; i++ {
		h.Reset()
		h.Write(data)
		h.Sum(nil)
	}
	elapsed := time.Since(start).Seconds()
	timePerOp := elapsed / float64(iterations)

	operations := math.Pow(2, float64(bits))
	seconds := operations * timePerOp

	return FormatDuration(seconds)
}

// TestAESSpeed измеряет скорость шифрования/дешифрования AES с заданным размером ключа.
// Возвращает: время шифрования, время дешифрования, пропускную способность в МБ/с.
func TestAESSpeed(keySize int, dataSizeMB int) (encryptTime, decryptTime, throughput float64) {
	// AES требует ключи 16, 24 или 32 байта (128, 192, 256 бит)
	// Для маленьких ключей дополняем до 16 байт (AES-128)
	effectiveKeySize := keySize
	if keySize < 16 {
		effectiveKeySize = 16
	} else if keySize > 16 && keySize < 24 {
		effectiveKeySize = 24
	} else if keySize > 24 && keySize < 32 {
		effectiveKeySize = 32
	}

	key := make([]byte, effectiveKeySize)
	rand.Read(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, 0, 0
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, 0, 0
	}
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	data := make([]byte, dataSizeMB*1024*1024)
	rand.Read(data)

	// Измеряем шифрование
	start := time.Now()
	ciphertext := gcm.Seal(nil, nonce, data, nil)
	encryptTime = time.Since(start).Seconds()

	// Измеряем дешифрование
	start = time.Now()
	_, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, 0, 0
	}
	decryptTime = time.Since(start).Seconds()

	throughput = float64(dataSizeMB) / encryptTime
	return
}

// EstimateBruteForceTimeAES возвращает строку с оценкой времени brute-force для AES.
func EstimateBruteForceTimeAES(keyBits int) string {
	// Минимальный ключ для AES - 16 байт (128 бит)
	effectiveKeyBits := keyBits
	if keyBits < 128 {
		effectiveKeyBits = 128
	}

	// Измеряем время одной операции
	key := make([]byte, effectiveKeyBits/8)
	rand.Read(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "н/д"
	}
	data := make([]byte, 16)
	rand.Read(data)

	start := time.Now()
	iterations := 10000
	for i := 0; i < iterations; i++ {
		block.Encrypt(data, data)
	}
	elapsed := time.Since(start).Seconds()
	timePerOp := elapsed / float64(iterations)

	// 2^keyBits операций
	operations := math.Pow(2, float64(keyBits))
	seconds := operations * timePerOp

	return FormatDuration(seconds)
}

// FormatDuration форматирует секунды в человекочитаемую строку.
func FormatDuration(seconds float64) string {
	// Format based on magnitude
	if seconds < 0.001 {
		return fmt.Sprintf("%.2f мкс", seconds*1e6)
	} else if seconds < 1 {
		return fmt.Sprintf("%.2f мс", seconds*1e3)
	} else if seconds < 60 {
		return fmt.Sprintf("%.2f сек", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.2f мин", seconds/60)
	} else if seconds < 24*3600 {
		return fmt.Sprintf("%.2f час", seconds/3600)
	} else if seconds < 365.25*24*3600 {
		return fmt.Sprintf("%.2f дней", seconds/(24*3600))
	} else {
		years := seconds / (365.25 * 24 * 3600)
		if years > 1e9 {
			return fmt.Sprintf("%.2e лет (вечность)", years)
		}
		return fmt.Sprintf("%.2e лет", years)
	}
}

// TestRSASpeed измеряет скорость шифрования/дешифрования RSA.
// Возвращает: время шифрования, время дешифрования, пропускную способность в МБ/с.
func TestRSASpeed(keyBits int, dataSizeMB int) (encryptTime, decryptTime, throughput float64) {
	// Минимальный размер ключа для RSA - 32 бита (технически возможно для тестов, но небезопасно)
	if keyBits < 32 {
		keyBits = 32
	}

	priv, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return 0, 0, 0
	}
	pub := &priv.PublicKey

	// RSA шифрует блоками
	blockSize := (keyBits / 8) - 11 // PKCS1v15 padding
	if blockSize <= 0 {
		blockSize = 1
	}
	numBlocks := (dataSizeMB * 1024 * 1024) / blockSize
	if numBlocks == 0 {
		numBlocks = 1
	}

	message := make([]byte, blockSize)
	rand.Read(message)

	// Измеряем шифрование
	start := time.Now()
	for i := 0; i < numBlocks; i++ {
		rsa.EncryptPKCS1v15(rand.Reader, pub, message)
	}
	encryptTime = time.Since(start).Seconds()

	// Подготавливаем зашифрованный блок для дешифрования
	ciphertext, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, message)

	// Измеряем дешифрование (значительно медленнее)
	start = time.Now()
	for i := 0; i < numBlocks; i++ {
		rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	}
	decryptTime = time.Since(start).Seconds()

	throughput = float64(dataSizeMB) / encryptTime
	return
}

// EstimateBruteForceTimeRSA возвращает строку с оценкой времени brute-force для RSA.
func EstimateBruteForceTimeRSA(keyBits int) string {
	// Минимальный размер ключа для RSA
	if keyBits < 32 {
		return "н/д (<32 бит)"
	}

	// Генерируем ключ для измерения
	priv, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return "ошибка генерации"
	}
	message := []byte("test")
	pub := &priv.PublicKey

	// Проверяем возможность шифрования (размер сообщения)
	blockSize := (keyBits / 8) - 11
	if blockSize <= 0 {
		message = []byte("t")
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, message)
	if err != nil {
		return "ошибка шифр."
	}

	start := time.Now()
	iterations := 10
	for i := 0; i < iterations; i++ {
		rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	}
	elapsed := time.Since(start).Seconds()
	timePerOp := elapsed / float64(iterations)

	// Для RSA используем упрощенную оценку
	var operations float64
	switch {
	case keyBits <= 64:
		operations = math.Pow(2, float64(keyBits))
	case keyBits <= 512:
		operations = 1e9
	case keyBits <= 1024:
		operations = 1e18
	case keyBits <= 2048:
		operations = 1e38
	case keyBits <= 4096:
		operations = 1e77
	default:
		operations = math.Pow(2, float64(keyBits)/2)
	}

	seconds := operations * timePerOp

	return FormatDuration(seconds)
}

// TestHashSpeed измеряет скорость хеширования.
// Возвращает время хеширования в секундах.
func TestHashSpeed(hasher hash.Hash, dataSizeMB int) float64 {
	data := make([]byte, dataSizeMB*1024*1024)
	rand.Read(data)

	start := time.Now()
	hasher.Write(data)
	hasher.Sum(nil)
	return time.Since(start).Seconds()
}

// TestEd25519Speed измеряет скорость подписи и верификации Ed25519.
// Возвращает время подписи, время верификации.
func TestEd25519Speed(dataSizeMB int) (signTime, verifyTime float64) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	message := make([]byte, dataSizeMB*1024*1024)
	rand.Read(message)

	// Измеряем подпись
	start := time.Now()
	sig := ed25519.Sign(priv, message)
	signTime = time.Since(start).Seconds()

	// Измеряем верификацию
	start = time.Now()
	ed25519.Verify(pub, message, sig)
	verifyTime = time.Since(start).Seconds()

	return
}

// TestHMACSpeed измеряет скорость HMAC-SHA256.
func TestHMACSpeed(keySize int, dataSizeMB int) float64 {
	key := make([]byte, keySize/8)
	rand.Read(key)
	data := make([]byte, dataSizeMB*1024*1024)
	rand.Read(data)

	h := hmac.New(sha256.New, key)
	start := time.Now()
	h.Write(data)
	h.Sum(nil)
	return time.Since(start).Seconds()
}

// TestChaCha20Speed измеряет скорость ChaCha20-Poly1305.
func TestChaCha20Speed(keySize int, dataSizeMB int) (encryptTime, decryptTime, throughput float64) {
	// ChaCha20 требует ключ 32 байта (256 бит)
	effectiveKeySize := 32
	if keySize/8 > 32 {
		effectiveKeySize = keySize / 8
	}

	key := make([]byte, effectiveKeySize)
	rand.Read(key)
	data := make([]byte, dataSizeMB*1024*1024)
	rand.Read(data)

	// Создаем ChaCha20-Poly1305 cipher
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return 0, 0, 0
	}

	nonce := make([]byte, aead.NonceSize())
	rand.Read(nonce)

	// Измеряем шифрование
	start := time.Now()
	ciphertext := aead.Seal(nil, nonce, data, nil)
	encryptTime = time.Since(start).Seconds()

	// Измеряем дешифрование
	start = time.Now()
	_, err = aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, 0, 0
	}
	decryptTime = time.Since(start).Seconds()

	throughput = float64(dataSizeMB) / encryptTime
	return
}

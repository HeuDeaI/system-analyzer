package memory

import (
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

const (
	testSize    = 64 * 1024 * 1024 // 64 MB
	bufferSize  = 4 * 1024         // 4 KB
	randomReads = 10000
)

func CreateTestFile() (string, error) {
	file, err := ioutil.TempFile("", "flash_test_*")
	if err != nil {
		return "", err
	}
	defer file.Close()

	data := make([]byte, bufferSize)
	rand.Read(data)
	for i := 0; i < testSize/bufferSize; i++ {
		_, err := file.Write(data)
		if err != nil {
			return "", err
		}
	}
	return file.Name(), nil
}

func CleanupTestFile(filePath string) {
	os.Remove(filePath)
}

func FlashWriteSpeedBenchmark(filePath string) (float64, string, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return 0, "", err
	}
	defer file.Close()

	data := make([]byte, bufferSize)
	rand.Read(data)
	start := time.Now()
	for i := 0; i < testSize/bufferSize; i++ {
		_, err := file.Write(data)
		if err != nil {
			return 0, "", err
		}
	}
	duration := time.Since(start).Seconds()
	return float64(testSize) / duration / (1024 * 1024), "MB/s", nil
}

func FlashReadSpeedBenchmark(filePath string) (float64, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, "", err
	}
	defer file.Close()

	data := make([]byte, bufferSize)
	start := time.Now()
	for {
		_, err := file.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, "", err
		}
	}
	duration := time.Since(start).Seconds()
	return float64(testSize) / duration / (1024 * 1024), "МБ/с", nil
}

func FlashRandomReadSpeedBenchmark(filePath string) (float64, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, "", err
	}
	defer file.Close()

	data := make([]byte, bufferSize)
	start := time.Now()
	for i := 0; i < randomReads; i++ {
		offset := rand.Int63n(int64(testSize - bufferSize))
		_, err := file.ReadAt(data, offset)
		if err != nil {
			return 0, "", err
		}
	}
	duration := time.Since(start).Seconds()
	return float64(randomReads*bufferSize) / duration / (1024 * 1024), "МБ/с", nil
}

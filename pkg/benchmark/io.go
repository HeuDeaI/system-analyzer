package benchmark

import (
	"os"
	"path/filepath"
	"time"
)

// DiskIOBenchmark measures sequential file I/O performance.
func DiskIOBenchmark() (float64, string) {
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "sys_analyzer_bench.tmp")
	defer os.Remove(tmpFile)
	
	data := make([]byte, 64*1024*1024) // 64MB
	
	start := time.Now()
	f, _ := os.Create(tmpFile)
	f.Write(data)
	f.Sync()
	f.Seek(0, 0)
	readData := make([]byte, 64*1024*1024)
	f.Read(readData)
	f.Close()
	
	elapsed := time.Since(start).Seconds()
	
	mbPerSec := 128.0 / elapsed // Write 64MB + Read 64MB
	return mbPerSec, "MB/s"
}

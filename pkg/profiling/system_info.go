package profiling

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var (
	lastNetIO  net.IOCountersStat
	lastDiskIO disk.IOCountersStat
	lastTime   time.Time
)

func init() {
	netIOs, _ := net.IOCounters(false)
	if len(netIOs) > 0 {
		lastNetIO = netIOs[0]
	}
	diskIOs, _ := disk.IOCounters()
	for _, stat := range diskIOs {
		lastDiskIO.ReadBytes += stat.ReadBytes
		lastDiskIO.WriteBytes += stat.WriteBytes
	}
	lastTime = time.Now()
}

func GetCPUInfo() (string, error) {
	info, err := cpu.Info()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s (%.2f ГГц, %d ядер)", info[0].ModelName, info[0].Mhz/1000, info[0].Cores), nil
}

func GetMemoryUsage() (string, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.2f/%.2f ГБ (%.2f%%)", float64(v.Used)/1e9, float64(v.Total)/1e9, v.UsedPercent), nil
}

func GetCPUUsage() ([]float64, error) {
	return cpu.Percent(time.Second, true)
}

func GetDiskIO() (string, error) {
	ios, err := disk.IOCounters()
	if err != nil {
		return "", err
	}
	var currentRead, currentWrite uint64
	for _, stat := range ios {
		currentRead += stat.ReadBytes
		currentWrite += stat.WriteBytes
	}
	now := time.Now()
	duration := now.Sub(lastTime).Seconds()
	if duration < 1 {
		duration = 1
	}
	readSpeed := float64(currentRead-lastDiskIO.ReadBytes) / duration / (1024 * 1024)
	writeSpeed := float64(currentWrite-lastDiskIO.WriteBytes) / duration / (1024 * 1024)
	lastDiskIO.ReadBytes = currentRead
	lastDiskIO.WriteBytes = currentWrite
	return fmt.Sprintf("Чтение: %.2f МБ/с / Запись: %.2f МБ/с", readSpeed, writeSpeed), nil
}

func GetNetIO() (string, error) {
	ios, err := net.IOCounters(false)
	if err != nil || len(ios) == 0 {
		return "", err
	}
	currentIO := ios[0]
	now := time.Now()
	duration := now.Sub(lastTime).Seconds()
	if duration < 1 {
		duration = 1
	}
	sentSpeed := float64(currentIO.BytesSent-lastNetIO.BytesSent) / duration / (1024 * 1024)
	recvSpeed := float64(currentIO.BytesRecv-lastNetIO.BytesRecv) / duration / (1024 * 1024)
	lastNetIO = currentIO
	lastTime = now
	return fmt.Sprintf("Получение: %.2f МБ/с / Отправка: %.2f МБ/с", recvSpeed, sentSpeed), nil
}

func GetHostInfo() (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	uptime := time.Duration(info.Uptime) * time.Second
	days := int(uptime.Hours() / 24)
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	return fmt.Sprintf("%s %s (Работает: %d д, %d ч, %d м)", info.Platform, info.PlatformVersion, days, hours, minutes), nil
}

func GetLoadAvg() (string, error) {
	avg, err := load.Avg()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.2f (1м), %.2f (5м), %.2f (15м)", avg.Load1, avg.Load5, avg.Load15), nil
}

func GetDiskUsage() (string, error) {
	usage, err := disk.Usage("/")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/: %.1f ГБ из %.1f ГБ", float64(usage.Used)/1e9, float64(usage.Total)/1e9), nil
}

// GetCoreTypes определяет количество P-cores и E-cores на Apple Silicon.
func GetCoreTypes() (pCores, eCores int, err error) {
	cmd := exec.Command("sysctl", "-n", "hw.perflevel0.physicalcpu")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get P-cores: %w", err)
	}
	pCores, err = strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse P-cores: %w", err)
	}

	cmd = exec.Command("sysctl", "-n", "hw.perflevel1.physicalcpu")
	output, err = cmd.Output()
	if err != nil {
		// Если нет E-cores, команда может завершиться ошибкой. Это нормально.
		return pCores, 0, nil
	}
	eCores, err = strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return pCores, 0, fmt.Errorf("failed to parse E-cores: %w", err)
	}

	return pCores, eCores, nil
}

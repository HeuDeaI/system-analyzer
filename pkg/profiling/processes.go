package profiling

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	PID    int32
	Name   string
	CPU    float64
	Memory float64
}

func GetTopProcesses(n int) ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var results []ProcessInfo
	for _, p := range processes {
		name, _ := p.Name()
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()
		if cpu > 0.1 || mem > 0.1 {
			results = append(results, ProcessInfo{
				PID:    p.Pid,
				Name:   name,
				CPU:    cpu,
				Memory: float64(mem),
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CPU > results[j].CPU
	})

	if len(results) > n {
		return results[:n], nil
	}
	return results, nil
}

func FormatProcessList(procs []ProcessInfo) string {
	if len(procs) == 0 {
		return "Нет активных процессов"
	}

	var sb strings.Builder
	// Заголовок с разделителями - увеличена ширина PID
	sb.WriteString("┌────────┬────────────────────┬──────────┬──────────┐\n")
	sb.WriteString(fmt.Sprintf("│ %-6s │ %-18s │ %-8s │ %-8s │\n", "PID", "Название", "CPU%", "Пам%"))
	sb.WriteString("├────────┼────────────────────┼──────────┼──────────┤\n")

	for _, p := range procs {
		name := p.Name
		if len(name) > 18 {
			name = name[:15] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ %-6d │ %-18s │ %7.1f%% │ %7.1f%% │\n", p.PID, name, p.CPU, p.Memory))
	}
	sb.WriteString("└────────┴────────────────────┴──────────┴──────────┘")
	return sb.String()
}

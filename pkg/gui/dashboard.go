package gui

import (
	"fmt"
	"log"
	"system-analyzer/pkg/profiling"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewDashboardPanel создает и возвращает панель мониторинга.
func NewDashboardPanel() fyne.CanvasObject {
	// Labels for dynamic data
	cpuInfoLabel := widget.NewLabel("...")
	memUsageLabel := widget.NewLabel("...")
	diskIOLabel := widget.NewLabel("...")
	netIOLabel := widget.NewLabel("...")
	hostLabel := widget.NewLabel("...")
	loadLabel := widget.NewLabel("...")
	diskUsageLabel := widget.NewLabel("...")

	// Cards for grouping information
	cpuCard := widget.NewCard("⚙️ Процессор", "", cpuInfoLabel)
	memCard := widget.NewCard("📚 ОЗУ", "", memUsageLabel)
	diskCard := widget.NewCard("💾 Диск", "", container.NewVBox(diskUsageLabel, diskIOLabel))
	netCard := widget.NewCard("🌐 Сеть", "", netIOLabel)
	sysCard := widget.NewCard("🖥️ Система", "", container.NewVBox(hostLabel, loadLabel))

	// Core usage bars
	coreUsageContainer := container.NewVBox()
	coreUsageCard := widget.NewCard("Загрузка ядер ЦП", "", coreUsageContainer)
	var coreBars []*widget.ProgressBar
	var coreLabels []*widget.Label

	// Get core types once
	pCores, _, err := profiling.GetCoreTypes()
	if err != nil {
		log.Println("Could not get core types:", err)
		pCores = -1 // Disable core type display
	}

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			cpuInfo, _ := profiling.GetCPUInfo()
			memUsage, _ := profiling.GetMemoryUsage()
			cpuUsage, _ := profiling.GetCPUUsage()
			diskIO, _ := profiling.GetDiskIO()
			netIO, _ := profiling.GetNetIO()
			host, _ := profiling.GetHostInfo()
			load, _ := profiling.GetLoadAvg()
			diskUsage, _ := profiling.GetDiskUsage()

			// Update labels
			cpuInfoLabel.SetText(cpuInfo)
			memUsageLabel.SetText(memUsage)
			diskIOLabel.SetText(fmt.Sprintf("Скорость: %s", diskIO))
			netIOLabel.SetText(fmt.Sprintf("Скорость: %s", netIO))
			hostLabel.SetText(host)
			loadLabel.SetText(fmt.Sprintf("Средняя нагрузка: %s", load))
			diskUsageLabel.SetText(diskUsage)

			// Update core usage bars
			if len(coreBars) != len(cpuUsage) {
				coreUsageContainer.RemoveAll()
				coreBars = nil
				coreLabels = nil
				for i := 0; i < len(cpuUsage); i++ {
					bar := widget.NewProgressBar()
					var label *widget.Label
					if pCores != -1 {
						if i < pCores {
							label = widget.NewLabel(fmt.Sprintf("Ядро %d (Производительное):", i+1))
						} else {
							label = widget.NewLabel(fmt.Sprintf("Ядро %d (Энергосберегающее):", i+1))
						}
					} else {
						label = widget.NewLabel(fmt.Sprintf("Ядро %d:", i+1))
					}
					coreUsageContainer.Add(label)
					coreUsageContainer.Add(bar)
					coreBars = append(coreBars, bar)
					coreLabels = append(coreLabels, label)
				}
			}
			for i, p := range cpuUsage {
				if i < len(coreBars) {
					coreBars[i].SetValue(p / 100)
				}
			}
		}
	}()

	// Layout
	grid := container.NewGridWithColumns(2, cpuCard, memCard, diskCard, netCard, sysCard)
	return container.NewScroll(container.NewVBox(grid, coreUsageCard))
}

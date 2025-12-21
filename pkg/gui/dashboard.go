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
// Панель отображает системную информацию в реальном времени.
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

	// Контейнер и карточка для отображения загрузки ядер ЦП.
	coreUsageContainer := container.NewVBox()
	coreUsageCard := widget.NewCard("Загрузка ядер ЦП", "", coreUsageContainer)
	var coreBars []*widget.ProgressBar // Слайсы для хранения виджетов прогресс-баров.
	var coreLabels []*widget.Label     // Слайсы для хранения меток ядер.

	// Получаем информацию о типах ядер (P-cores и E-cores) один раз при создании панели.
	pCores, _, err := profiling.GetCoreTypes()
	if err != nil {
		// Если получить информацию не удалось (например, на ОС, отличной от macOS),
		// логируем ошибку и отключаем отображение типов ядер.
		log.Println("Could not get core types:", err)
		pCores = -1
	}

	// Запускаем фоновую горутину для периодического обновления данных.
	go func() {
		for {
			// Обновляем данные каждые 500 миллисекунд.
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

			// Динамически обновляем количество прогресс-баров в соответствии с количеством ядер.
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

	// Собираем макет панели: сетка с информационными карточками и карточка с загрузкой ядер.
	// Все оборачивается в Scroll контейнер для возможности прокрутки.
	grid := container.NewGridWithColumns(2, cpuCard, memCard, diskCard, netCard, sysCard)
	return container.NewScroll(container.NewVBox(grid, coreUsageCard))
}

package gui

import (
	"fmt"
	"system-analyzer/pkg/profiling"
	"time"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewDashboardPanel() fyne.CanvasObject {
	cpuInfoLabel := widget.NewLabel("...")
	memUsageLabel := widget.NewLabel("...")
	diskIOLabel := widget.NewLabel("...")
	netIOLabel := widget.NewLabel("...")
	hostLabel := widget.NewLabel("...")
	loadLabel := widget.NewLabel("...")
	diskUsageLabel := widget.NewLabel("...")
	topProcsLabel := widget.NewLabel("...")

	cpuInfoLabel.TextStyle = fyne.TextStyle{Monospace: true}
	memUsageLabel.TextStyle = fyne.TextStyle{Monospace: true}
	diskIOLabel.TextStyle = fyne.TextStyle{Monospace: true}
	netIOLabel.TextStyle = fyne.TextStyle{Monospace: true}
	hostLabel.TextStyle = fyne.TextStyle{Monospace: true}
	loadLabel.TextStyle = fyne.TextStyle{Monospace: true}
	diskUsageLabel.TextStyle = fyne.TextStyle{Monospace: true}
	topProcsLabel.TextStyle = fyne.TextStyle{Monospace: true}

	coreUsageContainer := container.NewVBox()
	var coreBars []*widget.ProgressBar

	go func() {
		for {
			time.Sleep(1 * time.Second)
			cpuInfo, _ := profiling.GetCPUInfo()
			memUsage, _ := profiling.GetMemoryUsage()
			cpuUsage, _ := profiling.GetCPUUsage()
			diskIO, _ := profiling.GetDiskIO()
			netIO, _ := profiling.GetNetIO()
			host, _ := profiling.GetHostInfo()
			load, _ := profiling.GetLoadAvg()
			diskUsage, _ := profiling.GetDiskUsage()
			procs, _ := profiling.GetTopProcesses(5)

			cpuInfoLabel.SetText(cpuInfo)
			memUsageLabel.SetText(memUsage)
			diskIOLabel.SetText(diskIO)
			netIOLabel.SetText(netIO)
			hostLabel.SetText(host)
			loadLabel.SetText(load)
			diskUsageLabel.SetText(diskUsage)
			topProcsLabel.SetText(profiling.FormatProcessList(procs))

			if len(coreBars) != len(cpuUsage) {
				coreUsageContainer.RemoveAll()
				coreBars = nil
				for i := 0; i < len(cpuUsage); i++ {
					bar := widget.NewProgressBar()
					coreBars = append(coreBars, bar)
					coreUsageContainer.Add(container.NewBorder(nil, nil, widget.NewLabel(fmt.Sprintf("Ядро %d", i+1)), nil, bar))
				}
			}
			for i, p := range cpuUsage {
				if i < len(coreBars) {
					coreBars[i].SetValue(p / 100)
				}
			}
		}
	}()

	sysGrid := container.New(layout.NewGridLayout(2),
		createStatCard("Система", hostLabel),
		createStatCard("Средняя нагрузка", loadLabel),
		createStatCard("Процессор", cpuInfoLabel),
		createStatCard("Оперативная память", memUsageLabel),
		createStatCard("Дисковое пространство", diskUsageLabel),
		createStatCard("Дисковые операции", diskIOLabel),
		createStatCard("Сетевые операции", netIOLabel),
	)

	return container.NewPadded(container.NewVBox(
		canvas.NewText("Системная телеметрия", theme.ForegroundColor()),
		widget.NewSeparator(),
		sysGrid,
		widget.NewSeparator(),
		canvas.NewText("Загрузка ядер", theme.ForegroundColor()),
		coreUsageContainer,
		widget.NewSeparator(),
		canvas.NewText("Топ-5 процессов по ЦП", theme.ForegroundColor()),
		container.NewPadded(topProcsLabel),
	))
}

func createStatCard(title string, content fyne.CanvasObject) fyne.CanvasObject {
	t := canvas.NewText(title, theme.PrimaryColor())
	t.TextSize = 12
	t.TextStyle = fyne.TextStyle{Bold: true}
	bg := canvas.NewRectangle(color.Transparent)
	bg.StrokeColor = theme.DisabledColor()
	bg.StrokeWidth = 1
	return container.NewStack(bg, container.NewPadded(container.NewVBox(t, content)))
}

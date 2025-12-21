package main

import (
	"system-analyzer/pkg/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("Анализатор Производительности Системы")

	tabs := container.NewAppTabs(
		container.NewTabItem("📊 Мониторинг", gui.NewDashboardPanel()),
		container.NewTabItem("🚀 Производительность", gui.NewBenchmarkPanel()),
		container.NewTabItem("💾 Память", gui.NewMemoryPanel()),
		container.NewTabItem("🔀 Многозадачность", gui.NewConcurrencyPanel()),
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1280, 960))
	w.ShowAndRun()
}

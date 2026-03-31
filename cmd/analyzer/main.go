package main

import (
	"system-analyzer/pkg/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

// main is the entry point of the application.
func main() {
	// Create a new Fyne application.
	a := app.New()
	// Create the main application window with a title.
	w := a.NewWindow("Системный Анализатор")

	// Создаем контейнер с вкладками для различных панелей.
	tabs := container.NewAppTabs(
		container.NewTabItem("Мониторинг", container.NewScroll(gui.NewDashboardPanel())),
		container.NewTabItem("Процессор", container.NewScroll(gui.NewBenchmarkPanel("Процессор"))),
		container.NewTabItem("Память", container.NewScroll(gui.NewBenchmarkPanel("Память"))),
		container.NewTabItem("Математика", container.NewScroll(gui.NewBenchmarkPanel("Математика"))),
		container.NewTabItem("Криптография", container.NewScroll(gui.NewInteractiveCryptoPanel())),
	)

	// Set the tab container as the main window content.
	w.SetContent(tabs)
	// Set window to maximum available size
	w.Resize(fyne.NewSize(1600, 1000))
	w.CenterOnScreen()
	w.ShowAndRun()
}

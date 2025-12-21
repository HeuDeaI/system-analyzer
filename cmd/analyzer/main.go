package main

import (
	"system-analyzer/pkg/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

// main является точкой входа в приложение.
func main() {
	// Создаем новое приложение Fyne.
	a := app.New()
	// Создаем главное окно приложения с заголовком.
	w := a.NewWindow("Анализатор Производительности Системы")

	// Создаем контейнер с вкладками для различных панелей.
	tabs := container.NewAppTabs(
		// Вкладка "Мониторинг" для отображения системной информации в реальном времени.
		container.NewTabItem("📊 Мониторинг", gui.NewDashboardPanel()),
		// Вкладка "Производительность" для тестов процессора и пропускной способности ОЗУ.
		container.NewTabItem("🚀 Производительность", gui.NewBenchmarkPanel()),
		// Вкладка "Память" для тестов задержки кэшей, ОЗУ и скорости флеш-накопителей.
		container.NewTabItem("💾 Память", gui.NewMemoryPanel()),
		// Вкладка "Многозадачность" для тестов производительности горутин, каналов и мьютексов.
		container.NewTabItem("🔀 Многозадачность", gui.NewConcurrencyPanel()),
	)

	// Устанавливаем контейнер с вкладками как основное содержимое окна.
	w.SetContent(tabs)
	// Устанавливаем начальный размер окна.
	w.Resize(fyne.NewSize(1280, 960))
	// Центрируем окно на экране.
	w.CenterOnScreen()
	// Отображаем окно и запускаем главный цикл приложения.
	w.ShowAndRun()
}

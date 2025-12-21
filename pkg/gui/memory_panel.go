package gui

import (
	"fmt"
	"log"
	"system-analyzer/pkg/benchmark"
	"system-analyzer/pkg/memory"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewMemoryPanel создает и возвращает панель для тестирования подсистемы памяти.
func NewMemoryPanel() fyne.CanvasObject {
	// Количество итераций для каждого теста.
	const iterations = 3

	// Слайс, определяющий тесты задержки памяти.
	tests := []struct {
		name        string
		description string
		fn          benchmark.BenchmarkFunc
		resultLabel *widget.Label
	}{
		{"💡 Задержка L1 кэша", "Измеряет время доступа к кэшу первого уровня, самой быстрой памяти процессора.", memory.L1CacheLatencyBenchmark, widget.NewLabel("Результат: ...")},
		{"💻 Задержка L2 кэша", "Измеряет время доступа к кэшу второго уровня.", memory.L2CacheLatencyBenchmark, widget.NewLabel("Результат: ...")},
		{"💾 Задержка L3 кэша", "Измеряет время доступа к кэшу третьего уровня, общему для всех ядер.", memory.L3CacheLatencyBenchmark, widget.NewLabel("Результат: ...")},
		{"🏃 Задержка ОЗУ", "Измеряет время доступа к основной оперативной памяти (RAM).", memory.RAMLatencyBenchmark, widget.NewLabel("Результат: ...")},
	}

	// Слайс, определяющий тесты флеш-памяти (SSD/HDD).
	flashTests := []struct {
		name        string
		description string
		fn          func(string) (float64, string, error)
		resultLabel *widget.Label
	}{
		{"💿 Скорость записи на флеш-память", "Тестирует скорость последовательной записи большого файла.", memory.FlashWriteSpeedBenchmark, widget.NewLabel("Результат: ...")},
		{"📀 Скорость чтения с флеш-памяти", "Тестирует скорость последовательного чтения большого файла.", memory.FlashReadSpeedBenchmark, widget.NewLabel("Результат: ...")},
		{"🎲 Случайное чтение с флеш-памяти", "Тестирует скорость чтения случайных блоков данных из файла.", memory.FlashRandomReadSpeedBenchmark, widget.NewLabel("Результат: ...")},
	}

	// Создаем виджеты для отображения прогресса и запуска тестов.
	progressBar := widget.NewProgressBar()
	startButton := widget.NewButton("Запустить все тесты", nil)

	// Создаем вертикальный контейнер для карточек с тестами.
	content := container.NewVBox()
	for _, t := range tests {
		content.Add(widget.NewCard(t.name, t.description, t.resultLabel))
	}
	for _, t := range flashTests {
		content.Add(widget.NewCard(t.name, t.description, t.resultLabel))
	}

	// Определяем действие при нажатии на кнопку "Запустить все тесты".
	startButton.OnTapped = func() {
		// Запускаем тесты в отдельной горутине, чтобы не блокировать UI.
		go func() {
			// Блокируем кнопку на время выполнения тестов.
			startButton.Disable()
			defer startButton.Enable()

			// Канал для отслеживания общего прогресса.
			progressChan := make(chan float64)
			go func() {
				for p := range progressChan {
					progressBar.SetValue(p)
				}
			}()

			totalProgress := 0.0
			numTests := float64(len(tests) + len(flashTests))

			// Последовательно выполняем тесты задержки.
			for _, t := range tests {
				t.resultLabel.SetText("Выполняется...")
				individualProgress := make(chan float64)
				go func() {
					for p := range individualProgress {
						progressChan <- (totalProgress + p) / numTests
					}
				}()

				stats := benchmark.TestRunner(t.fn, iterations, individualProgress)
				close(individualProgress)
				t.resultLabel.SetText(fmt.Sprintf("Результат: %.2f %s", stats.Avg, stats.Unit))
				totalProgress++
			}

			// Создаем временный файл для тестов флеш-памяти.
			testFile, err := memory.CreateTestFile()
			if err != nil {
				log.Println("Failed to create test file:", err)
				return
			}
			defer memory.CleanupTestFile(testFile)

			// Последовательно выполняем тесты флеш-памяти.
			for _, t := range flashTests {
				t.resultLabel.SetText("Выполняется...")
				// Создаем обертку, так как функция теста флеш-памяти имеет другую сигнатуру.
				flashFn := func() (float64, string) {
					val, unit, err := t.fn(testFile)
					if err != nil {
						log.Printf("Flash test error: %v", err)
						return 0, ""
					}
					return val, unit
				}
				individualProgress := make(chan float64)
				go func() {
					for p := range individualProgress {
						progressChan <- (totalProgress + p) / numTests
					}
				}()
				stats := benchmark.TestRunner(flashFn, iterations, individualProgress)
				close(individualProgress)
				t.resultLabel.SetText(fmt.Sprintf("Результат: %.2f %s", stats.Avg, stats.Unit))
				totalProgress++
			}

			close(progressChan)
			progressBar.SetValue(1.0) // Устанавливаем прогресс в 100% по завершении.
		}()
	}

	// Возвращаем контейнер с прокруткой, содержащий все элементы панели.
	return container.NewScroll(container.NewVBox(
		startButton,
		progressBar,
		content,
	))
}

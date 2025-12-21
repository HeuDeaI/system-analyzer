package gui

import (
	"fmt"
	"system-analyzer/pkg/benchmark"
	"system-analyzer/pkg/concurrency"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewConcurrencyPanel создает и возвращает панель для тестирования производительности механизмов параллелизма.
func NewConcurrencyPanel() fyne.CanvasObject {
	// Количество итераций для каждого теста.
	const iterations = 3

	// Слайс, определяющий тесты для этой панели.
	tests := []struct {
		name        string
		description string
		fn          benchmark.BenchmarkFunc
		resultLabel *widget.Label
	}{
		{"🚀 Накладные расходы на горутины", "Измеряет время, необходимое для создания и завершения большого количества горутин.", concurrency.GoroutineOverheadBenchmark, widget.NewLabel("Результат: ...")},
		{"🔒 Блокировка мьютекса", "Тестирует производительность синхронизации с использованием мьютексов для защиты общего ресурса.", concurrency.MutexBenchmark, widget.NewLabel("Результат: ...")},
		{"⚛️ Атомарные операции", "Тестирует производительность атомарных операций, которые являются альтернативой мьютексам.", concurrency.AtomicBenchmark, widget.NewLabel("Результат: ...")},
		{"📡 Пропускная способность каналов", "Измеряет, как быстро данные могут быть отправлены и получены через каналы Go.", concurrency.ChannelBenchmark, widget.NewLabel("Результат: ...")},
		{"⛓️ Конвейерная обработка", "Оценивает производительность многоступенчатого конвейера, где каждая стадия выполняется в отдельной горутине.", concurrency.PipelineBenchmark, widget.NewLabel("Результат: ...")},
	}

	// Создаем виджеты для отображения прогресса и запуска тестов.
	progressBar := widget.NewProgressBar()
	startButton := widget.NewButton("Запустить все тесты", nil)

	// Создаем вертикальный контейнер для карточек с тестами.
	content := container.NewVBox()
	for _, t := range tests {
		card := widget.NewCard(t.name, t.description, t.resultLabel)
		content.Add(card)
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
			numTests := float64(len(tests))

			// Последовательно выполняем все тесты.
			for _, t := range tests {
				t.resultLabel.SetText("Выполняется...")

				// Канал для отслеживания прогресса отдельного теста.
				individualProgress := make(chan float64)
				go func() {
					for p := range individualProgress {
						// Обновляем общий прогресс на основе прогресса текущего теста.
						progressChan <- (totalProgress + p) / numTests
					}
				}()

				// Запускаем тест и получаем статистику.
				stats := benchmark.TestRunner(t.fn, iterations, individualProgress)
				close(individualProgress)

				// Обновляем метку с результатом.
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

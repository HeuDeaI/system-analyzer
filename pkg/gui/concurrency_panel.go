package gui

import (
	"fmt"
	"system-analyzer/pkg/benchmark"
	"system-analyzer/pkg/concurrency"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewConcurrencyPanel() fyne.CanvasObject {
	const iterations = 3

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

	progressBar := widget.NewProgressBar()
	startButton := widget.NewButton("Запустить все тесты", nil)

	content := container.NewVBox()
	for _, t := range tests {
		card := widget.NewCard(t.name, t.description, t.resultLabel)
		content.Add(card)
	}

	startButton.OnTapped = func() {
		go func() {
			startButton.Disable()
			defer startButton.Enable()

			progressChan := make(chan float64)
			go func() {
				for p := range progressChan {
					progressBar.SetValue(p)
				}
			}()

			totalProgress := 0.0
			numTests := float64(len(tests))

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

			close(progressChan)
			progressBar.SetValue(1.0)
		}()
	}

	return container.NewScroll(container.NewVBox(
		startButton,
		progressBar,
		content,
	))
}

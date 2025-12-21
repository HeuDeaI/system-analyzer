package gui

import (
	"fmt"
	"system-analyzer/pkg/benchmark"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewBenchmarkPanel() fyne.CanvasObject {
	const iterations = 3

	tests := []struct {
		name        string
		description string
		fn          benchmark.BenchmarkFunc
		resultLabel *widget.Label
	}{
		{"⚙️ Целочисленная арифметика", "Тестирует скорость выполнения базовых целочисленных операций (сложение, вычитание, умножение, деление).", benchmark.IntegerBenchmark, widget.NewLabel("Результат: ...")},
		{"✨ Арифметика с плавающей запятой", "Тестирует скорость выполнения операций с плавающей запятой, которые важны для научных вычислений и графики.", benchmark.FloatBenchmark, widget.NewLabel("Результат: ...")},
		{"📥 Пропускная способность чтения из ОЗУ", "Измеряет скорость последовательного чтения данных из оперативной памяти.", benchmark.ReadBandwidthBenchmark, widget.NewLabel("Результат: ...")},
		{"📤 Пропускная способность записи в ОЗУ", "Измеряет скорость последовательной записи данных в оперативную память.", benchmark.WriteBandwidthBenchmark, widget.NewLabel("Результат: ...")},
		{"🔄 Случайный доступ к ОЗУ", "Тестирует скорость доступа к случайным ячейкам памяти, что критично для баз данных и других приложений.", benchmark.RandomBandwidthBenchmark, widget.NewLabel("Результат: ...")},
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

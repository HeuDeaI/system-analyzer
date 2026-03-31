package gui

import (
	"fmt"
	"runtime"
	"system-analyzer/pkg/benchmark"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type BenchmarkTask struct {
	Name                 string
	Category             string
	Description          string
	Scientific           string
	ResultInterpretation func(float64) string // Функция для формирования описания результата
	Fn                   benchmark.BenchmarkFunc
	ResultLabel          *widget.Label
}

func NewBenchmarkPanel(category string) fyne.CanvasObject {
	allTasks := []BenchmarkTask{
		{
			Name:        "QuickSort",
			Category:    "Процессор",
			Description: "Сортировка массива (1 млн элементов).",
			Scientific:  "Тест эффективности кеша и предсказателя ветвлений. Показывает, насколько хорошо ЦП справляется с нелинейным выполнением кода.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.2f млн/с | 1 млн за %.2f сек | O(n log n)", result, 1/result)
			},
			Fn: benchmark.SortingBenchmark,
		},
		{
			Name:        "Concurrency",
			Category:    "Процессор",
			Description: "Параллельный расчет числа Пи (Монте-Карло).",
			Scientific:  "Тест масштабируемости системы. Показывает эффективность работы планировщика на всех доступных ядрах ЦП одновременно.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.2f млн/с | Ядро: %.0f тыс/с | Линейная масштабируемость", result, result/float64(runtime.NumCPU()))
			},
			Fn: benchmark.ConcurrencyBenchmark,
		},
		{
			Name:        "Fibonacci",
			Category:    "Процессор",
			Description: "Рекурсивное вычисление чисел Фибоначчи.",
			Scientific:  "Тест глубины стека и скорости вызова функций. Оценивает накладные расходы на рекурсию.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.2f оп/с | fib(35)", result)
			},
			Fn: benchmark.FibonacciBenchmark,
		},
		{
			Name:        "RAM Copy",
			Category:    "Память",
			Description: "Копирование больших блоков памяти (128 МБ).",
			Scientific:  "Измеряет чистую пропускную способность контроллера памяти (Memory Bandwidth). Критично для обработки видео и больших баз данных.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.2f ГБ/с | DDR4-3200: ~25-51 ГБ/с", result)
			},
			Fn: benchmark.RAMCopyBenchmark,
		},
		{
			Name:        "RAM Latency",
			Category:    "Память",
			Description: "Задержка доступа к памяти.",
			Scientific:  "Измеряет время доступа к случайным ячейкам памяти. Чем ниже задержка, тем быстрее работают приложения с нерегулярным доступом к данным.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.2f нс | Меньше - лучше", result)
			},
			Fn: benchmark.RAMLatencyBenchmark,
		},
		{
			Name:        "Disk I/O",
			Category:    "Память",
			Description: "Последовательная запись и чтение файла.",
			Scientific:  "Оценивает пропускную способность файловой системы и накопителя. Важно для баз данных и систем логирования.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.0f МБ/с | SSD SATA: ~500 | NVMe: ~3000-7000", result)
			},
			Fn: benchmark.DiskIOBenchmark,
		},
		{
			Name:        "Matrix Multiplication",
			Category:    "Математика",
			Description: "Умножение матриц (256x256).",
			Scientific:  "Классический тест производительности FPU (Floating Point Unit). Важен для графики, ИИ и научных расчетов.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.2f MFLOPS | Операции с плавающей точкой", result)
			},
			Fn: benchmark.MatrixMultiplicationBenchmark,
		},
		{
			Name:        "BigInt Math",
			Category:    "Математика",
			Description: "Арифметика сверхбольших чисел.",
			Scientific:  "Тестирует блоки умножения и деления процессора на числах произвольной точности. Важен для глубоких криптографических исследований.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.0f оп/с | Операции с 10^18", result)
			},
			Fn: benchmark.LargeIntArithmeticBenchmark,
		},
		{
			Name:        "Prime Search",
			Category:    "Математика",
			Description: "Поиск простых чисел (Miller-Rabin).",
			Scientific:  "Имитирует процесс генерации ключей безопасности. Требует интенсивных целочисленных вычислений и качественного энтропийного источника.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.2f чисел/с | RSA-2048: 2 числа | %.1f сек на ключ", result, 2/result)
			},
			Fn: benchmark.PrimeSearchBenchmark,
		},
		{
			Name:        "SHA-256",
			Category:    "Криптография",
			Description: "Тест хеширования SHA-256 (Secure Hash Algorithm 256-bit).",
			Scientific:  "Измеряет производительность логических операций (AND, XOR, ROT) над 32-битными словами. Важен для блокчейна и проверки целостности.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.0f МБ/с | %.0f МБ/сек | %.0f блоков/с", result, result, result*16)
			},
			Fn: benchmark.CryptoSHA256Benchmark,
		},
		{
			Name:        "SHA-512",
			Category:    "Криптография",
			Description: "Тест хеширования SHA-512.",
			Scientific:  "Использует 64-битные слова. На 64-битных системах может быть эффективнее SHA-256. Критичен для систем повышенной безопасности.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.0f МБ/с | Эффективнее SHA-256 на 64-бит CPU", result)
			},
			Fn: benchmark.CryptoSHA512Benchmark,
		},
		{
			Name:        "AES-256-GCM (Шифрование)",
			Category:    "Криптография",
			Description: "Шифрование AES-256 в режиме GCM.",
			Scientific:  "Проверяет наличие AES-NI инструкций. Используется в 90% трафика интернета (TLS/SSL). Прямо влияет на скорость VPN и защищенных соединений.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.0f МБ/с | 1 ГБ за %.1f сек | %.0f Мбит/с", result, 1024/result, result*8)
			},
			Fn: benchmark.CryptoAESBenchmark,
		},
		{
			Name:        "AES-256-GCM (Расшифрование)",
			Category:    "Криптография",
			Description: "Имитация времени расшифрования данных алгоритмом AES-256.",
			Scientific:  "Оценивает скорость обратной операции шифрования. Критично для серверов, обрабатывающих входящие защищенные потоки данных.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.0f МБ/с | ~%.0f клиентов на 100 Мбит/с", result, result*8/100)
			},
			Fn: benchmark.CryptoAESDecryptionBenchmark,
		},
		{
			Name:        "Ed25519",
			Category:    "Криптография",
			Description: "Подпись на эллиптических кривых.",
			Scientific:  "Тест производительности арифметики на кривой Ed25519. Важен для протоколов SSH и современных систем цифровой подписи.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.0f верификаций/с | %.0f SSH-соединений/с", result, result/2)
			},
			Fn: benchmark.CryptoEd25519Benchmark,
		},
		{
			Name:        "RSA-2048 (Расшифрование)",
			Category:    "Криптография",
			Description: "Имитация времени расшифрования RSA-2048.",
			Scientific:  "Оценивает скорость операций с закрытым ключом. Это самая ресурсоемкая операция при установке защищенного соединения HTTPS.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.0f оп/с | %.0f TLS рукопожатий/с | %.0f HTTPS коннектов/с", result, result, result/2)
			},
			Fn: benchmark.CryptoRSADecryptionBenchmark,
		},
		{
			Name:        "RSA-2048 (Генерация ключей)",
			Category:    "Криптография",
			Description: "Генерация новых пар ключей RSA-2048.",
			Scientific:  "Тестирует производительность генерации простых чисел и вычислений модульной арифметики. Критично для PKI и сертификатов.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.2f ключей/с | 1 ключ за %.1f сек | %.0f сертификатов/час", result, 1/result, result*3600)
			},
			Fn: benchmark.CryptoRSAKeyGenBenchmark,
		},
		{
			Name: "RSA-2048 (Оценка brute-force)", Category: "Криптография",
			Description: "Теоретическая оценка времени полного перебора ключа RSA-2048.",
			Scientific:  "Рассчитывает время подбора ключа методом грубой силы на основе измеренной производительности. RSA-2048 считается квантово-устойчивым до ~2030 года.",
			ResultInterpretation: func(result float64) string {
				return benchmark.FormatDuration(result)
			},
			Fn: benchmark.RSABruteForceEstimate,
		},
		{
			Name: "AES-128 (Оценка brute-force)", Category: "Криптография",
			Description: "Теоретическая оценка времени перебора ключа AES-128.",
			Scientific:  "AES-128 имеет 2^128 возможных ключей. Демонстрирует абсолютную криптографическую стойкость к классическому перебору.",
			ResultInterpretation: func(result float64) string {
				return benchmark.FormatDuration(result)
			},
			Fn: benchmark.AES128BruteForceEstimate,
		},
		{
			Name: "Шифр Цезаря (Перебор)", Category: "Криптография",
			Description: "Время перебора всех 25 ключей шифра Цезаря.",
			Scientific:  "Простейший моноалфавитный шифр. Показывает насколько быстро взламываются примитивные алгоритмы шифрования.",
			ResultInterpretation: func(result float64) string {
				return benchmark.FormatDuration(result)
			},
			Fn: benchmark.CaesarCipherBruteForce,
		},
		{
			Name: "XOR 8-бит (Перебор)", Category: "Криптография",
			Description: "Время перебора 256 вариантов 8-битного XOR ключа.",
			Scientific:  "Одноразовый блокнот с коротким ключом. Демонстрирует уязвимость XOR при недостаточной длине ключа.",
			ResultInterpretation: func(result float64) string {
				return benchmark.FormatDuration(result)
			},
			Fn: benchmark.XOR8BitBruteForce,
		},
		{
			Name: "4-значный PIN (Перебор)", Category: "Криптография",
			Description: "Время перебора 4-значного PIN кода (0000-9999).",
			Scientific:  "Имитирует атаку на простые PIN коды. 10000 комбинаций подбираются мгновенно даже на слабом процессоре.",
			ResultInterpretation: func(result float64) string {
				return benchmark.FormatDuration(result)
			},
			Fn: benchmark.PIN4DigitBruteForce,
		},
		{
			Name:        "Disk I/O",
			Category:    "Хранилище",
			Description: "Последовательная запись и чтение файла.",
			Scientific:  "Оценивает пропускную способность файловой системы и накопителя. Важно для баз данных и систем логирования.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.0f МБ/с | SSD SATA: ~500 | NVMe: ~3000-7000", result)
			},
			Fn: benchmark.DiskIOBenchmark,
		},
		{
			Name:        "String Ops",
			Category:    "Данные",
			Description: "Обработка текстовых данных (Replace, Split).",
			Scientific:  "Типичная нагрузка веб-серверов и текстовых анализаторов. Тестирует работу с динамической памятью и строковыми примитивами.",
			ResultInterpretation: func(result float64) string {
				return fmt.Sprintf("%.0f оп/с | %.0f МБ текста/с | %.0f JSON/с", result, result/1000, result/10)
			},
			Fn: benchmark.StringManipulationBenchmark,
		},
	}

	var tasks []BenchmarkTask
	for _, t := range allTasks {
		if category == "" || t.Category == category {
			t.ResultLabel = widget.NewLabel("Ожидание")
			t.ResultLabel.TextStyle = fyne.TextStyle{Monospace: true}
			tasks = append(tasks, t)
		}
	}

	progressBar := widget.NewProgressBar()
	btnText := fmt.Sprintf("Запустить тесты: %s", category)
	if category == "" {
		btnText = "Запустить все тесты"
	}
	startButton := widget.NewButtonWithIcon(btnText, theme.MediaPlayIcon(), nil)

	content := container.NewVBox()
	for i := range tasks {
		t := &tasks[i]
		descLabel := widget.NewLabel(fmt.Sprintf("Описание: %s\n\nДетали: %s", t.Description, t.Scientific))
		descLabel.Wrapping = fyne.TextWrapWord
		t.ResultLabel.Wrapping = fyne.TextWrapWord
		content.Add(widget.NewCard(t.Name, t.Category, container.NewVBox(descLabel, t.ResultLabel)))
	}

	startButton.OnTapped = func() {
		go func() {
			startButton.Disable()
			defer startButton.Enable()
			for i := range tasks {
				t := &tasks[i]
				t.ResultLabel.SetText("Вычисление...")
				stats := benchmark.TestRunner(t.Fn, 10, nil)
				resultText := fmt.Sprintf("%.2f %s (σ=%.1f%%)", stats.Avg, stats.Unit, stats.Variance)
				if t.ResultInterpretation != nil {
					resultText = t.ResultInterpretation(stats.Avg)
				}
				t.ResultLabel.SetText(resultText)
				progressBar.SetValue(float64(i+1) / float64(len(tasks)))
			}
		}()
	}

	return container.NewPadded(container.NewVBox(
		startButton, progressBar, widget.NewSeparator(), content,
	))
}

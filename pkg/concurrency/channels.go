package concurrency

import (
	"sync"
	"time"
)

const channelOps = 2000000

// ChannelBenchmark измеряет пропускную способность каналов.
func ChannelBenchmark() (float64, string) {
	ch := make(chan int, 100)
	var wg sync.WaitGroup
	wg.Add(2)

	start := time.Now()
	go func() {
		defer wg.Done()
		for i := 0; i < channelOps; i++ {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		defer wg.Done()
		for range ch {
			// Просто читаем из канала
		}
	}()

	wg.Wait()
	elapsed := time.Since(start)
	return float64(elapsed.Nanoseconds()) / float64(channelOps), "ns/op"
}

const pipelineStages = 5
const pipelineItems = 200000

// PipelineBenchmark измеряет производительность конвейерной обработки.
func PipelineBenchmark() (float64, string) {
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(pipelineStages)

	// Создаем каналы для стадий
	chans := make([]chan int, pipelineStages+1)
	for i := range chans {
		chans[i] = make(chan int, 100)
	}

	// Запускаем стадии
	for i := 0; i < pipelineStages; i++ {
		go func(in <-chan int, out chan<- int) {
			defer wg.Done()
			for item := range in {
				// Симулируем обработку
				out <- item * 2
			}
			close(out)
		}(chans[i], chans[i+1])
	}

	// Отправляем данные в первую стадию
	go func() {
		for i := 0; i < pipelineItems; i++ {
			chans[0] <- i
		}
		close(chans[0])
	}()

	// Потребляем результаты из последней стадии
	for range chans[pipelineStages] {
	}

	wg.Wait() // Ожидаем завершения последней стадии
	for i := 0; i < pipelineItems; i++ {
		<-chans[pipelineStages]
	}

	elapsed := time.Since(start)
	return float64(elapsed.Nanoseconds()) / float64(pipelineItems), "ns/op"
}

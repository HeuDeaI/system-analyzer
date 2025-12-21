package memory

import (
	"math/rand"
	"runtime"
	"sort"
	"time"
	"unsafe"
)

type node struct {
	next *node
}

func makeRandomList(size int) *node {
	nodes := make([]node, size)

	perm := rand.Perm(size)
	for i := 0; i < size-1; i++ {
		nodes[perm[i]].next = &nodes[perm[i+1]]
	}
	nodes[perm[size-1]].next = &nodes[perm[0]]

	return &nodes[perm[0]]
}

func chase(start *node, steps int) *node {
	p := start
	for i := 0; i < steps; i++ {
		p = p.next
	}
	return p
}

func measure(sizeBytes int, steps int, repeats int) float64 {
	size := sizeBytes / int(unsafe.Sizeof(node{}))

	results := make([]float64, repeats)

	for r := 0; r < repeats; r++ {
		head := makeRandomList(size)

		start := time.Now()
		_ = chase(head, steps)
		elapsed := time.Since(start)

		results[r] = float64(elapsed.Nanoseconds()) / float64(steps)
	}

	sort.Float64s(results)
	return results[repeats/2]
}

func runLatencyTest(sizeBytes int) (float64, string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	steps := 5_000_000
	repeats := 3

	ns := measure(sizeBytes, steps, repeats)
	return ns, "ns/доступ"
}

func L1CacheLatencyBenchmark() (float64, string) {
	return runLatencyTest(64 * 1024) // 64 KB
}

func L2CacheLatencyBenchmark() (float64, string) {
	return runLatencyTest(2 * 1024 * 1024) // 2 MB
}

func L3CacheLatencyBenchmark() (float64, string) {
	return runLatencyTest(16 * 1024 * 1024) // 16 MB
}

func RAMLatencyBenchmark() (float64, string) {
	return runLatencyTest(128 * 1024 * 1024) // 128 MB
}

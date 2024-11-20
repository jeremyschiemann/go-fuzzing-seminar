package increment

import (
	"sync"
	"testing"
	"time"
)

// SharedResource simulates a resource protected by a mutex.
type SharedResource struct {
	Value int
	mu    sync.Mutex
}

// Increment safely increments the value but introduces a subtle delay
// that may cause deadlock in tests with specific goroutine schedules.
func (r *SharedResource) Increment() {
	r.mu.Lock()
	defer r.mu.Unlock()
	time.Sleep(10 * time.Millisecond) // Simulate subtle contention
	r.Value++
}

func TestConcurrentIncrement(t *testing.T) {
	resource := &SharedResource{}
	const numWorkers = 4
	const incrementsPerWorker = 10

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Start workers that increment the shared resource.
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < incrementsPerWorker; j++ {
				// Introduce a blocking call that could lead to deadlock if workers are unsynchronized.
				if workerID == 1 && j == incrementsPerWorker/2 {
					// Worker 1 intentionally delays mid-workload.
					time.Sleep(100 * time.Millisecond)
				}
				resource.Increment()
			}
		}(i)
	}

	// Wait for all workers to finish.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for a bounded time to detect potential deadlocks.
	select {
	case <-done:
		// Verify the result after all increments.
		expected := numWorkers * incrementsPerWorker
		if resource.Value != expected {
			t.Errorf("Expected %d, got %d", expected, resource.Value)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Test timed out due to potential deadlock or blocking behavior")
	}
}

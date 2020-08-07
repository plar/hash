package pool

// TBD
import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	assert := assert.New(t)

	taskQueue := make(chan Task, 10)
	dispatcher := NewDispatcher(taskQueue, 8)
	dispatcher.Run()

	var wg sync.WaitGroup
	var m sync.Map

	total := 10000

	wg.Add(total)
	for i := 0; i < total; i++ {
		dispatcher.Dispatch(func(id int64) Task {
			return Task{
				ID: int64(i),
				Handler: func() {
					m.Store(id, id)
					wg.Done()
				},
			}
		}(int64(i)))
	}

	wg.Wait()

	for i := 0; i < total; i++ {
		id := int64(i)
		mv, ok := m.Load(id)
		assert.True(ok)
		assert.Equal(id, mv)
	}

	dispatcher.Stop()
}

func BenchmarkPoolDispatch(b *testing.B) {
	taskQueue := make(chan Task, 10)
	dispatcher := NewDispatcher(taskQueue, 8)
	dispatcher.Run()

	var wg sync.WaitGroup
	var m sync.Map

	total := b.N

	b.ResetTimer()

	wg.Add(total)
	for i := 0; i < total; i++ {
		dispatcher.Dispatch(func(id int64) Task {
			return Task{
				ID: int64(i),
				Handler: func() {
					m.Store(id, id)
					wg.Done()
				},
			}
		}(int64(i)))
	}

	wg.Wait()
	dispatcher.Stop()
}

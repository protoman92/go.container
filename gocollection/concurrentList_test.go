package gocollection

import (
	"testing"
)

func testListConcurrentOps(tb testing.TB, cl ConcurrentCollection) {
	params := &ConcurrentCollectionOpsParams{
		concurrentCollection: cl,
		log:                  false,
		keyCount:             500,
	}

	setupConcurrentCollectionOps(params)
}

func benchmarkListConcurrentOps(b *testing.B, clFn func() ConcurrentCollection) {
	for i := 0; i < b.N; i++ {
		testListConcurrentOps(b, clFn())
	}
}

func BenchmarkLockListConcurrentOps(b *testing.B) {
	benchmarkListConcurrentOps(b, func() ConcurrentCollection {
		sl := NewSliceList()
		return NewLockConcurrentList(sl)
	})
}

func TestLockListConcurrentOps(t *testing.T) {
	t.Parallel()
	sl := NewSliceList()
	cl := NewLockConcurrentList(sl)
	testListConcurrentOps(t, cl)
}

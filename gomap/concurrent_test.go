package gomap

import (
	"testing"
	"time"
)

func testConcurrentMapConcurrentOps(tb testing.TB, cm ConcurrentMap) {
	params := &ConcurrentMapOpsParams{
		concurrentMap:     cm,
		log:               false,
		keyCount:          500,
		setupWaitDuration: time.Second,
	}

	SetupConcurrentMapOps(params)
}

func benchmarkConcurrentMapConcurrentOps(b *testing.B, cmFn func() ConcurrentMap) {
	for i := 0; i < b.N; i++ {
		testConcurrentMapConcurrentOps(b, cmFn())
	}
}

func BenchmarkChannelConcurrentMapConcurrentOps(b *testing.B) {
	benchmarkConcurrentMapConcurrentOps(b, func() ConcurrentMap {
		bm := NewDefaultBasicMap()
		return NewChannelConcurrentMap(bm)
	})
}

func BenchmarkLockConcurrentMapConcurrentOps(b *testing.B) {
	benchmarkConcurrentMapConcurrentOps(b, func() ConcurrentMap {
		bm := NewDefaultBasicMap()
		return NewLockConcurrentMap(bm)
	})
}

func TestChannelConcurrentMapConcurrentOps(t *testing.T) {
	t.Parallel()
	bm := NewDefaultBasicMap()
	cm := NewChannelConcurrentMap(bm)
	testConcurrentMapConcurrentOps(t, cm)
}

func TestLockConcurrentMapConcurrentOps(t *testing.T) {
	t.Parallel()
	bm := NewDefaultBasicMap()
	cm := NewLockConcurrentMap(bm)
	testConcurrentMapConcurrentOps(t, cm)
}

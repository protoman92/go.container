package gomap

import (
	"fmt"
	"testing"
	"time"
)

func testConcurrentMapBasicOps(tb testing.TB, cm ConcurrentMap) {
	/// Setup
	key := "Key"
	value := "Value"

	/// When & Then
	cm.Set(key, value)

	if !cm.Contains(key) {
		tb.Errorf("Should contain %v", key)
	}

	getValue1, getFound1 := cm.Get(key)

	if getValue1 != value || !getFound1 {
		tb.Errorf("Should contain %v with value %v", key, value)
	}

	deletedFound := cm.Delete(key)
	getValue2, getFound2 := cm.Get(key)

	if (getFound2 && cm.Contains(key)) || getValue2 != nil || !deletedFound {
		tb.Errorf("Should not contain %v", key)
	}

	setPrev, setFound := cm.Set(key, value)
	length := cm.Length()

	if length != 1 || setPrev != nil || setFound {
		tb.Errorf("Should have length 1, but got %d", length)
	}

	cm.Clear()

	if !cm.IsEmpty() {
		tb.Errorf("Should be empty")
	}

	fmt.Printf("Final storage: %v\n", cm.UnderlyingStorage())
}

func TestChannelConcurrentMapBasicOps(t *testing.T) {
	t.Parallel()
	bm := NewDefaultBasicMap()
	cm := NewChannelConcurrentMap(bm)
	testConcurrentMapBasicOps(t, cm)
}

func TestLockConcurrentMapBasicOps(t *testing.T) {
	t.Parallel()
	bm := NewDefaultBasicMap()
	cm := NewLockConcurrentMap(bm)
	testConcurrentMapBasicOps(t, cm)
}

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

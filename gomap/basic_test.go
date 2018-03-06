package gomap

import (
	"fmt"
	"testing"
)

func testMapBasicOps(tb testing.TB, m Map) {
	/// Setup
	key := "Key"
	value := "Value"

	/// When & Then
	m.Set(key, value)

	if !m.Contains(key) {
		tb.Errorf("Should contain %v", key)
	}

	getValue1, getFound1 := m.Get(key)

	if getValue1 != value || !getFound1 {
		tb.Errorf("Should contain %v with value %v", key, value)
	}

	deletedFound := m.Delete(key)
	getValue2, getFound2 := m.Get(key)

	if (getFound2 && m.Contains(key)) || getValue2 != nil || !deletedFound {
		tb.Errorf("Should not contain %v", key)
	}

	setPrev, setFound := m.Set(key, value)
	length := m.Length()

	if length != 1 || setPrev != nil || setFound {
		tb.Errorf("Should have length 1, but got %d", length)
	}

	m.Clear()

	if !m.IsEmpty() {
		tb.Errorf("Should be empty")
	}

	fmt.Printf("Final storage: %v\n", m.UnderlyingStorage())
}

func testMapAllOps(t *testing.T, mapFn func() Map) {
	testMapBasicOps(t, mapFn())
}

func testBasicMapAllOps(t *testing.T) {
	t.Parallel()

	testMapAllOps(t, func() Map {
		return NewDefaultBasicMap()
	})
}

func TestChannelConcurrentMapAllOps(t *testing.T) {
	t.Parallel()

	testMapAllOps(t, func() Map {
		bm := NewDefaultBasicMap()
		return NewChannelConcurrentMap(bm)
	})
}

func TestLockConcurrentMapAllOps(t *testing.T) {
	t.Parallel()

	testMapAllOps(t, func() Map {
		bm := NewDefaultBasicMap()
		return NewLockConcurrentMap(bm)
	})
}

package gomap

import (
	"fmt"
	"testing"

	gl "github.com/protoman92/gocontainer/pkg/gocollection"
)

func testMapBasicOps(t *testing.T, m Map) {
	/// Setup
	key := "Key"
	value := "Value"

	/// When & Then
	m.Set(key, value)

	if !m.Contains(key) {
		t.Errorf("Should contain %v", key)
	}

	if val, found := m.Get(key); val != value || !found {
		t.Errorf("Should contain %v with value %v", key, value)
	}

	if prev, found := m.Delete(key); !found || prev == nil {
		t.Errorf("Should delete key")
	}

	if val, found := m.Get(key); found || val != nil {
		t.Errorf("Should not contain %v", key)
	}

	if prev, found := m.Set(key, value); found || prev != nil {
		t.Errorf("Should not have any previous value")
	}

	if length := m.Length(); length != 1 {
		t.Errorf("Should have length 1, but got %d", length)
	}

	m.Clear()

	if length := m.Length(); length > 0 {
		t.Errorf("Should not contain anything")
	}

	fmt.Printf("Final map %v\n", m)
}

func testMapKeys(t *testing.T, m Map) {
	/// Setup
	keys := []interface{}{1, 2, 3, 4, 5}

	/// When
	for ix := range keys {
		key := keys[ix]
		m.Set(key, key)
	}

	/// Then
	mapKeys := m.Keys()
	bl := gl.NewSliceList(mapKeys...)

	for ix := range keys {
		if contains := bl.Contains(keys[ix]); !contains {
			t.Errorf("Should contain key")
		}
	}
}

func testMapAllOps(t *testing.T, mapFn func() Map) {
	testMapBasicOps(t, mapFn())
	testMapKeys(t, mapFn())
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

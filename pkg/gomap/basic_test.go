package gomap

import (
	"fmt"
	"strconv"
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

	if getValue, getFound := m.Get(key); getValue != value || !getFound {
		t.Errorf("Should contain %v with value %v", key, value)
	}

	deletedFound := m.Delete(key)

	if getValue, getFound := m.Get(key); (getFound && m.Contains(key)) || getValue != nil || !deletedFound {
		t.Errorf("Should not contain %v", key)
	}

	setPrev, setFound := m.Set(key, value)

	if length := m.Length(); length != 1 || setPrev != nil || setFound {
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
	bl := gl.NewSliceList(mapKeys)

	for ix := range keys {
		key := strconv.Itoa(keys[ix].(int))

		if contains := bl.Contains(key); !contains {
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
package golist

import (
	"strconv"
	"testing"
)

func testListBasicOps(t *testing.T, list List) {
	/// Setup & When & Then
	add1 := list.Add(1)
	add2 := list.Add(2)
	add3 := list.Add(3)
	add4 := list.Add(4)

	if add1 != 1 || add2 != 1 || add3 != 1 || add4 != 1 {
		t.Errorf("Added wrong element count")
	}

	getValue1, getFound1 := list.Get(0)
	getValue2, getFound2 := list.Get(list.Length())
	getValue3, getFound3 := list.Get(list.Length() - 1)

	if getValue1 != 1 || !getFound1 {
		t.Errorf("Wrong element")
	}

	if getValue2 != nil || getFound2 {
		t.Errorf("Should not have found element")
	}

	if getValue3 != 4 || !getFound3 {
		t.Errorf("Wrong element")
	}

	if list.Length() != 4 {
		t.Errorf("Should have 4 elements")
	}

	deletedFound1 := list.Remove(2)

	if list.Length() != 3 || !deletedFound1 {
		t.Errorf("Should have 3 elements")
	}

	slice := make([]Element, 1000)

	for ix := range slice {
		slice[ix] = strconv.Itoa(ix)
	}

	addAll1 := list.AddAll(slice...)

	if addAll1 != 1000 {
		t.Errorf("Added wrong element count")
	}

	if list.Length() != 1003 {
		t.Errorf("Should have 1003 elements")
	}

	deletedFound2 := list.Remove("Not existent")

	if deletedFound2 {
		t.Errorf("Should not have found element")
	}

	list.Clear()

	if list.Length() != 0 {
		t.Errorf("Should not have any element")
	}
}

func testListAllOps(t *testing.T, listFn func() List) {
	testListBasicOps(t, listFn())
}

func TestSliceListAllOps(t *testing.T) {
	t.Parallel()

	testListAllOps(t, func() List {
		return NewSliceList()
	})
}

func TestLockConcurrentListAllOps(t *testing.T) {
	t.Parallel()

	testListAllOps(t, func() List {
		sl := NewSliceList()
		return NewLockConcurrentList(sl)
	})
}

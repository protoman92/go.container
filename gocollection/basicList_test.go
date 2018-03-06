package gocollection

import (
	"fmt"
	"testing"
)

func testListBasicOps(t *testing.T, list List) {
	/// Setup & When & Then
	getValue1, getFound1 := list.Get(0)

	if getValue1 != nil || getFound1 {
		t.Errorf("Should not have any element")
	}

	newElements := []Element{1, 2, 3, 4}
	list.AddAll(newElements...)

	for ix := range newElements {
		e := newElements[ix]
		value, found := list.Get(ix)

		if value != e || !found {
			t.Errorf("Should have found element")
		}
	}

	fmt.Printf("Final list: %v\n", list)
}

func testListAllOps(t *testing.T, listFn func() List) {
	testCollectionBasicOps(t, listFn())
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

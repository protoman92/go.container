package gocollection

import (
	"fmt"
	"testing"
)

func testListBasicOps(t *testing.T, list List) {
	/// Setup & When & Then
	if getValue, getFound := list.GetAt(0); getValue != nil || getFound {
		t.Errorf("Should not have any element")
	}

	newElements := []Element{1, 2, 3, 4}

	if added := list.AddAll(newElements...); added != 4 {
		t.Errorf("Wrong number of elements added")
	}

	for ix := range newElements {
		e := newElements[ix]

		if value, found := list.GetAt(ix); value != e || !found {
			t.Errorf("Should have found element")
		}
	}

	fmt.Printf("Final list: %v\n", list)
}

func testListRemoveAt(t *testing.T, list List) {
	/// Setup & When & Then
	if _, found := list.RemoveAt(0); found {
		t.Errorf("Should not remove anything")
	}

	list.AddAll(1, 2, 3, 4, 5, 6)

	if e, found := list.RemoveAt(3); !found || e != 4 {
		t.Errorf("Should remove something")
	}

	if e, found := list.RemoveAt(0); !found || e != 1 {
		t.Errorf("Should remove something")
	}

	if e, found := list.RemoveAt(3); !found || e != 6 {
		t.Errorf("Should remove something")
	}

	if removed := list.RemoveAllAt(0); removed != 1 {
		t.Errorf("Removed wrong number of elements")
	}

	if removed := list.RemoveAllAt(-1, 1000); removed != 0 {
		t.Errorf("Should not remove anything")
	}

	if contains := list.Contains(2); contains {
		t.Errorf("Should not contain element")
	}
}

func testListSetAt(t *testing.T, list List) {
	/// Setup & When
	if prev, found := list.SetAt(0, "Unable to set"); prev != nil || found {
		t.Errorf("Should not be set")
	}

	list.AddAll(1, 2, 3, 4, 5)

	if prev, found := list.SetAt(4, 6); !found || prev != 5 {
		t.Errorf("Set wrong element")
	}

	if prev, found := list.SetAt(1000, 1000); prev != nil || found {
		t.Errorf("Should not be set")
	}
}

func testListAllOps(t *testing.T, listFn func() List) {
	testCollectionAllOps(t, func() Collection {
		return listFn()
	})

	testListBasicOps(t, listFn())
	testListRemoveAt(t, listFn())
	testListSetAt(t, listFn())
}

func TestSliceListAllOps(t *testing.T) {
	t.Parallel()

	testListAllOps(t, func() List {
		return NewSliceList()
	})
}

func TestLockConcurrentSliceListAllOps(t *testing.T) {
	t.Parallel()

	testListAllOps(t, func() List {
		sl := NewSliceList()
		return NewLockConcurrentList(sl)
	})
}

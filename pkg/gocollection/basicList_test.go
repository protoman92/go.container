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

	newElements := []interface{}{1, 2, 3, 4}

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

func testListGetFirst(t *testing.T, list List) {
	/// Setup & When & Then
	if first, ok := list.GetFirst(); ok || first != nil {
		t.Errorf("Should get correct element")
	}

	list.AddAll(1, 2, 3, 4, 5, 6, 7)

	if first, ok := list.GetFirst(); !ok || first != 1 {
		t.Errorf("Should get correct element")
	}

	if first, ok := list.GetFirstFunc(func(e interface{}) bool {
		return e == 2
	}); !ok || first != 2 {
		t.Errorf("Should get correct element")
	}

	if first, ok := list.GetFirstFunc(func(e interface{}) bool {
		return e == "Should not exist"
	}); ok || first != nil {
		t.Errorf("Should get correct element")
	}
}

func testListIndexOf(t *testing.T, list List) {
	/// Setup & When & Then
	if index, found := list.IndexOf(0); found || index >= 0 {
		t.Errorf("Should not have found index")
	}

	elements := []interface{}{1, 2, 3, 4, 5}
	list.AddAll(elements...)

	for ix := range elements {
		element := elements[ix]

		if index, found := list.IndexOf(element); !found || index < 0 {
			t.Errorf("Should have found index")
		}
	}
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
	testListGetFirst(t, listFn())
	testListIndexOf(t, listFn())
	testListRemoveAt(t, listFn())
	testListSetAt(t, listFn())
}

func TestSliceListAllOps(t *testing.T) {
	t.Parallel()

	testListAllOps(t, func() List {
		return NewDefaultSliceList()
	})
}

func TestLockConcurrentSliceListAllOps(t *testing.T) {
	t.Parallel()

	testListAllOps(t, func() List {
		sl := NewDefaultSliceList()
		return NewLockConcurrentList(sl)
	})
}

func TestSliceListWithExistingList(t *testing.T) {
	/// Setup
	array1 := [4]interface{}{1, 2, 3, 4}
	slice1 := array1[:]
	list := NewSliceList(slice1...)

	/// When
	array1[3] = 5

	/// Then
	if element := slice1[3]; element != 5 {
		t.Errorf("Should have changed element")
	}

	if contains := list.Contains(5); contains {
		t.Errorf("Should not retain array reference")
	}

	for _, e := range slice1[:2] {
		if contains := list.Contains(e); !contains {
			t.Errorf("Should contain element")
		}
	}
}

func TestSliceRange(t *testing.T) {
	/// Setup
	list1 := NewSliceListForRange(10, -1, 2)
	list2 := NewSliceListForRange(1, 10, 1)

	/// When & Then
	if length := list1.Length(); length != 0 {
		t.Errorf("Should not have any element")
	}

	if contains := list2.Contains(10); contains {
		t.Errorf("Should not contain element")
	}

	if length := list2.Length(); length != 9 {
		t.Errorf("Wrong element count")
	}
}

package gocollection

import (
	"fmt"
	"strconv"
	"testing"
)

type removable struct {
	index int
}

func (rm *removable) String() string {
	return strconv.Itoa(rm.index)
}

func testCollectionBasicOps(t *testing.T, c Collection) {
	/// Setup & When & Then
	if contains := c.Contains(0); contains {
		t.Errorf("Should not contain anything")
	}

	add1 := c.Add(1)
	add2 := c.Add(2)
	add3 := c.Add(3)
	add4 := c.Add(4)

	if add1 != 1 || add2 != 1 || add3 != 1 || add4 != 1 {
		t.Errorf("Added wrong element count")
	}

	if contains := c.Contains(1); !contains {
		t.Errorf("Should contain element")
	}

	if c.Length() != 4 {
		t.Errorf("Should have 4 elements")
	}

	if containsAll := c.ContainsAll(4, 3, 2, 1); !containsAll {
		t.Errorf("Should contain all")
	}

	if deleted := c.Remove(2); deleted != 1 || c.Length() != 3 {
		t.Errorf("Should have 3 elements")
	}

	slice := make([]interface{}, 1000)

	for ix := range slice {
		slice[ix] = strconv.Itoa(ix)
	}

	if addAll := c.AddAll(slice...); addAll != 1000 {
		t.Errorf("Added wrong element count")
	}

	if c.Length() != 1003 {
		t.Errorf("Should have 1003 elements")
	}

	if deleted := c.Remove("Not existent"); deleted > 0 {
		t.Errorf("Should not have found element")
	}

	c.Clear()

	if c.Length() != 0 {
		t.Errorf("Should not have any element")
	}

	fmt.Printf("Final collection %v\n", c)
}

func testCollectionAdd(t *testing.T, c Collection) {
	/// Setup
	addCount := 1000
	addAllCount := 1000
	added := 0

	addSlice := make([]interface{}, addCount)
	addAllSlice := make([]interface{}, addAllCount)

	for ix := range addSlice {
		addSlice[ix] = ix
	}

	for ix := range addAllSlice {
		addAllSlice[ix] = -ix
	}

	combinedSlice := append(addSlice, addAllSlice...)
	combinedSlice1 := combinedSlice[addCount/2 : addAllCount+addAllCount/2]
	combinedSlice2 := append(addSlice, "Not existent")

	/// When & Then
	for ix := range addSlice {
		added += c.Add(ix)
	}

	if containsAll := c.ContainsAll(addSlice...); !containsAll {
		t.Errorf("Should contain all")
	}

	added += c.AddAll(addAllSlice...)

	if containsAll := c.ContainsAll(combinedSlice...); !containsAll {
		t.Errorf("Should contain all")
	}

	if containsAll := c.ContainsAll(combinedSlice1...); !containsAll {
		t.Errorf("Should contain all")
	}

	if containsAll := c.ContainsAll(combinedSlice2...); containsAll {
		t.Errorf("Should not contain all")
	}

	if added != addCount+addAllCount {
		t.Errorf("Wrong number of elements added")
	}
}

func testCollectionGetAllFunc(t *testing.T, c Collection) {
	/// Setup
	c.AddAll(1, 2, 3, 4, 5, 6, 7)

	/// When
	results := c.GetAllFunc(func(e interface{}) bool {
		if e, ok := e.(int); ok && e%2 == 0 {
			return true
		}

		return false
	})

	list := NewSliceList(results...)

	/// Then
	if contains := list.ContainsAll(2, 4, 6); !contains {
		t.Errorf("Should contain correct elements")
	}
}

func testCollectionRemove(t *testing.T, c Collection) {
	/// Setup
	addCount := 1000
	addSlice := make([]interface{}, addCount)

	for ix := range addSlice {
		addSlice[ix] = &removable{index: ix}
	}

	addSlice1 := addSlice[addCount/4 : addCount*3/4]

	/// When
	c.AddAll(addSlice...)

	if length := c.Length(); length != addCount {
		t.Errorf("Added wrong element count")
	}

	if removed := c.RemoveAll(addSlice1...); removed != len(addSlice1) {
		t.Errorf("Removed wrong element count")
	}

	if removed := c.RemoveAll(addSlice...); removed != addCount-len(addSlice1) {
		t.Errorf("Removed wrong element count")
	}

	if removed := c.RemoveAll(addSlice...); removed != 0 {
		t.Errorf("Should not remove anything")
	}
}

func testCollectionAllOps(t *testing.T, colFn func() Collection) {
	testCollectionBasicOps(t, colFn())
	testCollectionAdd(t, colFn())
	testCollectionGetAllFunc(t, colFn())
	testCollectionRemove(t, colFn())
}

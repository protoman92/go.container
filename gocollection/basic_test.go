package gocollection

import (
	"strconv"
	"testing"
)

func testCollectionBasicOps(t *testing.T, c Collection) {
	/// Setup & When & Then
	add1 := c.Add(1)
	add2 := c.Add(2)
	add3 := c.Add(3)
	add4 := c.Add(4)

	if add1 != 1 || add2 != 1 || add3 != 1 || add4 != 1 {
		t.Errorf("Added wrong element count")
	}

	if c.Length() != 4 {
		t.Errorf("Should have 4 elements")
	}

	deletedFound1 := c.Remove(2)

	if c.Length() != 3 || !deletedFound1 {
		t.Errorf("Should have 3 elements")
	}

	slice := make([]Element, 1000)

	for ix := range slice {
		slice[ix] = strconv.Itoa(ix)
	}

	addAll1 := c.AddAll(slice...)

	if addAll1 != 1000 {
		t.Errorf("Added wrong element count")
	}

	if c.Length() != 1003 {
		t.Errorf("Should have 1003 elements")
	}

	deletedFound2 := c.Remove("Not existent")

	if deletedFound2 {
		t.Errorf("Should not have found element")
	}

	c.Clear()

	if c.Length() != 0 {
		t.Errorf("Should not have any element")
	}
}

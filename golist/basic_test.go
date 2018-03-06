package golist

import (
	"strconv"
	"testing"
)

func TestBasicList(t *testing.T) {
	/// Setup
	list := NewBasicList()

	/// When & Then
	list.Add(1)
	list.Add(2)
	list.Add(3)
	list.Add(4)

	if list.Length() != 4 {
		t.Errorf("Should have 4 elements")
	}

	deleted := list.Remove(2)

	if list.Length() != 3 || !deleted {
		t.Errorf("Should have 3 elements")
	}

	slice := make([]Element, 1000)

	for ix := range slice {
		slice[ix] = strconv.Itoa(ix)
	}

	list.AddAll(slice...)

	if list.Length() != 1003 {
		t.Errorf("Should have 1003 elements")
	}
}

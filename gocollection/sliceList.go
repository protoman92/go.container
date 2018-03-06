package gocollection

import (
	"fmt"
)

type sliceList struct {
	slice []Element
}

func (sl *sliceList) String() string {
	return fmt.Sprint(sl.slice)
}

func (sl *sliceList) Add(element Element) int {
	sl.slice = append(sl.slice, element)
	return 1
}

func (sl *sliceList) AddAll(elements ...Element) int {
	for _, element := range elements {
		sl.slice = append(sl.slice, element)
	}

	return len(elements)
}

func (sl *sliceList) Clear() {
	sl.slice = make([]Element, 0)
}

func (sl *sliceList) Contains(element Element) bool {
	for _, e := range sl.slice {
		if e == element {
			return true
		}
	}

	return false
}

func (sl *sliceList) Get(index int) (Element, bool) {
	if index >= 0 && index < len(sl.slice) {
		return sl.slice[index], true
	}

	return nil, false
}

func (sl *sliceList) Length() int {
	return len(sl.slice)
}

func (sl *sliceList) Remove(element Element) bool {
	for ix := range sl.slice {
		e := sl.slice[ix]

		if e == element {
			slice1 := make([]Element, ix)
			copy(slice1, sl.slice[:ix])

			for jx := ix + 1; jx < len(sl.slice); jx++ {
				slice1 = append(slice1, sl.slice[jx])
			}

			sl.slice = slice1
			return true
		}
	}

	return false
}

// NewSliceList returns a new SliceList.
func NewSliceList() List {
	return &sliceList{slice: make([]Element, 0)}
}

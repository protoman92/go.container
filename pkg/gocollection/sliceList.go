package gocollection

import (
	"fmt"
	"sort"
)

type sliceList struct {
	slice []interface{}
}

func (sl *sliceList) String() string {
	return fmt.Sprint(sl.slice)
}

func (sl *sliceList) Add(element interface{}) int {
	sl.slice = append(sl.slice, element)
	return 1
}

func (sl *sliceList) AddAll(elements ...interface{}) int {
	for _, element := range elements {
		sl.slice = append(sl.slice, element)
	}

	return len(elements)
}

func (sl *sliceList) Clear() {
	sl.slice = make([]interface{}, 0)
}

func (sl *sliceList) Contains(element interface{}) bool {
	for ix := range sl.slice {
		if sl.slice[ix] == element {
			return true
		}
	}

	return false
}

func (sl *sliceList) ContainsAll(elements ...interface{}) bool {
	tempMap := make(map[interface{}]bool, len(sl.slice))

	for ix := range sl.slice {
		element := sl.slice[ix]

		if !tempMap[element] {
			tempMap[element] = true
		}
	}

	for ix := range elements {
		if _, found := tempMap[elements[ix]]; !found {
			return false
		}
	}

	return true
}

func (sl *sliceList) GetAt(index int) (interface{}, bool) {
	if index >= 0 && index < len(sl.slice) {
		return sl.slice[index], true
	}

	return nil, false
}

func (sl *sliceList) IndexOf(element interface{}) (int, bool) {
	for ix := range sl.slice {
		if sl.slice[ix] == element {
			return ix, true
		}
	}

	return -1, false
}

func (sl *sliceList) Length() int {
	return len(sl.slice)
}

func (sl *sliceList) Remove(element interface{}) int {
	for ix := range sl.slice {
		e := sl.slice[ix]

		if e == element {
			slice1 := make([]interface{}, ix)
			copy(slice1, sl.slice[:ix])

			for jx := ix + 1; jx < len(sl.slice); jx++ {
				slice1 = append(slice1, sl.slice[jx])
			}

			sl.slice = slice1
			return 1
		}
	}

	return 0
}

func (sl *sliceList) RemoveAt(index int) (interface{}, bool) {
	e, found := sl.GetAt(index)

	if found {
		length := sl.Length()
		first := make([]interface{}, index)
		second := make([]interface{}, length-index-1)
		copy(first, sl.slice[:index])
		copy(second, sl.slice[index+1:])
		first = append(first, second...)
		sl.slice = first
		return e, found
	}

	return nil, false
}

func (sl *sliceList) RemoveAllAt(indexes ...int) int {
	sort.Ints(indexes)
	length := sl.Length()
	lastIndex := 0
	newSlice := make([]interface{}, 0)

	for ix := range indexes {
		index := indexes[ix]

		if index < 0 {
			continue
		} else if index < length {
			newSlice = append(newSlice, sl.slice[lastIndex:index]...)
			lastIndex = index + 1
		} else {
			break
		}
	}

	newSlice = append(newSlice, sl.slice[lastIndex:length]...)
	sl.slice = newSlice
	return length - len(newSlice)
}

func (sl *sliceList) RemoveAll(elements ...interface{}) int {
	tempMap := make(map[interface{}][]int, 0)

	for ix := range sl.slice {
		element := sl.slice[ix]
		indexSlice, found := tempMap[element]

		if !found {
			indexSlice = make([]int, 0)
		}

		indexSlice = append(indexSlice, ix)
		tempMap[element] = indexSlice
	}

	removables := make([]int, 0)

	for ix := range elements {
		if indexes, found := tempMap[elements[ix]]; found && indexes != nil {
			removables = append(removables, indexes...)
		}
	}

	return sl.RemoveAllAt(removables...)
}

func (sl *sliceList) SetAt(index int, element interface{}) (interface{}, bool) {
	prev, found := sl.GetAt(index)

	if found {
		sl.slice[index] = element
		return prev, found
	}

	return nil, false
}

// NewSliceList returns a new SliceList with some default data. Note that the
// data will be copied before storage.
func NewSliceList(elements ...interface{}) List {
	slice1 := make([]interface{}, len(elements))
	copy(slice1, elements)
	return &sliceList{slice: slice1}
}

// NewDefaultSliceList returns a new default SliceList.
func NewDefaultSliceList() List {
	return &sliceList{slice: make([]interface{}, 0)}
}

// NewSliceListForRange returns a new SliceList containing all integers lying
// between an inclusive lower bound and an exclusive upper bound.
func NewSliceListForRange(inclusive int, exclusive int, step int) List {
	if exclusive < inclusive {
		exclusive = inclusive
	}

	slice1 := make([]interface{}, 0)

	for i := inclusive; i < exclusive; i += step {
		slice1 = append(slice1, i)
	}

	return NewSliceList(slice1...)
}

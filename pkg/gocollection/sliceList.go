package gocollection

import (
	"fmt"
	"sort"
)

type sliceList struct {
	slice []interface{}
}

func (l *sliceList) String() string {
	return fmt.Sprint(l.slice)
}

func (l *sliceList) Add(element interface{}) int {
	l.slice = append(l.slice, element)
	return 1
}

func (l *sliceList) AddAll(elements ...interface{}) int {
	for _, element := range elements {
		l.slice = append(l.slice, element)
	}

	return len(elements)
}

func (l *sliceList) Clear() {
	l.slice = make([]interface{}, 0)
}

func (l *sliceList) Contains(element interface{}) bool {
	for ix := range l.slice {
		if l.slice[ix] == element {
			return true
		}
	}

	return false
}

func (l *sliceList) ContainsAll(elements ...interface{}) bool {
	tempMap := make(map[interface{}]bool, len(l.slice))

	for ix := range l.slice {
		element := l.slice[ix]

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

func (l *sliceList) GetAllFunc(selector func(interface{}) bool) []interface{} {
	results := make([]interface{}, 0)

	for ix := range l.slice {
		if selector(l.slice[ix]) {
			results = append(results, l.slice[ix])
		}
	}

	return results
}

func (l *sliceList) GetFirst() (interface{}, bool) {
	return l.GetAt(0)
}

func (l *sliceList) GetFirstFunc(selector func(interface{}) bool) (interface{}, bool) {
	for ix := range l.slice {
		if selector(l.slice[ix]) {
			return l.slice[ix], true
		}
	}

	return nil, false
}

func (l *sliceList) GetAt(index int) (interface{}, bool) {
	if index >= 0 && index < len(l.slice) {
		return l.slice[index], true
	}

	return nil, false
}

func (l *sliceList) IndexOf(element interface{}) (int, bool) {
	for ix := range l.slice {
		if l.slice[ix] == element {
			return ix, true
		}
	}

	return -1, false
}

func (l *sliceList) Length() int {
	return len(l.slice)
}

func (l *sliceList) Remove(element interface{}) int {
	for ix := range l.slice {
		e := l.slice[ix]

		if e == element {
			slice1 := make([]interface{}, ix)
			copy(slice1, l.slice[:ix])

			for jx := ix + 1; jx < len(l.slice); jx++ {
				slice1 = append(slice1, l.slice[jx])
			}

			l.slice = slice1
			return 1
		}
	}

	return 0
}

func (l *sliceList) RemoveAt(index int) (interface{}, bool) {
	e, found := l.GetAt(index)

	if found {
		length := l.Length()
		first := make([]interface{}, index)
		second := make([]interface{}, length-index-1)
		copy(first, l.slice[:index])
		copy(second, l.slice[index+1:])
		first = append(first, second...)
		l.slice = first
		return e, found
	}

	return nil, false
}

func (l *sliceList) RemoveAllAt(indexes ...int) int {
	sort.Ints(indexes)
	length := l.Length()
	lastIndex := 0
	newSlice := make([]interface{}, 0)

	for ix := range indexes {
		index := indexes[ix]

		if index < 0 {
			continue
		} else if index < length {
			newSlice = append(newSlice, l.slice[lastIndex:index]...)
			lastIndex = index + 1
		} else {
			break
		}
	}

	newSlice = append(newSlice, l.slice[lastIndex:length]...)
	l.slice = newSlice
	return length - len(newSlice)
}

func (l *sliceList) RemoveAll(elements ...interface{}) int {
	tempMap := make(map[interface{}][]int, 0)

	for ix := range l.slice {
		element := l.slice[ix]
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

	return l.RemoveAllAt(removables...)
}

func (l *sliceList) SetAt(index int, element interface{}) (interface{}, bool) {
	prev, found := l.GetAt(index)

	if found {
		l.slice[index] = element
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

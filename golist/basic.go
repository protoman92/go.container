package golist

type basicList struct {
	slice []Element
}

func (bl *basicList) Add(element Element) {
	bl.slice = append(bl.slice, element)
}

func (bl *basicList) AddAll(elements ...Element) {
	for _, element := range elements {
		bl.slice = append(bl.slice, element)
	}
}

func (bl *basicList) Clear() {
	bl.slice = make([]Element, 0)
}

func (bl *basicList) Get(index int) (Element, bool) {
	if index >= 0 && index < len(bl.slice) {
		return bl.slice[index], true
	}

	return nil, false
}

func (bl *basicList) Length() int {
	return len(bl.slice)
}

func (bl *basicList) Remove(element Element) bool {
	for ix, e := range bl.slice {
		if e == element {
			slice1 := bl.slice[:ix]

			for jx := ix + 1; jx < len(bl.slice); jx++ {
				slice1 = append(slice1, bl.slice[jx])
			}

			bl.slice = slice1
			return true
		}
	}

	return false
}

// NewBasicList returns a new BasicList.
func NewBasicList() List {
	return &basicList{slice: make([]Element, 0)}
}

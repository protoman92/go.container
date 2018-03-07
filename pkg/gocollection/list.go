package gocollection

// This represents the minimal interface with additional List methods.
type listExtra interface {
	GetAt(index int) (interface{}, bool)
	IndexOf(element interface{}) (int, bool)
	RemoveAt(index int) (interface{}, bool)
	RemoveAllAt(indexes ...int) int
	SetAt(index int, element interface{}) (interface{}, bool)
}

// List represents an indexed Collection of elements.
type List interface {
	Collection
	listExtra
}

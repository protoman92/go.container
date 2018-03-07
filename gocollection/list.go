package gocollection

// This represents the minimal interface with additional List methods.
type listExtra interface {
	GetAt(index int) (Element, bool)
	RemoveAt(index int) (Element, bool)
	RemoveAllAt(indexes ...int) int
	SetAt(index int, element Element) (Element, bool)
}

// List represents an indexed Collection of elements.
type List interface {
	Collection
	listExtra
}

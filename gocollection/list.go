package gocollection

// This represents the minimal interface with additional List methods.
type listExtra interface {
	Get(index int) (Element, bool)
}

// List represents an indexed Collection of elements.
type List interface {
	Collection
	listExtra
}

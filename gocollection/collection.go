package gocollection

// Element represents a Collection element.
type Element interface{}

// Collection represents a collection of elements which may be ordered (list) or
// unordered (set).
type Collection interface {
	Add(element Element) int
	AddAll(elements ...Element) int
	Contains(element Element) bool
	ContainsAll(elements ...Element) bool
	Clear()
	Length() int
	Remove(element Element) int
	RemoveAll(elements ...Element) int
}

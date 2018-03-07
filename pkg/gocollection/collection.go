package gocollection

// Collection represents a collection of elements which may be ordered (list) or
// unordered (set).
type Collection interface {
	Add(element interface{}) int
	AddAll(elements ...interface{}) int
	Contains(element interface{}) bool
	ContainsAll(elements ...interface{}) bool
	Clear()
	Length() int
	Remove(element interface{}) int
	RemoveAll(elements ...interface{}) int
}

package golist

// Element represents a List element.
type Element interface{}

// List represents a indexed collection of elements.
type List interface {
	Add(element Element) int
	AddAll(elements ...Element) int
	Contains(element Element) bool
	Clear()
	Get(index int) (Element, bool)
	Length() int
	Remove(element Element) bool
}

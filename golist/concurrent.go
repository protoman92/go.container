package golist

// ConcurrentList represents a thread-safe List.
type ConcurrentList interface {
	List
	AddAsync(element Element, callback func(int))
	AddAllAsync(callback func(int), elements ...Element)
	ClearAsync(callback func())
	ContainsAsync(element Element, callback func(bool))
	GetAsync(index int, callback func(Element, bool))
	LengthAsync(callback func(int))
	RemoveAsync(element Element, callback func(bool))
}

type concurrentList struct {
	List
}

func (cl *concurrentList) AddAsync(element Element, callback func(int)) {
	go func() {
		added := cl.Add(element)
		callback(added)
	}()
}

func (cl *concurrentList) AddAllAsync(callback func(int), elements ...Element) {
	go func() {
		added := cl.AddAll(elements...)
		callback(added)
	}()
}

func (cl *concurrentList) ClearAsync(callback func()) {
	go func() {
		cl.Clear()
		callback()
	}()
}

func (cl *concurrentList) ContainsAsync(element Element, callback func(bool)) {
	go func() {
		contains := cl.Contains(element)
		callback(contains)
	}()
}

func (cl *concurrentList) GetAsync(index int, callback func(Element, bool)) {
	go func() {
		e, found := cl.Get(index)
		callback(e, found)
	}()
}

func (cl *concurrentList) LengthAsync(callback func(int)) {
	go func() {
		length := cl.Length()
		callback(length)
	}()
}

func (cl *concurrentList) RemoveAsync(element Element, callback func(bool)) {
	go func() {
		found := cl.Remove(element)
		callback(found)
	}()
}

func newConcurrentList(list List) ConcurrentList {
	return &concurrentList{List: list}
}

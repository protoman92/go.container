package gocollection

import (
	"fmt"
)

// ConcurrentCollection represents a thread-safe Collection.
type ConcurrentCollection interface {
	Collection
	AddAsync(element Element, callback func(int))
	AddAllAsync(callback func(int), elements ...Element)
	ClearAsync(callback func())
	ContainsAsync(element Element, callback func(bool))
	LengthAsync(callback func(int))
	RemoveAsync(element Element, callback func(int))
}

type concurrentCollection struct {
	Collection
}

func (cc *concurrentCollection) String() string {
	return fmt.Sprint(cc.Collection)
}

func (cc *concurrentCollection) AddAsync(element Element, callback func(int)) {
	go func() {
		added := cc.Add(element)
		callback(added)
	}()
}

func (cc *concurrentCollection) AddAllAsync(callback func(int), elements ...Element) {
	go func() {
		added := cc.AddAll(elements...)
		callback(added)
	}()
}

func (cc *concurrentCollection) ClearAsync(callback func()) {
	go func() {
		cc.Clear()
		callback()
	}()
}

func (cc *concurrentCollection) ContainsAsync(element Element, callback func(bool)) {
	go func() {
		contains := cc.Contains(element)
		callback(contains)
	}()
}

func (cc *concurrentCollection) LengthAsync(callback func(int)) {
	go func() {
		length := cc.Length()
		callback(length)
	}()
}

func (cc *concurrentCollection) RemoveAsync(element Element, callback func(int)) {
	go func() {
		found := cc.Remove(element)
		callback(found)
	}()
}

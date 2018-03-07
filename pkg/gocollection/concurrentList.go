package gocollection

import (
	"fmt"
)

// ConcurrentList represents a thread-safe List.
type ConcurrentList interface {
	ConcurrentCollection
	listExtra
	GetAtAsync(index int, callback func(interface{}, bool))
	RemoveAtAsync(index int, callback func(interface{}, bool))
	RemoveAllAtAsync(callback func(int), indexes ...int)
	SetAtAsync(index int, element interface{}, callback func(interface{}, bool))
}

type concurrentList struct {
	ConcurrentCollection
	listExtra
}

func (cl *concurrentList) String() string {
	return fmt.Sprint(cl.listExtra)
}

func (cl *concurrentList) GetAtAsync(index int, callback func(interface{}, bool)) {
	go func() {
		e, found := cl.GetAt(index)
		callback(e, found)
	}()
}

func (cl *concurrentList) RemoveAtAsync(index int, callback func(interface{}, bool)) {
	go func() {
		e, found := cl.RemoveAt(index)
		callback(e, found)
	}()
}

func (cl *concurrentList) RemoveAllAtAsync(callback func(int), indexes ...int) {
	go func() {
		removed := cl.listExtra.RemoveAllAt(indexes...)
		callback(removed)
	}()
}

func (cl *concurrentList) SetAtAsync(index int, element interface{}, callback func(interface{}, bool)) {
	go func() {
		prev, found := cl.SetAt(index, element)
		callback(prev, found)
	}()
}

func newConcurrentList(list List) ConcurrentList {
	collection := &concurrentCollection{Collection: list}
	return &concurrentList{ConcurrentCollection: collection, listExtra: list}
}

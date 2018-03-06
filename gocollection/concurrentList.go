package gocollection

import (
	"fmt"
)

// ConcurrentList represents a thread-safe List.
type ConcurrentList interface {
	ConcurrentCollection
	listExtra
	GetAsync(index int, callback func(Element, bool))
}

type concurrentList struct {
	ConcurrentCollection
	listExtra
}

func (cl *concurrentList) String() string {
	return fmt.Sprint(cl.listExtra)
}

func (cl *concurrentList) GetAsync(index int, callback func(Element, bool)) {
	go func() {
		e, found := cl.listExtra.Get(index)
		callback(e, found)
	}()
}

func newConcurrentList(list List) ConcurrentList {
	collection := &concurrentCollection{Collection: list}
	return &concurrentList{ConcurrentCollection: collection, listExtra: list}
}

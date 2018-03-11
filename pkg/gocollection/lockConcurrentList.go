package gocollection

import (
	"fmt"
	"sync"
)

type lockConcurrentList struct {
	*lockConcurrentCollection
	mutex *sync.RWMutex
	list  listExtra
}

// Do not return String() for the lock concurrent collection because these
// share the same mutex.
func (l *lockConcurrentList) String() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return fmt.Sprint(l.list)
}

func (l *lockConcurrentList) GetAt(index int) (interface{}, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	e, found := l.list.GetAt(index)
	return e, found
}

func (l *lockConcurrentList) GetFirst() (interface{}, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.list.GetFirst()
}

func (l *lockConcurrentList) GetFirstFunc(selector func(interface{}) bool) (interface{}, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.list.GetFirstFunc(selector)
}

func (l *lockConcurrentList) IndexOf(element interface{}) (int, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	index, found := l.list.IndexOf(element)
	return index, found
}

func (l *lockConcurrentList) RemoveAt(index int) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	e, found := l.list.RemoveAt(index)
	return e, found
}

func (l *lockConcurrentList) RemoveAllAt(indexes ...int) int {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.list.RemoveAllAt(indexes...)
}

func (l *lockConcurrentList) SetAt(index int, element interface{}) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.list.SetAt(index, element)
}

// NewLockConcurrentList returns a new LockConcurrentList.
func NewLockConcurrentList(list List) List {
	mutex := &sync.RWMutex{}
	lcc := &lockConcurrentCollection{mutex: mutex, storage: list}

	l := &lockConcurrentList{
		lockConcurrentCollection: lcc,
		list:  list,
		mutex: mutex,
	}

	return l
}

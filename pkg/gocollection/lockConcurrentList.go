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
func (lcl *lockConcurrentList) String() string {
	lcl.mutex.RLock()
	defer lcl.mutex.RUnlock()
	return fmt.Sprint(lcl.list)
}

func (lcl *lockConcurrentList) GetAt(index int) (interface{}, bool) {
	lcl.mutex.RLock()
	defer lcl.mutex.RUnlock()
	e, found := lcl.list.GetAt(index)
	return e, found
}

func (lcl *lockConcurrentList) IndexOf(element interface{}) (int, bool) {
	lcl.mutex.RLock()
	defer lcl.mutex.RUnlock()
	index, found := lcl.list.IndexOf(element)
	return index, found
}

func (lcl *lockConcurrentList) RemoveAt(index int) (interface{}, bool) {
	lcl.mutex.Lock()
	defer lcl.mutex.Unlock()
	e, found := lcl.list.RemoveAt(index)
	return e, found
}

func (lcl *lockConcurrentList) RemoveAllAt(indexes ...int) int {
	lcl.mutex.Lock()
	defer lcl.mutex.Unlock()
	return lcl.list.RemoveAllAt(indexes...)
}

func (lcl *lockConcurrentList) SetAt(index int, element interface{}) (interface{}, bool) {
	lcl.mutex.Lock()
	defer lcl.mutex.Unlock()
	return lcl.list.SetAt(index, element)
}

// NewLockConcurrentList returns a new LockConcurrentList.
func NewLockConcurrentList(list List) List {
	mutex := &sync.RWMutex{}
	lcc := &lockConcurrentCollection{mutex: mutex, storage: list}

	lcl := &lockConcurrentList{
		lockConcurrentCollection: lcc,
		list:  list,
		mutex: mutex,
	}

	return lcl
}

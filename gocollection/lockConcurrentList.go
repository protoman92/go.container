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

func (lcl *lockConcurrentList) String() string {
	return fmt.Sprint(lcl.list)
}

func (lcl *lockConcurrentList) Get(index int) (Element, bool) {
	lcl.mutex.RLock()
	defer lcl.mutex.RUnlock()
	e, found := lcl.list.Get(index)
	return e, found
}

// NewLockConcurrentList returns a new LockConcurrentList.
func NewLockConcurrentList(list List) ConcurrentList {
	mutex := &sync.RWMutex{}
	lcc := &lockConcurrentCollection{mutex: mutex, storage: list}

	lcl := &lockConcurrentList{
		lockConcurrentCollection: lcc,
		list:  list,
		mutex: mutex,
	}

	return newConcurrentList(lcl)
}

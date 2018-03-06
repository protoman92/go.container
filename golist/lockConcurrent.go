package golist

import "sync"

type lockConcurrentList struct {
	mutex   *sync.RWMutex
	storage List
}

func (lcl *lockConcurrentList) Add(element Element) int {
	lcl.mutex.Lock()
	defer lcl.mutex.Unlock()
	return lcl.storage.Add(element)
}

func (lcl *lockConcurrentList) AddAll(elements ...Element) int {
	lcl.mutex.Lock()
	defer lcl.mutex.Unlock()
	return lcl.storage.AddAll(elements...)
}

func (lcl *lockConcurrentList) Clear() {
	lcl.mutex.Lock()
	defer lcl.mutex.Unlock()
	lcl.storage.Clear()
}

func (lcl *lockConcurrentList) Contains(element Element) bool {
	lcl.mutex.RLock()
	defer lcl.mutex.RUnlock()
	return lcl.storage.Contains(element)
}

func (lcl *lockConcurrentList) Get(index int) (Element, bool) {
	lcl.mutex.RLock()
	defer lcl.mutex.RUnlock()
	e, found := lcl.storage.Get(index)
	return e, found
}

func (lcl *lockConcurrentList) Length() int {
	lcl.mutex.RLock()
	defer lcl.mutex.RUnlock()
	return lcl.storage.Length()
}

func (lcl *lockConcurrentList) Remove(element Element) bool {
	lcl.mutex.Lock()
	defer lcl.mutex.Unlock()
	return lcl.storage.Remove(element)
}

// NewLockConcurrentList returns a new LockConcurrentList.
func NewLockConcurrentList(list List) ConcurrentList {
	lcl := &lockConcurrentList{mutex: &sync.RWMutex{}, storage: list}
	return newConcurrentList(lcl)
}

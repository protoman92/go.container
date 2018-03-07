package gocollection

import (
	"sync"
)

type lockConcurrentCollection struct {
	mutex   *sync.RWMutex
	storage Collection
}

func (lcc *lockConcurrentCollection) Add(element interface{}) int {
	lcc.mutex.Lock()
	defer lcc.mutex.Unlock()
	return lcc.storage.Add(element)
}

func (lcc *lockConcurrentCollection) AddAll(elements ...interface{}) int {
	lcc.mutex.Lock()
	defer lcc.mutex.Unlock()
	return lcc.storage.AddAll(elements...)
}

func (lcc *lockConcurrentCollection) Clear() {
	lcc.mutex.Lock()
	defer lcc.mutex.Unlock()
	lcc.storage.Clear()
}

func (lcc *lockConcurrentCollection) Contains(element interface{}) bool {
	lcc.mutex.RLock()
	defer lcc.mutex.RUnlock()
	return lcc.storage.Contains(element)
}

func (lcc *lockConcurrentCollection) ContainsAll(elements ...interface{}) bool {
	lcc.mutex.RLock()
	defer lcc.mutex.RUnlock()
	return lcc.storage.ContainsAll(elements...)
}

func (lcc *lockConcurrentCollection) Length() int {
	lcc.mutex.RLock()
	defer lcc.mutex.RUnlock()
	return lcc.storage.Length()
}

func (lcc *lockConcurrentCollection) Remove(element interface{}) int {
	lcc.mutex.Lock()
	defer lcc.mutex.Unlock()
	return lcc.storage.Remove(element)
}

func (lcc *lockConcurrentCollection) RemoveAll(elements ...interface{}) int {
	lcc.mutex.Lock()
	defer lcc.mutex.Unlock()
	return lcc.storage.RemoveAll(elements...)
}

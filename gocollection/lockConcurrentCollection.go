package gocollection

import "sync"

type lockConcurrentCollection struct {
	mutex   *sync.RWMutex
	storage Collection
}

func (lcc *lockConcurrentCollection) Add(element Element) int {
	lcc.mutex.Lock()
	defer lcc.mutex.Unlock()
	return lcc.storage.Add(element)
}

func (lcc *lockConcurrentCollection) AddAll(elements ...Element) int {
	lcc.mutex.Lock()
	defer lcc.mutex.Unlock()
	return lcc.storage.AddAll(elements...)
}

func (lcc *lockConcurrentCollection) Clear() {
	lcc.mutex.Lock()
	defer lcc.mutex.Unlock()
	lcc.storage.Clear()
}

func (lcc *lockConcurrentCollection) Contains(element Element) bool {
	lcc.mutex.RLock()
	defer lcc.mutex.RUnlock()
	return lcc.storage.Contains(element)
}

func (lcc *lockConcurrentCollection) Length() int {
	lcc.mutex.RLock()
	defer lcc.mutex.RUnlock()
	return lcc.storage.Length()
}

func (lcc *lockConcurrentCollection) Remove(element Element) bool {
	lcc.mutex.Lock()
	defer lcc.mutex.Unlock()
	return lcc.storage.Remove(element)
}

package gomap

import "sync"

type lockConcurrentMap struct {
	mutex   *sync.RWMutex
	storage Map
}

func (lcm *lockConcurrentMap) Clear() {
	lcm.mutex.Lock()
	defer lcm.mutex.Unlock()
	lcm.storage.Clear()
}

func (lcm *lockConcurrentMap) Contains(key Key) bool {
	lcm.mutex.RLock()
	defer lcm.mutex.RUnlock()
	return lcm.storage.Contains(key)
}

func (lcm *lockConcurrentMap) Delete(key Key) int {
	lcm.mutex.Lock()
	defer lcm.mutex.Unlock()
	return lcm.storage.Delete(key)
}

func (lcm *lockConcurrentMap) Get(key Key) (Value, bool) {
	lcm.mutex.RLock()
	defer lcm.mutex.RUnlock()
	return lcm.storage.Get(key)
}

func (lcm *lockConcurrentMap) IsEmpty() bool {
	lcm.mutex.RLock()
	defer lcm.mutex.RUnlock()
	return lcm.storage.IsEmpty()
}

func (lcm *lockConcurrentMap) Length() int {
	lcm.mutex.RLock()
	defer lcm.mutex.RUnlock()
	return lcm.storage.Length()
}

func (lcm *lockConcurrentMap) Set(key Key, value Value) int {
	lcm.mutex.Lock()
	defer lcm.mutex.Unlock()
	return lcm.storage.Set(key, value)
}

func (lcm *lockConcurrentMap) UnderlyingStorage() map[Key]Value {
	lcm.mutex.RLock()
	defer lcm.mutex.RUnlock()
	return lcm.storage.UnderlyingStorage()
}

// NewLockConcurrentMap returns a new lock-based ConcurrentMap.
func NewLockConcurrentMap(storage Map) ConcurrentMap {
	cm := &lockConcurrentMap{mutex: &sync.RWMutex{}, storage: storage}
	return newConcurrentMap(cm)
}
